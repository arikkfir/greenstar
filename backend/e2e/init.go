package e2e

import (
	"fmt"
	"github.com/descope/go-sdk/descope/client"
	"github.com/iancoleman/strcase"
	"os"
	"path/filepath"
	"strings"
)

var (
	descopeProjectID          string
	descopeManagementKeyToken string
	descopeClient             *client.DescopeClient
	processID                 = os.Getpid()
)

func init() {
	descopeProjectID = readConfigValue("DESCOPE_PROJECT_ID")
	descopeManagementKeyToken = readConfigValue("DESCOPE_MANAGEMENT_KEY_TOKEN")

	if c, err := client.NewWithConfig(&client.Config{ProjectID: descopeProjectID, ManagementKey: descopeManagementKeyToken}); err != nil {
		panic(err)
	} else {
		descopeClient = c
	}
}

func readConfigValue(key string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}

	fileName, err := filepath.Abs("../../hack/" + strcase.ToKebab(key))
	if err != nil {
		panic(fmt.Errorf("failed to find file for environment variable '%s': %w", key, err))
	}

	b, err := os.ReadFile(fileName)
	if err != nil {
		panic(fmt.Errorf("%s was not found in neither an environment variable '%s' nor in '%s': %w", key, key, fileName, err))
	}

	return strings.TrimSpace(string(b))
}
