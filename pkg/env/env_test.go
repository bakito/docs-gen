package env

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bakito/docs-gen/internal/tests"
)

func Test_buildCombinedTag(t *testing.T) {
	docstests := []struct {
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
	for _, tt := range docstests {
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

type baseStructEnv struct {
	BaseField string `docs:"Doc Base" env:"BASE"`
}

type intermediateStructEnv struct {
	baseStructEnv
	InterField int `docs:"Doc Inter" env:"INTER"`
}

type extendedStructEnv struct {
	intermediateStructEnv
	TopField bool `docs:"Doc Top" env:"TOP"`
}

type inlineTaggedStructEnv struct {
	Base  baseStructEnv `docs:"ignored"   yaml:",inline"`
	Other string        `docs:"Doc Other" env:"OTHER"`
}

type pointerEmbeddedStructEnv struct {
	*baseStructEnv
	Other string `docs:"Doc Other" env:"OTHER"`
}

type unexportedFieldsBaseEnv struct {
	Exported   string `docs:"Doc Exported" env:"EXPORTED"`
	unexported string //nolint:unused // unexported field is used to verify docs generator ignores unexported struct fields
}

type unexportedFieldsExtendedEnv struct {
	unexportedFieldsBaseEnv
	Top string `docs:"Doc Top" env:"TOP"`
}

type testStructEnv struct {
	Field1 string `docs:"Doc 1" env:"FIELD1"`
	Field2 int    `docs:"Doc 2" env:"FIELD2"`
	Nested innerStruct
	Origin string `docs:"Doc Origin"` // Special case in code
}

func Test_writeEnvDocumentation_ExtendedStruct(t *testing.T) {
	var buf bytes.Buffer
	writeEnvDocumentation(&buf, reflect.TypeFor[extendedStructEnv](), "APP", nil)
	got := buf.String()

	expected := "| APP_BASE (string) | string | Doc Base  |\n" +
		"| APP_INTER (int)   | int    | Doc Inter |\n" +
		"| APP_TOP (bool)    | bool   | Doc Top   |\n"

	if got != expected {
		t.Errorf("writeEnvDocumentation() mismatch\nGot:\n%s\nExpected:\n%s", got, expected)
	}
}

func Test_writeEnvDocumentation_PointerEmbedded(t *testing.T) {
	var buf bytes.Buffer
	writeEnvDocumentation(&buf, reflect.TypeFor[pointerEmbeddedStructEnv](), "APP", nil)
	got := buf.String()

	expected := "| APP_BASE (string)  | string | Doc Base  |\n" +
		"| APP_OTHER (string) | string | Doc Other |\n"

	if got != expected {
		t.Errorf("writeEnvDocumentation() mismatch\nGot:\n%s\nExpected:\n%s", got, expected)
	}
}

func Test_writeEnvDocumentation_UnexportedFields(t *testing.T) {
	var buf bytes.Buffer
	writeEnvDocumentation(&buf, reflect.TypeFor[unexportedFieldsExtendedEnv](), "APP", nil)
	got := buf.String()

	expected := "| APP_EXPORTED (string) | string | Doc Exported |\n" +
		"| APP_TOP (string)      | string | Doc Top      |\n"

	if got != expected {
		t.Errorf("writeEnvDocumentation() mismatch\nGot:\n%s\nExpected:\n%s", got, expected)
	}
}

func Test_writeEnvDocumentation_InlineTag(t *testing.T) {
	var buf bytes.Buffer
	writeEnvDocumentation(&buf, reflect.TypeFor[inlineTaggedStructEnv](), "APP", nil)
	got := buf.String()

	expected := "| APP_BASE (string)  | string | Doc Base  |\n" +
		"| APP_OTHER (string) | string | Doc Other |\n"

	if got != expected {
		t.Errorf("writeEnvDocumentation() mismatch\nGot:\n%s\nExpected:\n%s", got, expected)
	}
}

func Test_writeEnvDocumentation(t *testing.T) {
	var buf bytes.Buffer
	writeEnvDocumentation(&buf, reflect.TypeOf(testStructEnv{}), "PRE", nil)
	got := buf.String()

	expectedSubstrings := []string{
		"| PRE_FIELD1 (string) | string | Doc 1     |",
		"| PRE_FIELD2 (int)    | int    | Doc 2     |",
		"| PRE_INNER (string)  | string | Doc Inner |",
	}

	for _, s := range expectedSubstrings {
		if !strings.Contains(got, s) {
			t.Errorf("writeEnvDocumentation() output missing substring: %v\nGot:\n%v", s, got)
		}
	}
}

func Test_writeEnvDocumentation_WithCustomizer(t *testing.T) {
	var buf bytes.Buffer
	writeEnvDocumentation(&buf, reflect.TypeFor[baseStructEnv](), "APP", func(envTag string, _ reflect.StructField) string {
		return "CUSTOM_" + envTag
	})
	got := buf.String()

	expected := "| APP_CUSTOM_BASE (string) | string | Doc Base |\n"
	if got != expected {
		t.Errorf("writeEnvDocumentation() mismatch\nGot:\n%s\nExpected:\n%s", got, expected)
	}
}

func Test_UpdateDocumentationWithCustomizer(t *testing.T) {
	input := "<!-- env-doc-start -->\nOld\n<!-- env-doc-end -->"
	fn := UpdateDocumentationWithCustomizer[baseStructEnv](
		"<!-- env-doc-start -->",
		"<!-- env-doc-end -->",
		func(envTag string, _ reflect.StructField) string {
			return "CUSTOM_" + envTag
		},
	)
	got := fn(input)

	expected := "<!-- env-doc-start -->\n" +
		"| Name                 | Type   | Description |\n" +
		"| :------------------- | ------ | :---------- |\n" +
		"| CUSTOM_BASE (string) | string | Doc Base    |\n" +
		"<!-- env-doc-end -->"

	if got != expected {
		t.Errorf("UpdateDocumentationWithCustomizer() mismatch\nGot:\n%s\nExpected:\n%s", got, expected)
	}
}

func Test_UpdateDocumentation_GoldenFile(t *testing.T) {
	docTests := []struct {
		name         string
		inputFile    string
		expectedFile string
		startMarker  string
		endMarker    string
	}{
		{
			name:         "env golden file",
			inputFile:    "README.md",
			expectedFile: "README.golden.md",
			startMarker:  "<!-- env-doc-start -->",
			endMarker:    "<!-- env-doc-end -->",
		},
	}

	for _, tt := range docTests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := os.ReadFile(filepath.Join("..", "..", "testdata", "env", tt.inputFile))
			if err != nil {
				t.Fatalf("read input golden file: %v", err)
			}

			got := UpdateDocumentation[testStructEnv](
				tt.startMarker,
				tt.endMarker,
			)(string(input))

			expectedPath := filepath.Join("..", "..", "testdata", "env", tt.expectedFile)

			if *tests.UpdateGoldenFiles {
				if err := os.WriteFile(expectedPath, []byte(got), 0o644); err != nil {
					t.Fatalf("update expected golden file: %v", err)
				}
			}

			expected, err := os.ReadFile(expectedPath)
			if err != nil {
				t.Fatalf("read expected golden file: %v", err)
			}

			if got != string(expected) {
				t.Errorf("UpdateDocumentation() output mismatch\nGot:\n%s\nExpected:\n%s", got, string(expected))
			}
		})
	}
}
