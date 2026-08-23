package productboard

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/grokify/productboard-go/internal/api"
)

// Feature represents a ProductBoard feature or subfeature.
type Feature struct {
	ID          string
	Type        string // "feature" or "subfeature"
	Name        string
	Description string
	Status      *Status
	Owner       *Member
	Tags        []SelectOption
	Timeframe   *Timeframe
	ParentID    string // component/feature this belongs to, if any
	Archived    bool
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

// FeatureList is a page of features, with a cursor for the next page.
type FeatureList struct {
	Features   []Feature
	NextCursor string
}

// GetFeature retrieves a feature (or subfeature) by ID.
func (c *Client) GetFeature(ctx context.Context, id string) (*Feature, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, &APIError{StatusCode: 400, Message: "invalid feature id: " + err.Error()}
	}

	res, err := c.apiClient.GetEntity(ctx, api.GetEntityParams{
		ID:     api.UUID(parsedID),
		Fields: []string{"all"},
	})
	if err != nil {
		return nil, wrapError("GetFeature", err)
	}

	resp, ok := res.(*api.EntityResponse)
	if !ok {
		return nil, &APIError{StatusCode: 404, Message: "feature not found"}
	}
	return featureFromEntity(resp.Data), nil
}

// ListFeatures lists features and subfeatures.
func (c *Client) ListFeatures(ctx context.Context, opts ...ListOption) (*FeatureList, error) {
	resp, err := c.listEntities(ctx, "ListFeatures",
		[]api.EntityType{api.EntityTypeFeature, api.EntityTypeSubfeature}, applyListOptions(opts...))
	if err != nil {
		return nil, err
	}

	list := &FeatureList{Features: make([]Feature, len(resp.Data)), NextCursor: nextCursor(resp.Links)}
	for i, e := range resp.Data {
		list.Features[i] = *featureFromEntity(e)
	}
	return list, nil
}

func featureFromEntity(e api.Entity) *Feature {
	f := &Feature{
		ID:   uuidString(e.ID),
		Type: string(e.Type),
	}
	fields, ok := e.Fields.Get()
	if !ok {
		return f
	}
	f.Name = fieldString(fields, "name")
	f.Description = fieldString(fields, "description")
	f.Status = fieldStatus(fields)
	f.Owner = fieldMember(fields, "owner")
	f.Tags = fieldSelectOptions(fields, "tags")
	f.Timeframe = fieldTimeframe(fields, "timeframe")
	f.Archived = fieldBool(fields, "archived")
	f.CreatedAt = fieldDateTime(fields, "createdAt")
	f.UpdatedAt = fieldDateTime(fields, "updatedAt")

	if rel, ok := e.Relationships.Get(); ok {
		f.ParentID = parentIDFromRelationships(rel)
	}
	return f
}
