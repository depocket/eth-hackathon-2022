package model

import (
	"math/big"
	"time"
)

type ResponseAddress struct {
	Addresses []AddressDgraph `json:"addresses"`
}

type ResponseTxn struct {
	Txns []TransactionDgraph `json:"txns"`
}

type ResponseFullFlow struct {
	FullFlow []TransactionDgraph `json:"full_flow"`
}

type FullFlow struct {
	Recipient []TransactionDgraph `json:"~recipient"`
	Sender    []TransactionDgraph `json:"~sender"`
}

type TransferAction struct {
	Amount    big.Int `json:"amount"`
	Method    string  `json:"method"`
	Recipient string  `json:"recipient"`
}

type AddressDgraph struct {
	UID     string   `json:"uid"`
	Address string   `json:"address"`
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Chain   string   `json:"chain"`
	DType   []string `json:"dgraph.type,omitempty"`
}

type TransactionDgraph struct {
	Chain        string        `json:"chain"`
	UID          string        `json:"uid"`
	Amount       big.Float     `json:"amount"`
	Sender       AddressDgraph `json:"sender"`
	Recipient    AddressDgraph `json:"recipient"`
	Name         string        `json:"name"`
	TokenAddress string        `json:"token_address"`
	TxnId        string        `json:"txn_id"`
	TxnTime      time.Time     `json:"txn_time"`
	DType        []string      `json:"dgraph.type,omitempty"`
}

type AddressDgraphResponse struct {
	UID       string                      `json:"uid,omitempty"`
	Address   string                      `json:"address,omitempty"`
	Name      string                      `json:"name,omitempty"`
	Type      string                      `json:"type,omitempty"`
	Recipient []TransactionDgraphResponse `json:"~recipient,omitempty"`
	Sender    []TransactionDgraphResponse `json:"~sender,omitempty"`
}

type TransactionDgraphResponse struct {
	UID          string                `json:"uid,omitempty"`
	Amount       float64               `json:"amount,omitempty"`
	Sender       AddressDgraphResponse `json:"sender,omitempty"`
	Recipient    AddressDgraphResponse `json:"recipient,omitempty"`
	Name         string                `json:"name,omitempty"`
	TokenAddress string                `json:"token_address,omitempty"`
	TxnId        string                `json:"txn_id,omitempty"`
	TxnTime      time.Time             `json:"txn_time,omitempty"`
}

type ResponseFlow struct {
	Data []AddressDgraphResponse `json:"data"`
}

type ResponsePath struct {
	Path []Path       `json:"_path_"`
	Node []NodeDgraph `json:"node"`
}

type NodeDgraph struct {
	UID     string `json:"uid,omitempty"`
	Address string `json:"address,omitempty"`
	Name    string `json:"name,omitempty"`
}

type Path struct {
	Weight    float64                `json:"_weight_"`
	Uid       string                 `json:"uid"`
	Sender    map[string]interface{} `json:"~sender"`
	Recipient map[string]interface{} `json:"~recipient"`
}
