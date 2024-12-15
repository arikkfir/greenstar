package sample

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/auth"
	"github.com/arikkfir/greenstar/backend/internal/server/middleware"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/account"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/tenant"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/transaction"
	"github.com/arikkfir/greenstar/backend/internal/server/util"
	"github.com/arikkfir/greenstar/backend/internal/util/db"
	"github.com/descope/go-sdk/descope/client"
	pgxorig "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/yaml.v3"
	"log/slog"
)

var (
	//go:embed acme-accounts.yaml
	acmeAccountsYAML string
	g                generator
)

func init() {
	if err := yaml.Unmarshal([]byte(acmeAccountsYAML), &g); err != nil {
		panic(fmt.Errorf("failed parsing acme-accounts.yaml: %w", err))
	} else if g.DefaultCurrency == "" {
		panic(fmt.Errorf("failed parsing acme-accounts.yaml: no default currency"))
	}
}

type generator struct {
	DefaultCurrency string                `yaml:"defaultCurrency"`
	Accounts        []accountWithChildren `yaml:"accounts"`
}

func Generate(ctx context.Context, descopeClient *client.DescopeClient, accessKey string, pool *pgxpool.Pool, th tenant.Handler, tenantID, tenantDisplayName string, txh transaction.Handler, ah account.Handler) error {
	slog.InfoContext(ctx, "Generating sample data")
	ctx = middleware.WithTenantID(ctx, tenantID)

	// Exchange an access key for a token and expose it in a context, simulating the way it's done in real HTTP requests
	_, token, err := descopeClient.Auth.ExchangeAccessKey(ctx, accessKey, nil)
	if err != nil {
		return fmt.Errorf("failed exchanging access key: %w", err)
	}
	ctx = auth.NewContextWithToken(ctx, &auth.Token{Token: *token, ID: token.ID})

	// Start a transaction & expose it in a context, simulating the way it's done in real HTTP requests
	txOptions := pgxorig.TxOptions{
		IsoLevel:       pgxorig.Serializable,
		AccessMode:     pgxorig.ReadWrite,
		DeferrableMode: pgxorig.NotDeferrable,
	}
	tx, err := pool.BeginTx(ctx, txOptions)
	if err != nil {
		return fmt.Errorf("failed starting a database transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	ctx = db.NewContextWithTx(ctx, tx)

	accountIDs := make(map[string]string)

	// Generate the data
	if err := generateTenant(ctx, th, tenantID, tenantDisplayName, txh, ah); err != nil {
		return fmt.Errorf("failed generating tenant: %w", err)
	}
	if err := generateAccounts(ctx, accountIDs, ah, tenantID); err != nil {
		return fmt.Errorf("failed generating accounts: %w", err)
	}
	if err := generateTransactions(ctx, accountIDs, txh, tenantID); err != nil {
		return fmt.Errorf("failed generating transactions: %w", err)
	}

	// Commit
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed committing transaction: %w", err)
	}

	return nil
}

func generateTenant(ctx context.Context, th tenant.Handler, tenantID, tenantDisplayName string, txh transaction.Handler, ah account.Handler) error {
	slog.InfoContext(ctx, "Generating sample tenant")

	if _, err := th.Get(ctx, tenant.GetRequest{ID: tenantID}); err != nil {
		if errors.Is(err, util.ErrNotFound) {
			if _, err = th.Create(ctx, tenant.CreateRequest{ID: tenantID, DisplayName: tenantDisplayName}); err != nil {
				return fmt.Errorf("failed creating tenant: %w", err)
			}
		} else {
			return fmt.Errorf("failed looking up tenant '%s': %w", tenantID, err)
		}
	} else if err := txh.DeleteAll(ctx, transaction.DeleteAllRequest{}); err != nil {
		return fmt.Errorf("failed deleting transactions for tenant '%s': %w", tenantID, err)
	} else if err := ah.DeleteAll(ctx, account.DeleteAllRequest{}); err != nil {
		return fmt.Errorf("failed deleting accounts for tenant '%s': %w", tenantID, err)
	}

	return nil
}

func generateAccounts(ctx context.Context, ids map[string]string, ah account.Handler, tenantID string) error {
	slog.InfoContext(ctx, "Generating sample accounts")
	for _, a := range g.Accounts {
		if err := a.applyAccount(ctx, ids, ah, tenantID, nil); err != nil {
			return fmt.Errorf("failed applying root account '%s': %w", a.ID, err)
		}
	}
	return nil
}

func generateTransactions(ctx context.Context, ids map[string]string, th transaction.Handler, tenantID string) error {
	slog.InfoContext(ctx, "Generating sample transactions")
	for _, a := range g.Accounts {
		if err := a.applyOutgoingTransactions(ctx, ids, th, tenantID, g.DefaultCurrency); err != nil {
			return fmt.Errorf("failed applying transactions of root account '%s': %w", a.ID, err)
		}
	}
	return nil
}
