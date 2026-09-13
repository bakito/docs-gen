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
	var counter int
	writeYAMLDocumentationHelper(w, prefix, pc, &counter)
}

func writeYAMLDocumentationHelper(w io.Writer, prefix Prefix, pc PrefixCustomizer, counter *int) {
	if prefix.FieldType.Kind() == reflect.Pointer {
		prefix.FieldType = prefix.FieldType.Elem()
	}
	if prefix.FieldType.Kind() != reflect.Struct {
		return
	}

	for i := range prefix.FieldType.NumField() {
		field := prefix.FieldType.Field(i)
		if field.PkgPath != "" && !field.Anonymous {
			continue
		}

		yamlTag := field.Tag.Get(common.TagYAML)
		if yamlTag == "-" {
			continue
		}

		ft := field.Type
		if ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
		}

		if common.IsInline(field) && ft.Kind() == reflect.Struct && ft.Name() != "Time" {
			inlinePrefix := Prefix{
				FieldType: ft,
				First:     prefix.First,
				Other:     prefix.Other,
			}
			writeYAMLDocumentationHelper(w, inlinePrefix, pc, counter)
			continue
		}

		parts := strings.Split(yamlTag, ",")
		tagFieldName := parts[0]

		newPrefix := Prefix{
			FieldType: ft,
		}

		pf := prefix.Other
		if *counter == 0 {
			pf = prefix.First
		}

		newPrefix.First = pf + "  "
		newPrefix.Other = prefix.Other + "  "

		if tagFieldName != "" {
			if pc != nil {
				pc(tagFieldName, &newPrefix)
			}

			docs := field.Tag.Get(common.TagDocs)
			fieldType := fieldTypeString(newPrefix.FieldType)
			fmt.Fprintf(w, "%s# %s (%s)\n", strings.ReplaceAll(pf, prefix.First, prefix.Other), docs, fieldType)
			fmt.Fprintf(w, "%s%s:\n", pf, tagFieldName)
			(*counter)++
		}

		if newPrefix.FieldType.Kind() == reflect.Struct && newPrefix.FieldType.Name() != "Time" {
			var nestedCounter int
			writeYAMLDocumentationHelper(w, newPrefix, pc, &nestedCounter)
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
