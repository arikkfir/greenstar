package e2e

import (
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/tenant"
	"github.com/arikkfir/greenstar/backend/internal/util/lang"
	"testing"
)

func TestTenants(t *testing.T) {
	admin := newSession(t, baseURL, exchangeAccessKeyForSessionKey(t, createAdminAccessKey(t)))

	t.Run("Create & Get Tenant", func(t *testing.T) {
		res := admin.createTenant(tenant.CreateRequest{DisplayName: fmt.Sprintf("%s-%d", t.Name(), processID)})
		t.Cleanup(func() { admin.deleteTenant(tenant.DeleteRequest{ID: res.ID}) })

		user := newSession(t, baseURL, exchangeAccessKeyForSessionKey(t, createTenantUserAccessKey(t, res.ID)))
		getRes := user.getTenant(tenant.GetRequest{ID: res.ID})
		assertEquals(t, res.ID, getRes.ID)
		assertEquals(t, res.CreatedAt, getRes.CreatedAt)
		assertEquals(t, res.UpdatedAt, getRes.UpdatedAt)
		assertEquals(t, res.DisplayName, getRes.DisplayName)
	})

	t.Run("List Tenants", func(t *testing.T) {
		baseTenantName := fmt.Sprintf("%s-%d", t.Name(), processID)

		var tenantIDs []string
		var createTenantResps []tenant.CreateResponse
		for i := 0; i < 5; i++ {
			displayName := fmt.Sprintf("%s-%d", baseTenantName, i)
			res := admin.createTenant(tenant.CreateRequest{DisplayName: displayName})
			tenantIDs = append(tenantIDs, res.ID)
			createTenantResps = append(createTenantResps, res)
			t.Cleanup(func() { admin.deleteTenant(tenant.DeleteRequest{ID: res.ID}) })
		}

		user := newSession(t, baseURL, exchangeAccessKeyForSessionKey(t, createTenantUserAccessKey(t, tenantIDs...)))

		t.Run("Sorting ASC", func(t *testing.T) {
			listResponse := user.listTenants(tenant.ListRequest{
				Count:       lang.PtrOf(uint(len(createTenantResps)) * 2),
				DisplayName: &baseTenantName,
				Sort:        []string{"displayName:ASC"},
			})
			assertEquals(t, len(createTenantResps), listResponse.TotalCount)
			assertEquals(t, len(createTenantResps), len(listResponse.Items))
			for i := 0; i < len(listResponse.Items); i++ {
				assertEquals(t, createTenantResps[i].ID, listResponse.Items[i].ID)
				assertEquals(t, createTenantResps[i].CreatedAt, listResponse.Items[i].CreatedAt)
				assertEquals(t, createTenantResps[i].UpdatedAt, listResponse.Items[i].UpdatedAt)
				assertEquals(t, createTenantResps[i].DisplayName, listResponse.Items[i].DisplayName)
			}
		})

		t.Run("Sorting DESC", func(t *testing.T) {
			listResponse := user.listTenants(tenant.ListRequest{
				Count:       lang.PtrOf(uint(len(createTenantResps)) * 2),
				DisplayName: &baseTenantName,
				Sort:        []string{"displayName:DESC"},
			})
			assertEquals(t, len(createTenantResps), listResponse.TotalCount)
			assertEquals(t, len(createTenantResps), len(listResponse.Items))
			for i := 0; i < len(listResponse.Items); i++ {
				assertEquals(t, createTenantResps[4-i].ID, listResponse.Items[i].ID)
				assertEquals(t, createTenantResps[4-i].CreatedAt, listResponse.Items[i].CreatedAt)
				assertEquals(t, createTenantResps[4-i].UpdatedAt, listResponse.Items[i].UpdatedAt)
				assertEquals(t, createTenantResps[4-i].DisplayName, listResponse.Items[i].DisplayName)
			}
		})

		t.Run("Paging", func(t *testing.T) {
			const pageSize = 3
			firstPageRes := user.listTenants(tenant.ListRequest{
				Count:       lang.PtrOf(uint(pageSize)),
				Offset:      lang.PtrOf(uint(0)),
				DisplayName: &baseTenantName,
				Sort:        []string{"displayName:ASC"},
			})
			assertEquals(t, len(createTenantResps), firstPageRes.TotalCount)
			assertEquals(t, pageSize, len(firstPageRes.Items))
			for i := 0; i < len(firstPageRes.Items); i++ {
				assertEquals(t, createTenantResps[i].ID, firstPageRes.Items[i].ID)
				assertEquals(t, createTenantResps[i].CreatedAt, firstPageRes.Items[i].CreatedAt)
				assertEquals(t, createTenantResps[i].UpdatedAt, firstPageRes.Items[i].UpdatedAt)
				assertEquals(t, createTenantResps[i].DisplayName, firstPageRes.Items[i].DisplayName)
			}

			secondPageRes := user.listTenants(tenant.ListRequest{
				Count:       lang.PtrOf(uint(pageSize)),
				Offset:      lang.PtrOf(uint(pageSize)),
				DisplayName: &baseTenantName,
				Sort:        []string{"displayName:ASC"},
			})
			assertEquals(t, len(createTenantResps), firstPageRes.TotalCount)
			assertEquals(t, pageSize-len(firstPageRes.Items), len(secondPageRes.Items))
			for i := 0; i < len(secondPageRes.Items); i++ {
				assertEquals(t, createTenantResps[i+pageSize].ID, firstPageRes.Items[i].ID)
				assertEquals(t, createTenantResps[i+pageSize].CreatedAt, firstPageRes.Items[i].CreatedAt)
				assertEquals(t, createTenantResps[i+pageSize].UpdatedAt, firstPageRes.Items[i].UpdatedAt)
				assertEquals(t, createTenantResps[i+pageSize].DisplayName, firstPageRes.Items[i].DisplayName)
			}
		})
	})

	t.Run("Patch Tenant", func(t *testing.T) {
		res := admin.createTenant(tenant.CreateRequest{DisplayName: fmt.Sprintf("%s-%d", t.Name(), processID)})
		t.Cleanup(func() { admin.deleteTenant(tenant.DeleteRequest{ID: res.ID}) })

		user := newSession(t, baseURL, exchangeAccessKeyForSessionKey(t, createTenantUserAccessKey(t, res.ID)))

		patchedDisplayName := fmt.Sprintf("Patched %s", res.DisplayName)
		patchRes := user.patchTenant(tenant.PatchRequest{
			ID:          res.ID,
			DisplayName: lang.PtrOf(patchedDisplayName),
		})
		assertEquals(t, res.ID, patchRes.ID)
		assertEquals(t, res.CreatedAt, patchRes.CreatedAt)
		assertNotBefore(t, patchRes.UpdatedAt, res.UpdatedAt)
		assertEquals(t, patchedDisplayName, patchRes.DisplayName)
	})

	t.Run("Update Tenant", func(t *testing.T) {
		res := admin.createTenant(tenant.CreateRequest{DisplayName: fmt.Sprintf("%s-%d", t.Name(), processID)})
		t.Cleanup(func() { admin.deleteTenant(tenant.DeleteRequest{ID: res.ID}) })

		user := newSession(t, baseURL, exchangeAccessKeyForSessionKey(t, createTenantUserAccessKey(t, res.ID)))

		updatedDisplayName := fmt.Sprintf("Updated %s", res.DisplayName)
		updateRes := user.updateTenant(tenant.UpdateRequest{
			ID:          res.ID,
			DisplayName: updatedDisplayName,
		})
		assertEquals(t, res.ID, updateRes.ID)
		assertEquals(t, res.CreatedAt, updateRes.CreatedAt)
		assertNotBefore(t, updateRes.UpdatedAt, res.UpdatedAt)
		assertEquals(t, updatedDisplayName, updateRes.DisplayName)
	})
}
