package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bakito/docs-gen/internal/tests"
)

func Test_UpdateDocumentation_GoldenFile(t *testing.T) {
	docTests := []struct {
		name         string
		inputFile    string
		expectedFile string
		startMarker  string
		endMarker    string
	}{
		{
			name:         "cli golden file",
			inputFile:    "README.md",
			expectedFile: "README.golden.md",
			startMarker:  "<!-- cli-doc-start -->",
			endMarker:    "<!-- cli-doc-end -->",
		},
	}

	for _, tt := range docTests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := os.ReadFile(filepath.Join("..", "..", "testdata", "cli", tt.inputFile))
			if err != nil {
				t.Fatalf("read input golden file: %v", err)
			}

			got := UpdateDocumentation(
				tt.startMarker,
				tt.endMarker,
				"../../",
				"go", "run", "internal/tests/cli/main.go", "--help",
			)(string(input))

			expectedPath := filepath.Join("..", "..", "testdata", "cli", tt.expectedFile)

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
