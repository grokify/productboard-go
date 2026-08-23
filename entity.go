package productboard

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/grokify/productboard-go/internal/api"
)

// uuidString converts an api.UUID (a defined type over uuid.UUID that
// doesn't inherit its String method) to its canonical string form.
func uuidString(id api.UUID) string {
	return uuid.UUID(id).String()
}

// parentIDFromRelationships finds the "parent" relationship, if any, and
// returns its target's ID.
func parentIDFromRelationships(rel api.EntityRelationships) string {
	for _, r := range rel.Data {
		if r.Type == api.EntityRelationshipTypeParent {
			return uuidString(r.Target.ID)
		}
	}
	return ""
}

// Status is a workspace-configurable status value (id + display name),
// shared by all entity types that have one (features, releases, etc.).
type Status struct {
	ID   string
	Name string
}

// Member is a workspace member reference (owner, assignee, etc.).
type Member struct {
	ID    string
	Email string
}

// SelectOption is a single value from a single/multi-select custom field.
type SelectOption struct {
	ID    string
	Name  string
	Color string
}

// Timeframe is a start/end date range with an optional granularity
// (year/quarter/month/day), used for initiative/feature/release planning.
type Timeframe struct {
	StartDate   *time.Time
	EndDate     *time.Time
	Granularity string
}

// field-value parsing helpers. api.EntityFieldValue is raw JSON (the spec's
// 18-way field-value union was patched to opaque at codegen time, since
// ogen can't synthesize a typed Go sum type for it) — these functions do
// the actual typing by shape, mirroring what ProductBoard's docs say each
// named field-value schema looks like.

func fieldString(fields api.EntityFields, key string) string {
	v, ok := fields[key]
	if !ok {
		return ""
	}
	var s string
	_ = json.Unmarshal(v, &s)
	return s
}

func fieldBool(fields api.EntityFields, key string) bool {
	v, ok := fields[key]
	if !ok {
		return false
	}
	var b bool
	_ = json.Unmarshal(v, &b)
	return b
}

// entityTimestamp converts an Entity's top-level createdAt/updatedAt
// (api.OptDateTimeFieldValue) to a *time.Time. Unlike other field values,
// these live outside the dynamic fields map — see api.Entity.
func entityTimestamp(v api.OptDateTimeFieldValue) *time.Time {
	t, ok := v.Get()
	if !ok {
		return nil
	}
	tt := time.Time(t)
	return &tt
}

func fieldStatus(fields api.EntityFields) *Status {
	v, ok := fields["status"]
	if !ok {
		return nil
	}
	var s Status
	if err := json.Unmarshal(v, &s); err != nil {
		return nil
	}
	return &s
}

func fieldMember(fields api.EntityFields, key string) *Member {
	v, ok := fields[key]
	if !ok {
		return nil
	}
	var m Member
	if err := json.Unmarshal(v, &m); err != nil {
		return nil
	}
	return &m
}

func fieldSelectOptions(fields api.EntityFields, key string) []SelectOption {
	v, ok := fields[key]
	if !ok {
		return nil
	}
	// tags/multi-select fields are arrays; single-select fields are a lone
	// object -- normalize both into a slice.
	var opts []SelectOption
	if err := json.Unmarshal(v, &opts); err == nil {
		return opts
	}
	var one SelectOption
	if err := json.Unmarshal(v, &one); err == nil {
		return []SelectOption{one}
	}
	return nil
}

func fieldTimeframe(fields api.EntityFields, key string) *Timeframe {
	v, ok := fields[key]
	if !ok {
		return nil
	}
	var raw struct {
		StartDate   string `json:"startDate"`
		EndDate     string `json:"endDate"`
		Granularity string `json:"granularity"`
	}
	if err := json.Unmarshal(v, &raw); err != nil {
		return nil
	}
	tf := &Timeframe{Granularity: raw.Granularity}
	if t, err := time.Parse("2006-01-02", raw.StartDate); err == nil {
		tf.StartDate = &t
	}
	if t, err := time.Parse("2006-01-02", raw.EndDate); err == nil {
		tf.EndDate = &t
	}
	return tf
}
