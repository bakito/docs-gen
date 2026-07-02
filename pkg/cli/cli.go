package cli

import (
	"context"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/bakito/docs-gen/pkg/common"
)

type Config struct {
	common.Config
	Command string
	Args    []string
	Dir     string
}

func UpdateDocumentation(cfg Config, fileContent string) string {
	var buf strings.Builder
	buf.WriteString("```\n")
	writeCliDocumentation(cfg, &buf)
	buf.WriteString("```\n")
	return common.UpdateDocumentationSection(cfg.Config, fileContent, buf.String())
}

func writeCliDocumentation(cfg Config, w io.Writer) {
	cmd := exec.CommandContext(context.Background(), cfg.Command, cfg.Args...)
	cmd.Dir = cfg.Dir
	output, err := cmd.Output()
	if err != nil {
		slog.With("executable", cfg.Command, "args", cfg.Args, "error", err).Error("Error executing cli")
		os.Exit(1)
	}
	if _, err := w.Write(output); err != nil {
		slog.Error("Error writing CLI documentation", "error", err)
		os.Exit(1)
	}
}
