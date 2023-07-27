package web

import (
	"context"
	"fmt"
	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/client"
	"github.com/rs/zerolog/log"
	"net/http"
)

const (
	contextTokenKey         = "$$$_descope_token_$$$"
	contextDescopeClientKey = "$$$_descope_client_$$$"
)

func DescopeAuthenticationMiddleware(descopeClient *client.DescopeClient, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), contextDescopeClientKey, descopeClient))
		if ok, token, err := descopeClient.Auth.ValidateAndRefreshSessionWithRequest(r, w); ok {
			r = r.WithContext(context.WithValue(r.Context(), contextTokenKey, token))
			next(w, r)
		} else {
			log.Ctx(r.Context()).Error().Err(err).Msg("Authentication token validation failed")
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}

func GetDescopeClient(ctx context.Context) *client.DescopeClient {
	v := ctx.Value(contextDescopeClientKey)
	if v == nil {
		panic(fmt.Errorf("no descope client found in context"))
	} else if descopeClient, ok := v.(*client.DescopeClient); ok {
		return descopeClient
	} else {
		panic(fmt.Sprintf("unexpected descope client type '%T' encountered: %+v", v, v))
	}
}

func GetToken(ctx context.Context) *descope.Token {
	v := ctx.Value(contextTokenKey)
	if v == nil {
		panic(fmt.Errorf("no token found in context"))
	} else if token, ok := v.(*descope.Token); ok {
		return token
	} else {
		panic(fmt.Sprintf("unexpected token type '%T' encountered: %+v", v, v))
	}
}
