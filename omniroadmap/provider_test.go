package omniroadmap

import (
	"net/http"
	"net/http/httptest"
	"testing"

	productboard "github.com/grokify/productboard-go"

	"github.com/grokify/omniroadmap-core/provider"
	"github.com/grokify/omniroadmap-core/provider/providertest"
)

func newTestProvider(t *testing.T, handler http.HandlerFunc) *Provider {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := productboard.NewClient(
		productboard.WithAPIToken("test-token"),
		productboard.WithBaseURL(server.URL),
	)
	if err != nil {
		t.Fatalf("productboard.NewClient: %v", err)
	}
	return NewProvider(client)
}

func TestConformance(t *testing.T) {
	p := newTestProvider(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[],"links":{"next":null}}`))
	})

	providertest.RunAll(t, providertest.Config{
		Provider:        p,
		SkipIntegration: true,
	})
}

func TestListItems_UnsupportedKind(t *testing.T) {
	p := newTestProvider(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("no HTTP call expected for an unsupported kind")
	})

	_, err := p.ListItems(t.Context(), &provider.ListItemsRequest{
		Kinds: []provider.ItemKind{provider.ItemKindInitiative},
	})
	if err == nil {
		t.Fatal("expected an error for an unsupported kind, got nil")
	}
}

func TestListStatuses_Unsupported(t *testing.T) {
	p := newTestProvider(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("no HTTP call expected — ListStatuses is unsupported")
	})

	_, err := p.ListStatuses(t.Context(), &provider.ListStatusesRequest{})
	if err == nil {
		t.Fatal("expected ErrUnsupportedOperation, got nil")
	}
}

func TestListCustomFieldDefinitions_Unsupported(t *testing.T) {
	p := newTestProvider(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("no HTTP call expected — ListCustomFieldDefinitions is unsupported")
	})

	_, err := p.ListCustomFieldDefinitions(t.Context(), &provider.ListCustomFieldDefinitionsRequest{})
	if err == nil {
		t.Fatal("expected ErrUnsupportedOperation, got nil")
	}
}

func TestListItems_FollowsAllPages(t *testing.T) {
	pages := map[string]string{
		"": `{
			"data": [{"id": "195a1cb2-728f-4be8-900f-aebbd84d7944", "type": "feature", "fields": {"name": "Page One Feature"}}],
			"links": {"next": "cursor-2"}
		}`,
		"cursor-2": `{
			"data": [{"id": "295a1cb2-728f-4be8-900f-aebbd84d7945", "type": "feature", "fields": {"name": "Page Two Feature"}}],
			"links": {"next": null}
		}`,
	}
	var requestedCursors []string
	p := newTestProvider(t, func(w http.ResponseWriter, r *http.Request) {
		cursor := r.URL.Query().Get("pageCursor")
		requestedCursors = append(requestedCursors, cursor)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(pages[cursor]))
	})

	resp, err := p.ListItems(t.Context(), &provider.ListItemsRequest{
		Kinds: []provider.ItemKind{provider.ItemKindFeature},
	})
	if err != nil {
		t.Fatalf("ListItems: %v", err)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("len(Items) = %d, want 2 (should follow the next-page cursor)", len(resp.Items))
	}
	if len(requestedCursors) != 2 || requestedCursors[1] != "cursor-2" {
		t.Errorf("requested cursors = %v, want [\"\", \"cursor-2\"]", requestedCursors)
	}
}
