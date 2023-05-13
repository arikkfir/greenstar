package auth

import (
	"fmt"
	"github.com/auth0/go-auth0/management"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func CreateGetConfigHandler(config Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		auth0API, err := management.New(
			config.Auth0Domain,
			management.WithContext(c),
			management.WithClientCredentials(config.APIClientID, config.APIClientSecret),
		)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to initialize the auth0 management API client: %w", err))
			return
		}

		orgList, err := auth0API.Organization.List()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to list organizations from the auth0 management API client: %w", err))
			return
		}

		tenant := c.Query("tenant")
		if c.Request.Host == "localhost" {
			if config.Auth0LocalhostOrgName != "" {
				tenant = config.Auth0LocalhostOrgName
			} else {
				c.AbortWithError(http.StatusBadRequest, fmt.Errorf("auth0 organization name to use for localhost is not configured"))
				return
			}
		} else {
			subdomain, _, found := strings.Cut(c.Request.Host, ".")
			if !found {
				c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse tenant (1st subdomain) from host '%s'", c.Request.Host))
				return
			} else {
				tenant = subdomain
			}
		}

		var org *management.Organization
		for _, curr := range orgList.Organizations {
			if curr.GetName() == tenant {
				org = curr
				break
			}
		}
		if org == nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("organization '%s' not found", c.Request.Host))
			return
		}

		type Auth0Config struct {
			Domain      string `json:"domain"`
			AppClientID string `json:"clientId"`
		}
		type AppConfig struct {
			Organization *management.Organization `json:"organization"`
			Auth0        Auth0Config              `json:"auth0"`
		}

		c.Negotiate(http.StatusOK, gin.Negotiate{
			Offered: []string{gin.MIMEJSON, gin.MIMEYAML},
			Data: AppConfig{
				Organization: &management.Organization{
					ID:          org.ID,
					Name:        org.Name,
					DisplayName: org.DisplayName,
				},
				Auth0: Auth0Config{
					Domain:      config.Auth0Domain,
					AppClientID: config.APPClientID,
				},
			},
		})
	}
}
