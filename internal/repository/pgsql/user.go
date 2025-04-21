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

func (u *UserRepo) Register(ctx context.Context, user model.User) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	query := `insert into users(email,name,password,salt) values ($1, $2, $3, $4)`
	stmt, err := u.Db.Prepare(query)
	if err != nil {
		slog.Error("error preparing statemet", "message", err)
		return csterr.ErrInternal
	}

	_, err = stmt.ExecContext(ctx, user.Email, user.Name, user.Password, user.Salt)
	if err != nil {
		slog.Error("Error inserting data", "message", err)
		return csterr.ErrInternal
	}
	return nil
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

func (u *UserRepo) DeleteProfile(ctx context.Context) error {
	return nil
}

func (u *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	query := `select id,email,name,password,salt,updated_at,created_at from users where email = $1`
	stmt, err := u.Db.Prepare(query)
	if err != nil {
		return nil, csterr.ErrInternal
	}
	defer stmt.Close()

	user := &model.User{}
	
	err = stmt.QueryRowContext(ctx, email).Scan(&user.Id, &user.Email, &user.Name,&user.Password, &user.Salt, &user.UpdatedAt, &user.CreatedAt)
	switch {
	case err == sql.ErrNoRows : 
		slog.Error(err.Error())
		return nil, csterr.ErrNotFound
	case err != nil : 
		slog.Error(err.Error())
		return nil, csterr.ErrInternal
	}
	return user, nil
}