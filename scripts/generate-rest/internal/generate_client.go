package internal

import (
	"context"
	"fmt"
	"github.com/iancoleman/strcase"
	"log/slog"
	"os/exec"
	"path/filepath"
)

type GenerateClientAction struct {
	APIFile string `flag:"true" required:"true"`
}

func (e *GenerateClientAction) Run(ctx context.Context) error {
	api, err := readAPI(e.APIFile)
	if err != nil {
		return fmt.Errorf("failed to read API spec: %w", err)
	}

	clientDir := "frontend/src/client"
	for _, model := range api.Models {
		tmplCtx := map[string]any{
			"api":   api,
			"model": model,
		}

		for _, t := range templates.Templates() {
			if isResourceClientTemplate(t.Name()) {
				filename := filepath.Join(clientDir, strcase.ToSnake(model.Name)+".ts")
				if err := writeToFile(filename, t.Name(), tmplCtx); err != nil {
					return fmt.Errorf("failed to write template '%s': %w", t.Name(), err)
				}
			}
		}
	}

	// if err := e.formatCode(ctx); err != nil {
	// 	return fmt.Errorf("failed to format code: %w", err)
	// }

	return nil
}

func (e *GenerateClientAction) formatCode(ctx context.Context) error {
	absPath, err := filepath.Abs("frontend")
	if err != nil {
		return fmt.Errorf("failed to determine absolute path: %w", err)
	} else {
		slog.Info("Running TypeScript prettier", "path", absPath)
	}

	cmd := exec.Command("npm", "run", "format-api-client")
	cmd.Dir = absPath
	cmd.Stdout = &LogWriter{
		Context: ctx,
		Logger:  slog.Default(),
		Level:   slog.LevelInfo,
	}
	cmd.Stderr = &LogWriter{
		Context: ctx,
		Logger:  slog.Default(),
		Level:   slog.LevelWarn,
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start TypeScript prettier: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("failed to wait for TypeScript prettier: %w", err)
	}
	if _, err := cmd.Stdout.(*LogWriter).Flush(); err != nil {
		return fmt.Errorf("failed to flush TypeScript prettier stdout: %w", err)
	}
	if _, err := cmd.Stderr.(*LogWriter).Flush(); err != nil {
		return fmt.Errorf("failed to flush TypeScript prettier stderr: %w", err)
	}
	return nil
}
