package e2e

import (
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/account"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/tenant"
	"github.com/arikkfir/greenstar/backend/internal/util/strings"
	"github.com/gavv/httpexpect/v2"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"testing"
	"time"
)

type session struct {
	t    *testing.T
	http *httpexpect.Expect
}

func newSession(t *testing.T, baseURL, jwt string) *session {
	return &session{
		t:    t,
		http: createAuthenticatedHTTPExpect(t, baseURL, jwt),
	}
}

func (s *session) createTenant(req tenant.CreateRequest) tenant.CreateResponse {
	res := tenant.CreateResponse{}
	s.http.POST("/tenants").
		WithHeader("Accept", "application/json").
		WithJSON(req).
		Expect().
		Status(http.StatusCreated).
		HasContentType("application/json").
		JSON().
		Decode(&res)
	assertNotZero(s.t, res.ID)
	assertEquals(s.t, time.Now(), res.CreatedAt, cmpopts.EquateApproxTime(5*time.Second))
	assertEquals(s.t, time.Now(), res.UpdatedAt, cmpopts.EquateApproxTime(5*time.Second))
	assertEquals(s.t, req.DisplayName, res.DisplayName)
	assertEquals(s.t, req.Slug, res.Slug)
	return res
}

func (s *session) createTestTenant() tenant.CreateResponse {
	displayName := fmt.Sprintf("Test %s-%d", s.t.Name(), processID)
	slug := strings.Slugify(displayName)
	res := s.createTenant(tenant.CreateRequest{
		DisplayName: displayName,
		Slug:        slug,
	})
	s.t.Cleanup(func() { s.deleteTenant(tenant.DeleteRequest{ID: res.ID}) })
	return res
}

func (s *session) getTenant(req tenant.GetRequest) tenant.GetResponse {
	res := tenant.GetResponse{}
	s.http.GET("/tenants/{id}").
		WithHeader("Accept", "application/json").
		WithPath("id", req.ID).
		Expect().
		Status(http.StatusOK).
		HasContentType("application/json").
		JSON().
		Decode(&res)
	assertNotZero(s.t, res.ID)
	assertNotZero(s.t, res.CreatedAt)
	assertNotZero(s.t, res.UpdatedAt)
	assertNotZero(s.t, res.DisplayName)
	assertNotZero(s.t, res.Slug)
	return res
}

func (s *session) listTenants(req tenant.ListRequest) tenant.ListResponse {
	res := tenant.ListResponse{}
	s.http.GET("/tenants").
		WithHeader("Accept", "application/json").
		WithQueryObject(req).
		Expect().
		Status(http.StatusOK).
		HasContentType("application/json").
		JSON().
		Decode(&res)
	return res
}

func (s *session) patchTenant(req tenant.PatchRequest) tenant.PatchResponse {
	res := tenant.PatchResponse{}
	s.http.PATCH("/tenants/{id}").
		WithHeader("Accept", "application/json").
		WithPath("id", req.ID).
		WithJSON(req).
		Expect().
		Status(http.StatusOK).
		HasContentType("application/json").
		JSON().
		Decode(&res)
	assertEquals(s.t, req.ID, res.ID)
	assertNotZero(s.t, res.CreatedAt)
	assertNotZero(s.t, res.UpdatedAt)
	if res.UpdatedAt.Before(res.CreatedAt) {
		s.t.Errorf("Update time should be greater than creation time after patching:\nCreatedAt: %s\nUpdatedAt: %s", res.CreatedAt, res.UpdatedAt)
	}
	return res
}

func (s *session) updateTenant(req tenant.UpdateRequest) tenant.UpdateResponse {
	res := tenant.UpdateResponse{}
	s.http.PUT("/tenants/{id}").
		WithHeader("Accept", "application/json").
		WithPath("id", req.ID).
		WithJSON(req).
		Expect().
		Status(http.StatusOK).
		HasContentType("application/json").
		JSON().
		Decode(&res)
	assertEquals(s.t, req.ID, res.ID)
	assertNotZero(s.t, res.CreatedAt)
	assertNotZero(s.t, res.UpdatedAt)
	assertEquals(s.t, req.DisplayName, res.DisplayName)
	assertEquals(s.t, req.Slug, res.Slug)
	if res.UpdatedAt.Before(res.CreatedAt) {
		s.t.Errorf("Update time should be greater than creation time after patching:\nCreatedAt: %s\nUpdatedAt: %s", res.CreatedAt, res.UpdatedAt)
	}
	return res
}

