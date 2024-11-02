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
		"Generate REST client",
		`This script generates TypeScript client for the REST server.`,
		&internal.GenerateClientAction{},
		nil,
	)

	// Execute the correct command
	os.Exit(int(command.Execute(os.Stderr, cmd, os.Args[1:], command.EnvVarsArrayToMap(os.Environ()))))
}
