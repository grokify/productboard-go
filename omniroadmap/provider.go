// Package omniroadmap implements omniroadmap-core's provider.Provider for
// ProductBoard, wrapping productboard-go's own *productboard.Client. Since
// productboard-go is itself a from-scratch SDK (no official ProductBoard Go
// SDK exists, even though ProductBoard publishes an OpenAPI spec), this
// adapter lives as an embedded omniroadmap/ subpackage rather than a
// separate repo — the same pattern as aha-go/omniroadmap.
package omniroadmap

import (
	productboard "github.com/grokify/productboard-go"

	omniroadmap "github.com/grokify/omniroadmap-core"
	"github.com/grokify/omniroadmap-core/provider"
)

const providerName = "productboard"

// Provider implements provider.Provider using productboard-go's
// *productboard.Client.
type Provider struct {
	client *productboard.Client
}

var _ provider.Provider = (*Provider)(nil)

// NewProvider wraps client as a provider.Provider.
func NewProvider(client *productboard.Client) *Provider {
	return &Provider{client: client}
}

func (p *Provider) Name() string { return providerName }

// Close is a no-op: productboard.Client owns no closable resources beyond
// its HTTP client, which the standard library manages for us.
func (p *Provider) Close() error { return nil }

func (p *Provider) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		// productboard-go v0.1 only wraps Feature/Subfeature, Release, and
		// Objective — Initiative/KeyResult are reachable via ProductBoard's
		// generic Entities API but have no typed wrapper yet.
		Kinds:                []provider.ItemKind{provider.ItemKindFeature, provider.ItemKindObjective},
		SupportsReleases:     true,
		SupportsObjectives:   true,
		SupportsCustomFields: false,
		SupportsWrite:        false,
	}
}

func init() {
	_ = omniroadmap.RegisterProvider(providerName, func(config any) (provider.Provider, error) {
		client, ok := config.(*productboard.Client)
		if !ok {
			return nil, omniroadmap.NewAPIError(providerName, 0, "invalid_config",
				"omniroadmap/productboard: expected *productboard.Client config")
		}
		return NewProvider(client), nil
	})
}
