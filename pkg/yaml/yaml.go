package yaml

import (
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"strings"

	"github.com/bakito/docs-gen/internal/common"
)

// UpdateDocumentation Updates the yaml documentation of the given type.
func UpdateDocumentation[T any](start, end string) common.UpdateDocsFunc {
	return func(fileContent string) string {
		slog.Info("Generating yaml documentation")
		cfg := common.Config{
			StartMarker: start,
			EndMarker:   end,
		}
		return updateDocumentationImpl[T](cfg, fileContent)
	}
}

func updateDocumentationImpl[T any](cfg common.Config, fileContent string) string {
	var buf strings.Builder
	buf.WriteString("```yaml\n")
	writeYAMLDocumentation(&buf, reflect.TypeFor[T](), "", "")
	buf.WriteString("```\n")

	return common.UpdateDocumentationSection(cfg, fileContent, buf.String())
}

func writeYAMLDocumentation(w io.Writer, t reflect.Type, firstPrefix, otherPrefix string) {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}

	var i int
	for _, field := range reflect.VisibleFields(t) {
		if field.PkgPath != "" {
			continue
		}

		yamlTag := field.Tag.Get(common.TagYAML)
		if yamlTag == "-" {
			continue
		}
		yamlTag = strings.TrimSuffix(yamlTag, ",omitempty")

		ft := field.Type
		if ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
		}

		pf := otherPrefix
		if i == 0 {
			pf = firstPrefix
		}

		newFirstPrefix := pf + "  "
		newOtherPrefix := otherPrefix + "  "

		if yamlTag == "replicas" && ft.Kind() == reflect.Slice {
			ft = ft.Elem()
			newFirstPrefix += "- "
			newOtherPrefix += "  "
		}

		if yamlTag != "" {
			docs := field.Tag.Get(common.TagDocs)
			fieldType := fieldTypeString(ft)
			fmt.Fprintf(w, "%s# %s (%s)\n", pf, docs, fieldType)
			fmt.Fprintf(w, "%s%s:\n", pf, yamlTag)
			i++
		}

		if ft.Kind() == reflect.Struct && ft.Name() != "Time" {
			writeYAMLDocumentation(w, ft, newFirstPrefix, newOtherPrefix)
		}
	}
}

func fieldTypeString(ft reflect.Type) string {
	if ft.Kind() == reflect.Map {
		return fmt.Sprintf("map[%s:%s]", ft.Key().Kind().String(), fieldTypeString(ft.Elem()))
	} else if ft.Kind() == reflect.Slice {
		return "[]" + fieldTypeString(ft.Elem())
	}
	return ft.Kind().String()
}
