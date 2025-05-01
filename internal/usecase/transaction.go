package usecase

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Fadil-Tao/gopher-pay/internal/model"
	csterr "github.com/Fadil-Tao/gopher-pay/utils/errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionRepo interface {
	WriteTransaction(ctx context.Context, transaction *model.Transaction, updateFn func() error) error
}

type TransactionUsecase struct {
	TransactionRepo
	WalletRepo
	UserRepo
}

func NewTransactionUsecase(transactionRepo TransactionRepo, walletRepo WalletRepo) *TransactionUsecase {
	return &TransactionUsecase{
		TransactionRepo: transactionRepo,
		WalletRepo:      walletRepo,
	}
}

func (t *TransactionUsecase) Transfer(ctx context.Context, userId int, phoneDestination string, amountStr string, description string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		slog.Error(err.Error())
		return csterr.ErrInternal
	}

	if amount.LessThanOrEqual(decimal.NewFromInt(10000)) {
		return errors.New("minimum amount needed is above 10000")
	}

	// check the user balance
	userWallet, err := t.WalletRepo.GetWalletInfoByUserId(ctx, userId)
	if err != nil {
		slog.Error("error initializing wallet repo", "message", err)
		return err
	}

	if userWallet.Balance.LessThanOrEqual(amount.Add(decimal.NewFromInt(30000))) {
		slog.Error("insufficient balance")
		return csterr.ErrInsufficientBalance
	}

	// get the source wallet info
	srcWallet, err := t.WalletRepo.GetWalletInfoByUserId(ctx, userId)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	// get the destination info
	destWallet, err := t.WalletRepo.GetWalletInfoByUserPhone(ctx, phoneDestination)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	// creaate new transaction data
	transaction := model.Transaction{
		Id:                  uuid.New(),
		SourceWalletId:      srcWallet.Id,
		DestinationWalletId: destWallet.Id,
		Amount:              amount,
		TransactionType:     model.Transfer,
		Description:         &description,
	}

	slog.Info("src wallet id :", "id",transaction.SourceWalletId)
	slog.Info("dest wallet id :", "id",transaction.DestinationWalletId)


	err = t.TransactionRepo.WriteTransaction(ctx, &transaction, func() error {
		// substract the sender balance
		err := t.WalletRepo.Substract(ctx, srcWallet.Id, amount)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		// add the receiveir balance
		err = t.WalletRepo.Add(ctx, destWallet.Id, amount)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}