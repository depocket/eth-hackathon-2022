package model

import (
	"github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	FromAddress   string         `json:"from_address" gorm:"index:transactions_from_idx;type:varchar(42)"`
	ToAddress     string         `json:"to_address" gorm:"index:transactions_to_idx;type:varchar(42)"`
	Hash          string         `gorm:"primaryKey;unique;type:varchar(66)"`
	Action        string         `json:"action" gorm:"index:transactions_action_idx;type:varchar(100)"`
	DecodedAction string         `json:"decoded_action" gorm:"index:transactions_decoded_action_idx;type:varchar(200)"`
	Input         string         `json:"input"`
	DecodedInput  postgres.Jsonb `json:"decoded_input,omitempty"  gorm:"type:jsonb"`
	BlockNumber   uint64         `json:"block_number" gorm:"index:transactions_block_number_idx"`
	Gas           uint64         `json:"gas"`
	GasPrice      uint64         `json:"gas_price"`
	GasFee        string         `json:"gas_fee"`
	Index         int            `gorm:"default: 0"`
	Type          int            `json:"type"`
	Value         int64          `json:"value"`
	Status        bool           `json:"status" gorm:"default: false"`
	Time          time.Time      `gorm:"type:time;index"`
}

type BaseModelID struct {
	ID uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
}

type BaseModelTime struct {
	CreatedAt time.Time       `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type Token struct {
	Address  string `gorm:"index;unique;not null"`
	Symbol   string `json:"symbol"`
	Decimals uint64 `json:"decimals"`
}

type TransactionResponse struct {
	Data           []Transaction `json:"data"`
	HasNextPage    bool          `json:"has_next_page"`
	NextPageCursor string        `json:"next_page_cursor"`
}
type TokenResponse struct {
	Data      []Token `json:"data"`
	ErrorCode int     `json:"error_code"`
}
