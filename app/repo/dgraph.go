package repo

import (
	"context"
	"depocket.io/app/model"
	"depocket.io/app/utils"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"strings"
	"time"
)

type Dgraph struct {
	DB    *dgo.Dgraph
	debug bool
}

func (r *Dgraph) GetRepo() *dgo.Dgraph {
	return r.DB
}

func NewRepoDgraph(db *dgo.Dgraph) *Dgraph {
	return &Dgraph{DB: db}
}

type DgraphInterface interface {
	GetByAddress(ctx context.Context, address string) (string, error)
	GetByTransaction(ctx context.Context, hash string) (string, error)
	CreateNode(ctx context.Context, field, value string, user interface{}) error
	UpdateNode(ctx context.Context, field, value string, user interface{}) error

	OutFlow(depth int, address string, token string, from time.Time, to time.Time) (*model.ResponseFlow, error)
	InFlow(depth int, address string, token string, from time.Time, to time.Time) (*model.ResponseFlow, error)
	FullFlow(depth int, address string, token string, from time.Time, to time.Time) (*model.ResponseFlow, error)
	Path(path int, from string, to string) (interface{}, error)
}

func (r *Dgraph) GetByAddress(ctx context.Context, address string) (string, error) {
	query := `{
		   addresses(func: eq(address, "` + address + `")) {uid}
		}`

	request := &api.Request{
		Query:     query,
		CommitNow: true,
	}

	resp, err := r.DB.NewTxn().Do(ctx, request)
	if err != nil {
		return "", err
	}

	res := model.ResponseAddress{}
	err = json.Unmarshal(resp.Json, &res)
	if err != nil {
		return "", err
	}

	if len(res.Addresses) == 0 {
		return "", utils.ErrRecordNotFound
	}
	return res.Addresses[0].UID, nil
}

func (r *Dgraph) GetByTransaction(ctx context.Context, hash string) (string, error) {
	query := `{
		   addresses(func: eq(txn_id, "` + hash + `")) {uid}
		}`

	request := &api.Request{
		Query:     query,
		CommitNow: true,
	}

	resp, err := r.DB.NewTxn().Do(ctx, request)
	if err != nil {
		return "", err
	}

	res := model.ResponseTxn{}
	err = json.Unmarshal(resp.Json, &res)
	if err != nil {
		return "", err
	}

	if len(res.Txns) == 0 {
		return "", utils.ErrRecordNotFound
	}
	return res.Txns[0].UID, nil
}

func (r *Dgraph) CreateNode(ctx context.Context, field, value string, user interface{}) error {
	pb, err := json.Marshal(user)
	if err != nil {
		return err
	}
	query := `{
		   userUid as var(func: eq(` + field + `, "` + value + `"))
		}`

	var mutations []*api.Mutation
	mutations = append(mutations, &api.Mutation{
		Cond:    `@if(eq(len(userUid), 0))`,
		SetJson: pb,
	})

	request := &api.Request{
		Query:     query,
		Mutations: mutations,
		CommitNow: true,
	}

	_, err = r.DB.NewTxn().Do(ctx, request)
	if err != nil {
		return err
	}
	return nil
}

