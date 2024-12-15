package middleware

import (
	"context"
	"fmt"
	"net/http"
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

func TokenMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxWithToken := context.WithValue(r.Context(), tokenKey, &Token{})
		r2 := r.WithContext(ctxWithToken)
		next.ServeHTTP(w, r2)
	}
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
