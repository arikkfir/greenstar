package main

import (
	"encoding/base64"
	"fmt"
	auth "github.com/arikkfir/greenstar/backend/auth/pkg"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

func CreateAuthGoogleLoginHandler(oauth *oauth2.Config, stateCookieName string, defaultPostLoginURL string) func(*gin.Context) {
	curl, err := url.Parse(oauth.RedirectURL)
	if err != nil {
		panic(fmt.Errorf("failed parsing oauth2 callback URL: %w", err))
	}

	return func(c *gin.Context) {
		consentNeeded := false
		if token := auth.GetClaims(c); token != nil {
			if token.RefreshToken == "" {
				consentNeeded = true
			} else if err := token.Valid(); err != nil {
				consentNeeded = true
			}
		} else {
			consentNeeded = true
		}

		b := make([]byte, 16)
		rand.Read(b)

		postLoginURL := c.Query("rurl")
		if postLoginURL == "" {
			postLoginURL = defaultPostLoginURL
		}

		state := base64.URLEncoding.EncodeToString(b) + "|" + base64.URLEncoding.EncodeToString([]byte(postLoginURL))
		c.SetCookie(stateCookieName, state, int((5 * time.Minute).Seconds()), "/", "", curl.Scheme == "https", true)

		authCodeOptions := []oauth2.AuthCodeOption{oauth2.AccessTypeOffline}
		if consentNeeded {
			authCodeOptions = append(authCodeOptions, oauth2.ApprovalForce)
		}
		c.Redirect(http.StatusTemporaryRedirect, oauth.AuthCodeURL(state, authCodeOptions...))
	}
}
