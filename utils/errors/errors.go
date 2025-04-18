package csterr
import "errors"

var (
	ErrNotFound = errors.New("error not found")
	ErrInternal = errors.New("internal server error")
	ErrIsAlreadyExist = errors.New("source already exist")
)