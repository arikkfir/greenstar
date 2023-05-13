package auth

import (
	"context"
	"fmt"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"strings"
)

const claimsContextKey = "$$$$____claims____$$$$"

type Claims struct {
	validator.RegisteredClaims
	CustomClaims
}

func (c *Claims) HasScope(expectedScope string) bool {
	result := strings.Split(c.CustomClaims.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}
	return false
}

func GetClaims(ctx context.Context) *Claims {
	v := ctx.Value(claimsContextKey)
	if v == nil {
		return nil
	} else if claims, ok := v.(*Claims); ok {
		return claims
	} else {
		panic(fmt.Sprintf("unexpected claims object encountered: %+v", v))
	}
}
