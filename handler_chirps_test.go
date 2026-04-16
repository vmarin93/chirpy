package main

import (
	"reflect"
	"testing"
)

func TestSanitize(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"nothing to sanitize": {input: "Clean", want: "Clean"},
		"censor all caps":     {input: "What the FORNAX", want: "What the ****"},
		"censor partly caps":  {input: "KerfufflE you too", want: "**** you too"},
		"censor all":          {input: "kerfuffle sharbert fornax", want: "**** **** ****"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := sanitize(tc.input)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
