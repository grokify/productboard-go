package productboard

import (
	"context"

	"github.com/grokify/productboard-go/internal/api"
)

// listEntities is the shared implementation behind ListFeatures/
// ListReleases/ListObjectives: build request params for the given entity
// type(s), call the generic ListEntities endpoint, and type-assert the
// success response.
func (c *Client) listEntities(ctx context.Context, operation string, types []api.EntityType, cfg *ListOptions) (*api.EntityListResponse, error) {
	params := api.ListEntitiesParams{
		Type:   types,
		Fields: []string{"all"},
	}
	if cfg.PageCursor != "" {
		params.PageCursor = api.NewOptString(cfg.PageCursor)
	}
	if cfg.Name != "" {
		params.Name = api.NewOptString(cfg.Name)
	}

	res, err := c.apiClient.ListEntities(ctx, params)
	if err != nil {
		return nil, wrapError(operation, err)
	}

	resp, ok := res.(*api.EntityListResponse)
	if !ok {
		return nil, &APIError{StatusCode: 500, Message: "unexpected response: entities not returned"}
	}
	return resp, nil
}

// nextCursor extracts the pagination cursor from a list response's links,
// or "" if there's no next page.
func nextCursor(links api.OptListLinks) string {
	l, ok := links.Get()
	if !ok || l.Next.Null {
		return ""
	}
	return l.Next.Value
}
