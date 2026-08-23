package productboard

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/grokify/productboard-go/internal/api"
)

// Release represents a ProductBoard release.
type Release struct {
	ID        string
	Name      string
	Status    *Status
	Timeframe *Timeframe
	Archived  bool
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

// ReleaseList is a page of releases, with a cursor for the next page.
type ReleaseList struct {
	Releases   []Release
	NextCursor string
}

// GetRelease retrieves a release by ID.
func (c *Client) GetRelease(ctx context.Context, id string) (*Release, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, &APIError{StatusCode: 400, Message: "invalid release id: " + err.Error()}
	}

	res, err := c.apiClient.GetEntity(ctx, api.GetEntityParams{
		ID:     api.UUID(parsedID),
		Fields: []string{"all"},
	})
	if err != nil {
		return nil, wrapError("GetRelease", err)
	}

	resp, ok := res.(*api.EntityResponse)
	if !ok {
		return nil, &APIError{StatusCode: 404, Message: "release not found"}
	}
	return releaseFromEntity(resp.Data), nil
}

// ListReleases lists releases.
func (c *Client) ListReleases(ctx context.Context, opts ...ListOption) (*ReleaseList, error) {
	resp, err := c.listEntities(ctx, "ListReleases",
		[]api.EntityType{api.EntityTypeRelease}, applyListOptions(opts...))
	if err != nil {
		return nil, err
	}

	list := &ReleaseList{Releases: make([]Release, len(resp.Data)), NextCursor: nextCursor(resp.Links)}
	for i, e := range resp.Data {
		list.Releases[i] = *releaseFromEntity(e)
	}
	return list, nil
}

func releaseFromEntity(e api.Entity) *Release {
	r := &Release{ID: uuidString(e.ID)}
	fields, ok := e.Fields.Get()
	if !ok {
		return r
	}
	r.Name = fieldString(fields, "name")
	r.Status = fieldStatus(fields)
	r.Timeframe = fieldTimeframe(fields, "timeframe")
	r.Archived = fieldBool(fields, "archived")
	r.CreatedAt = fieldDateTime(fields, "createdAt")
	r.UpdatedAt = fieldDateTime(fields, "updatedAt")
	return r
}
