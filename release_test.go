package productboard

import (
	"net/http"
	"testing"
)

const getReleaseFixture = `{
	"data": {
		"id": "395a1cb2-728f-4be8-900f-aebbd84d7946",
		"type": "release",
		"fields": {
			"name": "MVP3",
			"status": {"id": "st-2", "name": "Planned"},
			"timeframe": {"startDate": "2026-09-01", "endDate": "2026-09-30", "granularity": "month"},
			"archived": false
		}
	}
}`

func TestGetRelease(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(getReleaseFixture))
	})

	rel, err := client.GetRelease(t.Context(), "395a1cb2-728f-4be8-900f-aebbd84d7946")
	if err != nil {
		t.Fatalf("GetRelease: %v", err)
	}
	if rel.Name != "MVP3" {
		t.Errorf("Name = %q, want %q", rel.Name, "MVP3")
	}
	if rel.Status == nil || rel.Status.Name != "Planned" {
		t.Errorf("Status = %+v, want Name=Planned", rel.Status)
	}
}

const listReleasesFixture = `{
	"data": [
		{"id": "395a1cb2-728f-4be8-900f-aebbd84d7946", "type": "release", "fields": {"name": "MVP3"}}
	],
	"links": {"next": null}
}`

func TestListReleases(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(listReleasesFixture))
	})

	list, err := client.ListReleases(t.Context())
	if err != nil {
		t.Fatalf("ListReleases: %v", err)
	}
	if len(list.Releases) != 1 {
		t.Fatalf("len(Releases) = %d, want 1", len(list.Releases))
	}
	if list.NextCursor != "" {
		t.Errorf("NextCursor = %q, want empty (null next link)", list.NextCursor)
	}
}
