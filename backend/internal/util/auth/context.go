package auth

import (
	"context"
	"fmt"
)

type tokenKeyType struct{}

var tokenKey = &tokenKeyType{}

type Permission string

type Token struct{}

func (t *Token) IsPermittedGlobally(_ ...Permission) bool {
	return true
}

func (t *Token) IsPermittedForTenant(_ string, _ ...Permission) bool {
	return true
}

func GetToken(ctx context.Context) *Token {
	v := ctx.Value(tokenKey)
	if v == nil {
		panic(fmt.Errorf("no token found in context"))
	} else if token, ok := v.(*Token); ok {
		return token
	} else {
		panic(fmt.Sprintf("unexpected token type '%T' encountered: %+v", v, v))
	}
}

func WithToken(ctx context.Context, token *Token) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}
