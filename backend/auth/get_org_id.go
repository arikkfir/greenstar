package auth

import (
	"fmt"
	"github.com/auth0/go-auth0/management"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateGetOrgIDHandler(config Config) func(c *gin.Context) {
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

		for _, org := range orgList.Organizations {
			if org.GetName() == c.Param("tenant") {
				c.Negotiate(http.StatusOK, gin.Negotiate{
					Offered: []string{gin.MIMEJSON, gin.MIMEYAML},
					Data: struct {
						ID          string `json:"id"`
						Name        string `json:"name"`
						DisplayName string `json:"displayName"`
					}{ID: org.GetID(), Name: org.GetName(), DisplayName: org.GetDisplayName()},
				})
				return
			}
		}
		c.AbortWithStatus(http.StatusNotFound)
	}
}
