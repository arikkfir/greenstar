package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.35

import (
	"context"

	"github.com/arikkfir/greenstar/backend/gql"
	"github.com/arikkfir/greenstar/backend/model"
)

// Accounts is the resolver for the accounts field.
func (r *tenantResolver) Accounts(ctx context.Context, obj *model.Tenant) ([]*model.Account, error) {
	return r.AccountsService.Accounts(ctx, obj)
}

// Account is the resolver for the account field.
func (r *tenantResolver) Account(ctx context.Context, obj *model.Tenant, id string) (*model.Account, error) {
	return r.AccountsService.Account(ctx, obj, id)
}

// Transactions is the resolver for the transactions field.
func (r *tenantResolver) Transactions(ctx context.Context, obj *model.Tenant) ([]*model.Transaction, error) {
	return r.TransactionsService.Transactions(ctx, obj)
}

// Tenant returns gql.TenantResolver implementation.
func (r *Resolver) Tenant() gql.TenantResolver { return &tenantResolver{r} }

type tenantResolver struct{ *Resolver }
