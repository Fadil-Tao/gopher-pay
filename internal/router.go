package internal

import (
	"net/http"

	"github.com/Fadil-Tao/gopher-pay/internal/transport/rest"
)

type RouterInstance struct {
	UserUsecase string
	AuthUseCase string
}

// contain centralized api route
func (r *RouterInstance) NewRestRouter() *http.ServeMux {
	mux := http.NewServeMux()


	userHandler := rest.NewUserHandler(r.UserUsecase)
	authHandler := rest.NewAuthHandler(r.AuthUseCase)

	mux.HandleFunc("GET /healthcheck", rest.HealthCheck)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("PUT /profile", userHandler.EditProfile)

	return mux
}	
