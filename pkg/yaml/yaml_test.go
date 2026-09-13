package yaml

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bakito/docs-gen/internal/tests"
)

type innerStruct struct {
	Inner string `docs:"Doc Inner" yaml:"inner"`
}

type baseStructYAML struct {
	BaseField string `docs:"Doc Base" yaml:"base"`
}

type intermediateStructYAML struct {
	baseStructYAML
	InterField int `docs:"Doc Inter" yaml:"inter"`
}

type extendedStructYAML struct {
	intermediateStructYAML
	TopField bool `docs:"Doc Top" yaml:"top"`
}

type inlineTaggedStructYAML struct {
	Base  baseStructYAML `docs:"ignored"   yaml:",inline"`
	Other string         `docs:"Doc Other" yaml:"other"`
}

type pointerEmbeddedStructYAML struct {
	*baseStructYAML
	Other string `docs:"Doc Other" yaml:"other"`
}

type unexportedFieldsBaseYAML struct {
	Exported   string `docs:"Doc Exported" yaml:"exported"`
	unexported string //nolint:unused // unexported field is used to verify docs generator ignores unexported struct fields
}

type unexportedFieldsExtendedYAML struct {
	unexportedFieldsBaseYAML
	Top string `docs:"Doc Top" yaml:"top"`
}

type testStructYAML struct {
	Field1   string `docs:"Doc 1" yaml:"field1"`
	Nested   innerStruct
	Replicas []string `docs:"Doc Replicas" yaml:"replicas"`
}

func Test_writeYAMLDocumentation_ExtendedStruct(t *testing.T) {
	var buf bytes.Buffer
	writeYAMLDocumentation(&buf, Prefix{
		FieldType: reflect.TypeFor[extendedStructYAML](), First: "", Other: "",
	}, nil)
	got := buf.String()

	expected := "# Doc Base (string)\nbase:\n# Doc Inter (int)\ninter:\n# Doc Top (bool)\ntop:\n"
	if got != expected {
		t.Errorf("writeYAMLDocumentation() mismatch\nGot:\n%s\nExpected:\n%s", got, expected)
	}
}

func Test_writeYAMLDocumentation_PointerEmbedded(t *testing.T) {
	var buf bytes.Buffer
	writeYAMLDocumentation(&buf, Prefix{
		FieldType: reflect.TypeFor[pointerEmbeddedStructYAML](), First: "", Other: "",
	}, nil)
	got := buf.String()

	expected := "# Doc Base (string)\nbase:\n# Doc Other (string)\nother:\n"
	if got != expected {
		t.Errorf("writeYAMLDocumentation() mismatch\nGot:\n%s\nExpected:\n%s", got, expected)
	}
}

func Test_writeYAMLDocumentation_UnexportedFields(t *testing.T) {
	var buf bytes.Buffer
	writeYAMLDocumentation(&buf, Prefix{
		FieldType: reflect.TypeFor[unexportedFieldsExtendedYAML](), First: "", Other: "",
	}, nil)
	got := buf.String()

	expected := "# Doc Exported (string)\nexported:\n# Doc Top (string)\ntop:\n"
	if got != expected {
		t.Errorf("writeYAMLDocumentation() mismatch\nGot:\n%s\nExpected:\n%s", got, expected)
	}
}

func Test_writeYAMLDocumentation_InlineTag(t *testing.T) {
	var buf bytes.Buffer
	writeYAMLDocumentation(&buf, Prefix{
		FieldType: reflect.TypeFor[inlineTaggedStructYAML](), First: "", Other: "",
	}, nil)
	got := buf.String()

	expected := "# Doc Base (string)\nbase:\n# Doc Other (string)\nother:\n"
	if got != expected {
		t.Errorf("writeYAMLDocumentation() mismatch\nGot:\n%s\nExpected:\n%s", got, expected)
	}
}

func Test_writeYAMLDocumentation(t *testing.T) {
	var buf bytes.Buffer
	writeYAMLDocumentation(&buf, Prefix{
		FieldType: reflect.TypeFor[testStructYAML](), First: "", Other: "",
	}, nil)
	got := buf.String()

	expectedSubstrings := []string{
		"# Doc 1 (string)",
		"field1:",
		"# Doc Inner (string)",
		"inner:",
		"# Doc Replicas ([]string)",
		"replicas:",
	}

	for _, s := range expectedSubstrings {
		if !strings.Contains(got, s) {
			t.Errorf("writeYAMLDocumentation() output missing substring: %v\nGot:\n%v", s, got)
		}
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
			name:         "yaml golden file",
			inputFile:    "README.md",
			expectedFile: "README.golden.md",
			startMarker:  "<!-- yaml-doc-start -->",
			endMarker:    "<!-- yaml-doc-end -->",
		},
	}

	for _, tt := range docTests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := os.ReadFile(filepath.Join("..", "..", "testdata", "yaml", tt.inputFile))
			if err != nil {
				t.Fatalf("read input golden file: %v", err)
			}

			got := UpdateDocumentation[testStructYAML](
				tt.startMarker,
				tt.endMarker,
			)(string(input))

			expectedPath := filepath.Join("..", "..", "testdata", "yaml", tt.expectedFile)

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
