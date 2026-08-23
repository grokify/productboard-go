package omniroadmap

import (
	"context"

	omniroadmap "github.com/grokify/omniroadmap-core"
	"github.com/grokify/omniroadmap-core/provider"
)

// ListStatuses is unsupported: productboard-go v0.1 has no dedicated
// status-listing endpoint wrapper (ProductBoard's statuses are per-field
// configuration, reachable via GetEntityConfiguration/
// ListEntityConfigurations — not currently wrapped). Deriving statuses by
// paging every item would also be unreliable and expensive, so this
// reports unsupported rather than returning a partial/misleading list.
func (p *Provider) ListStatuses(ctx context.Context, req *provider.ListStatusesRequest) (*provider.ListStatusesResponse, error) {
	return nil, omniroadmap.ErrUnsupportedOperation
}
