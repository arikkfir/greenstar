package auth

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const claimsKey = "claims"

type Token struct {
	jwt.RegisteredClaims
	Tenant       string `json:"tenant,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

func GetClaims(ctx context.Context) *Token {
	v := ctx.Value(claimsKey)
	if v == nil {
		return nil
	} else {
		return v.(*Token)
	}
}

func setClaims(c *gin.Context, token *Token) {
	c.Set(claimsKey, token)
}
