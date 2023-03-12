package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/arikkfir/greenstar/backend/util/redisutil"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"net/http"
	"strings"
	"time"
)

func (a *Authenticator) HandleCallback(c *gin.Context) {
	expectedState, err := c.Cookie(a.Config.StateCookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("oauth state cookie is missing: %w", err))
		} else {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed getting oauth state cookie: %w", err))
		}
		return
	}

	stateFromProvider := c.Query("state")
	if stateFromProvider != expectedState {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("incorrect state detected in state cookie - expected '%s', found '%s'", expectedState, stateFromProvider))
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

	httpClient := a.OAuth.Client(c, token)
	defer httpClient.CloseIdleConnections()

	userInfoResponse, err := httpClient.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, gin.Error{
			Err:  fmt.Errorf("userinfo request failed: %w", err),
			Type: gin.ErrorTypePrivate,
			Meta: map[string]interface{}{
				"url":    userInfoResponse.Request.URL.String(),
				"status": userInfoResponse.Status,
			},
		})
		return
	}
	defer userInfoResponse.Body.Close()

	userInfo := GoogleAPIUserInfoResponse{}
	if userInfoResponse.StatusCode >= 200 && userInfoResponse.StatusCode <= 299 {
		responseBody := bytes.Buffer{}
		decoder := json.NewDecoder(io.TeeReader(userInfoResponse.Body, &responseBody))
		if err := decoder.Decode(&userInfo); err != nil {
			c.AbortWithError(http.StatusInternalServerError, gin.Error{
				Err:  fmt.Errorf("failed decoding google userinfo response: %w", err),
				Type: gin.ErrorTypePrivate,
				Meta: map[string]interface{}{"body": responseBody.String()},
			})
			return
		}
	} else if userInfoResponse.StatusCode == http.StatusUnauthorized {
		body, _ := io.ReadAll(userInfoResponse.Body)
		c.AbortWithError(http.StatusUnauthorized, gin.Error{
			Err:  fmt.Errorf("unauthorized to request userinfo"),
			Type: gin.ErrorTypePrivate,
			Meta: map[string]interface{}{"body": body},
		})
		return
	} else {
		body, _ := io.ReadAll(userInfoResponse.Body)
		c.AbortWithError(http.StatusInternalServerError, gin.Error{
			Err:  fmt.Errorf("failed to fetch google userinfo"),
			Type: gin.ErrorTypePrivate,
			Meta: map[string]interface{}{
				"url":    userInfoResponse.Request.URL.String(),
				"status": userInfoResponse.Status,
				"body":   body,
			},
		})
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
	} else if result := r.Do(c, r.B().Set().Key("session:"+sessionID).Value(string(sessionBytes)).Build()); result.Error() != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed persisting session to redis: %w", err))
		return
	}

	c.SetCookie(a.Config.StateCookieName, "", 0, "/", "", a.SecureCookies, true)
	c.SetCookie(a.Config.ClaimsCookieName, signedToken, int(a.Config.SessionDuration.Seconds()), "/", "", a.SecureCookies, true)
	c.Redirect(http.StatusFound, rurl)
}
