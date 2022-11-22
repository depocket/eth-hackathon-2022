package repo

import (
	"context"
	"depocket.io/app/model"
	"depocket.io/app/utils"
	"encoding/json"
	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
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
