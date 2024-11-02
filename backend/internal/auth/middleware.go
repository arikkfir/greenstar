package auth

import (
	"context"
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/server/util"
	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/client"
	"net/http"
)

type clientKeyType struct{}
type tokenKeyType struct{}

var clientKey = &clientKeyType{}
var tokenKey = &tokenKeyType{}

type Permission string

func NewPermission(scope, action string) Permission {
	return Permission(fmt.Sprintf("%s:%s", scope, action))
}

type Token struct {
	Token descope.Token
	ID    string `json:"id,omitempty"`
}

func (t *Token) IsPermittedGlobally(permissions ...Permission) bool {
	for _, permission := range permissions {
		if !t.Token.IsPermitted(string(permission)) {
			return false
		}
	}
	return true
}

func (t *Token) IsPermittedForTenant(tenant string, permissions ...Permission) bool {
	for _, permission := range permissions {
		hasGlobalPermission := t.Token.IsPermitted(string(permission))
		hasTenantPermission := t.Token.IsPermittedPerTenant(tenant, string(permission))
		if !hasGlobalPermission && !hasTenantPermission {
			return false
		}
	}
	return true
}

func (t *Token) GetTenants() []string {
	return t.Token.GetTenants()
}

func WithSDKMiddleware(descopeClient *client.DescopeClient, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), clientKey, descopeClient)))
	}
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		descopeClient := GetClient(r.Context())
		if ok, token, err := descopeClient.Auth.ValidateAndRefreshSessionWithRequest(r, w); err != nil {
			util.Logger(r.Context()).ErrorContext(r.Context(), "Authentication token validation failed", "err", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		} else if !ok {
			util.Logger(r.Context()).WarnContext(r.Context(), "Authentication token validation did not return an error, but ok=false", "err", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		} else {
			contextWithToken := NewContextWithToken(r.Context(), &Token{Token: *token, ID: token.ID})
			next.ServeHTTP(w, r.WithContext(contextWithToken))
		}
	})
}

func GetClient(ctx context.Context) *client.DescopeClient {
	v := ctx.Value(clientKey)
	if v == nil {
		panic(fmt.Errorf("no descope client found in context"))
	} else if descopeClient, ok := v.(*client.DescopeClient); ok {
		return descopeClient
	} else {
		panic(fmt.Sprintf("unexpected descope client type '%T' encountered: %+v", v, v))
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

func NewContextWithToken(ctx context.Context, token *Token) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}
