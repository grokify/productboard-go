package productboard

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := NewClient(
		WithAPIToken("test-token"),
		WithBaseURL(server.URL),
	)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return client
}

func TestNewClient_RequiresAPIToken(t *testing.T) {
	if _, err := NewClient(); err == nil {
		t.Fatal("NewClient() with no token: error = nil, want ErrMissingAPIToken")
	}
}

func TestNewClient_DefaultTimeout(t *testing.T) {
	client, err := NewClient(WithAPIToken("test-token"))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if got := client.config.HTTPClient.Timeout; got != 60*time.Second {
		t.Errorf("HTTPClient.Timeout = %v, want 60s", got)
	}
}

func TestNewClient_WithTimeout(t *testing.T) {
	client, err := NewClient(WithAPIToken("test-token"), WithTimeout(5*time.Second))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if got := client.config.HTTPClient.Timeout; got != 5*time.Second {
		t.Errorf("HTTPClient.Timeout = %v, want 5s", got)
	}
}

func TestNewClient_WithHTTPClient_NotOverridden(t *testing.T) {
	custom := &http.Client{Timeout: 3 * time.Second}
	client, err := NewClient(WithAPIToken("test-token"), WithHTTPClient(custom), WithTimeout(90*time.Second))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if client.config.HTTPClient != custom {
		t.Error("a caller-supplied HTTPClient should be used as-is, not replaced")
	}
	if got := client.config.HTTPClient.Timeout; got != 3*time.Second {
		t.Errorf("HTTPClient.Timeout = %v, want 3s (Timeout should not override a caller-supplied client)", got)
	}
}
