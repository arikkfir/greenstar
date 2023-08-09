package services

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/arikkfir/greenstar/backend/model"
	"github.com/arikkfir/greenstar/backend/util"
	"github.com/arikkfir/greenstar/backend/web"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog/log"
	"github.com/secureworks/errors"
	"strconv"
	"text/template"
	"time"
)

var (
	//go:embed scraper-job.tmpl.yaml
	scraperJobTmplString string
	scraperJobTmpl       *template.Template
)

func init() {
	tmpl, err := template.New("scraper-job").Parse(scraperJobTmplString)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse scraper-job.tmpl.yaml")
	} else {
		scraperJobTmpl = tmpl
	}
}

type TransactionsService struct {
	Service
}

func (s *TransactionsService) CreateTransaction(ctx context.Context, tenantID string, transaction model.TransactionChanges) (*model.Transaction, error) {
	if tenantID == GlobalTenantID {
		return nil, errors.New(util.ErrBadRequest, util.UserFacingTag)
	} else if !web.GetToken(ctx).IsPermittedPerTenant(tenantID, "Create transactions") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	session := s.getNeo4jSessionForTenant(ctx, neo4j.AccessModeWrite, tenantID)
	defer session.Close(ctx)

	v, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		id := uuid.NewString()

		createTxQuery := `// Create transaction
MATCH (sourceAccount:Account {accountID: $sourceAccountID})
MATCH (targetAccount:Account {accountID: $targetAccountID})
CREATE (sourceAccount)-[r:TransferredTo {txID: $id, date: $date, refID: $refID, amount: $amount, description: $description}]->(targetAccount)
RETURN sourceAccount.id, sourceAccount.displayName, r, targetAccount.id, targetAccount.displayName`
		createTxParams := map[string]interface{}{}
		createTxParams["sourceAccountID"] = transaction.SourceAccountID
		createTxParams["targetAccountID"] = transaction.TargetAccountID
		createTxParams["id"] = id
		createTxParams["date"] = transaction.Date
		createTxParams["refID"] = transaction.ReferenceID
		createTxParams["amount"] = transaction.Amount
		createTxParams["description"] = transaction.Description

		result, err := tx.Run(ctx, createTxQuery, createTxParams)
		if err != nil {
			return nil, errors.New("failed to execute query: %w\n%s", err, createTxQuery)
		}

		rec, err := result.Single(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to collect single record: %w", err)
		}

		txRel := rec.Values[2].(neo4j.Relationship)
		return &model.Transaction{
			ID:            txRel.Props["txID"].(string),
			Date:          txRel.Props["date"].(time.Time),
			TargetAccount: &model.Account{Tenant: &model.Tenant{ID: tenantID}, ID: rec.Values[3].(string), DisplayName: rec.Values[4].(string)},
			SourceAccount: &model.Account{Tenant: &model.Tenant{ID: tenantID}, ID: rec.Values[0].(string), DisplayName: rec.Values[1].(string)},
			ReferenceID:   txRel.Props["refID"].(string),
			Amount:        txRel.Props["amount"].(model.Money),
			Description:   txRel.Props["description"].(string),
		}, nil
	})
	return v.(*model.Transaction), err
}

func (s *TransactionsService) CreateTransactions(ctx context.Context, tenantID string, transactions []*model.TransactionChanges) (int, error) {
	if tenantID == GlobalTenantID {
		return 0, errors.New(util.ErrBadRequest, util.UserFacingTag)
	} else if !web.GetToken(ctx).IsPermittedPerTenant(tenantID, "Create transactions") {
		return 0, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	session := s.getNeo4jSessionForTenant(ctx, neo4j.AccessModeWrite, tenantID)
	defer session.Close(ctx)

	v, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		createTxQuery := `// Create transactions list`
		createTxParams := map[string]interface{}{}

		for i, input := range transactions {
			createTxQuery += fmt.Sprintf(`// Create transaction %d
MATCH (sourceAccount%d:Account {accountID: $sourceAccountID%d})
MATCH (targetAccount%d:Account {accountID: $targetAccountID%d})
CREATE (sourceAccount%d)-[:TransferredTo {
	txID: $id%d, 
	date: $date%d, 
	refID: $refID%d, 
	amount: $amount%d, 
	description: $description%d
}]->(targetAccount%d)`,
				i, i, i, i, i, i, i, i, i, i, i, i)
			createTxParams["sourceAccountID"+strconv.Itoa(i)] = input.SourceAccountID
			createTxParams["targetAccountID"+strconv.Itoa(i)] = input.TargetAccountID
			createTxParams["id"+strconv.Itoa(i)] = uuid.NewString()
			createTxParams["date"+strconv.Itoa(i)] = input.Date
			createTxParams["refID"+strconv.Itoa(i)] = input.ReferenceID
			createTxParams["amount"+strconv.Itoa(i)] = input.Amount
			createTxParams["description"+strconv.Itoa(i)] = input.Description
		}

		result, err := tx.Run(ctx, createTxQuery, createTxParams)
		if err != nil {
			return nil, errors.New("failed to execute query: %w\n%s", err, createTxQuery)
		}

		summary, err := result.Consume(ctx)
		if err != nil {
			return nil, errors.New("failed to consume query summary: %w", err)
		}

		relationshipsCreated := summary.Counters().RelationshipsCreated()
		if relationshipsCreated != len(transactions) {
			return nil, errors.New("incorrect number of transaction relationships created, expected %d, got %d", len(transactions), relationshipsCreated)
		}

		return relationshipsCreated, nil
	})
	return v.(int), err
}

func (s *TransactionsService) ScrapeIsraelBankYahav(ctx context.Context, tenantID, username, id, password string) (string, error) {
	panic(errors.New("not implemented: ScrapeIsraelBankYahav"))
}

func (s *TransactionsService) Transactions(ctx context.Context, tenant *model.Tenant) ([]*model.Transaction, error) {
	if tenant.ID == GlobalTenantID {
		return nil, errors.New(util.ErrBadRequest, util.UserFacingTag)
	} else if !web.GetToken(ctx).IsPermittedPerTenant(tenant.ID, "Read transactions") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	session := s.getNeo4jSessionForTenant(ctx, neo4j.AccessModeRead, tenant.ID)
	defer session.Close(ctx)

	v, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		const getTxQuery = `// Get databases representing tenants
MATCH (src:Account)-[tx:Transaction]->(dst:Account)
RETURN src.id, src.displayName, dst.id, dst.displayName, tx.id, tx.date, tx.referenceID, tx.amount, tx.description`

		result, err := tx.Run(ctx, getTxQuery, nil)
		if err != nil {
			return nil, errors.New("failed to execute query: %w\n%s", err, getTxQuery)
		}

		records, err := result.Collect(ctx)
		if err != nil {
			return nil, errors.New("failed to collect records: %w", err)
		}

		transactions := make([]*model.Transaction, 0)
		for _, rec := range records {
			src := &model.Account{Tenant: tenant, ID: rec.Values[0].(string), DisplayName: rec.Values[1].(string)}
			dst := &model.Account{Tenant: tenant, ID: rec.Values[2].(string), DisplayName: rec.Values[3].(string)}
			transactions = append(transactions, &model.Transaction{
				ID:            rec.Values[4].(string),
				Date:          rec.Values[5].(time.Time),
				TargetAccount: dst,
				SourceAccount: src,
				ReferenceID:   rec.Values[6].(string),
				Amount:        model.MustParseMoney(rec.Values[7].(string)),
				Description:   rec.Values[8].(string),
			})
		}
		return transactions, nil
	})
	return v.([]*model.Transaction), err
}
