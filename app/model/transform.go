package model

import "time"

type ResponseFlow struct {
	Data []AddressDgraphResponse `json:"data"`
}

type ResponseTransformed struct {
	Data  []AddressDgraphResponse `json:"data"`
	Nodes []Node                  `json:"nodes"`
	Edges []Edge                  `json:"edges"`
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

type Node struct {
	Id    string `json:"id"`
	Label string `json:"label"`
	Title string `json:"title"`
}

type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
}
