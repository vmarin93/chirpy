package main

import (
	"reflect"
	"testing"
)

func TestSanitize(t *testing.T) {
	censored := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	tests := map[string]struct {
		input    string
		censored map[string]struct{}
		want     string
	}{
		"nothing to sanitize": {input: "Clean", censored: censored, want: "Clean"},
		"censor all caps":     {input: "What the FORNAX", censored: censored, want: "What the ****"},
		"censor partly caps":  {input: "KerfufflE you too", censored: censored, want: "**** you too"},
		"censor all":          {input: "kerfuffle sharbert fornax", censored: censored, want: "**** **** ****"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := sanitize(tc.input, tc.censored)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
