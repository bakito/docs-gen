package common

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"
)

const (
	TagDocs = "docs"
	TagCLI  = "docs-cli"
	TagEnv  = "env"
	TagYAML = "yaml"
)

// IsInline checks if a struct field is embedded or marked with a yaml inline tag.
func IsInline(field reflect.StructField) bool {
	yamlTag := field.Tag.Get(TagYAML)
	if yamlTag == "-" {
		return false
	}
	parts := strings.Split(yamlTag, ",")
	for _, p := range parts {
		if strings.TrimSpace(p) == "inline" {
			return true
		}
	}
	if field.Anonymous && parts[0] == "" {
		return true
	}
	return false
}

type UpdateDocsFunc func(fileName string) string

type Config struct {
	StartMarker string
	EndMarker   string
}

// NewConfig creates a new Config.
func NewConfig(startMarker, endMarker string) Config {
	return Config{
		StartMarker: startMarker,
		EndMarker:   endMarker,
	}
}

// UpdateDocumentationSection updates the content between startMarker and endMarker in fileContent with newContent.
func UpdateDocumentationSection(cfg Config, fileContent, newContent string) string {
	startIdx := strings.Index(fileContent, cfg.StartMarker)
	endIdx := strings.Index(fileContent, cfg.EndMarker)

	if startIdx == -1 || endIdx == -1 {
		slog.Error(fmt.Sprintf("Could not find markers %s and %s", cfg.StartMarker, cfg.EndMarker))
		os.Exit(1)
	}

	return fileContent[:startIdx+len(cfg.StartMarker)] + "\n" + newContent + fileContent[endIdx:]
}
