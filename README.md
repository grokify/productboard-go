# Productboard Go SDK

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/grokify/productboard-go/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/grokify/productboard-go/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/grokify/productboard-go/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/grokify/productboard-go/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/grokify/productboard-go/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/grokify/productboard-go/actions/workflows/go-sast-codeql.yaml
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/productboard-go
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/productboard-go
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://grokify.github.io/productboard-go
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fproductboard-go
 [loc-svg]: https://tokei.rs/b1/github/grokify/productboard-go
 [repo-url]: https://github.com/grokify/productboard-go
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/productboard-go/blob/main/LICENSE

A thin, read-only Go SDK for the [ProductBoard REST API v2](https://developer.productboard.com/reference/introduction).

## Scope (v0.1)

ProductBoard publishes an official OpenAPI 3.1.1 spec — vendored here at
`openapi/entities.json` (see `openapi/SOURCE.md` for provenance) and used
to generate `internal/api` via [ogen](https://github.com/ogen-go/ogen).

This SDK covers just enough of ProductBoard's generic Entities API to read
Features, Releases, and Objectives:

- `GetFeature` / `ListFeatures`
- `GetRelease` / `ListReleases`
- `GetObjective` / `ListObjectives`

Read-only — no create/update/delete. ProductBoard's other APIs (Notes,
Teams, Members, Webhooks, Analytics, Jira/Plugin Integrations) aren't
covered. Built primarily to back a
[`github.com/grokify/omniroadmap-core`](https://github.com/grokify/omniroadmap-core)
`Provider` adapter (see the `omniroadmap/` subpackage) — grow it as real
need arises, same as any other package in this ecosystem.

## ProductBoard's dynamic field model

ProductBoard models Features/Releases/Objectives/etc. as generic, workspace-
configurable **Entities** with a `fields` map keyed by field name, where
each value's shape depends on the field's type (text, number, date, status,
member, single/multi-select, timeframe, ...). The spec's `EntityFieldValue`
schema is an 18-way undiscriminated `anyOf` across every possible shape,
which `ogen` can't synthesize a typed Go sum type for — so it's patched to
opaque JSON at codegen time (see `openapi/SOURCE.md`), and `entity.go`'s
`field*` helpers do the actual typing by inspecting each value's shape at
runtime. `Feature`/`Release`/`Objective` expose only the common fields
(name, description, status, owner, tags, timeframe, timestamps) as typed
struct fields — anything else in an entity's `fields` map isn't currently
surfaced.

## Quick start

```go
client, err := productboard.NewClient(
    productboard.WithAPIToken(os.Getenv("PRODUCTBOARD_API_TOKEN")),
)
if err != nil {
    log.Fatal(err)
}

features, err := client.ListFeatures(ctx)
if err != nil {
    log.Fatal(err)
}
for _, f := range features.Features {
    fmt.Println(f.Name, f.Status)
}
```

## Development

```bash
go build ./...
go vet ./...
golangci-lint run ./...
go test ./...
```

To regenerate `internal/api` from the currently-vendored spec:

```bash
make generate
```

To deliberately re-fetch a newer spec (review the diff before regenerating):

```bash
make update-spec
```
