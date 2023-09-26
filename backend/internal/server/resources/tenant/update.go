// Code generated by greenstar scripts; DO NOT EDIT.

package tenant

import (
	"encoding/json"
	"net/http"

	"github.com/arikkfir/greenstar/backend/internal/auth"
	"github.com/arikkfir/greenstar/backend/internal/server/util"
)

type UpdateRequest struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName,omitempty"`
	properties  []string
}

func (lr *UpdateRequest) UnmarshalJSON(data []byte) error {
	lr.properties = nil
	var tempMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return err
	}
	for key := range tempMap {
		lr.properties = append(lr.properties, key)
	}
	type typeAlias UpdateRequest
	alias := (*typeAlias)(lr)
	if err := json.Unmarshal(data, alias); err != nil {
		return err
	}
	return nil
}

type UpdateResponse Tenant

func (s *Server) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	l := util.Logger(ctx)
	if !auth.GetToken(ctx).IsPermittedForTenant(r.PathValue("id"), "tenants:update") {
		util.ServeError(w, r, util.ErrForbidden)
		l.With("tenantID", r.PathValue("TenantPathVariableName")).WarnContext(ctx, "Access denied", "permission", "tenants:update")
		return
	}

	req := UpdateRequest{}
	if err := util.UnmarshalBody(r, &req); err != nil {
		util.ServeError(w, r, err)
		return
	}
	req.ID = r.PathValue("id")
	if req.ID == "" {
		util.ServeError(w, r, util.ErrBadRequest)
		return
	}

	res, err := s.h.Update(ctx, req)
	if err != nil {
		if code := util.ServeError(w, r, err); code >= http.StatusInternalServerError {
			l.ErrorContext(ctx, "Failed updating tenant", "err", err)
		}
		return
	} else if res == nil {
		util.ServeError(w, r, util.ErrNotFound)
		return
	}

	if err := util.Marshal(w, r, http.StatusOK, res); err != nil {
		l.ErrorContext(ctx, "Failed marshaling tenant", "err", err)
		util.ServeError(w, r, err)
	}
}
