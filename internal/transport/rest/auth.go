package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/Fadil-Tao/gopher-pay/internal/model"
)


type AuthService interface {
	Register(ctx context.Context, user model.User) error
}

type AuthHandler struct {
	AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler{
	return &AuthHandler{
		AuthService : authService,
	}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	
}