func (s *session) deleteTenant(req tenant.DeleteRequest) tenant.DeleteResponse {
	res := tenant.DeleteResponse{}
	s.http.DELETE("/tenants/{id}").
		WithHeader("Accept", "application/json").
		WithPath("id", req.ID).
		Expect().
		Status(http.StatusNoContent)
	return res
}

func (s *session) createAccount(req account.CreateRequest) account.CreateResponse {
	res := account.CreateResponse{}
	s.http.POST("/tenants/{tenantID}").
		WithHeader("Accept", "application/json").
		WithPath("tenantID", req.TenantID).
		WithJSON(req).
		Expect().
		Status(http.StatusCreated).
		HasContentType("application/json").
		JSON().
		Decode(&res)
	assertNotZero(s.t, res.ID)
	assertEquals(s.t, time.Now(), res.CreatedAt, cmpopts.EquateApproxTime(5*time.Second))
	assertEquals(s.t, time.Now(), res.UpdatedAt, cmpopts.EquateApproxTime(5*time.Second))
	assertEquals(s.t, req.DisplayName, res.DisplayName)
	assertEquals(s.t, req.Icon, res.Icon)
	assertEquals(s.t, req.ParentID, res.ParentID)
	return res
}

func (s *session) getAccount(req account.GetRequest) account.GetResponse {
	res := account.GetResponse{}
	s.http.GET("/tenants/{tenantID}/accounts/{id}").
		WithHeader("Accept", "application/json").
		WithPath("tenantID", req.TenantID).
		WithPath("id", req.ID).
		Expect().
		Status(http.StatusOK).
		HasContentType("application/json").
		JSON().
		Decode(&res)
	assertNotZero(s.t, res.ID)
	assertNotZero(s.t, res.CreatedAt)
	assertNotZero(s.t, res.UpdatedAt)
	assertNotZero(s.t, res.DisplayName)
	return res
}

func (s *session) listAccounts(req account.ListRequest) account.ListResponse {
	res := account.ListResponse{}
	s.http.GET("/tenants/{tenantID}/accounts").
		WithHeader("Accept", "application/json").
		WithPath("tenantID", req.TenantID).
		WithQueryObject(req).
		Expect().
		Status(http.StatusOK).
		HasContentType("application/json").
		JSON().
		Decode(&res)
	return res
}

func (s *session) patchAccount(req account.PatchRequest) account.PatchResponse {
	res := account.PatchResponse{}
	s.http.PATCH("/tenants/{tenantID}/accounts/{id}").
		WithHeader("Accept", "application/json").
		WithPath("tenantID", req.TenantID).
		WithPath("id", req.ID).
		WithJSON(req).
		Expect().
		Status(http.StatusOK).
		HasContentType("application/json").
		JSON().
		Decode(&res)
	assertEquals(s.t, req.ID, res.ID)
	assertNotZero(s.t, res.CreatedAt)
	assertNotZero(s.t, res.UpdatedAt)
	if res.UpdatedAt.Before(res.CreatedAt) {
		s.t.Errorf("Update time should be greater than creation time after patching:\nCreatedAt: %s\nUpdatedAt: %s", res.CreatedAt, res.UpdatedAt)
	}
	return res
}

func (s *session) updateAccount(req account.UpdateRequest) account.UpdateResponse {
	res := account.UpdateResponse{}
	s.http.PUT("/tenants/{tenantID}/accounts/{id}").
		WithHeader("Accept", "application/json").
		WithPath("tenantID", req.TenantID).
		WithPath("id", req.ID).
		WithJSON(req).
		Expect().
		Status(http.StatusOK).
		HasContentType("application/json").
		JSON().
		Decode(&res)
	assertEquals(s.t, req.ID, res.ID)
	assertNotZero(s.t, res.CreatedAt)
	assertNotZero(s.t, res.UpdatedAt)
	assertEquals(s.t, req.DisplayName, res.DisplayName)
	assertEquals(s.t, req.Icon, res.Icon)
	assertEquals(s.t, req.ParentID, res.ParentID)
	if res.UpdatedAt.Before(res.CreatedAt) {
		s.t.Errorf("Update time should be greater than creation time after patching:\nCreatedAt: %s\nUpdatedAt: %s", res.CreatedAt, res.UpdatedAt)
	}
	return res
}

func (s *session) deleteAccount(req account.DeleteRequest) account.DeleteResponse {
	res := account.DeleteResponse{}
	s.http.DELETE("/tenants/{tenantID}/accounts/{id}").
		WithHeader("Accept", "application/json").
		WithPath("tenantID", req.TenantID).
		WithPath("id", req.ID).
		Expect().
		Status(http.StatusNoContent)
	return res
}
