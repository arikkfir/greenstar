package auth

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

func CreateVerifyTokenMiddleware(claimsCookieName, secret, expectedIssuer, expectedAudience string) func(*gin.Context) {
	return func(c *gin.Context) {
		signedJWTToken, err := c.Cookie(claimsCookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				c.AbortWithStatus(http.StatusUnauthorized)
			} else {
				c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed getting claims cookie: %w", err))
			}
			return
		}

		claims := &Token{}
		token, err := jwt.ParseWithClaims(signedJWTToken, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			} else {
				return []byte(secret), nil
			}
		})
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed parsing claims: %w", err))
			return
		}

		if !token.Valid {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid token"))
			return
		}

		if err := claims.Valid(); err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid claims: %w", err))
			return
		}

		if claims.Issuer != expectedIssuer {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("unexpected token issuer, expected '%s', got '%s'", expectedIssuer, claims.Issuer))
			return
		}

		if claims.Subject == "" {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("empty token subject"))
			return
		}

		if len(claims.Audience) == 0 {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("empty token audience"))
			return
		} else {
			found := false
			for _, aud := range claims.Audience {
				if aud == expectedAudience {
					found = true
					break
				}
			}
			if !found {
				c.AbortWithError(http.StatusForbidden, fmt.Errorf("inappropriate token audience: %v", claims.Audience))
				return
			}
		}

		if claims.ExpiresAt == nil {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("missing token expiration"))
			return
		} else if claims.ExpiresAt.Before(time.Now()) {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("token expired at '%s'", claims.ExpiresAt.Time.String()))
			return
		}

		if claims.NotBefore == nil {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("missing token start time"))
			return
		} else if claims.NotBefore.After(time.Now()) {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("token will only be active at '%s'", claims.NotBefore.Time.String()))
			return
		}

		if claims.Tenant == "" {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("missing token tenant"))
			return
		}

		if claims.AccessToken == "" {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("missing token OAuth access token"))
			return
		}

		if claims.RefreshToken == "" {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("missing token OAuth refresh token"))
			return
		}

		if claims.ID == "" {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("missing token ID"))
			return
		}

		c.Set(claimsKey, claims)
		c.Next()
	}
}
