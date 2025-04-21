package httpresponse

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type BaseResponse struct {
	Code    int    `json:"statusCode"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse[T any] struct {
	BaseResponse
	Data T `json:"data,omitempty"`
}

type ErrorResponse struct {
	BaseResponse
	Errors map[string][]string `json:"errors,omitempty"`
}

type PaginationInfo struct {
	TotalPage    int `json:"totalPage"`
	CurrentPage  int `json:"current"`
	Size         int `json:"size"`
	NextPage     int `json:"nextPage"`
	PreviousPage int `json:"prevPage"`
}

type PaginatedResponse[T any] struct {
	Items      T
	Pagination PaginationInfo
}

var fieldErrorMessages = map[string]string{
    "email":    "invalid email",
    "required": "cannot be empty",
    "min":      "value is too short",
    "max":      "value is too long",
}

func WriteSuccessResponse[T any](w http.ResponseWriter, status string, data T, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	response := SuccessResponse[T]{
		BaseResponse: BaseResponse{
			Code:   code,
			Status: status,
		},
		Data: data,
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func WriteErrorResponse(w http.ResponseWriter, status string, errors map[string][]string, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	response := ErrorResponse{
		BaseResponse: BaseResponse{
			Status:  status,
			Message: message,
			Code:    code,
		},
		Errors: errors,
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func WithPagination[T any](items T, totalItems, currentPage, pageSize int) *PaginatedResponse[T] {
	totalPage := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	nextPage := 0
	if currentPage < totalPage {
		nextPage = currentPage + 1
	}

	previousPage := 0
	if currentPage > 1 {
		previousPage = currentPage - 1
	}

	return &PaginatedResponse[T]{
		Items: items,
		Pagination: PaginationInfo{
			TotalPage:    totalPage,
			CurrentPage:  currentPage,
			Size:         pageSize,
			NextPage:     nextPage,
			PreviousPage: previousPage,
		},
	}
}

func SuccessResponseWithPagination[T any](w http.ResponseWriter, status string, items T, code, totalItems, currentPage, pageSize int) {
	paginatedData := WithPagination(items, totalItems, currentPage, pageSize)
	WriteSuccessResponse(w, status, paginatedData, code)
}

func MapValidationError(err validator.ValidationErrors) *map[string][]string {
	fieldErrors := make(map[string][]string)
	for _, fieldError := range err {
		field := fieldError.Field()
		jsonField := strings.ToLower(field[:1]) + field[1:]
		tag := fieldError.Tag()

		errorMsg := "validation failed"

		if message, exists := fieldErrorMessages[tag]; exists {
			errorMsg = message

			if tag == "min" || tag == "max" {
				param := fieldError.Param()
				switch tag {
				case "min":
					errorMsg = fmt.Sprintf("must be at least %s characters", param)
				case "max":
					errorMsg = fmt.Sprintf("must not exceed %s characters", param)
				}
			}
		}

		if messages, exists := fieldErrors[jsonField]; exists {
			fieldErrors[jsonField] = append(messages, errorMsg)
		} else {
			fieldErrors[jsonField] = []string{errorMsg}
		}
	}
	return &fieldErrors
}