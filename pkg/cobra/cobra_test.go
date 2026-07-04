package cobra

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bakito/docs-gen/internal/tests"
)

type innerStruct struct {
	Inner string `docs:"Doc Inner" docs-cli:"inner"`
}

type testStructCobra struct {
	Field1   string `docs:"Doc 1" docs-cli:"field1"`
	Nested   innerStruct
	Replicas []string `docs:"Doc Replicas" docs-cli:"replicas"`
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
			name:         "cobra golden file",
			inputFile:    "docs_input.go",
			expectedFile: "docs_golden.go",
			startMarker:  "// cobra-doc-start",
			endMarker:    "// cobra-doc-end",
		},
	}

	for _, tt := range docTests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := os.ReadFile(filepath.Join("..", "..", "testdata", "cobra", tt.inputFile))
			if err != nil {
				t.Fatalf("read input golden file: %v", err)
			}

			got := UpdateDocumentation[testStructCobra](
				tt.startMarker,
				tt.endMarker,
			)(string(input))

			expectedPath := filepath.Join("..", "..", "testdata", "cobra", tt.expectedFile)

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
