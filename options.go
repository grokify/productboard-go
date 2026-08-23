package productboard

import (
	"net/http"
	"time"
)

// Option configures the Client.
type Option func(*Config)

// WithAPIToken sets the ProductBoard API token.
func WithAPIToken(token string) Option {
	return func(c *Config) { c.APIToken = token }
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) { c.HTTPClient = client }
}

// WithTimeout sets the request timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) { c.Timeout = timeout }
}

// WithBaseURL overrides the default API URL.
func WithBaseURL(baseURL string) Option {
	return func(c *Config) { c.BaseURL = baseURL }
}

// ListOptions configures list operations. ProductBoard's Entities API uses
// cursor-based pagination, not page numbers — PageCursor comes from the
// previous response's NextCursor.
type ListOptions struct {
	PageCursor string
	Name       string
}

// ListOption configures a list operation.
type ListOption func(*ListOptions)

// WithPageCursor sets the pagination cursor (from a previous ListResult's
// NextCursor).
func WithPageCursor(cursor string) ListOption {
	return func(o *ListOptions) { o.PageCursor = cursor }
}

// WithNameFilter filters results by entity name.
func WithNameFilter(name string) ListOption {
	return func(o *ListOptions) { o.Name = name }
}

func applyListOptions(opts ...ListOption) *ListOptions {
	o := &ListOptions{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}
