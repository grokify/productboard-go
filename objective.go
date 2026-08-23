package productboard

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/grokify/productboard-go/internal/api"
)

// Objective represents a ProductBoard objective (OKR-style goal).
type Objective struct {
	ID          string
	Name        string
	Description string
	Status      *Status
	Owner       *Member
	Timeframe   *Timeframe
	Archived    bool
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

// ObjectiveList is a page of objectives, with a cursor for the next page.
type ObjectiveList struct {
	Objectives []Objective
	NextCursor string
}

// GetObjective retrieves an objective by ID.
func (c *Client) GetObjective(ctx context.Context, id string) (*Objective, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, &APIError{StatusCode: 400, Message: "invalid objective id: " + err.Error()}
	}

	res, err := c.apiClient.GetEntity(ctx, api.GetEntityParams{
		ID:     api.UUID(parsedID),
		Fields: []string{"all"},
	})
	if err != nil {
		return nil, wrapError("GetObjective", err)
	}

	resp, ok := res.(*api.EntityResponse)
	if !ok {
		return nil, &APIError{StatusCode: 404, Message: "objective not found"}
	}
	return objectiveFromEntity(resp.Data), nil
}

// ListObjectives lists objectives.
func (c *Client) ListObjectives(ctx context.Context, opts ...ListOption) (*ObjectiveList, error) {
	resp, err := c.listEntities(ctx, "ListObjectives",
		[]api.EntityType{api.EntityTypeObjective}, applyListOptions(opts...))
	if err != nil {
		return nil, err
	}

	list := &ObjectiveList{Objectives: make([]Objective, len(resp.Data)), NextCursor: nextCursor(resp.Links)}
	for i, e := range resp.Data {
		list.Objectives[i] = *objectiveFromEntity(e)
	}
	return list, nil
}

func objectiveFromEntity(e api.Entity) *Objective {
	o := &Objective{ID: uuidString(e.ID)}
	fields, ok := e.Fields.Get()
	if !ok {
		return o
	}
	o.Name = fieldString(fields, "name")
	o.Description = fieldString(fields, "description")
	o.Status = fieldStatus(fields)
	o.Owner = fieldMember(fields, "owner")
	o.Timeframe = fieldTimeframe(fields, "timeframe")
	o.Archived = fieldBool(fields, "archived")
	o.CreatedAt = fieldDateTime(fields, "createdAt")
	o.UpdatedAt = fieldDateTime(fields, "updatedAt")
	return o
}