func (r *Dgraph) UpdateNode(ctx context.Context, field, value string, user interface{}) error {
	pb, err := json.Marshal(user)
	if err != nil {
		return err
	}
	query := `{
			userUid as var(func: eq(` + field + `, "` + value + `")) {uid}
		}`

	var mutations []*api.Mutation
	mutations = append(mutations, &api.Mutation{
		Cond:    `@if(NOT eq(len(userUid),0))`,
		SetJson: pb,
	})

	request := &api.Request{
		Query:     query,
		Mutations: mutations,
		CommitNow: true,
	}

	_, err = r.DB.NewTxn().Do(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func (r *Dgraph) FullFlow(depth int, address string, token string, from time.Time, to time.Time) (*model.ResponseFlow, error) {
	// people who have given money to anyone who has given money to me
	// people who have received money from anyone who has received money from me
	query := `
		{
		  data(func: eq(address,"[ADDRESS]")) @recurse (depth:[DEPTH])  {
			
			~recipient @filter(between(txn_time,"[FROM]","[TO]") AND eq(token_address,"[TOKEN]")) 
			recipient 
			
			~sender  @filter(between(txn_time,"[FROM]","[TO]") AND eq(token_address,"[TOKEN]")) 
			sender

			uid
			txn_id 
			txn_time
			name
		    token_address
			address
		  }
		}
	`
	replacer := strings.NewReplacer(
		"[ADDRESS]", address,
		"[DEPTH]", fmt.Sprintf("%v", depth),
		"[FROM]", from.String(),
		"[TO]", to.String(),
		"[TOKEN]", token,
	)
	query = replacer.Replace(query)

	resp, err := r.DB.NewTxn().Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	res := &model.ResponseFlow{}
	if err := json.Unmarshal(resp.Json, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Dgraph) InFlow(depth int, address string, token string, from time.Time, to time.Time) (*model.ResponseFlow, error) {
	// show me people who have given money to anyone who has given money to me
	query := `
		{
		  data(func: eq(address,"[ADDRESS]")) @recurse (depth:[DEPTH])  {
			~recipient @filter(between(txn_time,"[FROM]","[TO]") AND eq(token_address,"[TOKEN]")) 
			uid
			sender
			txn_id 
			txn_time
			name
		    token_address
			address
		  }
		}
	`
	replacer := strings.NewReplacer(
		"[ADDRESS]", address,
		"[DEPTH]", fmt.Sprintf("%v", depth),
		"[FROM]", from.String(),
		"[TO]", to.String(),
		"[TOKEN]", token,
	)
	query = replacer.Replace(query)

	resp, err := r.DB.NewTxn().Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	res := &model.ResponseFlow{}
	if err := json.Unmarshal(resp.Json, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Dgraph) OutFlow(depth int, address string, token string, from time.Time, to time.Time) (*model.ResponseFlow, error) {
	// show me people who have received money from anyone who has received money from me
	query := `
		{
		  data(func: eq(address,"[ADDRESS]")) @recurse (depth:[DEPTH])  {
			~sender @filter(between(txn_time,"[FROM]","[TO]") AND eq(token_address,"[TOKEN]")) 
			uid
			recipient
			txn_id 
			txn_time
			name
		    token_address
			address
		  }
		}
	`
	replacer := strings.NewReplacer(
		"[ADDRESS]", address,
		"[DEPTH]", fmt.Sprintf("%v", depth),
		"[FROM]", from.String(),
		"[TO]", to.String(),
		"[TOKEN]", token,
	)
	query = replacer.Replace(query)

	resp, err := r.DB.NewTxn().Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	res := &model.ResponseFlow{}
	if err := json.Unmarshal(resp.Json, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Dgraph) Path(path int, from string, to string) (interface{}, error) {
	// show me people who have received money from anyone who has received money from me
	query := ` 
	{
  
	 from as var(func: eq(address, "[FROM_ADDRESS]"))
	 to as var(func: eq(address, "[TO_ADDRESS]"))
	  
	 path as shortest(from: uid(from), to: uid(to),numpaths: [PATH]) {
	  ~recipient
	  ~sender
	  sender
	  recipient
	 }
	  
	 path(func: uid(path)) {
	   uid
	   address
	   name
	 }
	  
	}
	`
	replacer := strings.NewReplacer(
		"[FROM_ADDRESS]", from,
		"[TO_ADDRESS]", to,
		"[PATH]", fmt.Sprintf("%v", path),
	)
	query = replacer.Replace(query)

	resp, err := r.DB.NewTxn().Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	res := map[string]interface{}{}
	if err := json.Unmarshal(resp.Json, &res); err != nil {
		return nil, err
	}
	return res, nil
}
