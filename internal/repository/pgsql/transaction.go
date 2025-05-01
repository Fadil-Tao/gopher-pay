package pgsql

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/Fadil-Tao/gopher-pay/internal/model"
	csterr "github.com/Fadil-Tao/gopher-pay/utils/errors"
)

type TransactionRepo struct {
	DB *sql.DB
}

func NewTransactionRepo(db *sql.DB) *TransactionRepo {
	return &TransactionRepo{
		DB: db,
	}
}

func (t *TransactionRepo) GetUserTransactions(ctx context.Context, userId int, maxItem int, page int) ([]model.Transaction, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	query := `
	select (t.id, u.email as "sender", u.email as "receiver", t.amount, t.transaction_type, t.description, t.created_at) 
	from transaction as t
	join users as user_sender on t.sender_id = user_sender.id 
	join users as user_receiver on t.receiver_id =  user_receiver.id
	where t.senderName = $1 or t.receiverName = $1 order by t.created_at desc limit $2 OFFSET $3; 
	`
	_, err := t.DB.Prepare(query)
	if err != nil {
		slog.Error(err.Error())
		return nil, csterr.ErrInternal
	}

	return nil, nil
}

func (t *TransactionRepo) WriteTransaction(ctx context.Context, transaction *model.Transaction, updateFn func() error) error {
	return runInTx(t.DB, func(tx *sql.Tx) error {

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := updateFn()
		if err != nil {
			return err
		}

		query := `insert into transactions(id, source_wallet_id, destination_wallet_id, amount, transaction_type, description) values ($1,$2,$3,$4,$5,$6);`
		stmt, err := t.DB.Prepare(query)
		if err != nil {
			slog.Error(err.Error())
			return csterr.ErrInternal
		}

		defer stmt.Close()

		res, err := stmt.ExecContext(ctx, transaction.Id, transaction.SourceWalletId, transaction.DestinationWalletId, transaction.Amount, transaction.TransactionType, transaction.Description)
		if err != nil {
			slog.Error(err.Error())
			return csterr.ErrInternal
		}
		
		rowsAffected,err := res.RowsAffected()
		if err != nil { 
			slog.Error(err.Error())
			return csterr.ErrInternal
		}
		slog.Info("inserting transaction", "rows affected", rowsAffected)
		if rowsAffected < 1{
			return csterr.ErrInternal 
		}
		return nil
	})
}

func (w *WalletRepo) Transfer(ctx context.Context) error {
	return nil
}
