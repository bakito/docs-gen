package cli

import (
	"context"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/bakito/docs-gen/internal/common"
)

// UpdateDocumentation Updates the documentation of the environment variables of the given type.
func UpdateDocumentation(start, end, dir, command string, args ...string) common.UpdateDocsFunc {
	return func(fileContent string) string {
		ul := slog.With("command", command, "args", args)
		ul.Info("Generating cli documentation")
		cfg := config{
			Config:  common.NewConfig(start, end),
			Command: command,
			Dir:     dir,
			Args:    args,
		}

		return updateDocumentationImpl(ul, cfg, fileContent)
	}
}

type config struct {
	common.Config
	Command string
	Args    []string
	Dir     string
}

func updateDocumentationImpl(ul *slog.Logger, cfg config, fileContent string) string {
	var buf strings.Builder
	buf.WriteString("```\n")
	writeCliDocumentation(ul, cfg, &buf)
	buf.WriteString("```\n")
	return common.UpdateDocumentationSection(cfg.Config, fileContent, buf.String())
}

func writeCliDocumentation(ul *slog.Logger, cfg config, w io.Writer) {
	cmd := exec.CommandContext(context.Background(), cfg.Command, cfg.Args...)
	cmd.Dir = cfg.Dir
	output, err := cmd.Output()
	if err != nil {
		ul.Error("Error executing cli", "error", err)
		os.Exit(1)
	}
	if _, err := w.Write(output); err != nil {
		ul.Error("Error writing CLI documentation", "error", err)
		os.Exit(1)
	}
}
