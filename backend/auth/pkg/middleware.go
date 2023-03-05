package auth

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

func isAudiencePresent(audience string, audiences []string) bool {
	for _, a := range audiences {
		if a == audience {
			return true
		}
	}
	return false
}

func createTokenVerificationCallback(secret string) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		} else {
			return []byte(secret), nil
		}
	}
}

func CreateObtainTokenMiddleware(claimsCookieName, secret string, useSecureCookies bool) func(*gin.Context) {
	return func(c *gin.Context) {
		signedJWTToken, err := c.Cookie(claimsCookieName)
		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				c.Error(fmt.Errorf("failed getting claims cookie: %w", err))
				c.SetCookie(claimsCookieName, "", -1, "/", "", useSecureCookies, true)
			}
			setClaims(c, nil)
			c.Next()
			return
		}

		claims := &Token{}
		if token, err := jwt.ParseWithClaims(signedJWTToken, claims, createTokenVerificationCallback(secret)); err != nil {
			c.Error(fmt.Errorf("failed parsing JWT token: %w", err))
			c.SetCookie(claimsCookieName, "", -1, "/", "", useSecureCookies, true)
			claims = nil
		} else if !token.Valid {
			c.Error(fmt.Errorf("invalid token"))
			c.SetCookie(claimsCookieName, "", -1, "/", "", useSecureCookies, true)
			claims = nil
		} else if claims.AccessToken == "" {
			c.Error(fmt.Errorf("missing access token from JWT token"))
			c.SetCookie(claimsCookieName, "", -1, "/", "", useSecureCookies, true)
			claims = nil
		} else if claims.RefreshToken == "" {
			c.Error(fmt.Errorf("missing refresh token from JWT token"))
			c.SetCookie(claimsCookieName, "", -1, "/", "", useSecureCookies, true)
			claims = nil
		}
		setClaims(c, claims)
		c.Next()
	}
}

func CreateVerifyTokenMiddleware(expectedIssuer, expectedAudience string) func(*gin.Context) {
	return func(c *gin.Context) {
		claims := GetClaims(c)
		if claims == nil {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("missing claims"))
		} else if err := claims.Valid(); err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid claims: %w", err))
		} else if claims.Issuer != expectedIssuer {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("unexpected token issuer, expected '%s', got '%s'", expectedIssuer, claims.Issuer))
		} else if claims.Subject == "" {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("empty token subject"))
		} else if claims.ExpiresAt == nil {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("missing token expiration"))
		} else if claims.ExpiresAt.Before(time.Now()) {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("token expired at '%s'", claims.ExpiresAt.Time.String()))
		} else if claims.NotBefore == nil {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("missing token start time"))
		} else if claims.NotBefore.After(time.Now()) {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("token will only be active at '%s'", claims.NotBefore.Time.String()))
		} else if !isAudiencePresent(expectedAudience, claims.Audience) {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("inappropriate token audience: %v", claims.Audience))
		} else if claims.Tenant == "" {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("missing token tenant"))
		} else if claims.AccessToken == "" {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("missing token OAuth access token"))
		} else if claims.RefreshToken == "" {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("missing token OAuth refresh token"))
		} else if claims.ID == "" {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("missing token ID"))
		} else {
			setClaims(c, claims)
			c.Next()
		}
	}
}
