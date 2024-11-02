package internal

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/tools/imports"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func inferTypesFromAPIType(propertyType string) (goType, typeScriptType, oasType, oasTypeFormat string, err error) {
	switch propertyType {
	case "string":
		return propertyType, "string", "string", "", nil
	case "timestamp":
		return "time.Time", "Date", "string", "date-time", nil
	case "decimal":
		return "decimal.Decimal", "number", "number", "double", nil
	default:
		return "", "", "", "", fmt.Errorf("unsupported property type '%s'", propertyType)
	}
}

func writeToFile(filename, tmplName string, data any) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", dir, err)
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", filename, err)
	}
	if err := templates.ExecuteTemplate(file, tmplName, data); err != nil {
		file.Close()
		return fmt.Errorf("failed to execute OpenAPI spec template: %w", err)
	}
	file.Close()

	if strings.HasSuffix(filename, ".go") {
		optimizedSource, err := imports.Process(filename, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to optimize imports for '%s': %w", filename, err)
		} else if err := os.WriteFile(filename, optimizedSource, 0644); err != nil {
			return fmt.Errorf("failed to write back optimized source of '%s': %w", filename, err)
		}
	}

	return nil
}

type LogWriter struct {
	Context context.Context
	Logger  *slog.Logger
	Level   slog.Level
	Args    []any
	mutex   sync.Mutex
	buffer  bytes.Buffer
}

func (lw *LogWriter) Write(p []byte) (n int, err error) {
	lw.mutex.Lock()
	defer lw.mutex.Unlock()

	n, err = lw.buffer.Write(p)
	if err != nil {
		return n, err
	}

	// Try to process complete lines in the buffer.
	for {
		line, err := lw.buffer.ReadBytes('\n')
		if err == io.EOF {
			// partial line; write back to buffer until next call
			lw.buffer.Write(line)
			break
		} else if err != nil {
			return n, err
		}

		// Remove new-line at the end and write to logger
		if strings.HasSuffix(string(line), "\n") {
			line = line[:len(line)-1]
		}
		lw.Logger.Log(lw.Context, lw.Level, string(line), lw.Args...)
	}

	return n, nil
}

func (lw *LogWriter) Flush() (n int, err error) {
	lw.mutex.Lock()
	defer lw.mutex.Unlock()

	b := lw.buffer.Len()
	if b > 0 {
		lw.Logger.Log(lw.Context, lw.Level, lw.buffer.String(), lw.Args...)
		lw.buffer.Reset()
	}
	return b, nil
}
