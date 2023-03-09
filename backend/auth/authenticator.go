package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	Config              *Config
	SecureCookies       bool
	DefaultPostLoginURL string
	OAuth               *oauth2.Config
}

func (a *Authenticator) createTokenVerificationCallback(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	} else {
		return []byte(a.Config.Google.ClientSecret), nil
	}
}
