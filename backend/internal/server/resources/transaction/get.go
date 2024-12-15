// Code generated by greenstar scripts; DO NOT EDIT.

package transaction

import (
	"github.com/arikkfir/greenstar/backend/internal/auth"
	"github.com/arikkfir/greenstar/backend/internal/server/middleware"
	"github.com/arikkfir/greenstar/backend/internal/server/util"
	"github.com/shopspring/decimal"
	"net/http"
	"time"
)

var (
	_ = decimal.Decimal{}
	_ = time.Time{}
)

type GetRequest struct {
	ID       string  `json:"id"`
	Currency *string `url:"currency,omitempty"`
}

func (lr *GetRequest) UnmarshalFromRequest(r *http.Request) error {
	lr.ID = r.PathValue("id")
	if lr.ID == "" {
		return util.ErrBadRequest
	}
	return nil
}

type GetResponse Transaction

func (s *Server) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	l := util.Logger(ctx)

	tenantID := middleware.GetTenantID(ctx)
	if tenantID != "" {
		l = l.With("tenantID", tenantID)
	}
	authToken := auth.GetToken(ctx)
	if !authToken.IsPermittedGlobally("transactions:read") {
		if tenantID != "" {
			if !authToken.IsPermittedForTenant(tenantID, "transactions:read") {
				util.ServeError(w, r, util.ErrForbidden)
				l.WarnContext(ctx, "Access denied", "permission", "transactions:read")
				return
			}
		} else {
			util.ServeError(w, r, util.ErrForbidden)
			l.WarnContext(ctx, "Access denied", "permission", "transactions:read")
			return
		}
	}

	req := GetRequest{}
	if err := req.UnmarshalFromRequest(r); err != nil {
		util.ServeError(w, r, err)
		return
	}

	res, err := s.h.Get(ctx, req)
	if err != nil {
		if code := util.ServeError(w, r, err); code >= http.StatusInternalServerError {
			l.ErrorContext(ctx, "Failed getting transaction", "err", err)
		}
		return
	} else if res == nil {
		util.ServeError(w, r, util.ErrNotFound)
		return
	}

	if err := util.Marshal(w, r, http.StatusOK, res); err != nil {
		l.ErrorContext(ctx, "Failed marshaling transaction", "err", err)
		util.ServeError(w, r, err)
	}
}
