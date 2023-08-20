package services

import (
	"context"
	"github.com/arikkfir/greenstar/backend/model"
	"github.com/arikkfir/greenstar/backend/util"
	"github.com/arikkfir/greenstar/backend/web"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog/log"
	"github.com/secureworks/errors"
)

type TenantsService struct {
	Service
}

func (s *TenantsService) Tenants(ctx context.Context) ([]*model.Tenant, error) {
	descopeTenants, err := s.Descope.Management.Tenant().LoadAll()
	if err != nil {
		return nil, errors.New("failed loading tenants: %w", err)
	}

	tenants := make([]*model.Tenant, 0)
	for _, tenant := range descopeTenants {
		if web.GetToken(ctx).IsPermittedPerTenant(GlobalTenantID, "tenant:read") || web.GetToken(ctx).IsPermittedPerTenant(tenant.ID, "tenant:read") {
			tenants = append(tenants, &model.Tenant{
				ID:          tenant.ID,
				DisplayName: tenant.Name,
			})
		}
	}
	return tenants, nil
}

func (s *TenantsService) Tenant(ctx context.Context, id string) (*model.Tenant, error) {
	if !web.GetToken(ctx).IsPermittedPerTenant(id, "tenant:read") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	tenant, err := s.Descope.Management.Tenant().Load(id)
	if err != nil {
		return nil, errors.New("failed loading tenant: %w", err)
	}

	return &model.Tenant{
		ID:          tenant.ID,
		DisplayName: tenant.Name,
	}, nil
}

func (s *TenantsService) CreateTenant(ctx context.Context, tenantID *string, tenant model.TenantChanges) (*model.Tenant, error) {
	if !web.GetToken(ctx).IsPermittedPerTenant(GlobalTenantID, "tenant:write") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	var id string
	if tenantID == nil {
		id = util.RandomHash(7)
	} else if *tenantID == GlobalTenantID {
		return nil, errors.New("This tenant ID is reserved.", util.ErrBadRequest, util.UserFacingTag)
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
		if err := web.GetDescopeClient(ctx).Management.Tenant().Delete(id); err != nil {
			log.Ctx(ctx).Error().Stack().Err(err).Msg("failed to delete Descope tenant as a cleanup due to a failure to create its Neo4j DB")
		}
		return nil, errors.New("failed to create tenant: %w", err)
	}

	return &model.Tenant{
		ID:          id,
		DisplayName: tenant.DisplayName,
	}, nil
}

func (s *TenantsService) UpdateTenant(ctx context.Context, tenantID string, tenant model.TenantChanges) (*model.Tenant, error) {
	if tenantID == GlobalTenantID {
		return nil, errors.New("Updating the global tenant is not allowed.", util.ErrBadRequest, util.UserFacingTag)
	}

	if !web.GetToken(ctx).IsPermittedPerTenant(GlobalTenantID, "tenant:write") &&
		!web.GetToken(ctx).IsPermittedPerTenant(tenantID, "tenant:write") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
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
	if tenantID == GlobalTenantID {
		return "", errors.New("Deleting the global tenant is not allowed.", util.ErrBadRequest, util.UserFacingTag)
	}

	if !web.GetToken(ctx).IsPermittedPerTenant(GlobalTenantID, "tenant:write") &&
		!web.GetToken(ctx).IsPermittedPerTenant(tenantID, "tenant:write") {
		return "", errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	if err := s.Descope.Management.Tenant().Delete(tenantID); err != nil {
		return "", errors.New("failed to delete tenant: %w", err)
	}
	return tenantID, nil
}
