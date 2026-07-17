package mapdata

import (
	"log/slog"
	"maps"
	"slices"
	"strings"
	"text/template"

	"github.com/bakito/docs-gen/internal/common"
)

// UpdateDocumentation Updates the documentation by the provided map.
func UpdateDocumentation[T any](data map[string]T, start, end, tpl, tplPrefix, tplSuffix string) common.UpdateDocsFunc {
	t1 := template.New("tpl")
	t1, err := t1.Parse(tpl)
	if err != nil {
		slog.Error("Error parsing Template", "error", err)
		return nil
	}

	return func(fileContent string) string {
		slog.Info("Generating map docs")
		cfg := config{
			Config: common.Config{
				StartMarker: start,
				EndMarker:   end,
			},
			prefix:   tplPrefix,
			suffix:   tplSuffix,
			template: t1,
		}

		return updateDocumentationImpl(cfg, fileContent, data)
	}
}

func updateDocumentationImpl[T any](cfg config, fileContent string, data map[string]T) string {
	var buf strings.Builder

	if cfg.prefix != "" {
		buf.WriteString(cfg.prefix)
	}

	keys := slices.Collect(maps.Keys(data))
	slices.Sort(keys)

	for _, key := range keys {

		err := cfg.template.Execute(&buf, map[string]any{"Key": key, "Value": data[key]})
		if err != nil {
			slog.Error("Error rendering Template", "error", err)
			return ""
		}
	}
	if cfg.suffix != "" {
		buf.WriteString(cfg.suffix)
	}
	return common.UpdateDocumentationSection(cfg.Config, fileContent, buf.String())
}

type config struct {
	common.Config
	prefix   string
	suffix   string
	template *template.Template
}
