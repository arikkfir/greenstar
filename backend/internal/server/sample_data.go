package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/auth"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/account"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/tenant"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/transaction"
	"github.com/arikkfir/greenstar/backend/internal/server/util"
	"github.com/arikkfir/greenstar/backend/internal/util/db"
	"github.com/arikkfir/greenstar/backend/internal/util/lang"
	"github.com/arikkfir/greenstar/backend/internal/util/strings"
	pgxorig "github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"log/slog"
	"math/rand/v2"
	"time"
)

type SampleDataError struct {
	msg   string
	cause error
}

func (e *SampleDataError) Error() string {
	return e.msg
}

func (e *SampleDataError) Unwrap() error {
	return e.cause
}

func (s *Server) newAccount(ctx context.Context, tenantID, displayName, icon string, parent *account.Account) account.Account {
	var parentID *string
	if parent != nil {
		parentID = lang.PtrOf(parent.ID)
	}
	if res, err := s.AccountsHandler.Create(ctx, account.CreateRequest{
		TenantID:    tenantID,
		DisplayName: displayName,
		Icon:        strings.NilIfEmpty(icon),
		ParentID:    parentID,
	}); err != nil {
		panic(&SampleDataError{msg: fmt.Sprintf("failed creating '%s' account: %s", displayName, err), cause: err})
	} else {
		return account.Account(*res)
	}
}

func (s *Server) newTx(ctx context.Context, tenantID string, date time.Time, source, target account.Account, description string, amount decimal.Decimal, currency string) transaction.Transaction {
	if t, err := s.TransactionsHandler.Create(ctx, transaction.CreateRequest{
		TenantID:        tenantID,
		Date:            date,
		ReferenceID:     strings.RandomHash(7),
		Amount:          amount,
		Currency:        currency,
		Description:     lang.PtrOf(description),
		SourceAccountID: source.ID,
		TargetAccountID: target.ID,
	}); err != nil {
		panic(&SampleDataError{msg: fmt.Sprintf("failed creating transaction: %s", err), cause: err})
	} else {
		return transaction.Transaction(*t)
	}
}

func (s *Server) newDecimal(ds string) decimal.Decimal {
	d, err := decimal.NewFromString(ds)
	if err != nil {
		panic(err)
	}
	return d
}

func (s *Server) randomDate() time.Time {
	return time.Now().AddDate(0, rand.IntN(12)*-1, rand.IntN(30)*-1)
}

