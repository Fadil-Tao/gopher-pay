package pgsql_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Fadil-Tao/gopher-pay/internal/model"
	"github.com/Fadil-Tao/gopher-pay/internal/repository/pgsql"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestShouldWriteRegister(t *testing.T){
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error : %s", err)
	}
	defer db.Close()

	transactRepo := pgsql.NewTransactionRepo(db)

	amount, err := decimal.NewFromString("10000.23")
	if err != nil {
		t.Fatalf("unexpected error string to decimal : %s", err)
	} 
	Description := "money from zambia"

	mockTransact := &model.Transaction{
		Id: uuid.New() ,
		SourceWalletId: uuid.New(),
		DestinationWalletId: uuid.New(),
		Amount:  amount,
		Description:  &Description,
	}
	
	mock.ExpectBegin() 
	mock.ExpectPrepare("insert into transactions\\(id, source_wallet_id, destination_wallet_id, amount, transaction_type, description\\) values").
	ExpectExec().
	WithArgs(
		mockTransact.Id, 
		mockTransact.SourceWalletId, 
		mockTransact.DestinationWalletId, 
		mockTransact.Amount,
		mockTransact.TransactionType,
		mockTransact.Description).
	WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	updateFnCalled := false
	mockUpdatefn := func () error {
		updateFnCalled = true
		return nil
	}

	err = transactRepo.WriteTransaction(context.Background(), mockTransact, mockUpdatefn)

	assert.NoError(t, err)
	assert.True(t, updateFnCalled, "update fn should be called")
	assert.NoError(t, mock.ExpectationsWereMet())
}