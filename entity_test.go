package productboard

import (
	"testing"

	"github.com/grokify/productboard-go/internal/api"
)

func fieldsFromJSON(t *testing.T, jsonStr string) api.EntityFields {
	t.Helper()
	return api.EntityFields{"raw": api.EntityFieldValue(jsonStr)}
}

func TestFieldSelectOptions_SingleAndMulti(t *testing.T) {
	tests := []struct {
		name string
		json string
		want int
	}{
		{"single object", `{"id":"1","name":"High","color":"red"}`, 1},
		{"array", `[{"id":"1","name":"api"},{"id":"2","name":"web"}]`, 2},
		{"empty array", `[]`, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := fieldsFromJSON(t, tt.json)
			got := fieldSelectOptions(fields, "raw")
			if len(got) != tt.want {
				t.Errorf("len = %d, want %d", len(got), tt.want)
			}
		})
	}
}

func TestFieldStatus_Missing(t *testing.T) {
	fields := api.EntityFields{}
	if got := fieldStatus(fields); got != nil {
		t.Errorf("fieldStatus on missing key = %+v, want nil", got)
	}
}