func (s *Server) GenerateSampleTenant(ctx context.Context, tenantID, tenantDisplayName, accessKey string) (result error) {
	slog.Default().InfoContext(ctx, "Generating sample data")

	_, token, err := s.Descope.Auth.ExchangeAccessKey(ctx, accessKey, nil)
	if err != nil {
		return fmt.Errorf("failed exchanging access key: %w", err)
	}
	ctx = auth.NewContextWithToken(ctx, &auth.Token{Token: *token, ID: token.ID})

	// Start a transaction
	txOptions := pgxorig.TxOptions{
		IsoLevel:       pgxorig.Serializable,
		AccessMode:     pgxorig.ReadWrite,
		DeferrableMode: pgxorig.NotDeferrable,
	}
	tx, err := s.Pool.BeginTx(ctx, txOptions)
	if err != nil {
		return fmt.Errorf("failed beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Make transaction available in the context we pass to the model handlers
	ctx = db.NewContextWithTx(ctx, tx)

	// Delete the tenant, ignoring errors
	if _, err := s.TenantsHandler.Get(ctx, tenant.GetRequest{ID: tenantID}); err != nil {
		if errors.Is(err, util.ErrNotFound) {
			if _, err = s.TenantsHandler.Create(ctx, tenant.CreateRequest{ID: tenantID, DisplayName: tenantDisplayName}); err != nil {
				return fmt.Errorf("failed creating tenant: %w", err)
			}
		} else {
			return fmt.Errorf("failed looking up tenant '%s': %w", tenantID, err)
		}
	} else if err := s.TransactionsHandler.DeleteAll(ctx, transaction.DeleteAllRequest{TenantID: tenantID}); err != nil {
		return fmt.Errorf("failed deleting transactions for tenant '%s': %w", tenantID, err)
	} else if err := s.AccountsHandler.DeleteAll(ctx, account.DeleteAllRequest{TenantID: tenantID}); err != nil {
		return fmt.Errorf("failed deleting accounts for tenant '%s': %w", tenantID, err)
	}

	// Allows utility functions to panic instead of returning errors - and we transform OUR panics back into errors
	// This allows for a more fluid code control
	defer func() {
		if result == nil {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					var sampleDataError *SampleDataError
					if errors.As(err, &sampleDataError) {
						result = err
						return
					}
				}
				panic(err)
			}
		}
	}()

	now := time.Now()

	// Create data
	employments := s.newAccount(ctx, tenantID, "Employers", "MapsHomeWork", nil)
	microsoft := s.newAccount(ctx, tenantID, "Microsoft", "Microsoft", &employments)
	amazon := s.newAccount(ctx, tenantID, "Amazon", "Amazon", &employments)
	google := s.newAccount(ctx, tenantID, "Google", "Google", &employments)

	checkingAccounts := s.newAccount(ctx, tenantID, "Checking Accounts", "BankTransfer", nil)
	bankOfAmerica := s.newAccount(ctx, tenantID, "Bank of America", "BankOfAmerica", &checkingAccounts)
	deutscheBank := s.newAccount(ctx, tenantID, "Deutsche Bank", "DeutscheBank", &checkingAccounts)

	assets := s.newAccount(ctx, tenantID, "Assets", "Savings", nil)
	realEstate := s.newAccount(ctx, tenantID, "Real Estate", "Apartment", &assets)
	cars := s.newAccount(ctx, tenantID, "Cars", "DirectionsCarFilled", &assets)

	leisure := s.newAccount(ctx, tenantID, "Leisure", "BeachAccess", nil)

	entertainment := s.newAccount(ctx, tenantID, "Entertainment", "Movie", &leisure)
	netflix := s.newAccount(ctx, tenantID, "Netflix", "Netflix", &entertainment)
	amazonPrime := s.newAccount(ctx, tenantID, "Amazon Prime", "AmazonPrime", &entertainment)
	spotify := s.newAccount(ctx, tenantID, "Spotify", "Spotify", &entertainment)
	youtube := s.newAccount(ctx, tenantID, "YouTube", "YouTube", &entertainment)
	foodAndDrink := s.newAccount(ctx, tenantID, "Food and Drink", "Fastfood", &leisure)

	household := s.newAccount(ctx, tenantID, "Household", "Home", &leisure)
	householdSupplies := s.newAccount(ctx, tenantID, "Household Supplies", "HouseholdSupplies", &household)
	householdGroceries := s.newAccount(ctx, tenantID, "Groceries", "LocalGroceryStore", &householdSupplies)
	householdMeat := s.newAccount(ctx, tenantID, "Meat", "FoodTurkey", &householdSupplies)
	householdElectronics := s.newAccount(ctx, tenantID, "Household Electronics", "ElectricBolt", &household)

	clothing := s.newAccount(ctx, tenantID, "Clothing", "Checkroom", nil)
	shoes := s.newAccount(ctx, tenantID, "Shoes", "ShoeSneaker", &clothing)
	nike := s.newAccount(ctx, tenantID, "Nike", "Nike", &shoes)
	adidas := s.newAccount(ctx, tenantID, "Adidas", "Adidas", &shoes)
	zara := s.newAccount(ctx, tenantID, "Zara", "Zara", &clothing)

	for i := range 6 {
		s.newTx(ctx, tenantID, now.AddDate(0, -i, 0), microsoft, bankOfAmerica, "Salary", decimal.NewFromInt(25000).Round(2), "ILS")
	}
	for i := range 3 {
		s.newTx(ctx, tenantID, now.AddDate(0, -6-i, 0), amazon, deutscheBank, "Salary", decimal.NewFromInt(30000).Round(2), "ILS")
	}
	for i := range 3 {
		s.newTx(ctx, tenantID, now.AddDate(0, -6-3-i, 0), google, deutscheBank, "Salary", decimal.NewFromInt(35000).Round(2), "ILS")
	}
	for i := range 12 {
		s.newTx(ctx, tenantID, now.AddDate(0, -i, -now.Day()+1), bankOfAmerica, realEstate, "Mortgage payment", decimal.NewFromInt(5000).Round(2), "ILS")
		s.newTx(ctx, tenantID, now.AddDate(0, -i, -now.Day()+1), bankOfAmerica, cars, "Car payment", decimal.NewFromInt(5000).Round(2), "ILS")
		s.newTx(ctx, tenantID, now.AddDate(0, -i, -now.Day()+1), deutscheBank, netflix, "Netflix Subscription", s.newDecimal("29.99").Round(2), "ILS")
		s.newTx(ctx, tenantID, now.AddDate(0, -i, -now.Day()+1), deutscheBank, amazonPrime, "AmazonPrime Subscription", s.newDecimal("14.98").Round(2), "ILS")
		s.newTx(ctx, tenantID, now.AddDate(0, -i, -now.Day()+1), deutscheBank, spotify, "Spotify Subscription", s.newDecimal("9.99").Round(2), "ILS")
	}

	for i := range 3 {
		s.newTx(ctx, tenantID, now.AddDate(0, -i, -now.Day()+1), deutscheBank, youtube, "YouTube Premium Subscription", s.newDecimal("12.99").Round(2), "ILS")
	}

	for range 2 {
		s.newTx(ctx, tenantID, s.randomDate(), bankOfAmerica, foodAndDrink, "Going out", decimal.NewFromFloat(rand.Float64()*500).Round(2), "ILS")
		s.newTx(ctx, tenantID, s.randomDate(), bankOfAmerica, householdGroceries, "Groceries", decimal.NewFromFloat(rand.Float64()*200).Round(2), "ILS")
		s.newTx(ctx, tenantID, s.randomDate(), bankOfAmerica, householdMeat, "Meat", decimal.NewFromFloat(rand.Float64()*700).Round(2), "ILS")
		s.newTx(ctx, tenantID, s.randomDate(), bankOfAmerica, householdElectronics, "Electronics", decimal.NewFromFloat(rand.Float64()*2000).Round(2), "ILS")
		s.newTx(ctx, tenantID, s.randomDate(), bankOfAmerica, nike, "Nike Sneakers", decimal.NewFromFloat(rand.Float64()*900).Round(2), "ILS")
		s.newTx(ctx, tenantID, s.randomDate(), bankOfAmerica, adidas, "Adidas running shoes", decimal.NewFromFloat(rand.Float64()*1200).Round(2), "ILS")
		s.newTx(ctx, tenantID, s.randomDate(), bankOfAmerica, zara, "Clothes from Zara", decimal.NewFromFloat(rand.Float64()*1200).Round(2), "ILS")
	}

	// Commit
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed committing transaction: %w", err)
	}

	return nil
}
