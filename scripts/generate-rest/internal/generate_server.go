package internal

import (
	"context"
	"fmt"
	"github.com/iancoleman/strcase"
	"os"
	"path/filepath"
)

type GenerateServerAction struct {
	APIFile string `flag:"true" required:"true"`
}

func (e *GenerateServerAction) Run(_ context.Context) error {
	api, err := readAPI(e.APIFile)
	if err != nil {
		return fmt.Errorf("failed to read API spec: %w", err)
	}

	serverDir := "backend/internal/server"
	for _, model := range api.Models {
		tmplCtx := map[string]any{
			"api":   api,
			"model": model,
		}

		for _, t := range templates.Templates() {
			if resourceName, ok := isResourceServerTemplate(t.Name()); ok {
				filename := filepath.Join(serverDir, "resources", strcase.ToSnake(model.Name), resourceName)
				if filepath.Base(filename) == "handler_impl.go" {
					if _, err := os.Stat(filename); err == nil {
						continue // Never override the `handler_impl.go` file
					} else if !os.IsNotExist(err) {
						return fmt.Errorf("failed to check if file exists '%s': %w", filename, err)
					}
				}
				if err := writeToFile(filename, t.Name(), tmplCtx); err != nil {
					return fmt.Errorf("failed to write template '%s': %w", t.Name(), err)
				}
			}
		}
	}

	if err := writeToFile(filepath.Join(serverDir, "resources/openapi/openapi.yaml"), "openapi.yaml.tmpl", api); err != nil {
		return fmt.Errorf("failed to write OpenAPI spec: %w", err)
	}
	if err := writeToFile(filepath.Join(serverDir, "resources/openapi/server.go"), "openapi.go.tmpl", api); err != nil {
		return fmt.Errorf("failed to write OpenAPI server: %w", err)
	}

	if err := writeToFile(filepath.Join(serverDir, "server.go"), "server.go.tmpl", api); err != nil {
		return fmt.Errorf("failed to write server: %w", err)
	}

	return nil
}
