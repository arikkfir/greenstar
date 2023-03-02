package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	auth "github.com/arikkfir/greenstar/backend/auth/pkg"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"strings"
	"time"
)

type GoogleAPIUserInfoResponse struct {
	Email        string `json:"email"`
	FamilyName   string `json:"family_name"`
	GivenName    string `json:"given_name"`
	HostedDomain string `json:"hd"`
	ID           string `json:"id"`
	Link         string `json:"link"`
	Locale       string `json:"locale"`
	Name         string `json:"name"`
	PictureURL   string `json:"picture"`
}

func CreateAuthGoogleCallbackHandler(oauth *oauth2.Config, stateCookieName, claimsCookieName string, useSecureCookies bool, issuer string) func(*gin.Context) {
	return func(c *gin.Context) {
		expectedState, err := c.Cookie(stateCookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				c.AbortWithError(http.StatusBadRequest, fmt.Errorf("oauth state cookie is missing: %w", err))
			} else {
				c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed getting oauth state cookie: %w", err))
			}
			return
		}

		stateFromProvider, found := c.GetPostForm("state")
		if !found {
			stateFromProvider = c.Query("state")
		}
		if stateFromProvider != expectedState {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("incorrect state detected in state cookie - expected '%s', found '%s'", expectedState, stateFromProvider))
			return
		}

		var rurl string
		if _, encodedRURL, found := strings.Cut(stateFromProvider, "|"); !found {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid state, missing '|' sign: %s", stateFromProvider))
		} else if rurlBytes, err := base64.URLEncoding.DecodeString(encodedRURL); err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed URL-decoding the redirect URL inside state: %s", stateFromProvider))
		} else {
			rurl = string(rurlBytes)
		}

		code, found := c.GetPostForm("code")
		if !found || code == "" {
			code = c.Query("code")
		}

		token, err := oauth.Exchange(c, code)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed exchanging OAuth code for OAuth token: %w", err))
			return
		}

		userInfoResponse, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			c.Error(gin.Error{
				Err:  fmt.Errorf("userinfo request failed: %w", err),
				Type: gin.ErrorTypePrivate,
				Meta: map[string]interface{}{
					"url":    userInfoResponse.Request.URL.String(),
					"status": userInfoResponse.Status,
					"token":  token.AccessToken,
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
				c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed decoding userinfo response: %w", err))
				return
			}
		} else if userInfoResponse.StatusCode == http.StatusUnauthorized {
			body, _ := io.ReadAll(userInfoResponse.Body)
			c.AbortWithStatus(http.StatusUnauthorized)
			c.Error(gin.Error{
				Err:  fmt.Errorf("userinfo request failed with status code %d", userInfoResponse.StatusCode),
				Type: gin.ErrorTypePrivate,
				Meta: map[string]interface{}{
					"url":    userInfoResponse.Request.URL.String(),
					"status": userInfoResponse.Status,
					"token":  token.AccessToken,
					"body":   string(body),
				},
			})
			return
		} else {
			body, _ := io.ReadAll(userInfoResponse.Body)
			c.AbortWithStatus(http.StatusInternalServerError)
			c.Error(gin.Error{
				Err:  fmt.Errorf("userinfo request failed with status code %d", userInfoResponse.StatusCode),
				Type: gin.ErrorTypePrivate,
				Meta: map[string]interface{}{
					"url":    userInfoResponse.Request.URL.String(),
					"status": userInfoResponse.Status,
					"token":  token.AccessToken,
					"body":   string(body),
				},
			})
			return
		}

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Token{
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        uuid.NewString(),
				Issuer:    issuer,
				Subject:   "google|" + userInfo.ID,
				Audience:  []string{"greenstar.admin", "greenstar.auth", "greenstar.operations", "greenstar.public"},
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now().Add(time.Minute * -5)),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			},
			Tenant:       userInfo.HostedDomain,
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		})

		signedToken, err := jwtToken.SignedString([]byte(oauth.ClientSecret))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed signing JWT token: %w", err))
			return
		}

		// TODO: encrypt claims cookie
		c.SetCookie(stateCookieName, "", -1, "/", "", useSecureCookies, true)
		c.SetCookie(claimsCookieName, signedToken, int((24 * time.Hour).Seconds()), "/", "", useSecureCookies, true)
		c.Redirect(http.StatusFound, rurl)
	}
}
