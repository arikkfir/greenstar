package util

import (
	"errors"
	"net/http"
)

var (
	ErrBadRequest                   = errors.New(http.StatusText(http.StatusBadRequest))                   // ErrBadRequest indicates malformed syntax, etc.
	ErrUnauthorized                 = errors.New(http.StatusText(http.StatusUnauthorized))                 // ErrUnauthorized indicates missing credentials.
	ErrForbidden                    = errors.New(http.StatusText(http.StatusForbidden))                    // ErrForbidden indicates lack of necessary permissions.
	ErrNotFound                     = errors.New(http.StatusText(http.StatusNotFound))                     // ErrNotFound indicates resource could not be found.
	ErrMethodNotAllowed             = errors.New(http.StatusText(http.StatusMethodNotAllowed))             // ErrMethodNotAllowed indicates method not supported by resource.
	ErrNotAcceptable                = errors.New(http.StatusText(http.StatusNotAcceptable))                // ErrNotAcceptable indicates no client-acceptable media type for resource.
	ErrConflict                     = errors.New(http.StatusText(http.StatusConflict))                     // ErrConflict indicates a conflict the user must resolve and resubmit.
	ErrGone                         = errors.New(http.StatusText(http.StatusGone))                         // ErrGone indicates the resource is no longer available.
	ErrPreconditionFailed           = errors.New(http.StatusText(http.StatusPreconditionFailed))           // ErrPreconditionFailed indicates one or more request conditions are false.
	ErrUnsupportedMediaType         = errors.New(http.StatusText(http.StatusUnsupportedMediaType))         // ErrUnsupportedMediaType indicates payload media type is not supported.
	ErrRequestedRangeNotSatisfiable = errors.New(http.StatusText(http.StatusRequestedRangeNotSatisfiable)) // ErrRequestedRangeNotSatisfiable indicates the given Range cannot be satisfied.
	ErrUnprocessableEntity          = errors.New(http.StatusText(http.StatusUnprocessableEntity))          // ErrUnprocessableEntity indicates request payload is well-formed syntactically but not semantically.
	ErrLocked                       = errors.New(http.StatusText(http.StatusLocked))                       // ErrLocked indicates that the resource is currently locked.
	ErrPreconditionRequired         = errors.New(http.StatusText(http.StatusPreconditionRequired))         // ErrPreconditionRequired indicates request must be conditional but isn't.
	ErrInternalError                = errors.New(http.StatusText(http.StatusInternalServerError))          // ErrInternalError indicates an unexpected internal server error.
	ErrNotImplemented               = errors.New(http.StatusText(http.StatusNotImplemented))               // ErrNotImplemented indicates that a requested operation is not implemented.
	ErrServiceUnavailable           = errors.New(http.StatusText(http.StatusServiceUnavailable))           // ErrServiceUnavailable indicates service is currently unavailable.

	httpCodesMap = map[error]int{
		ErrBadRequest:                   http.StatusBadRequest,
		ErrUnauthorized:                 http.StatusUnauthorized,
		ErrForbidden:                    http.StatusForbidden,
		ErrNotFound:                     http.StatusNotFound,
		ErrMethodNotAllowed:             http.StatusMethodNotAllowed,
		ErrNotAcceptable:                http.StatusNotAcceptable,
		ErrConflict:                     http.StatusConflict,
		ErrGone:                         http.StatusGone,
		ErrPreconditionFailed:           http.StatusPreconditionFailed,
		ErrUnsupportedMediaType:         http.StatusUnsupportedMediaType,
		ErrRequestedRangeNotSatisfiable: http.StatusRequestedRangeNotSatisfiable,
		ErrUnprocessableEntity:          http.StatusUnprocessableEntity,
		ErrLocked:                       http.StatusLocked,
		ErrPreconditionRequired:         http.StatusPreconditionRequired,
		ErrInternalError:                http.StatusInternalServerError,
		ErrNotImplemented:               http.StatusNotImplemented,
		ErrServiceUnavailable:           http.StatusServiceUnavailable,
	}
)

func ServeError(w http.ResponseWriter, r *http.Request, err error) (httpCode int) {
	for httpErr, code := range httpCodesMap {
		if errors.Is(err, httpErr) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			http.Error(w, err.Error(), code)
			return code
		}
	}
	return ServeError(w, r, ErrInternalError)
}
