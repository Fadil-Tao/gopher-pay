package usecase

import (
	"context"

	"github.com/Fadil-Tao/gopher-pay/internal/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type WalletRepo interface{
	CreateWallet(ctx context.Context, userId int )error
	GetWalletInfoByUserId(ctx context.Context, userId int)(*model.Wallet, error)
	Add(ctx context.Context, id uuid.UUID, amount decimal.Decimal) error
	Substract(ctx context.Context,id uuid.UUID ,amount decimal.Decimal) error 
	GetWalletInfoByUserPhone(ctx context.Context, userPhone string)(*model.Wallet, error)
}	

type WalletUseCase struct {
	WalletRepo
}

func NewWalletUsecase(walletRepo WalletRepo) *WalletUseCase{
	return &WalletUseCase{
		WalletRepo:  walletRepo,
	}
}