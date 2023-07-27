package services

import (
	"context"
	"github.com/arik-kfir/greenstar/backend/model"
	"github.com/arik-kfir/greenstar/backend/util"
	"github.com/arik-kfir/greenstar/backend/web"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/secureworks/errors"
)

type TenantsService struct {
	Service
}

func (s *TenantsService) Tenants(ctx context.Context) ([]*model.Tenant, error) {
	if !web.GetToken(ctx).IsPermittedPerTenant(GlobalTenantID, "Manage tenants") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	descopeTenants, err := s.Descope.Management.Tenant().LoadAll()
	if err != nil {
		return nil, errors.New("failed loading tenants: %w", err)
	}

	tenants := make([]*model.Tenant, 0)
	for _, tenant := range descopeTenants {
		tenants = append(tenants, &model.Tenant{
			ID:          tenant.ID,
			DisplayName: tenant.Name,
		})
	}
	return tenants, nil
}

func (s *TenantsService) Tenant(ctx context.Context, id string) (*model.Tenant, error) {
	if !web.GetToken(ctx).IsPermittedPerTenant(GlobalTenantID, "Manage tenants") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	descopeTenants, err := s.Descope.Management.Tenant().LoadAll()
	if err != nil {
		return nil, errors.New("failed loading tenant: %w", err)
	}
	for _, tenant := range descopeTenants {
		if tenant.ID == id {
			return &model.Tenant{
				ID:          tenant.ID,
				DisplayName: tenant.Name,
			}, nil
		}
	}
	return nil, errors.New("tenant not found")
}

func (s *TenantsService) CreateTenant(ctx context.Context, tenantID *string, tenant model.TenantChanges) (*model.Tenant, error) {
	if !web.GetToken(ctx).IsPermittedPerTenant(GlobalTenantID, "Manage tenants") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	var id string
	if tenantID == nil {
		id = util.RandomHash(7)
	} else if *tenantID == GlobalTenantID {
		return nil, errors.New("Tenant already exists.", util.ErrBadRequest, util.UserFacingTag)
	} else {
		id = *tenantID
	}

	if err := s.Descope.Management.Tenant().CreateWithID(id, tenant.DisplayName, nil); err != nil {
		return nil, errors.New("failed to create tenant: %w", err)
	}

	createDBFunc := func(tx neo4j.ManagedTransaction) (any, error) {
		//goland:noinspection SqlNoDataSourceInspection
		const createTenantDBQuery = `CREATE DATABASE $id IF NOT EXISTS`
		_, err := tx.Run(ctx, createTenantDBQuery, map[string]any{"id": id})
		if err != nil {
			return nil, errors.New("failed to create tenant database: %w\n%s", err, createTenantDBQuery)
		}
		return nil, nil
	}

	createDBSession := s.getNeo4jSessionForSystem(ctx, neo4j.AccessModeRead)
	defer createDBSession.Close(ctx)

	if _, err := createDBSession.ExecuteWrite(ctx, createDBFunc); err != nil {
		// TODO: delete tenant in Descope via Descope Management API
		return nil, errors.New("failed to create tenant: %w", err)
	}

	return &model.Tenant{
		ID:          id,
		DisplayName: tenant.DisplayName,
	}, nil
}

func (s *TenantsService) UpdateTenant(ctx context.Context, tenantID string, tenant model.TenantChanges) (*model.Tenant, error) {
	if !web.GetToken(ctx).IsPermittedPerTenant(GlobalTenantID, "Manage tenants") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	if tenantID == GlobalTenantID {
		return nil, errors.New("Updating the global tenant is not allowed.", util.ErrBadRequest, util.UserFacingTag)
	}

	if err := s.Descope.Management.Tenant().Update(tenantID, tenant.DisplayName, nil); err != nil {
		return nil, errors.New("failed to update tenant: %w", err)
	}
	return &model.Tenant{
		ID:          tenantID,
		DisplayName: tenant.DisplayName,
	}, nil
}

func (s *TenantsService) DeleteTenant(ctx context.Context, tenantID string) (string, error) {
	if !web.GetToken(ctx).IsPermittedPerTenant(GlobalTenantID, "Manage tenants") {
		return "", errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	if tenantID == GlobalTenantID {
		return "", errors.New("Deleting the global tenant is not allowed.", util.ErrBadRequest, util.UserFacingTag)
	}

	if err := s.Descope.Management.Tenant().Delete(tenantID); err != nil {
		return "", errors.New("failed to delete tenant: %w", err)
	}
	return tenantID, nil
}
