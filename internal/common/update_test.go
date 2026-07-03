package common

import (
	"testing"
)

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
