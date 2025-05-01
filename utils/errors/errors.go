package csterr
import "errors"

var (
	ErrNotFound = errors.New("error not found")
	ErrInternal = errors.New("internal server error")
	ErrEmailAlreadyUsed = errors.New("email already used")
	ErrUnauthorized = errors.New("unauthorized")
	ErrPhoneAlreadyUsed = errors.New("phone already used")
	ErrInsufficientAmount = errors.New("insufficient amount")
	ErrInsufficientBalance = errors.New("insufficient balance")

)