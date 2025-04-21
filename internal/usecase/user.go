package usecase

import (
	"context"

	"github.com/Fadil-Tao/gopher-pay/internal/model"
)

type UserRepo interface {
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type UserUsecase struct {
	userRepo UserRepo	
}

func NewUserUsecase(userRepo UserRepo) *UserUsecase {
	return &UserUsecase{
		userRepo:  userRepo,
	}
}