package e2e

import (
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/account"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/tenant"
	"github.com/arikkfir/greenstar/backend/internal/util/lang"
	"slices"
	"testing"
)

func TestAccountsEndpoint(t *testing.T) {
	admin := newSession(t, baseURL, exchangeAccessKeyForSessionKey(t, createAdminAccessKey(t)))

	t.Run("Create & Get Account", func(t *testing.T) {
		tnt := admin.createTestTenant()
		user := newSession(t, baseURL, exchangeAccessKeyForSessionKey(t, createTenantUserAccessKey(t, tnt.ID)))

		displayName := fmt.Sprintf("%s-%d Root", t.Name(), processID)
		res := user.createAccount(account.CreateRequest{
			TenantID:    tnt.ID,
			DisplayName: displayName,
		})

		a := user.getAccount(account.GetRequest{TenantID: tnt.ID, ID: res.ID})
		assertEquals(t, displayName, a.DisplayName)
		assertIsNil(t, a.Icon)
		assertIsNil(t, a.ParentID)
	})

	t.Run("List Accounts", func(t *testing.T) {
		tnt := admin.createTestTenant()
		user := newSession(t, baseURL, exchangeAccessKeyForSessionKey(t, createTenantUserAccessKey(t, tnt.ID)))

		const count = 5

		a := user.createAccount(account.CreateRequest{
			TenantID:    tnt.ID,
			DisplayName: fmt.Sprintf("%s-%d Root", t.Name(), processID),
			Icon:        lang.PtrOf("root"),
		})
		a1 := user.createAccount(account.CreateRequest{
			TenantID:    tnt.ID,
			DisplayName: fmt.Sprintf("%s -> %s", a.DisplayName, "A1"),
			Icon:        lang.PtrOf("a1"),
			ParentID:    lang.PtrOf(a.ID),
		})
		a2 := user.createAccount(account.CreateRequest{
			TenantID:    tnt.ID,
			DisplayName: fmt.Sprintf("%s -> %s", a.DisplayName, "A2"),
			Icon:        lang.PtrOf("a2"),
			ParentID:    lang.PtrOf(a.ID),
		})
		b1 := user.createAccount(account.CreateRequest{
			TenantID:    tnt.ID,
			DisplayName: fmt.Sprintf("%s -> %s", a1.DisplayName, "B1"),
			Icon:        lang.PtrOf("b1"),
			ParentID:    lang.PtrOf(a1.ID),
		})
		b2 := user.createAccount(account.CreateRequest{
			TenantID:    tnt.ID,
			DisplayName: fmt.Sprintf("%s -> %s", a2.DisplayName, "B2"),
			Icon:        lang.PtrOf("b2"),
			ParentID:    lang.PtrOf(a2.ID),
		})

		t.Run("Sorting ASC", func(t *testing.T) {
			accounts := []account.Account{account.Account(a), account.Account(a1), account.Account(a2), account.Account(b1), account.Account(b2)}

			listResponse := user.listAccounts(account.ListRequest{
				Count:       lang.PtrOf(uint(count) * 2),
				DisplayName: lang.PtrOf(a.DisplayName),
				Sort:        []string{"displayName:ASC"},
			})
			assertEquals(t, len(accounts), listResponse.TotalCount)
			assertEquals(t, len(accounts), len(listResponse.Items))
			for i := 0; i < len(accounts); i++ {
				assertEquals(t, accounts[i].ID, listResponse.Items[i].ID)
				assertEquals(t, accounts[i].CreatedAt, listResponse.Items[i].CreatedAt)
				assertEquals(t, accounts[i].UpdatedAt, listResponse.Items[i].UpdatedAt)
				assertEquals(t, accounts[i].DisplayName, listResponse.Items[i].DisplayName)
				assertEquals(t, accounts[i].Icon, listResponse.Items[i].Icon)
				assertEquals(t, accounts[i].ParentID, listResponse.Items[i].ParentID)
			}
		})

		t.Run("Sorting DESC", func(t *testing.T) {
			accounts := []account.Account{account.Account(a), account.Account(a1), account.Account(a2), account.Account(b1), account.Account(b2)}
			slices.Reverse(accounts)

			listResponse := user.listAccounts(account.ListRequest{
				Count:       lang.PtrOf(uint(count) * 2),
				DisplayName: lang.PtrOf(a.DisplayName),
				Sort:        []string{"displayName:DESC"},
			})
			assertEquals(t, len(accounts), listResponse.TotalCount)
			assertEquals(t, len(accounts), len(listResponse.Items))
			for i := 0; i < len(accounts); i++ {
				assertEquals(t, accounts[i].ID, listResponse.Items[i].ID)
				assertEquals(t, accounts[i].CreatedAt, listResponse.Items[i].CreatedAt)
				assertEquals(t, accounts[i].UpdatedAt, listResponse.Items[i].UpdatedAt)
				assertEquals(t, accounts[i].DisplayName, listResponse.Items[i].DisplayName)
				assertEquals(t, accounts[i].Icon, listResponse.Items[i].Icon)
				assertEquals(t, accounts[i].ParentID, listResponse.Items[i].ParentID)
			}
		})

		t.Run("Paging", func(t *testing.T) {
			const pageSize = 3
			accounts := []account.Account{account.Account(a), account.Account(a1), account.Account(a2), account.Account(b1), account.Account(b2)}
			firstPageAccounts := []account.Account{account.Account(a), account.Account(a1), account.Account(a2)}
			secondPageAccounts := []account.Account{account.Account(b1), account.Account(b2)}

			firstPageRes := user.listAccounts(account.ListRequest{
				Count:       lang.PtrOf(uint(pageSize)),
				Offset:      lang.PtrOf(uint(0)),
				DisplayName: lang.PtrOf(a.DisplayName),
				Sort:        []string{"displayName:ASC"},
			})
			assertEquals(t, len(accounts), firstPageRes.TotalCount)
			assertEquals(t, len(firstPageAccounts), len(firstPageRes.Items))
			for i := 0; i < len(firstPageAccounts); i++ {
				assertEquals(t, firstPageAccounts[i].ID, firstPageRes.Items[i].ID)
				assertEquals(t, firstPageAccounts[i].CreatedAt, firstPageRes.Items[i].CreatedAt)
				assertEquals(t, firstPageAccounts[i].UpdatedAt, firstPageRes.Items[i].UpdatedAt)
				assertEquals(t, firstPageAccounts[i].DisplayName, firstPageRes.Items[i].DisplayName)
				assertEquals(t, firstPageAccounts[i].Icon, firstPageRes.Items[i].Icon)
				assertEquals(t, firstPageAccounts[i].ParentID, firstPageRes.Items[i].ParentID)
			}

			secondPageRes := user.listAccounts(account.ListRequest{
				Count:       lang.PtrOf(uint(pageSize)),
				Offset:      lang.PtrOf(uint(pageSize)),
				DisplayName: lang.PtrOf(a.DisplayName),
				Sort:        []string{"displayName:ASC"},
			})
			assertEquals(t, len(accounts), secondPageRes.TotalCount)
			assertEquals(t, len(secondPageAccounts), len(secondPageRes.Items))
			for i := 0; i < len(secondPageAccounts); i++ {
				assertEquals(t, secondPageAccounts[i].ID, secondPageRes.Items[i].ID)
				assertEquals(t, secondPageAccounts[i].CreatedAt, secondPageRes.Items[i].CreatedAt)
				assertEquals(t, secondPageAccounts[i].UpdatedAt, secondPageRes.Items[i].UpdatedAt)
				assertEquals(t, secondPageAccounts[i].DisplayName, secondPageRes.Items[i].DisplayName)
				assertEquals(t, secondPageAccounts[i].Icon, secondPageRes.Items[i].Icon)
				assertEquals(t, secondPageAccounts[i].ParentID, secondPageRes.Items[i].ParentID)
			}
		})
	})

	t.Run("Patch Account", func(t *testing.T) {
		tnt := admin.createTestTenant()
		user := newSession(t, baseURL, exchangeAccessKeyForSessionKey(t, createTenantUserAccessKey(t, tnt.ID)))

		res := user.createAccount(account.CreateRequest{
			TenantID:    tnt.ID,
			DisplayName: fmt.Sprintf("%s-%d Root", t.Name(), processID),
		})

		patchedDisplayName := fmt.Sprintf("Patched %s", res.DisplayName)
		patchRes := user.patchAccount(account.PatchRequest{
			TenantID:    tnt.ID,
			ID:          res.ID,
			DisplayName: lang.PtrOf(patchedDisplayName),
		})
		assertEquals(t, patchedDisplayName, patchRes.DisplayName)
	})

	t.Run("Update Account", func(t *testing.T) {
		tnt := admin.createTestTenant()
		user := newSession(t, baseURL, exchangeAccessKeyForSessionKey(t, createTenantUserAccessKey(t, tnt.ID)))

		res := user.createAccount(account.CreateRequest{
			TenantID:    tnt.ID,
			DisplayName: fmt.Sprintf("%s-%d Root", t.Name(), processID),
		})

		updatedDisplayName := fmt.Sprintf("Updated %s", res.DisplayName)
		updateRes := user.updateTenant(tenant.UpdateRequest{
			ID:          res.ID,
			DisplayName: updatedDisplayName,
		})
		assertEquals(t, updatedDisplayName, updateRes.DisplayName)
	})
}
