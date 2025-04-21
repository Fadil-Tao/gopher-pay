package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type UserHandler struct {
	userService string
}

func NewUserHandler(userService string) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) EditProfile(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")	
	if err := json.NewEncoder(w).Encode(map[string]string{
		"mesasge": u.userService,
	}); err != nil {
		slog.Error("error encoding message")
	}
} 

func (u *UserHandler) SomethingSecret(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"mesasge": "this suppossed be a secret",
	}); err != nil {
		slog.Error("error encoding message")
	}
}