package productboard

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ogen-go/ogen/validate"
)

// Sentinel errors for configuration and validation.
var (
	ErrMissingAPIToken = errors.New("productboard: api token is required")
)

// APIError represents an error response from the ProductBoard API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("productboard: API error (status %d): %s", e.StatusCode, e.Message)
}

// IsNotFound returns true if the error indicates a 404 Not Found response.
func IsNotFound(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsUnauthorized returns true if the error indicates a 401 Unauthorized response.
func IsUnauthorized(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusUnauthorized
	}
	return false
}

// IsRateLimited returns true if the error indicates a 429 Too Many Requests response.
func IsRateLimited(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusTooManyRequests
	}
	return false
}

// wrapError wraps an error from the API client with operation context.
func wrapError(operation string, err error) error {
	if err == nil {
		return nil
	}

	var statusErr *validate.UnexpectedStatusCodeError
	if errors.As(err, &statusErr) {
		return &APIError{
			StatusCode: statusErr.StatusCode,
			Message:    fmt.Sprintf("%s: %v", operation, err),
		}
	}

	errMsg := err.Error()
	statusCode := 0
	switch {
	case strings.Contains(errMsg, "404") || strings.Contains(errMsg, "not found"):
		statusCode = http.StatusNotFound
	case strings.Contains(errMsg, "401") || strings.Contains(errMsg, "unauthorized"):
		statusCode = http.StatusUnauthorized
	case strings.Contains(errMsg, "403") || strings.Contains(errMsg, "forbidden"):
		statusCode = http.StatusForbidden
	case strings.Contains(errMsg, "429") || strings.Contains(errMsg, "rate limit"):
		statusCode = http.StatusTooManyRequests
	case strings.Contains(errMsg, "500") || strings.Contains(errMsg, "internal server"):
		statusCode = http.StatusInternalServerError
	}

	if statusCode > 0 {
		return &APIError{StatusCode: statusCode, Message: fmt.Sprintf("%s: %v", operation, err)}
	}

	return fmt.Errorf("%s: %w", operation, err)
}
