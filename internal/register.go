package internal

import (
	"database/sql"
	"net/http"

	"github.com/Fadil-Tao/gopher-pay/internal/repository/pgsql"
	"github.com/Fadil-Tao/gopher-pay/internal/usecase"
)

func Register(conn *sql.DB) *http.ServeMux{

	// repository layer injection
	userRepo := pgsql.NewUserRepo(conn)
	walletRepo := pgsql.NewWalletRepo(conn)
	transactRepo := pgsql.NewTransactionRepo(conn)

	// use case layer injection
	authUseCase := usecase.NewAuthUsecase(userRepo,walletRepo)
	userUseCase := usecase.NewUserUsecase(userRepo)
	transactUsecase := usecase.NewTransactionUsecase(transactRepo, walletRepo)

	router :=  NewHttpRouterInstance(
		*userUseCase, 
		*authUseCase, 
		*transactUsecase,
	)
	return router.NewRestRouter()
}