package main

import (
	"github.com/arikkfir/command"
	"github.com/arikkfir/greenstar/scripts/generate-rest/internal"
	"os"
	"path/filepath"
)

func main() {
	cmd := command.MustNew(
		filepath.Base(os.Args[0]),
		"Generate REST backend",
		`This script generates REST handlers for a given model set.`,
		&internal.GenerateServerAction{},
		nil,
	)

	// Execute the correct command
	os.Exit(int(command.Execute(os.Stderr, cmd, os.Args[1:], command.EnvVarsArrayToMap(os.Environ()))))
}
