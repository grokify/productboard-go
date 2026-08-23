package productboard

import (
	"net/http"
	"testing"
)

const getObjectiveFixture = `{
	"data": {
		"id": "495a1cb2-728f-4be8-900f-aebbd84d7947",
		"type": "objective",
		"fields": {
			"name": "Grow enterprise adoption",
			"description": "Increase enterprise customer count",
			"owner": {"id": "ow-2", "email": "pm@example.com"}
		}
	}
}`

func TestGetObjective(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(getObjectiveFixture))
	})

	o, err := client.GetObjective(t.Context(), "495a1cb2-728f-4be8-900f-aebbd84d7947")
	if err != nil {
		t.Fatalf("GetObjective: %v", err)
	}
	if o.Name != "Grow enterprise adoption" {
		t.Errorf("Name = %q, want %q", o.Name, "Grow enterprise adoption")
	}
	if o.Owner == nil || o.Owner.Email != "pm@example.com" {
		t.Errorf("Owner = %+v, want Email=pm@example.com", o.Owner)
	}
}

const listObjectivesFixture = `{
	"data": [
		{"id": "495a1cb2-728f-4be8-900f-aebbd84d7947", "type": "objective", "fields": {"name": "Grow enterprise adoption"}}
	],
	"links": {"next": null}
}`

func TestListObjectives(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(listObjectivesFixture))
	})

	list, err := client.ListObjectives(t.Context())
	if err != nil {
		t.Fatalf("ListObjectives: %v", err)
	}
	if len(list.Objectives) != 1 {
		t.Fatalf("len(Objectives) = %d, want 1", len(list.Objectives))
	}
}
