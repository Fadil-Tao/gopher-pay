package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Fadil-Tao/gopher-pay/internal/model"
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
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	validate := validator.New()

	err := validate.Struct(newUser)
	if err != nil {
		slog.Error(err.Error())
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("validation error : %s", errors), http.StatusBadRequest)
		return
	}

	err = a.svc.Register(ctx, newUser)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "user successfully registered"})
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	loginReq := struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required,min=8"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		slog.Error("error decoding request body", "message", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err := validate.Struct(loginReq)
	if err != nil {
		slog.Error(err.Error())
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("request body error : %s", errors), http.StatusBadRequest)
		return
	}

	jwtToken, err := a.svc.Login(ctx, loginReq.Email, loginReq.Password)
	if err != nil {
		http.Error(w, "An unexpected error occurred on the server", http.StatusInternalServerError)
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
