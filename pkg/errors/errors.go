package errors

import "net/http"

// AppError represents an application error
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// Common error constructors
func NotFound(message string) *AppError {
	return NewAppError("NOT_FOUND", message, http.StatusNotFound)
}

func BadRequest(message string) *AppError {
	return NewAppError("BAD_REQUEST", message, http.StatusBadRequest)
}

func Unauthorized(message string) *AppError {
	return NewAppError("UNAUTHORIZED", message, http.StatusUnauthorized)
}

func Forbidden(message string) *AppError {
	return NewAppError("FORBIDDEN", message, http.StatusForbidden)
}

func InternalError(message string) *AppError {
	return NewAppError("INTERNAL_ERROR", message, http.StatusInternalServerError)
}

func Conflict(message string) *AppError {
	return NewAppError("CONFLICT", message, http.StatusConflict)
}
