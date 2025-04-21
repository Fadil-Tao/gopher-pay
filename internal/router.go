package internal

import (
	"encoding/json"
	"net/http"

	"github.com/Fadil-Tao/gopher-pay/internal/transport/middleware"
	"github.com/Fadil-Tao/gopher-pay/internal/transport/rest"
	"github.com/Fadil-Tao/gopher-pay/internal/usecase"
)

type RouterInstance struct {
	UserUsecase usecase.UserUsecase
	AuthUseCase usecase.AuthUsecase
}

func NewHttpRouterInstance(UserUsecase usecase.UserUsecase, AuthUseCase usecase.AuthUsecase) *RouterInstance {
	return &RouterInstance{
		UserUsecase: UserUsecase,
		AuthUseCase: AuthUseCase,
	}
}

func (r *RouterInstance) NewRestRouter() *http.ServeMux {
	mux := http.NewServeMux()
	authHandler := rest.NewAuthHandler(&r.AuthUseCase)

	// public 
	mux.HandleFunc("GET /healthcheck", rest.HealthCheck)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/logout", authHandler.Logout)

	// require auth
    mux.Handle("GET /profile", middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]string{"message": "this is suppossed to be a secret!!!"})
    })))
	mux.Handle("GET /secret", middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]string{"message": "this is suppossed to be a secret!!!"})
    })))
	
	return mux
}