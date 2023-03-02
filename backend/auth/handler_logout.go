package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateAuthGoogleLogoutHandler(stateCookieName, claimsCookieName string, useSecureCookies bool, postLogoutURL string) func(*gin.Context) {
	return func(c *gin.Context) {
		c.SetCookie(stateCookieName, "", -1, "/", "", useSecureCookies, true)
		c.SetCookie(claimsCookieName, "", -1, "/", "", useSecureCookies, true)
		c.Redirect(http.StatusFound, postLogoutURL)
	}
}
