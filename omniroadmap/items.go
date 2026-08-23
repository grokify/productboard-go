package omniroadmap

import (
	"context"
	"slices"

	omniroadmap "github.com/grokify/omniroadmap-core"
	"github.com/grokify/omniroadmap-core/provider"
)

// ListItems lists items across the kinds requested (or all supported kinds
// if req.Kinds is empty). ProductBoard's Entities API is cursor-paginated,
// not page-numbered — req.Page is ignored; the response's TotalCount is
// left at 0 (ProductBoard's list endpoint doesn't return a total count,
// only a next-page cursor) rather than fabricated.
func (p *Provider) ListItems(ctx context.Context, req *provider.ListItemsRequest) (*provider.ListItemsResponse, error) {
	kinds := req.Kinds
	if len(kinds) == 0 {
		kinds = p.Capabilities().Kinds
	}

	resp := &provider.ListItemsResponse{PerPage: req.PerPage}
	for _, kind := range kinds {
		items, err := p.listItemsByKind(ctx, kind)
		if err != nil {
			return nil, err
		}
		resp.Items = append(resp.Items, items...)
	}
	return resp, nil
}

func (p *Provider) listItemsByKind(ctx context.Context, kind provider.ItemKind) ([]provider.Item, error) {
	switch kind {
	case provider.ItemKindFeature:
		list, err := p.client.ListFeatures(ctx)
		if err != nil {
			return nil, wrapErr("ListFeatures", err)
		}
		items := make([]provider.Item, len(list.Features))
		for i := range list.Features {
			items[i] = *itemFromFeature(&list.Features[i])
		}
		return items, nil

	case provider.ItemKindObjective:
		list, err := p.client.ListObjectives(ctx)
		if err != nil {
			return nil, wrapErr("ListObjectives", err)
		}
		items := make([]provider.Item, len(list.Objectives))
		for i := range list.Objectives {
			items[i] = *itemFromObjective(&list.Objectives[i])
		}
		return items, nil

	default:
		if !slices.Contains(p.Capabilities().Kinds, kind) {
			return nil, omniroadmap.ErrUnsupportedOperation
		}
		return nil, nil
	}
}

// GetItem fetches a single, fully-populated item by kind and source ID.
func (p *Provider) GetItem(ctx context.Context, req *provider.GetItemRequest) (*provider.Item, error) {
	switch req.Kind {
	case provider.ItemKindFeature:
		f, err := p.client.GetFeature(ctx, req.ID)
		if err != nil {
			return nil, wrapErr("GetFeature", err)
		}
		return itemFromFeature(f), nil

	case provider.ItemKindObjective:
		o, err := p.client.GetObjective(ctx, req.ID)
		if err != nil {
			return nil, wrapErr("GetObjective", err)
		}
		return itemFromObjective(o), nil

	default:
		return nil, omniroadmap.ErrUnsupportedOperation
	}
}

func wrapErr(op string, err error) error {
	return &opError{op: op, err: err}
}

type opError struct {
	op  string
	err error
}

func (e *opError) Error() string { return "omniroadmap/productboard: " + e.op + ": " + e.err.Error() }
func (e *opError) Unwrap() error { return e.err }
