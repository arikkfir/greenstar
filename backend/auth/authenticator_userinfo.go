package auth

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
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
	session := GetSession(c)

	httpClient := a.OAuth.Client(c, session.Token)
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

	if userInfoResponse.StatusCode >= 200 && userInfoResponse.StatusCode <= 299 {
		userInfo := GoogleAPIUserInfoResponse{}
		decoder := json.NewDecoder(userInfoResponse.Body)
		if err := decoder.Decode(&userInfo); err != nil {
			c.AbortWithError(http.StatusInternalServerError, gin.Error{
				Err:  fmt.Errorf("failed decoding google userinfo response: %w", err),
				Type: gin.ErrorTypePrivate,
			})
		}

		c.Negotiate(http.StatusOK, gin.Negotiate{
			Offered: []string{gin.MIMEJSON},
			Data:    userInfo,
		})
	} else if userInfoResponse.StatusCode == http.StatusUnauthorized {
		body, _ := io.ReadAll(userInfoResponse.Body)
		c.AbortWithError(http.StatusUnauthorized, gin.Error{
			Err:  fmt.Errorf("unauthorized to request userinfo"),
			Type: gin.ErrorTypePrivate,
			Meta: map[string]interface{}{"body": body},
		})
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
	}
}