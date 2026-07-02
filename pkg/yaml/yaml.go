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
		return updateDocumentationImpl[T](cfg, fileContent, nil)
	}
}

// UpdateDocumentationWithCustomizer Updates the yaml documentation of the given type.
func UpdateDocumentationWithCustomizer[T any](start, end string, pc PrefixCustomizer) common.UpdateDocsFunc {
	return func(fileContent string) string {
		slog.Info("Generating yaml documentation")
		cfg := common.Config{
			StartMarker: start,
			EndMarker:   end,
		}
		return updateDocumentationImpl[T](cfg, fileContent, pc)
	}
}

func updateDocumentationImpl[T any](cfg common.Config, fileContent string, pc PrefixCustomizer) string {
	var buf strings.Builder
	buf.WriteString("```yaml\n")
	writeYAMLDocumentation(&buf, Prefix{
		FieldType: reflect.TypeFor[T](), First: "", Other: "",
	}, pc)
	buf.WriteString("```\n")

	return common.UpdateDocumentationSection(cfg, fileContent, buf.String())
}

func writeYAMLDocumentation(w io.Writer, prefix Prefix, pc PrefixCustomizer) {
	if prefix.FieldType.Kind() == reflect.Pointer {
		prefix.FieldType = prefix.FieldType.Elem()
	}
	if prefix.FieldType.Kind() != reflect.Struct {
		return
	}

	var i int
	for _, field := range reflect.VisibleFields(prefix.FieldType) {
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

		pf := prefix.Other
		if i == 0 {
			pf = prefix.First
		}

		fieldPrefix := Prefix{
			FieldType: ft,
			First:     pf + "  ",
			Other:     prefix.Other + "  ",
		}

		if pc != nil {
			pc(yamlTag, &fieldPrefix)
		}
		if yamlTag != "" {
			docs := field.Tag.Get(common.TagDocs)
			fieldType := fieldTypeString(ft)
			fmt.Fprintf(w, "%s# %s (%s)\n", pf, docs, fieldType)
			fmt.Fprintf(w, "%s%s:\n", pf, yamlTag)
			i++
		}

		if ft.Kind() == reflect.Struct && ft.Name() != "Time" {
			writeYAMLDocumentation(w, fieldPrefix, pc)
		}
	}
}

type Prefix struct {
	FieldType reflect.Type
	First     string
	Other     string
}

// PrefixCustomizer is a function that can be used to customize the prefix of a field.
type PrefixCustomizer = func(yamlTag string, prefix *Prefix)

func fieldTypeString(ft reflect.Type) string {
	if ft.Kind() == reflect.Map {
		return fmt.Sprintf("map[%s:%s]", ft.Key().Kind().String(), fieldTypeString(ft.Elem()))
	} else if ft.Kind() == reflect.Slice {
		return "[]" + fieldTypeString(ft.Elem())
	}
	return ft.Kind().String()
}
