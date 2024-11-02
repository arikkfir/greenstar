package auth

import (
	"context"
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/server/util"
	"github.com/descope/go-sdk/descope/client"
	descopeLogger "github.com/descope/go-sdk/descope/logger"
)

type DescopeConfig struct {
	DescopeLogLevel              string `required:"true"`
	DescopeProjectID             string `required:"true"`
	DescopeManagementKeyToken    string `required:"true"`
	DescopeBackendAccessKeyToken string `required:"true"`
}

func (c *DescopeConfig) NewSDK(ctx context.Context) (*client.DescopeClient, error) {
	var ll descopeLogger.LogLevel
	switch c.DescopeLogLevel {
	case "none":
		ll = descopeLogger.LogNone
	case "debug":
		ll = descopeLogger.LogDebugLevel
	default:
		ll = descopeLogger.LogNone
	}
	descopeClient, err := client.NewWithConfig(&client.Config{
		ProjectID:     c.DescopeProjectID,
		ManagementKey: c.DescopeManagementKeyToken,
		LogLevel:      ll,
		Logger:        &DescopeLogger{Logger: util.Logger(ctx), LogLevel: ll},
	})
	if err != nil {
		return nil, fmt.Errorf("failed creating Descope client: %w", err)
	}
	return descopeClient, nil
}
