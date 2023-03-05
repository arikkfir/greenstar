package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	auth "github.com/arikkfir/greenstar/backend/auth/pkg"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func CreateAuthGoogleUserInfoHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		token := auth.GetClaims(c)

		response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?prettyPrint=true&oauth_token=" + token.AccessToken)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			c.Error(gin.Error{
				Err:  fmt.Errorf("userinfo request failed: %w", err),
				Type: gin.ErrorTypePrivate,
				Meta: map[string]interface{}{
					"url":    response.Request.URL.String(),
					"status": response.Status,
					"token":  token.AccessToken,
				},
			})
			return
		}
		defer response.Body.Close()

		if response.StatusCode >= 200 && response.StatusCode <= 299 {
			userInfo := GoogleAPIUserInfoResponse{}
			responseBody := bytes.Buffer{}
			decoder := json.NewDecoder(io.TeeReader(response.Body, &responseBody))
			if err := decoder.Decode(&userInfo); err != nil {
				c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed decoding userinfo response: %w", err))
				return
			}
			c.Negotiate(http.StatusOK, gin.Negotiate{Offered: []string{gin.MIMEJSON, gin.MIMEYAML}, Data: userInfo})
		} else if response.StatusCode == http.StatusUnauthorized {
			body, _ := io.ReadAll(response.Body)
			c.AbortWithStatus(http.StatusUnauthorized)
			c.Error(gin.Error{
				Err:  fmt.Errorf("userinfo request failed with status code %d", response.StatusCode),
				Type: gin.ErrorTypePrivate,
				Meta: map[string]interface{}{
					"url":    response.Request.URL.String(),
					"status": response.Status,
					"token":  token.AccessToken,
					"body":   string(body),
				},
			})
			return
		} else {
			body, _ := io.ReadAll(response.Body)
			c.AbortWithStatus(http.StatusInternalServerError)
			c.Error(gin.Error{
				Err:  fmt.Errorf("userinfo request failed with status code %d", response.StatusCode),
				Type: gin.ErrorTypePrivate,
				Meta: map[string]interface{}{
					"url":    response.Request.URL.String(),
					"status": response.Status,
					"token":  token.AccessToken,
					"body":   string(body),
				},
			})
			return
		}
	}
}
