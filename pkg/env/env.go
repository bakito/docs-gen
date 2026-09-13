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
		return updateDocumentationImpl[T](cfg, fileContent, nil)
	}
}

// UpdateDocumentationWithCustomizer Updates the documentation of the environment variables of the given type.
func UpdateDocumentationWithCustomizer[T any](start, end string, etc TagCustomizer) common.UpdateDocsFunc {
	return func(fileContent string) string {
		slog.Info("Generating environment variables documentation")
		cfg := common.Config{
			StartMarker: start,
			EndMarker:   end,
		}
		return updateDocumentationImpl[T](cfg, fileContent, etc)
	}
}

func updateDocumentationImpl[T any](cfg common.Config, fileContent string, etc TagCustomizer) string {
	entries := collectEnvEntries(reflect.TypeFor[T](), "", etc)

	w0, w1, w2 := len("Name"), len("Type"), len("Description")
	for _, e := range entries {
		if len(e.name) > w0 {
			w0 = len(e.name)
		}
		if len(e.typ) > w1 {
			w1 = len(e.typ)
		}
		if len(e.docs) > w2 {
			w2 = len(e.docs)
		}
	}

	var buf strings.Builder
	fmt.Fprintf(&buf, "| %-*s | %-*s | %-*s |\n", w0, "Name", w1, "Type", w2, "Description")
	fmt.Fprintf(
		&buf,
		"| :%-*s | %-*s | :%-*s |\n",
		w0-1,
		strings.Repeat("-", w0-1),
		w1,
		strings.Repeat("-", w1),
		w2-1,
		strings.Repeat("-", w2-1),
	)

	for _, e := range entries {
		fmt.Fprintf(&buf, "| %-*s | %-*s | %-*s |\n", w0, e.name, w1, e.typ, w2, e.docs)
	}

	return common.UpdateDocumentationSection(cfg, fileContent, buf.String())
}

type envEntry struct {
	name string
	typ  string
	docs string
}

func writeEnvDocumentation(w io.Writer, t reflect.Type, prefix string, etc TagCustomizer) {
	entries := collectEnvEntries(t, prefix, etc)
	w0, w1, w2 := 0, 0, 0
	for _, e := range entries {
		if len(e.name) > w0 {
			w0 = len(e.name)
		}
		if len(e.typ) > w1 {
			w1 = len(e.typ)
		}
		if len(e.docs) > w2 {
			w2 = len(e.docs)
		}
	}

	for _, e := range entries {
		fmt.Fprintf(w, "| %-*s | %-*s | %-*s |\n", w0, e.name, w1, e.typ, w2, e.docs)
	}
}

func collectEnvEntries(t reflect.Type, prefix string, etc TagCustomizer) []envEntry {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}

	var entries []envEntry
	for i := range t.NumField() {
		field := t.Field(i)
		if field.PkgPath != "" && !field.Anonymous {
			continue
		}

		envTag := field.Tag.Get(common.TagEnv)
		if etc != nil {
			envTag = etc(envTag, field)
		}

		ft := field.Type
		if ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
		}

		if common.IsInline(field) && ft.Kind() == reflect.Struct && ft.Name() != "Time" {
			entries = append(entries, collectEnvEntries(ft, prefix, etc)...)
			continue
		}

		combinedTag := buildCombinedTag(prefix, envTag)

		if ft.Kind() == reflect.Struct && ft.Name() != "Time" {
			entries = append(entries, collectEnvEntries(ft, strings.TrimSuffix(combinedTag, "_"), etc)...)
		} else if envTag != "" {
			envVar := strings.Trim(combinedTag, "_") + " (" + ft.Kind().String() + ")"
			docs := field.Tag.Get(common.TagDocs)
			entries = append(entries, envEntry{
				name: envVar,
				typ:  ft.Kind().String(),
				docs: docs,
			})
		}
	}
	return entries
}

type TagCustomizer = func(envTag string, field reflect.StructField) string

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
