package pgsql

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/Fadil-Tao/gopher-pay/internal/model"
	csterr "github.com/Fadil-Tao/gopher-pay/utils/errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type WalletRepo struct {
	DB *sql.DB
}

func NewWalletRepo(db *sql.DB) *WalletRepo {
	return &WalletRepo{
		DB: db,
	}
}

func (w *WalletRepo) GetWalletInfoByUserId(ctx context.Context, userId int) (*model.Wallet, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	query := `select id, balance from wallets where user_id = $1;`
	stmt, err := w.DB.Prepare(query)
	if err != nil {
		return nil, csterr.ErrInternal
	}
	defer stmt.Close()

	wallet := &model.Wallet{}
	err = stmt.QueryRowContext(ctx, userId).Scan(&wallet.Id, &wallet.Balance)
	if err != nil {
		slog.Error(err.Error())
		if err == sql.ErrNoRows {
			return nil, csterr.ErrNotFound
		}
		return nil, csterr.ErrInternal
	}
	return wallet, nil
}

func (w *WalletRepo) GetWalletInfoByUserPhone(ctx context.Context, userPhone string) (*model.Wallet, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	query := `select w.id, w.balance from wallets as w join users on users.id = wallets.user_id where users.phone = $1`

	stmt, err := w.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "message", err)
		return nil, csterr.ErrInternal
	}
	defer stmt.Close()

	wallet := &model.Wallet{}
	err = stmt.QueryRowContext(ctx, userPhone).Scan(&wallet.Id, &wallet.Balance)
	if err != nil {
		slog.Error(err.Error())
		if err == sql.ErrNoRows {
			return nil, csterr.ErrNotFound
		}
		return nil, csterr.ErrInternal
	}
	return wallet, nil
}

func (w *WalletRepo) Add(ctx context.Context, id uuid.UUID, amount decimal.Decimal) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	query := `update wallets set balance = balance - $2 where id = $1;`
	stmt, err := w.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return csterr.ErrInternal
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id, amount)
	if err != nil {
		slog.Error("error executing context", "message", err)
		return csterr.ErrInternal
	}

	return nil
}

func (w *WalletRepo) Substract(ctx context.Context, id uuid.UUID, amount decimal.Decimal) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	query := `update wallets set balance = balance - $2 where id = $1`
	stmt, err := w.DB.Prepare(query)
	if err != nil {
		return csterr.ErrInternal
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id, amount)
	if err != nil {
		slog.Error(err.Error())
		return csterr.ErrInternal
	}
	return nil
}

func (w *WalletRepo) CreateWallet(ctx context.Context, userId int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	slog.Debug("new user id", "message", userId)

	newWalletId := uuid.New()
	query := `insert into wallets(id, user_id) values ($1, $2);`
	stmt, err := w.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "message", err)
		return csterr.ErrInternal
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, newWalletId, userId)
	if err != nil {
		slog.Error("error creating wallet", "message", err.Error())
		return csterr.ErrInternal
	}
	return nil
}
