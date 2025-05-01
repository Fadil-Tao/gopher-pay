package model

import (
	"time"
	"github.com/shopspring/decimal"
	"github.com/google/uuid"
)

type Wallet struct {
	Id uuid.UUID `json:"id"`
	UserId int `json:"userId,omitempty"`
	Balance decimal.Decimal `json:"balance" validation:"min=10000, max=10000000"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}