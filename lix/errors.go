package lix

// LixError is the base error type for all SDK errors.
type LixError struct {
	msg string
}

func (e *LixError) Error() string { return e.msg }

// HttpClientError is returned when an underlying HTTP transport error occurs.
type HttpClientError struct{ LixError }

// UnauthorizedError is returned on HTTP 401.
type UnauthorizedError struct{ LixError }

// NotFoundException is returned on HTTP 404.
type NotFoundException struct{ LixError }

// RateLimitError is returned on HTTP 429.
type RateLimitError struct{ LixError }

// ServerError is returned on HTTP 500.
type ServerError struct{ LixError }

// ValidationError is returned on HTTP 400 and contains field-level parameter errors.
type ValidationError struct {
	LixError
	Data map[string]interface{}
}
