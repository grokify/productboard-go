package omniroadmap

import (
	"context"

	omniroadmap "github.com/grokify/omniroadmap-core"
	"github.com/grokify/omniroadmap-core/provider"
)

// ListCustomFieldDefinitions is unsupported: productboard-go v0.1 doesn't
// wrap ProductBoard's field-configuration endpoints
// (GetEntityConfiguration/ListEntityConfigurations, which are generated in
// internal/api but have no hand-written wrapper yet). Capabilities().
// SupportsCustomFields is false, matching this.
func (p *Provider) ListCustomFieldDefinitions(ctx context.Context, req *provider.ListCustomFieldDefinitionsRequest) (*provider.ListCustomFieldDefinitionsResponse, error) {
	return nil, omniroadmap.ErrUnsupportedOperation
}
