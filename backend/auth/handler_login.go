package main

import (
	"encoding/base64"
	"fmt"
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
		b := make([]byte, 16)
		rand.Read(b)

		rurl, found := c.GetPostForm("rurl")
		if !found || rurl == "" {
			rurl = c.Query("rurl")
		}
		if rurl == "" {
			rurl = defaultPostLoginURL
		}

		state := base64.URLEncoding.EncodeToString(b) + "|" + base64.URLEncoding.EncodeToString([]byte(rurl))

		c.SetCookie(stateCookieName, state, int((5 * time.Minute).Seconds()), "/", "", curl.Scheme == "https", true)

		// TODO: only send "oauth2.ApprovalForce" option if missing or partial cookie (esp. if missing refresh token)
		c.Redirect(http.StatusTemporaryRedirect, oauth.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce))
	}
}
