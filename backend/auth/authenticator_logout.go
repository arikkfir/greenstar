package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Authenticator) HandleLogout(c *gin.Context) {
	c.SetCookie(a.Config.StateCookieName, "", -1, "/", "", a.SecureCookies, true)
	c.SetCookie(a.Config.ClaimsCookieName, "", -1, "/", "", a.SecureCookies, true)
	c.Redirect(http.StatusFound, a.DefaultPostLoginURL)
}
