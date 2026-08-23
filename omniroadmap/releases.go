package omniroadmap

import (
	"context"

	productboard "github.com/grokify/productboard-go"

	"github.com/grokify/omniroadmap-core/provider"
)

// ListReleases lists all releases. Cursor-paginated like ListItems —
// req.Page is ignored, every page is fetched internally, and TotalCount is
// left at 0 (not returned by ProductBoard's API).
func (p *Provider) ListReleases(ctx context.Context, req *provider.ListReleasesRequest) (*provider.ListReleasesResponse, error) {
	releases, err := paginate(func(cursor string) ([]provider.Release, string, error) {
		list, err := p.client.ListReleases(ctx, productboard.WithPageCursor(cursor))
		if err != nil {
			return nil, "", wrapErr("ListReleases", err)
		}
		batch := make([]provider.Release, len(list.Releases))
		for i := range list.Releases {
			batch[i] = releaseFromProductboard(&list.Releases[i])
		}
		return batch, list.NextCursor, nil
	})
	if err != nil {
		return nil, err
	}
	return &provider.ListReleasesResponse{Releases: releases, PerPage: req.PerPage}, nil
}
