// Code generated by greenstar scripts; DO NOT EDIT.

package tenant

import (
	"encoding/json"
	"net/http"
	"slices"
	"time"

	"github.com/arikkfir/greenstar/backend/internal/auth"
	"github.com/arikkfir/greenstar/backend/internal/server/util"
	"github.com/shopspring/decimal"
)

var (
	_ = decimal.Decimal{}
	_ = time.Time{}
)

type CreateRequest struct {
	ID          string `json:"id,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	properties  []string
}

func (lr *CreateRequest) HasID() bool { return slices.Contains(lr.properties, "id") }
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

type CreateResponse Tenant

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	l := util.Logger(ctx)
	if !auth.GetToken(ctx).IsPermittedGlobally("tenants:create") {
		util.ServeError(w, r, util.ErrForbidden)
		l.WarnContext(ctx, "Access denied", "permission", "tenants:create")
		return
	}

	req := CreateRequest{}
	if err := util.UnmarshalBody(r, &req); err != nil {
		util.ServeError(w, r, err)
		return
	}

	res, err := s.h.Create(ctx, req)
	if err != nil {
		if code := util.ServeError(w, r, err); code >= http.StatusInternalServerError {
			l.ErrorContext(ctx, "Failed creating tenant", "err", err)
		}
		return
	}

	if err := util.Marshal(w, r, http.StatusCreated, res); err != nil {
		l.ErrorContext(ctx, "Failed marshaling tenant", "err", err)
		util.ServeError(w, r, err)
	}
}
