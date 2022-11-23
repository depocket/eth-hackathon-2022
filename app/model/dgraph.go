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
	DType   []string `json:"dgraph.type,omitempty"`
}

type TransactionDgraph struct {
	UID          string        `json:"uid"`
	Amount       float64       `json:"amount"`
	Sender       AddressDgraph `json:"sender"`
	Recipient    AddressDgraph `json:"recipient"`
	Name         string        `json:"name"`
	TokenAddress string        `json:"token_address"`
	TxnId        string        `json:"txn_id"`
	TxnTime      time.Time     `json:"txn_time"`
	DType        []string      `json:"dgraph.type,omitempty"`
}
