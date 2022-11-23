package model

import "time"

type FlowRequest struct {
	Depth   int       `json:"depth" binding:"required"`
	Token   string    `json:"token" binding:"required"`
	Address string    `json:"address" binding:"required"`
	From    time.Time `json:"from" binding:"required"`
	To      time.Time `json:"to" binding:"required"`
}

type PathRequest struct {
	Path        int    `json:"path" binding:"required"`
	FromAddress string `json:"from_address" binding:"required"`
	ToAddress   string `json:"to_address" binding:"required"`
}
