package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionType string

const (
	Withdrawal TransactionType = "withdrawal"
	TopUp TransactionType = "topup"
	Transfer TransactionType = "transfer"
)

type Transaction struct {
	Id uuid.UUID `json:"id"`
	SourceWalletId uuid.UUID `json:"sender,omitempty"`
	DestinationWalletId uuid.UUID `json:"receiver,omitempty"`
	Amount decimal.Decimal `json:"amount"`
	TransactionType TransactionType `json:"transactionType"`
	Description *string `json:"description,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}


type TransferRequest struct {
	AmountStr   string	`json:"amount" validate:"required,amount"`
	Description string	`json:"description,omitempty"`
	PhoneDest   string	`json:"destination" validate:"required"`
}

