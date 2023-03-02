package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
)

const claimsKey = "claims"

type Token struct {
	jwt.RegisteredClaims
	Tenant       string `json:"tenant,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

func GetToken(ctx context.Context) *Token {
	return ctx.Value(claimsKey).(*Token)
}
