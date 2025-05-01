package internal

import (
	"encoding/json"
	"net/http"

	"github.com/Fadil-Tao/gopher-pay/internal/transport/middleware"
	"github.com/Fadil-Tao/gopher-pay/internal/transport/rest"
	"github.com/Fadil-Tao/gopher-pay/internal/usecase"
)

type RouterInstance struct {
	UserUsecase     usecase.UserUsecase
	AuthUseCase     usecase.AuthUsecase
	TransactUseCase usecase.TransactionUsecase
}

func NewHttpRouterInstance(UserUsecase usecase.UserUsecase, AuthUseCase usecase.AuthUsecase, transactUsecase usecase.TransactionUsecase) *RouterInstance {
	return &RouterInstance{
		UserUsecase:     UserUsecase,
		AuthUseCase:     AuthUseCase,
		TransactUseCase: transactUsecase,
	}
}

func (r *RouterInstance) NewRestRouter() *http.ServeMux {
	mux := http.NewServeMux()
	authHandler := rest.NewAuthHandler(&r.AuthUseCase)
	transactHandler := rest.NewTransactionHandler(&r.TransactUseCase)

	// public
	mux.HandleFunc("GET /healthcheck", rest.HealthCheck)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/logout", authHandler.Logout)

	// require auth
	mux.Handle("GET /profile", middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "this is suppossed to be a secret!!!"})
	})))
	mux.Handle("POST /pay", middleware.RequireAuth(http.HandlerFunc(transactHandler.Transfer)))
	mux.Handle("GET /secret", middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "this is suppossed to be a secret!!!"})
	})))

	return mux
}
