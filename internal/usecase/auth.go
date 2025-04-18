package usecase

import (
	"context"
	"log/slog"

	"github.com/Fadil-Tao/gopher-pay/internal/model"
	csterr "github.com/Fadil-Tao/gopher-pay/utils/errors"
)

type AuthRepo interface {
	Register(ctx context.Context, user model.User) error
	IsUserExist(ctx context.Context, email string) (bool, error)
}

type AuthService struct {
	AuthRepo
}

func NewAuthSerice(authRepo AuthRepo) *AuthService {
	return &AuthService{
		authRepo,
	}
}

func (a *AuthService) Register(ctx context.Context, user model.User) error {
	select {
	case <-ctx.Done():
		slog.Error("context passed")
		return ctx.Err()
	default:
	}

	isUserExist, err := a.AuthRepo.IsUserExist(ctx, user.Email)
	if err != nil {
		slog.Error(err.Error())
		return csterr.ErrInternal
	}
	if isUserExist {
		return csterr.ErrIsAlreadyExist
	}
	err = a.AuthRepo.Register(ctx, user)
	if err != nil {
		return csterr.ErrInternal
	}
	return nil
}
