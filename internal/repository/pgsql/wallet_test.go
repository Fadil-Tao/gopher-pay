package pgsql_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Fadil-Tao/gopher-pay/internal/repository/pgsql"
	"github.com/stretchr/testify/assert"
)

func TestCreateWallet(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("unexpected error : %s", err)
    }
    defer db.Close()
    
    walletRepo := pgsql.NewWalletRepo(db)
    userId := 1
    
    query := `insert into wallets\(id, user_id\) values \(\$1, \$2\);`
    
    mock.ExpectPrepare(query).
        ExpectExec().
        WithArgs(sqlmock.AnyArg(), userId).
        WillReturnResult(sqlmock.NewResult(0, 1))
    
    err = walletRepo.CreateWallet(context.Background(), userId)
    
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}
