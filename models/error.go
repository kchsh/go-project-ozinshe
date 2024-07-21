package models

type ApiError struct {
	Error string `json:"error"`
}

func NewApiError(msg string) ApiError {
	return ApiError{msg}
}
