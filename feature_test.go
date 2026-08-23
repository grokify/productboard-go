package productboard

import (
	"net/http"
	"testing"
)

const getFeatureFixture = `{
	"data": {
		"id": "195a1cb2-728f-4be8-900f-aebbd84d7944",
		"type": "feature",
		"fields": {
			"name": "Awesome Public API",
			"description": "A public API for integrations",
			"status": {"id": "st-1", "name": "In Progress"},
			"owner": {"id": "ow-1", "email": "john@doe.com"},
			"tags": [{"id": "tag-1", "name": "api", "color": "blue"}],
			"timeframe": {"startDate": "2023-10-01", "endDate": "2023-10-31", "granularity": "month"},
			"archived": false
		},
		"relationships": {
			"data": [
				{
					"type": "parent",
					"target": {"id": "318de52f-4e38-4c94-a550-a0d47a1f212e", "type": "component"}
				}
			]
		},
		"createdAt": "2023-01-01T00:00:00Z",
		"updatedAt": "2023-01-02T00:00:00Z"
	}
}`

func TestGetFeature(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(getFeatureFixture))
	})

	f, err := client.GetFeature(t.Context(), "195a1cb2-728f-4be8-900f-aebbd84d7944")
	if err != nil {
		t.Fatalf("GetFeature: %v", err)
	}

	if f.Name != "Awesome Public API" {
		t.Errorf("Name = %q, want %q", f.Name, "Awesome Public API")
	}
	if f.Status == nil || f.Status.Name != "In Progress" {
		t.Errorf("Status = %+v, want Name=In Progress", f.Status)
	}
	if f.Owner == nil || f.Owner.Email != "john@doe.com" {
		t.Errorf("Owner = %+v, want Email=john@doe.com", f.Owner)
	}
	if len(f.Tags) != 1 || f.Tags[0].Name != "api" {
		t.Errorf("Tags = %+v, want [{Name: api}]", f.Tags)
	}
	if f.Timeframe == nil || f.Timeframe.Granularity != "month" {
		t.Errorf("Timeframe = %+v, want Granularity=month", f.Timeframe)
	}
	if f.Timeframe.StartDate == nil || f.Timeframe.StartDate.Format("2006-01-02") != "2023-10-01" {
		t.Errorf("Timeframe.StartDate = %v, want 2023-10-01", f.Timeframe.StartDate)
	}
	if f.ParentID != "318de52f-4e38-4c94-a550-a0d47a1f212e" {
		t.Errorf("ParentID = %q, want %q", f.ParentID, "318de52f-4e38-4c94-a550-a0d47a1f212e")
	}
	if f.CreatedAt == nil || f.CreatedAt.Format("2006-01-02") != "2023-01-01" {
		t.Errorf("CreatedAt = %v, want 2023-01-01", f.CreatedAt)
	}
	if f.UpdatedAt == nil || f.UpdatedAt.Format("2006-01-02") != "2023-01-02" {
		t.Errorf("UpdatedAt = %v, want 2023-01-02", f.UpdatedAt)
	}
}

func TestGetFeature_InvalidID(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("no HTTP call expected for an invalid UUID")
	})

	if _, err := client.GetFeature(t.Context(), "not-a-uuid"); err == nil {
		t.Fatal("expected an error for an invalid UUID, got nil")
	}
}

const listFeaturesFixture = `{
	"data": [
		{
			"id": "195a1cb2-728f-4be8-900f-aebbd84d7944",
			"type": "feature",
			"fields": {"name": "Feature One"}
		},
		{
			"id": "295a1cb2-728f-4be8-900f-aebbd84d7945",
			"type": "subfeature",
			"fields": {"name": "Feature Two"}
		}
	],
	"links": {"next": "cursor-abc"}
}`

func TestListFeatures(t *testing.T) {
	var gotQuery string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(listFeaturesFixture))
	})

	list, err := client.ListFeatures(t.Context())
	if err != nil {
		t.Fatalf("ListFeatures: %v", err)
	}

	if len(list.Features) != 2 {
		t.Fatalf("len(Features) = %d, want 2", len(list.Features))
	}
	if list.Features[0].Name != "Feature One" {
		t.Errorf("Features[0].Name = %q, want %q", list.Features[0].Name, "Feature One")
	}
	if list.Features[1].Type != "subfeature" {
		t.Errorf("Features[1].Type = %q, want %q", list.Features[1].Type, "subfeature")
	}
	if list.NextCursor != "cursor-abc" {
		t.Errorf("NextCursor = %q, want %q", list.NextCursor, "cursor-abc")
	}
	if gotQuery == "" {
		t.Error("expected request to include query parameters (type filter)")
	}
}
