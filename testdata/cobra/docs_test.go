package cobra

import (
	"testing"
)

func Test_cflagVar(t *testing.T) {
	tests := []struct {
		name          string
		p             *int
		flagName      string
		value         int
		expectedName  string
		expectedValue int
		expectedDoc   string
	}{
		{
			name:          "field1 mapping",
			p:             new(int),
			flagName:      "field1",
			value:         42,
			expectedName:  "field1",
			expectedValue: 42,
			expectedDoc:   "Doc 1",
		},
		{
			name:          "replicas mapping",
			p:             new(int),
			flagName:      "replicas",
			value:         3,
			expectedName:  "replicas",
			expectedValue: 3,
			expectedDoc:   "Doc Replicas",
		},
		{
			name:          "unmapped field",
			p:             new(int),
			flagName:      "unknown",
			value:         99,
			expectedName:  "unknown",
			expectedValue: 99,
			expectedDoc:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotP, gotName, gotValue, gotDoc := cflagVar(tt.p, tt.flagName, tt.value)
			if gotP != tt.p {
				t.Errorf("cflagVar() gotP = %v, want %v", gotP, tt.p)
			}
			if gotName != tt.expectedName {
				t.Errorf("cflagVar() gotName = %v, want %v", gotName, tt.expectedName)
			}
			if gotValue != tt.expectedValue {
				t.Errorf("cflagVar() gotValue = %v, want %v", gotValue, tt.expectedValue)
			}
			if gotDoc != tt.expectedDoc {
				t.Errorf("cflagVar() gotDoc = %v, want %v", gotDoc, tt.expectedDoc)
			}
		})
	}
}

func Test_cflag(t *testing.T) {
	tests := []struct {
		name          string
		flagName      string
		value         string
		expectedName  string
		expectedValue string
		expectedDoc   string
	}{
		{
			name:          "field1 mapping",
			flagName:      "field1",
			value:         "test",
			expectedName:  "field1",
			expectedValue: "test",
			expectedDoc:   "Doc 1",
		},
		{
			name:          "replicas mapping",
			flagName:      "replicas",
			value:         "5",
			expectedName:  "replicas",
			expectedValue: "5",
			expectedDoc:   "Doc Replicas",
		},
		{
			name:          "unmapped field",
			flagName:      "other",
			value:         "value",
			expectedName:  "other",
			expectedValue: "value",
			expectedDoc:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotValue, gotDoc := cflag(tt.flagName, tt.value)
			if gotName != tt.expectedName {
				t.Errorf("cflag() gotName = %v, want %v", gotName, tt.expectedName)
			}
			if gotValue != tt.expectedValue {
				t.Errorf("cflag() gotValue = %v, want %v", gotValue, tt.expectedValue)
			}
			if gotDoc != tt.expectedDoc {
				t.Errorf("cflag() gotDoc = %v, want %v", gotDoc, tt.expectedDoc)
			}
		})
	}
}

func Test_cflagP(t *testing.T) {
	tests := []struct {
		name              string
		flagName          string
		shorthand         string
		value             bool
		expectedName      string
		expectedShorthand string
		expectedValue     bool
		expectedDoc       string
	}{
		{
			name:              "field1 mapping",
			flagName:          "field1",
			shorthand:         "f",
			value:             true,
			expectedName:      "field1",
			expectedShorthand: "f",
			expectedValue:     true,
			expectedDoc:       "Doc 1",
		},
		{
			name:              "replicas mapping",
			flagName:          "replicas",
			shorthand:         "r",
			value:             false,
			expectedName:      "replicas",
			expectedShorthand: "r",
			expectedValue:     false,
			expectedDoc:       "Doc Replicas",
		},
		{
			name:              "unmapped field",
			flagName:          "verbose",
			shorthand:         "v",
			value:             true,
			expectedName:      "verbose",
			expectedShorthand: "v",
			expectedValue:     true,
			expectedDoc:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotShorthand, gotValue, gotDoc := cflagP(tt.flagName, tt.shorthand, tt.value)
			if gotName != tt.expectedName {
				t.Errorf("cflagP() gotName = %v, want %v", gotName, tt.expectedName)
			}
			if gotShorthand != tt.expectedShorthand {
				t.Errorf("cflagP() gotShorthand = %v, want %v", gotShorthand, tt.expectedShorthand)
			}
			if gotValue != tt.expectedValue {
				t.Errorf("cflagP() gotValue = %v, want %v", gotValue, tt.expectedValue)
			}
			if gotDoc != tt.expectedDoc {
				t.Errorf("cflagP() gotDoc = %v, want %v", gotDoc, tt.expectedDoc)
			}
		})
	}
}
