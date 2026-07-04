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

type testStructYAML struct {
	Field1   string `docs:"Doc 1" yaml:"field1"`
	Nested   innerStruct
	Replicas []string `docs:"Doc Replicas" yaml:"replicas"`
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
