// Code generated by greenstar scripts; DO NOT EDIT.

package account

import (
	"encoding/json"
	"github.com/arikkfir/greenstar/backend/internal/auth"
	"github.com/arikkfir/greenstar/backend/internal/server/middleware"
	"github.com/arikkfir/greenstar/backend/internal/server/util"
	"github.com/shopspring/decimal"
	"net/http"
	"slices"
	"time"
)

var (
	_ = decimal.Decimal{}
	_ = time.Time{}
)

type CreateRequest struct {
	DisplayName string  `json:"displayName,omitempty"`
	Icon        *string `json:"icon,omitempty"`
	ParentID    *string `json:"parentID,omitempty"`
	properties  []string
}

func (lr *CreateRequest) HasIcon() bool     { return slices.Contains(lr.properties, "icon") }
func (lr *CreateRequest) HasParentID() bool { return slices.Contains(lr.properties, "parentID") }
func (lr *CreateRequest) UnmarshalJSON(data []byte) error {
	lr.properties = nil
	var tempMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return err
	}
	for key := range tempMap {
		lr.properties = append(lr.properties, key)
	}
	type typeAlias CreateRequest
	alias := (*typeAlias)(lr)
	if err := json.Unmarshal(data, alias); err != nil {
		return err
	}
	return nil
}

type CreateResponse Account

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	l := util.Logger(ctx)

	tenantID := middleware.GetTenantID(ctx)
	if tenantID != "" {
		l = l.With("tenantID", tenantID)
	}
	authToken := auth.GetToken(ctx)
	if !authToken.IsPermittedGlobally("accounts:create") {
		if tenantID != "" {
			if !authToken.IsPermittedForTenant(tenantID, "accounts:create") {
				util.ServeError(w, r, util.ErrForbidden)
				l.WarnContext(ctx, "Access denied", "permission", "accounts:create")
				return
			}
		} else {
			util.ServeError(w, r, util.ErrForbidden)
			l.WarnContext(ctx, "Access denied", "permission", "accounts:create")
			return
		}
	}

	req := CreateRequest{}
	if err := util.UnmarshalBody(r, &req); err != nil {
		util.ServeError(w, r, err)
		return
	}

	res, err := s.h.Create(ctx, req)
	if err != nil {
		if code := util.ServeError(w, r, err); code >= http.StatusInternalServerError {
			l.ErrorContext(ctx, "Failed creating account", "err", err)
		}
		return
	}

	if err := util.Marshal(w, r, http.StatusCreated, res); err != nil {
		l.ErrorContext(ctx, "Failed marshaling account", "err", err)
		util.ServeError(w, r, err)
	}
}
