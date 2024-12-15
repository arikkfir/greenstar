package main

import (
	"bytes"
	"context"
	"embed"
	"flag"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v3"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

var (
	apiFile string

	//go:embed templates
	templatesFS embed.FS
	templates   []*template.Template
)

func init() {
	flag.StringVar(&apiFile, "api-file", "api.yaml", "Path to the API spec file")
	flag.Parse()

	strcase.ConfigureAcronym("ParentID", "parentID")
	pluralizeClient := pluralize.NewClient()
	functions := template.FuncMap{
		"fail": func(format string, args ...any) (string, error) { return "", fmt.Errorf(format, args...) },
		"strings": func(sa []string) string {
			b := bytes.Buffer{}
			for _, s := range sa {
				if b.Len() > 0 {
					b.WriteString(", ")
				}
				b.WriteString("\"" + strings.ReplaceAll(s, "\"", "\\\"") + "\"")
			}
			return b.String()
		},
		"toCamelCase":      strcase.ToCamel,
		"toLowerCamelCase": strcase.ToLowerCamel,
		"toSnake":          strcase.ToSnake,
		"toScreamingSnake": strcase.ToScreamingSnake,
		"toKebab":          strcase.ToKebab,
		"toScreamingKebab": strcase.ToScreamingKebab,
		"isPlural":         pluralizeClient.IsPlural,
		"toPlural":         pluralizeClient.Plural,
		"isSingular":       pluralizeClient.IsSingular,
		"toSingular":       pluralizeClient.Singular,
		"pluralize":        pluralizeClient.Pluralize,
	}

	if t, err := collectTemplates(sprig.FuncMap(), functions); err != nil {
		panic(err)
	} else {
		templates = t
	}
}

func collectTemplates(funcMaps ...template.FuncMap) ([]*template.Template, error) {
	var templates []*template.Template
	if err := fs.WalkDir(templatesFS, ".", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Check if the entry is a file and ends with .tmpl
		if !info.IsDir() && filepath.Ext(path) == ".tmpl" {
			t := template.New(path)
			for _, fm := range funcMaps {
				t = t.Funcs(fm)
			}
			if b, err := templatesFS.ReadFile(path); err != nil {
				return fmt.Errorf("failed reading template file '%s': %w", path, err)
			} else if tmpl, err := t.Parse(string(b)); err != nil {
				return fmt.Errorf("failed parsing template file '%s': %w", path, err)
			} else {
				templates = append(templates, tmpl)
			}
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to walk templates directory: %w", err)
	}
	return templates, nil
}

func main() {
	ctx := context.Background()

	slog.Info("Reading API specification", "file", apiFile)
	api, err := readAPISpec()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	errCh := make(chan error, len(templates))
	wg := sync.WaitGroup{}
	for _, tmpl := range templates {
		wg.Add(1)
		go func(ctx context.Context, api *API, tmpl *template.Template) {
			defer wg.Done()
			if err := processTemplate(ctx, *api, tmpl); err != nil {
				errCh <- fmt.Errorf("failed to process template '%s': %w", tmpl.Name(), err)
			}
		}(ctx, api, tmpl)
	}
	wg.Wait()
	close(errCh)

	failed := false
	for err := range errCh {
		slog.Error(err.Error())
		failed = true
	}

	if failed {
		os.Exit(1)
	}
}

func readAPISpec() (*API, error) {
	if err := validateAPIFile(apiFile); err != nil {
		return nil, err
	}

	f, err := os.Open(apiFile)
	if err != nil {
		return nil, fmt.Errorf("failed opening API specification file: %w", err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)

	var api API
	if err := decoder.Decode(&api); err != nil {
		return nil, fmt.Errorf("failed decoding API specification file: %w", err)
	}

	return &api, nil
}

func processTemplate(_ context.Context, api API, tmpl *template.Template) error {
	tmplName := tmpl.Name()
	if strings.Contains(tmplName, "$model") {
		for modelName, model := range api.Models {
			targetFile := strings.TrimPrefix(
				strings.TrimSuffix(
					strings.ReplaceAll(
						tmplName,
						"$model",
						strings.ToLower(modelName)),
					".tmpl",
				),
				"templates/",
			)

			slog.Info("Processing model template", "tmpl", tmplName, "model", modelName, "targetFile", targetFile)
			data := map[string]any{"api": api, "model": model}
			if err := writeTemplate(tmpl, data, targetFile); err != nil {
				return fmt.Errorf("failed to write template '%s': %w", tmplName, err)
			}
		}
	} else {
		targetFile := strings.TrimPrefix(strings.TrimSuffix(tmplName, ".tmpl"), "templates/")
		slog.Info("Processing global template", "tmpl", tmplName, "targetFile", targetFile)
		if err := writeTemplate(tmpl, api, targetFile); err != nil {
			return fmt.Errorf("failed to write template '%s': %w", tmplName, err)
		}
	}
	return nil
}

func writeTemplate(tmpl *template.Template, tmplData any, targetFile string) error {
	dir := filepath.Dir(targetFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", dir, err)
	}

	file, err := os.OpenFile(targetFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", targetFile, err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, tmplData); err != nil {
		return fmt.Errorf("failed executing template '%s': %w", tmpl.Name(), err)
	}

	return nil
}
