package pgsql

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/Fadil-Tao/gopher-pay/internal/model"
	csterr "github.com/Fadil-Tao/gopher-pay/utils/errors"
)

type UserRepo struct {
	Db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		Db: db,
	}
}

func (u *UserRepo) Register(ctx context.Context, user model.User, createFn func(userId int)error) error {
	return runInTx(u.Db, func(tx *sql.Tx) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		query := `insert into users(email,name,password,salt,phone) values ($1, $2, $3, $4,$5) returning id`
		stmt, err := u.Db.Prepare(query)
		if err != nil {
			slog.Error("error preparing statement", "message", err)
			return csterr.ErrInternal
		}
		defer stmt.Close()

		var newUserId *int
		err = stmt.QueryRowContext(ctx, user.Email, user.Name, user.Password, user.Salt, user.Phone).Scan(&newUserId)
		if err != nil {
			slog.Error("Error inserting data", "message", err)
			return csterr.ErrInternal
		}
		err = createFn(*newUserId)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		return nil
	})
}

func (u *UserRepo) IsUserExist(ctx context.Context, email string) (bool, error) {
	select {
	case <-ctx.Done():
		slog.Error("context passed")
		return false, ctx.Err()
	default:
	}

	query := `select count(id) from users where email = $1`
	stmt, err := u.Db.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return false, csterr.ErrInternal
	}
	defer stmt.Close()
	var count int
	err = stmt.QueryRowContext(ctx, email).Scan(&count)
	if err != nil {
		slog.Error(err.Error())
		return false, csterr.ErrInternal
	}
	return count > 0, nil
}

func (u *UserRepo) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	query := `select id,email,name,password,salt,phone,updated_at,created_at from users where phone = $1`
	stmt, err := u.Db.Prepare(query)
	if err != nil {
		slog.Error(err.Error())
		return nil, csterr.ErrInternal
	}
	defer stmt.Close()

	user := &model.User{}

	err = stmt.QueryRowContext(ctx, phone).Scan(&user.Id, &user.Email, &user.Name, &user.Password, &user.Salt, &user.Password, &user.UpdatedAt, &user.CreatedAt)
	switch {
	case err == sql.ErrNoRows:
		slog.Error(err.Error())
		return nil, csterr.ErrNotFound
	case err != nil:
		slog.Error(err.Error())
		return nil, csterr.ErrInternal
	}
	return user, nil
}

func (u *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	query := `select id,email,name,password,salt,phone,updated_at,created_at from users where email = $1`
	stmt, err := u.Db.Prepare(query)
	if err != nil {
		slog.Error(err.Error())
		return nil, csterr.ErrInternal
	}
	defer stmt.Close()

	user := &model.User{}

	err = stmt.QueryRowContext(ctx, email).Scan(&user.Id, &user.Email, &user.Name, &user.Password, &user.Salt, &user.Phone, &user.UpdatedAt, &user.CreatedAt)
	switch {
	case err == sql.ErrNoRows:
		slog.Error(err.Error())
		return nil, csterr.ErrNotFound
	case err != nil:
		slog.Error(err.Error())
		return nil, csterr.ErrInternal
	}
	return user, nil
}