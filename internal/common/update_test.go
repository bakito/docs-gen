package common

import (
	"reflect"
	"testing"
)

type InlineBase struct{}

type isInlineTestStruct struct {
	InlineBase
	AnonInline        InlineBase `yaml:",inline"`
	AnonInlineNoComma InlineBase `yaml:"inline"`
	NamedInline       InlineBase `yaml:",inline"`
	NamedInlineOpt    InlineBase `yaml:",inline,omitempty"`
	NamedNotInline    InlineBase `yaml:"named"`
	AnonNamed         InlineBase `yaml:"anonNamed"`
	Ignored           InlineBase `yaml:"-"`
	NoTagNamed        InlineBase
}

func Test_IsInline(t *testing.T) {
	st := reflect.TypeFor[isInlineTestStruct]()
	tests := []struct {
		fieldName string
		expected  bool
	}{
		{"InlineBase", true},
		{"AnonInline", true},
		{"AnonInlineNoComma", true},
		{"NamedInline", true},
		{"NamedInlineOpt", true},
		{"NamedNotInline", false},
		{"AnonNamed", false},
		{"Ignored", false},
		{"NoTagNamed", false},
	}

	for _, tt := range tests {
		t.Run(tt.fieldName, func(t *testing.T) {
			f, ok := st.FieldByName(tt.fieldName)
			if !ok {
				t.Fatalf("field %s not found", tt.fieldName)
			}
			if got := IsInline(f); got != tt.expected {
				t.Errorf("IsInline(%s) = %v, want %v", tt.fieldName, got, tt.expected)
			}
		})
	}
}

func Test_updateDocumentationSection(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		startMarker string
		endMarker   string
		newContent  string
		expected    string
	}{
		{
			"Standard update",
			"Before\n<!-- start -->\nOld\n<!-- end -->\nAfter",
			"<!-- start -->",
			"<!-- end -->",
			"New\n",
			"Before\n<!-- start -->\nNew\n<!-- end -->\nAfter",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UpdateDocumentationSection(
				NewConfig(tt.startMarker, tt.endMarker),
				tt.fileContent,
				tt.newContent,
			); got != tt.expected {
				t.Errorf("updateDocumentationSection() = %v, want %v", got, tt.expected)
			}
		})
	}
}
