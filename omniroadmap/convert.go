package omniroadmap

import (
	"fmt"

	productboard "github.com/grokify/productboard-go"

	"github.com/grokify/omniroadmap-core/provider"
)

func canonicalID(sourceID string) string {
	return fmt.Sprintf("%s:%s", providerName, sourceID)
}

// itemFromFeature converts a productboard.Feature to a canonical Item.
// Features and subfeatures both map to ItemKindFeature; the distinction is
// preserved in Metadata["productboard.type"].
func itemFromFeature(f *productboard.Feature) *provider.Item {
	item := &provider.Item{
		ID:          canonicalID(f.ID),
		Provider:    providerName,
		SourceID:    f.ID,
		Kind:        provider.ItemKindFeature,
		Name:        f.Name,
		Description: f.Description,
		Status:      statusFromProductboard(f.Status),
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
		Owner:       personFromMember(f.Owner),
		Tags:        tagNames(f.Tags),
		Metadata: map[string]any{
			"productboard.type":     f.Type,
			"productboard.archived": f.Archived,
		},
	}
	if f.ParentID != "" {
		item.ParentID = canonicalID(f.ParentID)
	}
	if f.Timeframe != nil {
		item.StartDate = f.Timeframe.StartDate
		item.DueDate = f.Timeframe.EndDate
	}
	return item
}

// itemFromObjective converts a productboard.Objective to a canonical Item.
func itemFromObjective(o *productboard.Objective) *provider.Item {
	item := &provider.Item{
		ID:          canonicalID(o.ID),
		Provider:    providerName,
		SourceID:    o.ID,
		Kind:        provider.ItemKindObjective,
		Name:        o.Name,
		Description: o.Description,
		Status:      statusFromProductboard(o.Status),
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
		Owner:       personFromMember(o.Owner),
		Metadata: map[string]any{
			"productboard.archived": o.Archived,
		},
	}
	if o.Timeframe != nil {
		item.StartDate = o.Timeframe.StartDate
		item.DueDate = o.Timeframe.EndDate
	}
	return item
}

// releaseFromProductboard converts a productboard.Release to a canonical Release.
func releaseFromProductboard(r *productboard.Release) provider.Release {
	rel := provider.Release{
		ID:       canonicalID(r.ID),
		Provider: providerName,
		SourceID: r.ID,
		Name:     r.Name,
		Status:   statusFromProductboard(r.Status),
		Metadata: map[string]any{
			"productboard.archived": r.Archived,
		},
	}
	if r.Timeframe != nil {
		rel.StartDate = r.Timeframe.StartDate
		rel.ReleaseDate = r.Timeframe.EndDate
	}
	return rel
}

// statusFromProductboard converts a productboard.Status to a canonical
// Status. ProductBoard's Status field value is just {id, name} — no
// completion flag or category is exposed via the API, unlike Aha's
// WorkflowStatus (Complete bool) or Jira's StatusCategory (todo/in-
// progress/done key) — so Category is deliberately left unset rather than
// guessed from the name.
func statusFromProductboard(s *productboard.Status) *provider.Status {
	if s == nil {
		return nil
	}
	return &provider.Status{ID: s.ID, Name: s.Name}
}

// personFromMember converts a productboard.Member to a canonical Person.
// ProductBoard's member field value is just {id, email} — no display name.
func personFromMember(m *productboard.Member) *provider.Person {
	if m == nil {
		return nil
	}
	return &provider.Person{ID: m.ID, Email: m.Email}
}

func tagNames(tags []productboard.SelectOption) []string {
	if len(tags) == 0 {
		return nil
	}
	names := make([]string, len(tags))
	for i, t := range tags {
		names[i] = t.Name
	}
	return names
}
