package gqlutil

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/rs/zerolog/log"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func ErrorPresenter(ctx context.Context, e error) *gqlerror.Error {
	if gqlErr, ok := e.(*gqlerror.Error); ok {
		return gqlErr
	}

	path := graphql.GetPath(ctx)
	log.Ctx(ctx).
		Error().
		Err(e).
		Interface("gqlPath", path).
		Msg("GraphQL handling error")

	return &gqlerror.Error{
		Message:    e.Error(),
		Path:       path,
		Extensions: map[string]interface{}{"code": "INTERNAL_SERVER_ERROR"},
	}
}

func PanicRecoverer(ctx context.Context, err interface{}) error {
	log.Ctx(ctx).Debug().Interface("panic", err).Msg("Recovered from panic")
	if e, ok := err.(error); ok {
		return e
	}
	return fmt.Errorf("internal server error: %v", err)
}
