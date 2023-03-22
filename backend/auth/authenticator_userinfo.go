package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
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

func (a *Authenticator) HandleUserInfo(c *gin.Context) {
	claims, err := a.parseClaimsFromCookie(c, "greenstar.auth")
	if err != nil {
		c.SetCookie(a.Config.StateCookieName, "", -1, "/", "", a.SecureCookies, true)
		c.SetCookie(a.Config.ClaimsCookieName, "", -1, "/", "", a.SecureCookies, true)
		c.AbortWithError(http.StatusUnauthorized, gin.Error{
			Err:  fmt.Errorf("failed getting claims from cookie: %w", err),
			Type: gin.ErrorTypePrivate,
		})
		return
	} else if claims == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	session, err := a.loadSession(c, claims.ID)
	if err != nil {
		c.SetCookie(a.Config.StateCookieName, "", -1, "/", "", a.SecureCookies, true)
		c.SetCookie(a.Config.ClaimsCookieName, "", -1, "/", "", a.SecureCookies, true)
		c.AbortWithError(http.StatusUnauthorized, gin.Error{
			Err:  fmt.Errorf("failed loading session for claims: %w", err),
			Type: gin.ErrorTypePrivate,
		})
		return
	} else if session == nil {
		c.SetCookie(a.Config.StateCookieName, "", -1, "/", "", a.SecureCookies, true)
		c.SetCookie(a.Config.ClaimsCookieName, "", -1, "/", "", a.SecureCookies, true)
		c.AbortWithError(http.StatusUnauthorized, gin.Error{
			Err:  fmt.Errorf("session not found for claims: %w", err),
			Type: gin.ErrorTypePrivate,
		})
		return
	} else if err := a.verifyClaimsAndSession(claims, session); err != nil {
		c.SetCookie(a.Config.StateCookieName, "", -1, "/", "", a.SecureCookies, true)
		c.SetCookie(a.Config.ClaimsCookieName, "", -1, "/", "", a.SecureCookies, true)
		c.AbortWithError(http.StatusUnauthorized, gin.Error{
			Err:  fmt.Errorf("failed verifying claims and session: %w", err),
			Type: gin.ErrorTypePrivate,
		})
		return
	}

	if session.MockUserInfo != nil {
		if session.HasPermission(PermissionAuthUserInfoMock) {
			log.Ctx(c).Warn().Interface("mock", session.MockUserInfo).Msg("Mock user info provided!")
			c.Negotiate(http.StatusOK, gin.Negotiate{
				Offered: []string{gin.MIMEJSON},
				Data:    *session.MockUserInfo,
			})
		} else {
			c.AbortWithError(http.StatusForbidden, gin.Error{
				Err:  fmt.Errorf("insufficient permissions to use mock userinfo"),
				Type: gin.ErrorTypePrivate,
			})
		}
		return
	}

	userInfo, err := a.loadUserInfo(c, session.Token)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, gin.Error{
			Err:  fmt.Errorf("failed loading userinfo: %w", err),
			Type: gin.ErrorTypePrivate,
		})
	} else {
		c.Negotiate(http.StatusOK, gin.Negotiate{
			Offered: []string{gin.MIMEJSON},
			Data:    userInfo,
		})
	}
}
