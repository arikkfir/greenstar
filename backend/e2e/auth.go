package e2e

import (
	"context"
	"fmt"
	"github.com/descope/go-sdk/descope"
	"testing"
)

var (
	AdminRoles = []string{"Admin"}
	UserRoles  = []string{"User"}
)

func exchangeAccessKeyForSessionKey(t *testing.T, accessKey string) string {
	_, token, err := descopeClient.Auth.ExchangeAccessKey(context.Background(), accessKey, nil)
	if err != nil {
		t.Fatalf("failed to exchange access key: %v", err)
	}
	return token.JWT
}

func createAdminAccessKey(t *testing.T) string {
	accessKeyMgr := descopeClient.Management.AccessKey()
	accessKeyName := fmt.Sprintf("%s-admin-%d", t.Name(), processID)

	key, resp, err := accessKeyMgr.Create(context.Background(), accessKeyName, 0, AdminRoles, nil, "", nil)
	if err != nil {
		t.Fatalf("failed to create admin access key: %v", err)
	}

	t.Cleanup(func() {
		if err := accessKeyMgr.Delete(context.Background(), resp.ID); err != nil {
			t.Errorf("failed to delete admin access key '%s': %v", resp.ID, err)
		}
	})

	return key
}

func createTenantUserAccessKey(t *testing.T, tenantIDs ...string) string {
	accessKeyMgr := descopeClient.Management.AccessKey()
	accessKeyName := fmt.Sprintf("%s-user-%d", t.Name(), processID)

	var tenants []*descope.AssociatedTenant
	for _, id := range tenantIDs {
		tenants = append(tenants, &descope.AssociatedTenant{TenantID: id, Roles: UserRoles})
	}
	key, resp, err := accessKeyMgr.Create(context.Background(), accessKeyName, 0, nil, tenants, "", nil)
	if err != nil {
		t.Fatalf("failed to create user access key: %v", err)
	}

	t.Cleanup(func() {
		if err := accessKeyMgr.Delete(context.Background(), resp.ID); err != nil {
			t.Errorf("failed to delete user access key '%s': %v", resp.ID, err)
		}
	})

	return key
}
