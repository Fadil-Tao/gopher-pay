package internal

import (
	"database/sql"
	"net/http"

	"github.com/Fadil-Tao/gopher-pay/internal/repository/pgsql"
	"github.com/Fadil-Tao/gopher-pay/internal/usecase"
)

func Register(conn *sql.DB) *http.ServeMux{
	userRepo := pgsql.NewUserRepo(conn)
	
	authUseCase := usecase.NewAuthUsecase(userRepo)
	userUseCase := usecase.NewUserUsecase(userRepo)

	router :=  NewHttpRouterInstance(*userUseCase, *authUseCase)
	return router.NewRestRouter()
}