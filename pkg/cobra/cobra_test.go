package cobra

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bakito/docs-gen/internal/tests"
)

type innerStruct struct {
	Inner string `docs:"Doc Inner" docs-cli:"inner"`
}

type baseStructCobra struct {
	BaseField string `docs:"Doc Base" docs-cli:"base"`
}

type intermediateStructCobra struct {
	baseStructCobra
	InterField int `docs:"Doc Inter" docs-cli:"inter"`
}

type extendedStructCobra struct {
	intermediateStructCobra
	TopField bool `docs:"Doc Top" docs-cli:"top"`
}

type inlineTaggedStructCobra struct {
	Base  baseStructCobra `docs:"ignored"   yaml:",inline"`
	Other string          `docs:"Doc Other" docs-cli:"other"`
}

type pointerEmbeddedStructCobra struct {
	*baseStructCobra
	Other string `docs:"Doc Other" docs-cli:"other"`
}

type unexportedFieldsBaseCobra struct {
	Exported   string `docs:"Doc Exported" docs-cli:"exported"`
	unexported string //nolint:unused // unexported field is used to verify docs generator ignores unexported struct fields
}

type unexportedFieldsExtendedCobra struct {
	unexportedFieldsBaseCobra
	Top string `docs:"Doc Top" docs-cli:"top"`
}

type testStructCobra struct {
	Field1   string `docs:"Doc 1" docs-cli:"field1"`
	Nested   innerStruct
	Replicas []string `docs:"Doc Replicas" docs-cli:"replicas"`
}

func Test_writeCobraMapping_ExtendedStruct(t *testing.T) {
	var buf bytes.Buffer
	writeCobraMapping(&buf, reflect.TypeFor[extendedStructCobra]())
	got := buf.String()

	expected := "\t`base`: `Doc Base`,\n" +
		"\t`inter`: `Doc Inter`,\n" +
		"\t`top`: `Doc Top`,\n"

	if got != expected {
		t.Errorf("writeCobraMapping() mismatch\nGot:\n%s\nExpected:\n%s", got, expected)
	}
}

func Test_writeCobraMapping_PointerEmbedded(t *testing.T) {
	var buf bytes.Buffer
	writeCobraMapping(&buf, reflect.TypeFor[pointerEmbeddedStructCobra]())
	got := buf.String()

	expected := "\t`base`: `Doc Base`,\n" +
		"\t`other`: `Doc Other`,\n"

	if got != expected {
		t.Errorf("writeCobraMapping() mismatch\nGot:\n%s\nExpected:\n%s", got, expected)
	}
}

func Test_writeCobraMapping_UnexportedFields(t *testing.T) {
	var buf bytes.Buffer
	writeCobraMapping(&buf, reflect.TypeFor[unexportedFieldsExtendedCobra]())
	got := buf.String()

	expected := "\t`exported`: `Doc Exported`,\n" +
		"\t`top`: `Doc Top`,\n"

	if got != expected {
		t.Errorf("writeCobraMapping() mismatch\nGot:\n%s\nExpected:\n%s", got, expected)
	}
}

func Test_writeCobraMapping_InlineTag(t *testing.T) {
	var buf bytes.Buffer
	writeCobraMapping(&buf, reflect.TypeFor[inlineTaggedStructCobra]())
	got := buf.String()

	expected := "\t`base`: `Doc Base`,\n" +
		"\t`other`: `Doc Other`,\n"

	if got != expected {
		t.Errorf("writeCobraMapping() mismatch\nGot:\n%s\nExpected:\n%s", got, expected)
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
