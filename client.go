// Package productboard is a Go SDK for the ProductBoard REST API v2
// (https://developer.productboard.com/reference/introduction). Unlike
// aha-go, ProductBoard publishes a real official OpenAPI spec, vendored at
// openapi/entities.json (see openapi/SOURCE.md for provenance) and used to
// generate internal/api via ogen — see generate.sh / make generate.
//
// Scope is deliberately thin for v0.1: read-only access to Features,
// Releases, and Objectives via ProductBoard's generic Entities API, just
// enough to back an github.com/grokify/omniroadmap-core Provider adapter
// (see the omniroadmap/ subpackage). Full CRUD across ProductBoard's other
// APIs (Notes, Teams, Members, Webhooks, Analytics, Integrations) is out of
// scope until there's real need for it.
package productboard

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grokify/productboard-go/internal/api"
)

const (
	// SDKVersion is the version of this SDK.
	SDKVersion = "0.1.0"

	// SDKName is the name of this SDK.
	SDKName = "productboard-go"
)

// Client provides access to the ProductBoard API.
type Client struct {
	config    *Config
	apiClient *api.Client
}

// NewClient creates a new ProductBoard client with the given options.
//
// Configuration is loaded in the following order (later values override
// earlier ones): defaults, PRODUCTBOARD_API_TOKEN environment variable,
// options passed to NewClient.
func NewClient(opts ...Option) (*Client, error) {
	cfg := &Config{}
	cfg.loadDefaults()
	cfg.loadEnv()

	for _, opt := range opts {
		opt(cfg)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	secSource := &securitySource{token: cfg.APIToken}
	httpClient := &sdkHTTPClient{client: cfg.HTTPClient}

	apiClient, err := api.NewClient(cfg.buildBaseURL(), secSource, api.WithClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("creating API client: %w", err)
	}

	return &Client{config: cfg, apiClient: apiClient}, nil
}

// BaseURL returns the API base URL.
func (c *Client) BaseURL() string {
	return c.config.buildBaseURL()
}

// API returns the low-level ogen-generated client for advanced use, for API
// operations not covered by the high-level wrapper.
func (c *Client) API() *api.Client {
	return c.apiClient
}

// securitySource implements api.SecuritySource for bearer token auth.
type securitySource struct {
	token string
}

func (s *securitySource) BearerAuth(ctx context.Context, operationName api.OperationName) (api.BearerAuth, error) {
	return api.BearerAuth{Token: s.token}, nil
}

// sdkHTTPClient wraps an HTTP client to add SDK identification headers.
type sdkHTTPClient struct {
	client *http.Client
}

func (c *sdkHTTPClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", fmt.Sprintf("%s/%s", SDKName, SDKVersion))
	if req.Header.Get("Content-Type") == "" && req.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json")
	}
	return c.client.Do(req) //nolint:gosec // G704: URL from client config, not user input
}
