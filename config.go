package productboard

import (
	"net/http"
	"os"
	"time"
)

// defaultBaseURL is ProductBoard's public API v2 base URL.
const defaultBaseURL = "https://api.productboard.com/v2"

// Config holds the configuration for the Client.
type Config struct {
	// APIToken is your ProductBoard API token (Personal Access Token).
	APIToken string

	// HTTPClient is the HTTP client to use for requests.
	// If nil, http.DefaultClient is used.
	HTTPClient *http.Client

	// Timeout is the request timeout. Default is 60 seconds.
	Timeout time.Duration

	// BaseURL overrides the default API URL. If empty, defaultBaseURL is used.
	BaseURL string
}

func (c *Config) loadDefaults() {
	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}
	if c.Timeout == 0 {
		c.Timeout = 60 * time.Second
	}
}

// loadEnv loads configuration from environment variables. Only used if the
// corresponding field is empty.
func (c *Config) loadEnv() {
	if c.APIToken == "" {
		c.APIToken = os.Getenv("PRODUCTBOARD_API_TOKEN")
	}
}

func (c *Config) validate() error {
	if c.APIToken == "" {
		return ErrMissingAPIToken
	}
	return nil
}

func (c *Config) buildBaseURL() string {
	if c.BaseURL != "" {
		return c.BaseURL
	}
	return defaultBaseURL
}
