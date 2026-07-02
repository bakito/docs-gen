package common

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

const (
	TagDocs = "docs"
	TagEnv  = "env"
	TagYAML = "yaml"
)

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
