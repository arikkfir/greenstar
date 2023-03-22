package gqlutil

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/rs/zerolog/log"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"reflect"
)

type UserError struct {
	err *gqlerror.Error
}

func (e *UserError) Error() string {
	return e.err.Error()
}

func (e *UserError) Is(err error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(e)
}

func (e *UserError) Unwrap() error {
	return e.err
}

func ErrorPresenter(ctx context.Context, e error) *gqlerror.Error {
	if errors.Is(e, &UserError{}) {
		return e.(*gqlerror.Error)
	}

	path := graphql.GetPath(ctx)
	log.Ctx(ctx).
		Error().
		Stack().
		Err(e).
		Interface("actualError", fmt.Sprintf("%T: %+v", e, e)).
		Interface("gqlPath", path).
		Msg("Internal error occurred")

	return &gqlerror.Error{
		Message:    "An internal error has occurred.",
		Path:       path,
		Extensions: map[string]interface{}{"code": "INTERNAL_SERVER_ERROR"},
	}
}

func PanicRecoverer(ctx context.Context, err interface{}) error {
	log.Ctx(ctx).Error().Stack().Interface("panic", fmt.Sprintf("%T: %+v", err, err)).Msg("Recovered from panic")
	if e, ok := err.(error); ok {
		return e
	}
	return fmt.Errorf("internal server error: %v", err)
}
