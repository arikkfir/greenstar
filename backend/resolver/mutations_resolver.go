package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.32

import (
	"context"

	"github.com/arik-kfir/greenstar/backend/gql"
	"github.com/arik-kfir/greenstar/backend/model"
)

// CreateTenant is the resolver for the createTenant field.
func (r *mutationResolver) CreateTenant(ctx context.Context, tenantID *string, tenant model.TenantChanges) (*model.Tenant, error) {
	return r.TenantsService.CreateTenant(ctx, tenantID, tenant)
}

// UpdateTenant is the resolver for the updateTenant field.
func (r *mutationResolver) UpdateTenant(ctx context.Context, tenantID string, tenant model.TenantChanges) (*model.Tenant, error) {
	return r.TenantsService.UpdateTenant(ctx, tenantID, tenant)
}

// DeleteTenant is the resolver for the deleteTenant field.
func (r *mutationResolver) DeleteTenant(ctx context.Context, tenantID string) (string, error) {
	return r.TenantsService.DeleteTenant(ctx, tenantID)
}

// CreateAccount is the resolver for the createAccount field.
func (r *mutationResolver) CreateAccount(ctx context.Context, accountID *string, account model.AccountChanges) (*model.Account, error) {
	return r.AccountsService.CreateAccount(ctx, accountID, account)
}

// UpdateAccount is the resolver for the updateAccount field.
func (r *mutationResolver) UpdateAccount(ctx context.Context, accountID string, account model.AccountChanges) (*model.Account, error) {
	return r.AccountsService.UpdateAccount(ctx, accountID, account)
}

// DeleteAccount is the resolver for the deleteAccount field.
func (r *mutationResolver) DeleteAccount(ctx context.Context, accountID string) (string, error) {
	return r.AccountsService.DeleteAccount(ctx, accountID)
}

// CreateTransaction is the resolver for the createTransaction field.
func (r *mutationResolver) CreateTransaction(ctx context.Context, transaction model.TransactionChanges) (*model.Transaction, error) {
	return r.TransactionsService.CreateTransaction(ctx, transaction)
}

// CreateTransactions is the resolver for the createTransactions field.
func (r *mutationResolver) CreateTransactions(ctx context.Context, transactions []*model.TransactionChanges) (int, error) {
	return r.TransactionsService.CreateTransactions(ctx, transactions)
}

// ScrapeIsraelBankYahav is the resolver for the scrapeIsraelBankYahav field.
func (r *mutationResolver) ScrapeIsraelBankYahav(ctx context.Context, username string, id string, password string) (string, error) {
	return r.TransactionsService.ScrapeIsraelBankYahav(ctx, username, id, password)
}

// UpdateOperation is the resolver for the updateOperation field.
func (r *mutationResolver) UpdateOperation(ctx context.Context, id string, op model.OperationChanges) (*model.Operation, error) {
	return r.OperationsService.UpdateOperation(ctx, id, op)
}

// Mutation returns gql.MutationResolver implementation.
func (r *Resolver) Mutation() gql.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
