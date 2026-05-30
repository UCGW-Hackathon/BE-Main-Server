package dto

type SuccessResponse[TResponse any] struct {
	Status  string    `json:"status"`
	Message string    `json:"message,omitempty"`
	Data    TResponse `json:"data,omitempty"`
	Meta    any       `json:"meta,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Status    string       `json:"status"`
	Message   string       `json:"message"`
	ErrorCode string       `json:"error_code,omitempty"`
	Errors    []FieldError `json:"errors,omitempty"`
}
