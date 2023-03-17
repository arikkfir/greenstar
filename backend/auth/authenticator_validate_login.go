package auth

import (
	"errors"
	"fmt"
	"github.com/arikkfir/greenstar/backend/util/redisutil"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

func (a *Authenticator) CreateMiddlewareForFilteringUnauthenticated(audience string) func(c *gin.Context) {
	return func(c *gin.Context) {
		claimsCookieValue, err := c.Cookie(a.Config.ClaimsCookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				c.AbortWithStatus(http.StatusUnauthorized)
			} else {
				c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed getting cookie: %w", err))
			}
			return
		}

		claims := jwt.RegisteredClaims{}
		session := Session{}
		r := redisutil.GetRedis(c)

		if _, err := jwt.ParseWithClaims(claimsCookieValue, &claims, a.createTokenVerificationCallback); err != nil {
			c.SetCookie(a.Config.StateCookieName, "", -1, "/", "", a.SecureCookies, true)
			c.SetCookie(a.Config.ClaimsCookieName, "", -1, "/", "", a.SecureCookies, true)
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed parsing JWT token '%s': %w", claimsCookieValue, err))
		} else if !claims.VerifyAudience(audience, true) {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("forbidden audience: %s", audience))
		} else if !claims.VerifyExpiresAt(time.Now(), true) {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("expired token"))
		} else if !claims.VerifyIssuedAt(time.Now(), true) {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("invalid issued-at value"))
		} else if !claims.VerifyIssuer("greenstar.auth", true) {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("unexpected issuer, expected '%s', got: %s", "greenstar.auth", claims.Issuer))
		} else if !claims.VerifyNotBefore(time.Now(), true) {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("token not yet active"))
		} else if result := r.Do(c, r.B().Get().Key("session:"+claims.ID).Build()); result.Error() != nil {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("failed getting session '%s' from redis: %w", claims.ID, err))
		} else if err := result.DecodeJSON(&session); err != nil {
			c.SetCookie(a.Config.StateCookieName, "", -1, "/", "", a.SecureCookies, true)
			c.SetCookie(a.Config.ClaimsCookieName, "", -1, "/", "", a.SecureCookies, true)
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed decoding session from redis: %w", err))
		} else if session.Claims.Issuer != claims.Issuer {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("issuer mismatch, expected '%s', got: %s", session.Claims.Issuer, claims.Issuer))
		} else if session.Claims.Subject != claims.Subject {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("subject mismatch, expected '%s', got: %s", session.Claims.Subject, claims.Subject))
		} else if strings.Join(session.Claims.Audience, ",") != strings.Join(claims.Audience, ",") {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("audience mismatch, expected '%s', got: %s", session.Claims.Audience, claims.Audience))
		} else if session.Claims.ExpiresAt.UnixNano() != claims.ExpiresAt.UnixNano() {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("expiresAt mismatch, expected '%s', got: %s", session.Claims.ExpiresAt, claims.ExpiresAt))
		} else if session.Claims.NotBefore.UnixNano() != claims.NotBefore.UnixNano() {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("notBefore mismatch, expected '%s', got: %s", session.Claims.NotBefore, claims.NotBefore))
		} else if session.Claims.IssuedAt.UnixNano() != claims.IssuedAt.UnixNano() {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("issuedAt mismatch, expected '%s', got: %s", session.Claims.IssuedAt, claims.IssuedAt))
		} else if session.Claims.ID != claims.ID {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("id mismatch, expected '%s', got: %s", session.Claims.ID, claims.ID))
		} else {
			setSession(c, &session)
			c.Next()
		}
	}
}
