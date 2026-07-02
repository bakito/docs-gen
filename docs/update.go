package docs

import (
	"log/slog"
	"os"

	"github.com/bakito/docs-gen/internal/common"
)

func UpdateDocumentation(fileName string, updateFuncs ...common.UpdateDocsFunc) {
	if len(updateFuncs) == 0 {
		slog.Error("No update functions provided")
		os.Exit(1)
	}

	fl := slog.With("fileName", fileName)
	fl.Info("Reading File")
	content, err := os.ReadFile(fileName)
	if err != nil {
		fl.Error("Error reading file", "error", err)
		os.Exit(1)
	}

	fileContent := string(content)
	for _, fn := range updateFuncs {
		fileContent = fn(fileContent)
	}

	fl.Info("Writing file")
	err = os.WriteFile(fileName, []byte(fileContent), 0o644)
	if err != nil {
		fl.Error("Error writing file", "error", err)
		os.Exit(1)
	}
}
