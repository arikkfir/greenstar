package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/arikkfir/greenstar/backend/util/redisutil"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
	"time"
)

func (a *Authenticator) InitiateLogin(c *gin.Context) {
	postLoginURL := c.Query("rurl")
	if postLoginURL == "" {
		postLoginURL = a.DefaultPostLoginURL
	}

	claimsCookieValue, err := c.Cookie(a.Config.ClaimsCookieName)
	if err != nil {
		if !errors.Is(err, http.ErrNoCookie) {
			// failed fetching cookie (and not because it's missing, something else went wrong)
			c.SetCookie(a.Config.ClaimsCookieName, "", 0, "/", "", a.SecureCookies, true)
		}
	} else {
		claims := jwt.RegisteredClaims{}
		_, err := jwt.ParseWithClaims(claimsCookieValue, &claims, a.createTokenVerificationCallback)
		if err != nil {
			c.SetCookie(a.Config.ClaimsCookieName, "", 0, "/", "", a.SecureCookies, true)
		} else if !claims.VerifyAudience("greenstar.auth", true) {
			c.SetCookie(a.Config.ClaimsCookieName, "", 0, "/", "", a.SecureCookies, true)
		} else if !claims.VerifyExpiresAt(time.Now(), true) {
			c.SetCookie(a.Config.ClaimsCookieName, "", 0, "/", "", a.SecureCookies, true)
		} else if !claims.VerifyIssuedAt(time.Now(), true) {
			c.SetCookie(a.Config.ClaimsCookieName, "", 0, "/", "", a.SecureCookies, true)
		} else if !claims.VerifyIssuer("greenstar.auth", true) {
			c.SetCookie(a.Config.ClaimsCookieName, "", 0, "/", "", a.SecureCookies, true)
		} else if !claims.VerifyNotBefore(time.Now(), true) {
			c.SetCookie(a.Config.ClaimsCookieName, "", 0, "/", "", a.SecureCookies, true)
		} else {
			session := Session{}
			r := redisutil.GetRedis(c)
			if result := r.Do(c, r.B().Get().Key("session:"+claims.ID).Build()); result.Error() != nil {
				c.SetCookie(a.Config.StateCookieName, "", 0, "/", "", a.SecureCookies, true)
				c.SetCookie(a.Config.ClaimsCookieName, "", 0, "/", "", a.SecureCookies, true)
				c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed getting session '%s' from redis: %w", claims.ID, err))
				return
			} else if !result.IsCacheHit() {
				c.SetCookie(a.Config.ClaimsCookieName, "", 0, "/", "", a.SecureCookies, true)
			} else if err := result.DecodeJSON(&session); err != nil {
				c.SetCookie(a.Config.StateCookieName, "", 0, "/", "", a.SecureCookies, true)
				c.SetCookie(a.Config.ClaimsCookieName, "", 0, "/", "", a.SecureCookies, true)
				c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed decoding session from redis: %w", err))
				return
			} else if session.Claims.Issuer != claims.Issuer {
				c.SetCookie(a.Config.ClaimsCookieName, "", 0, "/", "", a.SecureCookies, true)
			} else if session.Claims.Subject != claims.Subject {
				c.SetCookie(a.Config.ClaimsCookieName, "", 0, "/", "", a.SecureCookies, true)
			} else if strings.Join(session.Claims.Audience, ",") != strings.Join(claims.Audience, ",") {
				c.SetCookie(a.Config.ClaimsCookieName, "", 0, "/", "", a.SecureCookies, true)
			} else if session.Claims.ExpiresAt.UnixNano() != claims.ExpiresAt.UnixNano() {
				c.SetCookie(a.Config.ClaimsCookieName, "", 0, "/", "", a.SecureCookies, true)
			} else if session.Claims.NotBefore.UnixNano() != claims.NotBefore.UnixNano() {
				c.SetCookie(a.Config.ClaimsCookieName, "", 0, "/", "", a.SecureCookies, true)
			} else if session.Claims.IssuedAt.UnixNano() != claims.IssuedAt.UnixNano() {
				c.SetCookie(a.Config.ClaimsCookieName, "", 0, "/", "", a.SecureCookies, true)
			} else if session.Claims.ID != claims.ID {
				c.SetCookie(a.Config.ClaimsCookieName, "", 0, "/", "", a.SecureCookies, true)
			} else if !session.Token.Valid() {
				c.SetCookie(a.Config.ClaimsCookieName, "", 0, "/", "", a.SecureCookies, true)
			} else {
				c.SetCookie(a.Config.StateCookieName, "", 0, "/", "", a.SecureCookies, true)
				c.Redirect(http.StatusTemporaryRedirect, postLoginURL)
				return
			}
		}
	}

	state := uuid.New().String() + "|" + base64.URLEncoding.EncodeToString([]byte(postLoginURL))
	c.SetCookie(a.Config.StateCookieName, state, int((5 * time.Minute).Seconds()), "/", "", a.SecureCookies, true)
	c.SetCookie(a.Config.ClaimsCookieName, "", 0, "/", "", a.SecureCookies, true)
	c.Redirect(http.StatusTemporaryRedirect, a.OAuth.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce))
}
