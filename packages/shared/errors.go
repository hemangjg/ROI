package shared

import "net/http"

// AppError is a typed application error mapped to HTTP status codes.
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message}
}

func BadRequest(code, message string) *AppError {
	return NewAppError(http.StatusBadRequest, code, message)
}

func Internal(message string) *AppError {
	return NewAppError(http.StatusInternalServerError, "internal_error", message)
}

func Unauthorized(code, message string) *AppError {
	return NewAppError(http.StatusUnauthorized, code, message)
}

func Conflict(code, message string) *AppError {
	return NewAppError(http.StatusConflict, code, message)
}

func NotFound(code, message string) *AppError {
	return NewAppError(http.StatusNotFound, code, message)
}

func Forbidden(code, message string) *AppError {
	return NewAppError(http.StatusForbidden, code, message)
}

func RateLimited(code, message string) *AppError {
	return NewAppError(http.StatusTooManyRequests, code, message)
}