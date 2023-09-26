package middleware

import (
	"github.com/arikkfir/greenstar/backend/internal/server/util"
	"net/http"
	"strings"
)

type Handlers struct {
	OPTIONS http.Handler
	HEAD    http.Handler
	POST    http.Handler
	GET     http.Handler
	PUT     http.Handler
	PATCH   http.Handler
	DELETE  http.Handler
}

func (h *Handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch strings.ToUpper(r.Method) {
	case http.MethodOptions:
		if h.OPTIONS != nil {
			h.OPTIONS.ServeHTTP(w, r)
		} else {
			methods := []string{http.MethodOptions}
			if h.POST != nil {
				methods = append(methods, http.MethodPost)
			}
			if h.GET != nil {
				methods = append(methods, http.MethodGet, http.MethodHead)
			}
			if h.PUT != nil {
				methods = append(methods, http.MethodPut)
			}
			if h.PATCH != nil {
				methods = append(methods, http.MethodPatch)
			}
			if h.DELETE != nil {
				methods = append(methods, http.MethodDelete)
			}
			w.Header().Set("Allow", strings.Join(methods, ", "))
			w.WriteHeader(http.StatusNoContent)
		}
	case http.MethodHead:
		if h.HEAD != nil {
			h.HEAD.ServeHTTP(w, r)
		} else if h.GET != nil {
			h.GET.ServeHTTP(&responseBodyDiscarder{ResponseWriter: w}, r)
		} else {
			util.ServeError(w, r, util.ErrMethodNotAllowed)
		}
	case http.MethodPost:
		if h.POST != nil {
			h.POST.ServeHTTP(w, r)
		} else {
			util.ServeError(w, r, util.ErrMethodNotAllowed)
		}
	case http.MethodGet:
		if h.GET != nil {
			h.GET.ServeHTTP(w, r)
		} else {
			util.ServeError(w, r, util.ErrMethodNotAllowed)
		}
	case http.MethodPut:
		if h.PUT != nil {
			h.PUT.ServeHTTP(w, r)
		} else {
			util.ServeError(w, r, util.ErrMethodNotAllowed)
		}
	case http.MethodPatch:
		if h.PATCH != nil {
			h.PATCH.ServeHTTP(w, r)
		} else {
			util.ServeError(w, r, util.ErrMethodNotAllowed)
		}
	case http.MethodDelete:
		if h.DELETE != nil {
			h.DELETE.ServeHTTP(w, r)
		} else {
			util.ServeError(w, r, util.ErrMethodNotAllowed)
		}
	}
}
