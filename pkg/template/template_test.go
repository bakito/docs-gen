package template

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bakito/docs-gen/internal/tests"
)

func Test_UpdateDocumentation_GoldenFile(t *testing.T) {
	docTests := []struct {
		name         string
		data         map[string]string
		inputFile    string
		expectedFile string
		startMarker  string
		endMarker    string
		tpl          string
		tplPrefix    string
		tplSuffix    string
	}{
		{
			name:         "template golden file",
			inputFile:    "README.md",
			expectedFile: "README.golden.md",
			startMarker:  "<!-- template-doc-start -->",
			endMarker:    "<!-- template-doc-end -->",

			data:      map[string]string{"A": "a", "B": "b", "C": "c"},
			tpl:       "| {{ .Key }} | {{ .Value }} |\n",
			tplPrefix: "| a key | a value |\n| ------ | ----------- |\n",
		},
	}

	for _, tt := range docTests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := os.ReadFile(filepath.Join("..", "..", "testdata", "template", tt.inputFile))
			if err != nil {
				t.Fatalf("read input golden file: %v", err)
			}

			got := UpdateDocumentation[string](
				tt.data,
				tt.startMarker,
				tt.endMarker,
				tt.tpl,
				WithPrefix(tt.tplPrefix),
				WithSuffix(tt.tplSuffix),
			)(string(input))

			expectedPath := filepath.Join("..", "..", "testdata", "template", tt.expectedFile)

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
