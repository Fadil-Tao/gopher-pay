package usecase_test

import (
	"context"
	"testing"

	"github.com/Fadil-Tao/gopher-pay/internal/model"
	"github.com/Fadil-Tao/gopher-pay/internal/usecase"
	mocks "github.com/Fadil-Tao/gopher-pay/internal/usecase/mock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)



func TestShouldTransfer(t *testing.T) {
	mockTransactRepo := new(mocks.TransactionRepo)
	mockWalletRepo := new(mocks.WalletRepo)

	transactUsecase := usecase.NewTransactionUsecase(mockTransactRepo, mockWalletRepo)

	ctx := context.Background()

	newTransactReq := struct {
		userId      int
		phoneDest   string
		amountStr   string
		description string
	}{
		userId:      1,
		phoneDest:   "089765432128",
		amountStr:   "11000.32",
		description: "money from gambia",
	}

	amount, _ := decimal.NewFromString(newTransactReq.amountStr)

	srcWalletId := uuid.New()
	destWalletId := uuid.New()

	mockWalletRepo.On("GetWalletInfoByUserId", ctx, newTransactReq.userId).
		Return(&model.Wallet{
			Id:      srcWalletId,
			UserId:  newTransactReq.userId,
			Balance: decimal.NewFromInt(1000000),
		}, nil).Twice()

	mockWalletRepo.On("GetWalletInfoByUserPhone", ctx, newTransactReq.phoneDest).
		Return(&model.Wallet{
			Id:      destWalletId,
			UserId:  2,
			Balance: decimal.NewFromInt(500000),
		}, nil)

	mockWalletRepo.On("Substract", ctx, srcWalletId, amount).Return(nil)
	mockWalletRepo.On("Add", ctx, destWalletId, amount).Return(nil)

	mockTransactRepo.On("WriteTransaction", ctx, mock.AnythingOfType("*model.Transaction"), mock.AnythingOfType("func() error")).
		Return(nil).Run(func(args mock.Arguments) {
			transaction := args.Get(1).(*model.Transaction)

			assert.Equal(t, srcWalletId, transaction.SourceWalletId)
			assert.Equal(t, destWalletId, transaction.DestinationWalletId)
			assert.Equal(t, amount, transaction.Amount)
			assert.Equal(t, model.TransactionType("transfer"), transaction.TransactionType)

			txFunc := args.Get(2).(func() error)
			err := txFunc()
			assert.NoError(t, err)
		})

	err := transactUsecase.Transfer(ctx, newTransactReq.userId, newTransactReq.phoneDest, newTransactReq.amountStr, newTransactReq.description)
	assert.NoError(t, err)

	mockTransactRepo.AssertExpectations(t)
	mockWalletRepo.AssertExpectations(t)
}
