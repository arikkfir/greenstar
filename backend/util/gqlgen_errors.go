package util

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/rs/zerolog/log"
	"github.com/secureworks/errors"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

//goland:noinspection GoTypeAssertionOnErrors
func GraphErrorPresenter(ctx context.Context, e error) *gqlerror.Error {
	var path ast.Path
	err := e
	if ge, ok := err.(*gqlerror.Error); ok {
		// Unwrap one level of gqlerror.Error since gql always wraps returned errors from
		// resolvers with it - and we want to test the underlying error if it's a tag provider
		err = errors.Unwrap(ge)
		path = ge.Path
	} else {
		path = graphql.GetPath(ctx)
	}
	if te, ok := err.(errors.TagProvider); ok {
		if te.HasTag(UserFacingTagKey) {
			return e.(*gqlerror.Error)
		}
	}

	log.Ctx(ctx).
		Error().
		Stack().
		Err(e).
		Interface("gqlPath", path).
		Msg("Internal error occurred")

	return &gqlerror.Error{
		Message:    "An internal error has occurred.",
		Path:       path,
		Extensions: map[string]interface{}{"code": "INTERNAL_SERVER_ERROR"},
	}
}

//goland:noinspection GoTypeAssertionOnErrors
func GraphPanicRecoverer(_ context.Context, p interface{}) error {
	if e, ok := p.(error); ok {
		// We assume stack trace will be retained for the wrapped error
		return errors.New("recovered panic: %w", e)
	} else {
		// No stack trace will be available, so we format the panic object itself with as much information as possible
		return errors.New("recovered panic of type '%T': %+v", p, p)
	}
}
