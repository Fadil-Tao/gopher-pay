package rest

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/Fadil-Tao/gopher-pay/internal/model"
	csterr "github.com/Fadil-Tao/gopher-pay/utils/errors"
	"github.com/Fadil-Tao/gopher-pay/utils/httpresponse"
	"github.com/go-playground/validator/v10"
)

type AuthService interface {
	Register(ctx context.Context, user model.User) error
	Login(ctx context.Context, email string, password string) (response *string, err error)
}

type AuthHandler struct {
	svc AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{
		svc: authService,
	}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var newUser model.User

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		slog.Error(err.Error())
		httpresponse.WriteErrorResponse(w, "error", nil , "invalid request", http.StatusBadRequest)
		return
	}

	validate := validator.New()

	err := validate.Struct(newUser)
	if err != nil {
		slog.Error(err.Error())
		fieldErrors := httpresponse.MapValidationError(err.(validator.ValidationErrors))
		httpresponse.WriteErrorResponse(w,"badRequest", *fieldErrors, "invalid request body", http.StatusBadRequest)
		return
	}

	err = a.svc.Register(ctx, newUser)
	if err != nil {
		slog.Error(err.Error())
		if err == csterr.ErrIsAlreadyExist {
			httpresponse.WriteErrorResponse(w, "error", nil, "user already exist", http.StatusConflict)
			return 
		}
		httpresponse.WriteErrorResponse(w, "error", nil, "unexpected internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "user registered successfully",
	})
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	loginReq := struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		slog.Error("error decoding request body", "message", err)
		httpresponse.WriteErrorResponse(w,"badRequest",nil, "invalid request body", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err := validate.Struct(loginReq)
	if err != nil {
		slog.Error(err.Error())
		fieldErrors := httpresponse.MapValidationError(err.(validator.ValidationErrors))
		httpresponse.WriteErrorResponse(w,"badRequest", *fieldErrors, "invalid request body", http.StatusBadRequest)
		return
	}

	jwtToken, err := a.svc.Login(ctx, loginReq.Email, loginReq.Password)
	if err != nil {
		slog.Error(err.Error())
		httpresponse.WriteErrorResponse(w,"unauthorized", nil,"invalid credentials",http.StatusUnauthorized)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    *jwtToken,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Login Successfull",
	})
}


func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:     "token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge: -1,
	}

	http.SetCookie(w, &c)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"Message": "Logout Success"})	
}