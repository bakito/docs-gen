package env

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func Test_buildCombinedTag(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		envTag   string
		expected string
	}{
		{"Both empty", "", "", ""},
		{"Prefix only", "PREFIX", "", "PREFIX"},
		{"Tag only", "", "TAG", "TAG"},
		{"Both set", "PREFIX", "TAG", "PREFIX_TAG"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildCombinedTag(tt.prefix, tt.envTag); got != tt.expected {
				t.Errorf("buildCombinedTag() = %v, want %v", got, tt.expected)
			}
		})
	}
}

type innerStruct struct {
	Inner string `docs:"Doc Inner" env:"INNER"`
}

type testStructEnv struct {
	Field1 string `docs:"Doc 1" env:"FIELD1"`
	Field2 int    `docs:"Doc 2" env:"FIELD2"`
	Nested innerStruct
	Origin string `docs:"Doc Origin"` // Special case in code
}

func Test_writeEnvDocumentation(t *testing.T) {
	var buf bytes.Buffer
	writeEnvDocumentation(&buf, reflect.TypeOf(testStructEnv{}), "PRE", nil)
	got := buf.String()

	expectedSubstrings := []string{
		"| PRE_FIELD1 (string) | string | Doc 1 |",
		"| PRE_FIELD2 (int) | int | Doc 2 |",
		"| PRE_INNER (string) | string | Doc Inner |",
	}

	for _, s := range expectedSubstrings {
		if !strings.Contains(got, s) {
			t.Errorf("writeEnvDocumentation() output missing substring: %v\nGot:\n%v", s, got)
		}
	}
}
