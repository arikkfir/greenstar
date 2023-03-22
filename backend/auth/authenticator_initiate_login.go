package auth

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

func (a *Authenticator) HandleInitiateLogin(c *gin.Context) {
	postLoginURL := c.Query("rurl")
	if postLoginURL == "" {
		postLoginURL = a.DefaultPostLoginURL
	}

	state := uuid.New().String() + "|" + base64.URLEncoding.EncodeToString([]byte(postLoginURL))
	c.SetCookie(a.Config.StateCookieName, state, int((5 * time.Minute).Seconds()), "/", "", a.SecureCookies, true)
	c.SetCookie(a.Config.ClaimsCookieName, "", 0, "/", "", a.SecureCookies, true)
	c.Redirect(http.StatusTemporaryRedirect, a.OAuth.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce))
}
