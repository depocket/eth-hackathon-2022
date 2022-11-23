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
	Block         Block          `gorm:"foreignKey:BlockNumber"`
}

type Block struct {
	Number       int64         `gorm:"primaryKey;unique"`
	Hash         string        `gorm:"uniqueIndex;type:varchar(66)"`
	ParentHash   string        `gorm:"type:varchar(66)"`
	Time         time.Time     `gorm:"type:time"`
	Sid          string        `gorm:"uniqueIndex"`
	Transactions []Transaction `gorm:"foreignKey:BlockNumber;references:Number"`
}

func (Block) TableName() string {
	return "bsc.blocks"
}

func (Transaction) TableName() string {
	return "bsc.transactions"
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
	BaseModelID
	BaseModelTime
	Address     string  `gorm:"index;unique;not null"`
	Name        string  `json:"name"`
	Symbol      string  `json:"symbol"`
	IconUrl     *string `json:"icon_url"`
	SiteUrl     *string `json:"site_url"`
	ProjectCode string  `json:"project_code"`
	Chain       string  `json:"chain"`
	Type        string  `json:"type"`
	Decimals    uint64  `json:"decimals"`
	Price       float64 `json:"price"`
	TotalSupply uint64  `json:"total_supply"`
	IsCore      bool    `json:"is_core" gorm:"default:false"`
	IsVerified  bool    `json:"is_verified" gorm:"default:true"`
	IsActive    bool    `json:"is_active" gorm:"default:true"`
}

func (Token) TableName() string {
	return "public.tokens"
}
