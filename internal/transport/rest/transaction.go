package rest

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/Fadil-Tao/gopher-pay/internal/model"
	"github.com/Fadil-Tao/gopher-pay/internal/transport/middleware"
	csterr "github.com/Fadil-Tao/gopher-pay/utils/errors"
	"github.com/Fadil-Tao/gopher-pay/utils/httpresponse"
	customvalidator "github.com/Fadil-Tao/gopher-pay/utils/validator"
	"github.com/go-playground/validator/v10"
)

type TransactionUsecase interface {
	Transfer(ctx context.Context, userId int, phoneDestination string, amountStr string, description string) error
}

type TransactHandler struct{
	TransactionUsecase
}

func NewTransactionHandler(transactUseCase TransactionUsecase)*TransactHandler{
	return &TransactHandler{
		TransactionUsecase: transactUseCase,
	}
}

func (t *TransactHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var transferReq model.TransferRequest

	if err := json.NewDecoder(r.Body).Decode(&transferReq); err != nil {
		slog.Error("error decoding request body", "message", err)
		httpresponse.WriteErrorResponse(w, "badRequest", nil, "invalid request body", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	validate.RegisterValidation("amount", customvalidator.AmountValidator)
	err := validate.Struct(transferReq)
	if err != nil {
		slog.Error(err.Error())
		fieldErrors := httpresponse.MapValidationError(err.(validator.ValidationErrors))
		httpresponse.WriteErrorResponse(w, "badRequest", *fieldErrors, "invalid request body", http.StatusBadRequest)
		return
	}

	userId, err := middleware.GetUserId(w,r)
	if err != nil {
		slog.Error(err.Error())
		httpresponse.WriteErrorResponse(w, "authorized", nil, "invalid key", http.StatusBadRequest)
		return
	}
	
	err = t.TransactionUsecase.Transfer(ctx, userId, transferReq.PhoneDest,transferReq.AmountStr,transferReq.Description)
	if err != nil {
		slog.Error(err.Error())
		if err == csterr.ErrInsufficientBalance {
			httpresponse.WriteErrorResponse(w,"failed", nil, "insufficient funds", http.StatusPaymentRequired)
			return
		}
		if err == csterr.ErrInternal {
			httpresponse.WriteErrorResponse(w,"failed", nil, "payment failed", http.StatusPaymentRequired)
			return
		}
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "transfer successfull",
	})
}