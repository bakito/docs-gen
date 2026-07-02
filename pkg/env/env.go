package env

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
		slog.Info("Generating environment variables documentation")
		cfg := common.Config{
			StartMarker: start,
			EndMarker:   end,
		}
		return updateDocumentationImpl[T](cfg, fileContent)
	}
}

func updateDocumentationImpl[T any](cfg common.Config, fileContent string) string {
	var buf strings.Builder
	buf.WriteString("| Name | Type | Description |\n")
	buf.WriteString("| :--- | ---- |:----------- |\n")
	writeEnvDocumentation(&buf, reflect.TypeFor[T](), "")

	return common.UpdateDocumentationSection(cfg, fileContent, buf.String())
}

func writeEnvDocumentation(w io.Writer, t reflect.Type, prefix string) {
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

		envTag := field.Tag.Get(common.TagEnv)
		if envTag == "" {
			switch field.Name {
			case "Origin":
				envTag = "ORIGIN"
			case "Replica":
				envTag = "REPLICA#"
			}
		}

		combinedTag := buildCombinedTag(prefix, envTag)

		ft := field.Type
		if ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
		}

		if ft.Kind() == reflect.Struct && ft.Name() != "Time" {
			writeEnvDocumentation(w, ft, strings.TrimSuffix(combinedTag, "_"))
		} else if envTag != "" {
			envVar := strings.Trim(combinedTag, "_") + " (" + ft.Kind().String() + ")"
			docs := field.Tag.Get(common.TagDocs)
			fmt.Fprintf(w, "| %s | %s | %s |\n", envVar, ft.Kind().String(), docs)
		}
	}
}

func buildCombinedTag(prefix, envTag string) string {
	if prefix != "" && envTag != "" {
		if strings.HasPrefix(envTag, prefix+"_") {
			return envTag
		}
		return prefix + "_" + envTag
	} else if prefix != "" {
		return prefix
	}
	return envTag
}
