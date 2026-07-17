package cobra

import (
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"strings"

	"github.com/bakito/docs-gen/internal/common"
)

// UpdateDocumentation Updates the documentation of the environment variables of the given type.
func UpdateDocumentation[T any](start, end string) common.UpdateDocsFunc {
	return func(fileContent string) string {
		slog.Info("Generating cobra doc")
		cfg := common.Config{
			StartMarker: start,
			EndMarker:   end,
		}
		return updateDocumentationImpl[T](cfg, fileContent)
	}
}

func updateDocumentationImpl[T any](cfg common.Config, fileContent string) string {
	var buf strings.Builder
	buf.WriteString("var docsCobraMapping = map[string]string{\n")
	writeCobraMapping(&buf, reflect.TypeFor[T]())
	buf.WriteString("}\n")
	buf.WriteString(helper)
	return common.UpdateDocumentationSection(cfg, fileContent, buf.String())
}

func writeCobraMapping(w io.Writer, t reflect.Type) {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}

	for _, field := range reflect.VisibleFields(t) {
		if field.PkgPath != "" {
			continue
		}

		docsTag := field.Tag.Get(common.TagDocs)
		if docsTag == "" {
			continue
		}
		ft := field.Type
		if ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
		}
		if ft.Kind() == reflect.Struct && ft.Name() != "Time" {
			writeCobraMapping(w, ft)
		}
		cliTag := field.Tag.Get(common.TagCLI)
		if cliTag == "" {
			continue
		}

		fmt.Fprintf(w, "	`%s`: `%s`,\n", cliTag, docsTag)
	}
}

const helper = `
func cflagVar[T any](p *T, name string, value T) (pOut *T, nameOut string, valueOut T, reason string) {
	return p, name, value, docsCobraMapping[name]
}

func cflag[T any](name string, value T) (nameOut string, valueOut T, reason string) {
	return name, value, docsCobraMapping[name]
}

func cflagP[T any](name, shorthand string, value T) (nameOut, shorthandOut string, valueOut T, reason string) {
	return name, shorthand, value, docsCobraMapping[name]
}

`
