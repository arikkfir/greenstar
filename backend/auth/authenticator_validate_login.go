package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Authenticator) CreateRequireAuthenticationMiddleware(audience string) func(c *gin.Context) {
	return func(c *gin.Context) {
		claims, err := a.parseClaimsFromCookie(c, audience)
		if err != nil {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("failed parsing claims from cookie: %w", err))
			return
		} else if claims == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		session, err := a.loadSession(c, claims.ID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed loading session for claims: %w", err))
			return
		} else if session == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		} else if err := a.verifyClaimsAndSession(claims, session); err != nil {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("failed verifying claims and session: %w", err))
			return
		}

		setSession(c, session)
		c.Next()
	}
}
