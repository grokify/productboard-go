package productboard

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
