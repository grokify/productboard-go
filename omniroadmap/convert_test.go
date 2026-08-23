package omniroadmap

import (
	"testing"
	"time"

	productboard "github.com/grokify/productboard-go"

	"github.com/grokify/omniroadmap-core/provider"
)

func TestItemFromFeature(t *testing.T) {
	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)

	f := &productboard.Feature{
		ID:          "195a1cb2-728f-4be8-900f-aebbd84d7944",
		Type:        "feature",
		Name:        "Dark mode",
		Description: "Add a dark theme",
		Status:      &productboard.Status{ID: "st-1", Name: "In Progress"},
		Owner:       &productboard.Member{ID: "ow-1", Email: "pm@example.com"},
		Tags:        []productboard.SelectOption{{ID: "t1", Name: "ui"}},
		Timeframe:   &productboard.Timeframe{StartDate: &start, EndDate: &end},
		ParentID:    "318de52f-4e38-4c94-a550-a0d47a1f212e",
	}

	item := itemFromFeature(f)

	if item.ID != "productboard:195a1cb2-728f-4be8-900f-aebbd84d7944" {
		t.Errorf("ID = %q, want %q", item.ID, "productboard:195a1cb2-728f-4be8-900f-aebbd84d7944")
	}
	if item.Kind != provider.ItemKindFeature {
		t.Errorf("Kind = %q, want %q", item.Kind, provider.ItemKindFeature)
	}
	if item.Status == nil || item.Status.Name != "In Progress" {
		t.Errorf("Status = %+v, want Name=In Progress", item.Status)
	}
	if item.Status.Category != "" {
		t.Errorf("Status.Category = %q, want empty (ProductBoard doesn't expose a category)", item.Status.Category)
	}
	if item.Owner == nil || item.Owner.Email != "pm@example.com" {
		t.Errorf("Owner = %+v, want Email=pm@example.com", item.Owner)
	}
	if len(item.Tags) != 1 || item.Tags[0] != "ui" {
		t.Errorf("Tags = %v, want [ui]", item.Tags)
	}
	if item.StartDate == nil || !item.StartDate.Equal(start) {
		t.Errorf("StartDate = %v, want %v", item.StartDate, start)
	}
	if item.DueDate == nil || !item.DueDate.Equal(end) {
		t.Errorf("DueDate = %v, want %v", item.DueDate, end)
	}
	if item.ParentID != "productboard:318de52f-4e38-4c94-a550-a0d47a1f212e" {
		t.Errorf("ParentID = %q, want %q", item.ParentID, "productboard:318de52f-4e38-4c94-a550-a0d47a1f212e")
	}
	if item.Metadata["productboard.type"] != "feature" {
		t.Errorf("Metadata[productboard.type] = %v, want %q", item.Metadata["productboard.type"], "feature")
	}
}

func TestItemFromObjective(t *testing.T) {
	o := &productboard.Objective{
		ID:   "495a1cb2-728f-4be8-900f-aebbd84d7947",
		Name: "Grow enterprise adoption",
	}
	item := itemFromObjective(o)

	if item.Kind != provider.ItemKindObjective {
		t.Errorf("Kind = %q, want %q", item.Kind, provider.ItemKindObjective)
	}
	if item.Name != "Grow enterprise adoption" {
		t.Errorf("Name = %q, want %q", item.Name, "Grow enterprise adoption")
	}
}

func TestReleaseFromProductboard(t *testing.T) {
	r := &productboard.Release{
		ID:     "395a1cb2-728f-4be8-900f-aebbd84d7946",
		Name:   "MVP3",
		Status: &productboard.Status{ID: "st-2", Name: "Planned"},
	}
	rel := releaseFromProductboard(r)

	if rel.ID != "productboard:395a1cb2-728f-4be8-900f-aebbd84d7946" {
		t.Errorf("ID = %q, want %q", rel.ID, "productboard:395a1cb2-728f-4be8-900f-aebbd84d7946")
	}
	if rel.Name != "MVP3" {
		t.Errorf("Name = %q, want %q", rel.Name, "MVP3")
	}
	if rel.Status == nil || rel.Status.Name != "Planned" {
		t.Errorf("Status = %+v, want Name=Planned", rel.Status)
	}
}

func TestPersonFromMember_Nil(t *testing.T) {
	if got := personFromMember(nil); got != nil {
		t.Errorf("personFromMember(nil) = %+v, want nil", got)
	}
}
