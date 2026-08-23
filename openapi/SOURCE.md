# Vendored spec provenance

`entities.json` is a **vendored, pinned copy** of ProductBoard's own official
OpenAPI spec — not hand-authored.

- **Source URL**: https://developer.productboard.com/openapi/entities.yaml
  (despite the `.yaml` URL, the actual response body is JSON — saved here as
  `entities.json` to match its real content type)
- **Fetched**: 2026-08-12
- **OpenAPI version**: 3.1.1
- **Spec info.version**: 2.0.0 ("Entities")
- **Server / base URL**: `https://api.productboard.com/v2`

ProductBoard publishes 8 other spec files covering Notes, Teams, Members,
Jira Integrations, Plugin Integrations, Webhooks, and Analytics — only
`entities.json` is vendored here for v0.1 (covers product, component,
feature, subfeature, initiative, objective, keyResult, release,
releaseGroup, user, company). See the repo root README for scope.

To deliberately re-fetch a newer version, run `make update-spec` (updates
this file's URL target's content in place — review the diff before running
`make generate`, don't regenerate on an unreviewed spec change).

## Manual patches applied after fetching (must be reapplied after `make update-spec`)

The raw fetched file didn't generate cleanly with `ogen` as-is. Two patches
were applied by hand to the vendored copy — if you re-fetch, check whether
they're still needed and reapply if so:

1. **OpenAPI 3.1 nullable-type-array → 3.0-style `nullable: true`.** 6
   schema properties used JSON Schema 2020-12's `"type": ["string", "null"]`
   syntax, which this version of `ogen` fails to parse (`cannot unmarshal
   !!seq into string`). Converted each to `"type": "string", "nullable":
   true` (the OpenAPI 3.0-style `aha-go`'s own spec already uses).
2. **`EntityFieldValue` (the dynamic field-value union) → opaque schema.**
   The original schema is an 18-way undiscriminated `anyOf` across every
   possible field-value shape (Boolean/Date/Member/Status/MultiSelect/...),
   which `ogen` can't synthesize a typed Go sum type for
   (`complex anyOf not implemented` / `can't infer fields discriminator`).
   Replaced its body with an empty/no-constraint schema so `ogen` generates
   it as opaque JSON (`type EntityFieldValue jx.Raw`) — the hand-written
   wrapper (`entity.go`'s `field*` helpers) does the actual typing by
   inspecting each value's shape at runtime instead.

`ogen.yml`'s `ignore_not_implemented` list (`complex anyOf`,
`discriminator inference`) is also required for the *other*, lower-value
operations this patch doesn't touch (Create/Update/Search entities, List
field values) — those stay unimplemented in `internal/api` since they're
out of this SDK's read-only v0.1 scope anyway.
