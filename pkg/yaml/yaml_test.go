package yaml

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
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
