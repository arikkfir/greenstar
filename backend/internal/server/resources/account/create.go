// Code generated by greenstar scripts; DO NOT EDIT.

package account

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/arikkfir/greenstar/backend/internal/auth"
	"github.com/arikkfir/greenstar/backend/internal/server/util"
	"github.com/shopspring/decimal"
)

type CreateRequest struct {
	TenantID            string           `json:"-"`
	Balance             *decimal.Decimal `json:"balance,omitempty"`
	DisplayName         string           `json:"displayName,omitempty"`
	Icon                *string          `json:"icon,omitempty"`
	ParentID            *string          `json:"parentID,omitempty"`
	TotalIncomingAmount *decimal.Decimal `json:"totalIncomingAmount,omitempty"`
	TotalOutgoingAmount *decimal.Decimal `json:"totalOutgoingAmount,omitempty"`
	properties          []string
}

func (lr *CreateRequest) HasBalance() bool  { return slices.Contains(lr.properties, "balance") }
func (lr *CreateRequest) HasIcon() bool     { return slices.Contains(lr.properties, "icon") }
func (lr *CreateRequest) HasParentID() bool { return slices.Contains(lr.properties, "parentID") }
func (lr *CreateRequest) HasTotalIncomingAmount() bool {
	return slices.Contains(lr.properties, "totalIncomingAmount")
}
func (lr *CreateRequest) HasTotalOutgoingAmount() bool {
	return slices.Contains(lr.properties, "totalOutgoingAmount")
}
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
	if !auth.GetToken(ctx).IsPermittedForTenant(r.PathValue("tenantID"), "accounts:read") {
		util.ServeError(w, r, util.ErrForbidden)
		l.With("tenantID", r.PathValue("TenantPathVariableName")).WarnContext(ctx, "Access denied", "permission", "accounts:read")
		return
	}

	req := CreateRequest{}
	if err := util.UnmarshalBody(r, &req); err != nil {
		util.ServeError(w, r, err)
		return
	}
	req.TenantID = r.PathValue("tenantID")
	if req.TenantID == "" {
		util.ServeError(w, r, fmt.Errorf("%w: invalid tenant ID", util.ErrBadRequest))
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