package omniroadmap

import (
	"context"

	"github.com/grokify/omniroadmap-core/provider"
)

// ListReleases lists releases. Cursor-paginated like ListItems — req.Page
// is ignored, TotalCount is left at 0 (not returned by ProductBoard's API).
func (p *Provider) ListReleases(ctx context.Context, req *provider.ListReleasesRequest) (*provider.ListReleasesResponse, error) {
	list, err := p.client.ListReleases(ctx)
	if err != nil {
		return nil, wrapErr("ListReleases", err)
	}

	releases := make([]provider.Release, len(list.Releases))
	for i := range list.Releases {
		releases[i] = releaseFromProductboard(&list.Releases[i])
	}
	return &provider.ListReleasesResponse{Releases: releases, PerPage: req.PerPage}, nil
}
