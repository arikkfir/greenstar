package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/arikkfir/greenstar/backend/util/redisutil"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

func (a *Authenticator) HandleCallback(c *gin.Context) {
	expectedState, err := a.parseOAuthStateFromCookie(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed getting oauth state cookie: %w", err))
		return
	} else if expectedState == nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("oauth state cookie is missing: %w", err))
		return
	}

	stateFromProvider := c.Query("state")
	if stateFromProvider != *expectedState {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("incorrect state detected in state cookie - expected '%s', found '%s'", *expectedState, stateFromProvider))
		return
	}

	var sessionID, rurl string
	if id, encodedRURL, found := strings.Cut(stateFromProvider, "|"); !found {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid state, missing '|' sign: %s", stateFromProvider))
		return
	} else if rurlBytes, err := base64.URLEncoding.DecodeString(encodedRURL); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed URL-decoding the redirect URL inside state: %s", stateFromProvider))
		return
	} else {
		sessionID = id
		rurl = string(rurlBytes)
	}

	code := c.Query("code")
	if code == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("missing code query parameter"))
		return
	}

	token, err := a.OAuth.Exchange(c, code)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed exchanging OAuth code for OAuth token: %w", err))
		return
	}

	userInfo, err := a.loadUserInfo(c, token)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed loading user info: %w", err))
		return
	}

	claims := jwt.RegisteredClaims{
		ID:        sessionID,
		Issuer:    "greenstar.auth",
		Subject:   "google|" + userInfo.ID,
		Audience:  []string{"greenstar.admin", "greenstar.auth", "greenstar.operations", "greenstar.public"},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now().Add(time.Minute * -5)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	signedToken, err := jwtToken.SignedString([]byte(a.Config.Google.ClientSecret))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed signing JWT token: %w", err))
		return
	}

	// TODO: encrypt claims cookie
	// TODO: decide set of permissions from user repository
	session := Session{
		Claims: claims,
		Token:  token,
		Tenant: userInfo.HostedDomain,
		Permissions: []string{
			PermissionAuthUserInfo,
			PermissionAdminCreateTenant,
		},
	}

	r := redisutil.GetRedis(c)
	if sessionBytes, err := json.Marshal(session); err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed marshalling session before persisting to Redis: %w", err))
		return
	} else if result := r.Do(c, r.B().Set().Key("session:"+sessionID).Value(string(sessionBytes)).ExatTimestamp(claims.ExpiresAt.Unix()).Build()); result.Error() != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed persisting session to redis: %w", err))
		return
	}

	c.SetCookie(a.Config.StateCookieName, "", 0, "/", "", a.SecureCookies, true)
	c.SetCookie(a.Config.ClaimsCookieName, signedToken, int(a.Config.SessionDuration.Seconds()), "/", "", a.SecureCookies, true)
	c.Redirect(http.StatusFound, rurl)
}
