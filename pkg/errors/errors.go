package apperrors

import "net/http"

// AppError is a custom error type with HTTP status code
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

// Common errors
var (
	ErrNotFound        = &AppError{Code: http.StatusNotFound, Message: "resource not found"}
	ErrUnauthorized    = &AppError{Code: http.StatusUnauthorized, Message: "unauthorized"}
	ErrForbidden       = &AppError{Code: http.StatusForbidden, Message: "forbidden"}
	ErrBadRequest      = &AppError{Code: http.StatusBadRequest, Message: "bad request"}
	ErrInternalServer  = &AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	ErrConflict        = &AppError{Code: http.StatusConflict, Message: "resource already exists"}
	ErrInvalidPassword = &AppError{Code: http.StatusUnauthorized, Message: "invalid email or password"}
	ErrEmailTaken      = &AppError{Code: http.StatusConflict, Message: "email already taken"}
)

// New creates a new AppError
func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}
