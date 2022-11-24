package handler

import (
	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MigrationHandler struct {
	dgraph *dgo.Dgraph
	log    *zap.Logger
}

func NewMigrationHandler(r *gin.RouterGroup, log *zap.Logger, dgraph *dgo.Dgraph) {
	m := MigrationHandler{dgraph: dgraph, log: log}
	ar := r.Group("/dgraph")
	ar.GET("/migrate", m.MigrateDgraph)
}

func (h *MigrationHandler) MigrateDgraph(c *gin.Context) {
	err := h.dgraph.Alter(c, &api.Operation{
		Schema: `
			address: string @index(hash) .
			type: string @index(hash) .
			name: string @index(hash) .
			chain: string @index(hash) .
			txn_id: string @index(hash) .
			txn_time: datetime .
			amount: float .
			recipient: uid @reverse .
			sender: uid @reverse .
			symbol: string @index(hash) .
			token_address: string @index(hash) .

			type Address {
				address
				type
				name
				chain
			}

			type Transaction {
				txn_id
				txn_time
				amount	
				symbol
				token_address
				sender
				recipient
				chain
			}
		`,
	})
	if err != nil {
		h.log.Panic(err.Error())
	}
}
