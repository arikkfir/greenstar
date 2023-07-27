package services

import (
	"context"
	"fmt"
	"github.com/arik-kfir/greenstar/backend/model"
	"github.com/arik-kfir/greenstar/backend/util"
	"github.com/arik-kfir/greenstar/backend/web"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/secureworks/errors"
	"strconv"
	"time"
)

type AccountsService struct {
	Service
}

func (s *AccountsService) CreateAccount(ctx context.Context, tenantID string, accountID *string, account model.AccountChanges) (*model.Account, error) {
	if tenantID == GlobalTenantID {
		return nil, errors.New(util.ErrBadRequest, util.UserFacingTag)
	} else if !web.GetToken(ctx).IsPermittedPerTenant(tenantID, "Manage accounts") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	session := s.getNeo4jSessionForTenant(ctx, neo4j.AccessModeWrite, tenantID)
	defer session.Close(ctx)

	var id string
	if accountID == nil {
		id = util.RandomHash(7)
	} else if len(*accountID) != 7 {
		return nil, errors.New("account ID must be 7 letters long")
	} else {
		id = *accountID
	}

	createAccountParams := map[string]interface{}{}

	createAccountQuery := `// Create account
CREATE (account:Account {accountID: $accountID, displayName: $displayName})`
	createAccountParams["accountID"] = id
	createAccountParams["displayName"] = account.DisplayName

	if account.ParentID != nil {
		createAccountQuery += "WITH account\n"
		createAccountQuery += "MATCH (parent:Account {accountID: $parentAccountID}) \n"
		createAccountQuery += "CREATE (account)-[:ChildOf]->(parent)\n"
		createAccountParams["parentAccountID"] = account.ParentID
	}

	for i, kv := range account.Labels {
		createAccountQuery += fmt.Sprintf("MERGE (label:Label {name: $labelName%d})\n", i)
		createAccountQuery += fmt.Sprintf("MERGE (account)-[r:HasLabel {value: $labelValue%d}]->(label)\n", i)
		createAccountParams["labelName"+strconv.Itoa(i)] = kv.Key
		createAccountParams["labelValue"+strconv.Itoa(i)] = kv.Value
	}

	createAccountQuery += "RETURN account.id, account.displayName"

	v, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, createAccountQuery, createAccountParams)
		if err != nil {
			return nil, errors.New("failed to execute query: %w\n%s", err, createAccountQuery)
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, errors.New("failed to collect single record: %w", err)
		}

		return &model.Account{
			Tenant:      &model.Tenant{ID: tenantID},
			ID:          record.Values[0].(string),
			DisplayName: record.Values[1].(string),
		}, nil
	})
	if v == nil {
		return nil, err
	} else {
		return v.(*model.Account), err
	}
}

func (s *AccountsService) UpdateAccount(ctx context.Context, tenantID, accountID string, account model.AccountChanges) (*model.Account, error) {
	if tenantID == GlobalTenantID {
		return nil, errors.New(util.ErrBadRequest, util.UserFacingTag)
	} else if !web.GetToken(ctx).IsPermittedPerTenant(tenantID, "Manage accounts") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	session := s.getNeo4jSessionForTenant(ctx, neo4j.AccessModeWrite, tenantID)
	defer session.Close(ctx)

	updateAccountParams := map[string]interface{}{}

	updateAccountQuery := `// Update account
MATCH (account:Account {accountID: $accountID}) `
	updateAccountParams["accountID"] = accountID

	if account.DisplayName != nil {
		updateAccountQuery += "SET account.displayName = $displayName\n"
		updateAccountParams["displayName"] = *account.DisplayName
	}

	if account.ParentID != nil {
		updateAccountQuery += "WITH account\n"
		updateAccountQuery += "OPTIONAL MATCH (account)-[oldChildOfRel:ChildOf]->(oldParent:Account)\n"
		updateAccountQuery += "WHERE oldParent.accountID <> $parentAccountID\n"
		updateAccountQuery += "DELETE oldChildOfRel\n"
		updateAccountQuery += "WITH account\n"
		updateAccountQuery += "MERGE (parent:Account {accountID: $parentAccountID})\n"
		updateAccountQuery += "MERGE (account)-[:ChildOf]->(parent)\n"
		updateAccountParams["parentAccountID"] = *account.ParentID
	}

	for i, kv := range account.Labels {
		updateAccountQuery += "OPTIONAL MATCH (account)-[oldLabelRel:HasLabel]->(:Label)\n"
		updateAccountQuery += "DELETE oldLabelRel\n"
		updateAccountQuery += "WITH account\n"
		updateAccountQuery += fmt.Sprintf("MERGE (label:Label {name: $labelName%d})\n", i)
		updateAccountQuery += fmt.Sprintf("MERGE (account)-[r:HasLabel {value: $labelValue%d}]->(label)\n", i)
		updateAccountParams["labelName"+strconv.Itoa(i)] = kv.Key
		updateAccountParams["labelValue"+strconv.Itoa(i)] = kv.Value
	}

	updateAccountQuery += "RETURN account.id, account.displayName"

	v, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, updateAccountQuery, updateAccountParams)
		if err != nil {
			return nil, errors.New("failed to execute query: %w\n%s", err, updateAccountQuery)
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, errors.New("failed to collect single record: %w", err)
		}

		return &model.Account{
			Tenant:      &model.Tenant{ID: tenantID},
			ID:          record.Values[0].(string),
			DisplayName: record.Values[1].(string),
		}, nil
	})
	return v.(*model.Account), err
}

func (s *AccountsService) DeleteAccount(ctx context.Context, tenantID, accountID string) (string, error) {
	if tenantID == GlobalTenantID {
		return "", errors.New(util.ErrBadRequest, util.UserFacingTag)
	} else if !web.GetToken(ctx).IsPermittedPerTenant(tenantID, "Manage accounts") {
		return "", errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	session := s.getNeo4jSessionForTenant(ctx, neo4j.AccessModeWrite, tenantID)
	defer session.Close(ctx)

	v, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		const deleteAccountQuery = `// Delete account
MATCH (account:Account {accountID: $accountID})
DELETE account`

		result, err := tx.Run(ctx, deleteAccountQuery, map[string]any{"accountID": accountID})
		if err != nil {
			return nil, errors.New("failed to execute query: %w\n%s", err, deleteAccountQuery)
		}

		summary, err := result.Consume(ctx)
		if err != nil {
			return "", errors.New("failed to consume query summary: %w", err)
		}

		if summary.Counters().NodesDeleted() != 1 {
			return accountID, errors.New("no records deleted")
		}
		return accountID, nil
	})
	return v.(string), err
}

func (s *AccountsService) Accounts(ctx context.Context, tenant *model.Tenant, rootsOnly *bool) ([]*model.Account, error) {
	if tenant.ID == GlobalTenantID {
		return nil, errors.New(util.ErrBadRequest, util.UserFacingTag)
	} else if !web.GetToken(ctx).IsPermittedPerTenant(tenant.ID, "Read accounts") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	session := s.getNeo4jSessionForTenant(ctx, neo4j.AccessModeRead, tenant.ID)
	defer session.Close(ctx)

	v, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		const getRootAccountsCypher = `// Get root accounts
MATCH (a:Account) 
WHERE NOT exists ((a)-[:ChildOf]->(:Account)) 
RETURN a.id, a.displayName`
		const getAllAccountsCypher = `// Get all accounts
MATCH (a:Account) 
RETURN a.id, a.displayName`

		query := getAllAccountsCypher
		if rootsOnly != nil && *rootsOnly == true {
			query = getRootAccountsCypher
		}
		result, err := tx.Run(ctx, query, nil)
		if err != nil {
			return nil, errors.New("failed to execute query: %w\n%s", err, query)
		}

		records, err := result.Collect(ctx)
		if err != nil {
			return nil, errors.New("failed to collect records: %w", err)
		}

		accounts := make([]*model.Account, 0)
		for _, rec := range records {
			accounts = append(accounts, &model.Account{
				Tenant:      tenant,
				ID:          rec.Values[0].(string),
				DisplayName: rec.Values[1].(string),
			})
		}
		return accounts, nil
	})
	return v.([]*model.Account), err
}

func (s *AccountsService) Account(ctx context.Context, tenant *model.Tenant, accountID string) (*model.Account, error) {
	if tenant.ID == GlobalTenantID {
		return nil, errors.New(util.ErrBadRequest, util.UserFacingTag)
	} else if !web.GetToken(ctx).IsPermittedPerTenant(tenant.ID, "Read accounts") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	session := s.getNeo4jSessionForTenant(ctx, neo4j.AccessModeRead, tenant.ID)
	defer session.Close(ctx)

	v, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		const getAccountCypher = `// Get account by ID
MATCH (account:Account {id: $accountID})
RETURN account.id, account.displayName`

		result, err := tx.Run(ctx, getAccountCypher, map[string]any{"accountID": accountID})
		if err != nil {
			return nil, errors.New("failed to execute query: %w\n%s", err, getAccountCypher)
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, errors.New("failed to collect single record: %w", err)
		}

		return &model.Account{
			Tenant:      tenant,
			ID:          record.Values[0].(string),
			DisplayName: record.Values[1].(string),
		}, nil
	})
	return v.(*model.Account), err
}

func (s *AccountsService) Labels(ctx context.Context, obj *model.Account) ([]*model.KeyAndValue, error) {
	if !web.GetToken(ctx).IsPermittedPerTenant(obj.Tenant.ID, "Read accounts") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	session := s.getNeo4jSessionForTenant(ctx, neo4j.AccessModeRead, obj.Tenant.ID)
	defer session.Close(ctx)

	v, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		const getLabelsCypher = `// Get account labels
MATCH (acc:Account {accountID: $accountID})-[r:HasLabel]->(l:Label)
RETURN l.name, r.value`

		result, err := tx.Run(ctx, getLabelsCypher, map[string]any{"accountID": obj.ID})
		if err != nil {
			return nil, errors.New("failed to execute query: %w\n%s", err, getLabelsCypher)
		}

		records, err := result.Collect(ctx)
		if err != nil {
			return nil, errors.New("failed to collect records: %w", err)
		}

		var labels []*model.KeyAndValue
		for _, record := range records {
			labels = append(labels, &model.KeyAndValue{
				Key:   record.Values[0].(string),
				Value: record.Values[1].(string),
			})
		}
		return labels, nil
	})
	return v.([]*model.KeyAndValue), err
}

func (s *AccountsService) Children(ctx context.Context, obj *model.Account) ([]*model.Account, error) {
	if !web.GetToken(ctx).IsPermittedPerTenant(obj.Tenant.ID, "Read accounts") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	session := s.getNeo4jSessionForTenant(ctx, neo4j.AccessModeRead, obj.Tenant.ID)
	defer session.Close(ctx)

	v, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		const getChildrenCypher = `// Get account children
MATCH (acc:Account)-[:ChildOf]->(parent:Account {accountID: $accountID})
RETURN acc.id, acc.displayName`

		result, err := tx.Run(ctx, getChildrenCypher, map[string]any{"accountID": obj.ID})
		if err != nil {
			return nil, errors.New("failed to execute query: %w\n%s", err, getChildrenCypher)
		}

		records, err := result.Collect(ctx)
		if err != nil {
			return nil, errors.New("failed to collect records: %w", err)
		}

		accounts := make([]*model.Account, 0)
		for _, rec := range records {
			accounts = append(accounts, &model.Account{
				ID:          rec.Values[0].(string),
				DisplayName: rec.Values[1].(string),
			})
		}
		return accounts, nil
	})
	return v.([]*model.Account), err
}

func (s *AccountsService) Parent(ctx context.Context, obj *model.Account) (*model.Account, error) {
	if !web.GetToken(ctx).IsPermittedPerTenant(obj.Tenant.ID, "Read accounts") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	session := s.getNeo4jSessionForTenant(ctx, neo4j.AccessModeRead, obj.Tenant.ID)
	defer session.Close(ctx)

	v, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		const getParentCypher = `// Get account parent
MATCH (child:Account {accountID: $accountID})-[:ChildOf]->(parent:Account) 
RETURN parent.id, parent.displayName`

		result, err := tx.Run(ctx, getParentCypher, map[string]any{"accountID": obj.ID})
		if err != nil {
			return nil, errors.New("failed to execute query: %w\n%s", err, getParentCypher)
		}

		records, err := result.Collect(ctx)
		if err != nil {
			return nil, errors.New("failed to collect records: %w", err)
		}

		if len(records) == 0 {
			return nil, nil
		}

		if len(records) > 1 {
			return nil, errors.New("too many records matched")
		}

		return &model.Account{
			ID:          records[0].Values[0].(string),
			DisplayName: records[0].Values[1].(string),
		}, nil
	})
	return v.(*model.Account), err
}

func (s *AccountsService) OutgoingTransactions(ctx context.Context, obj *model.Account) ([]*model.Transaction, error) {
	if !web.GetToken(ctx).IsPermittedPerTenant(obj.Tenant.ID, "Read accounts") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	} else if !web.GetToken(ctx).IsPermittedPerTenant(obj.Tenant.ID, "Read transactions") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	session := s.getNeo4jSessionForTenant(ctx, neo4j.AccessModeRead, obj.Tenant.ID)
	defer session.Close(ctx)

	v, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		const getTxQuery = `// Get databases representing tenants
MATCH (origin:Account {id: $sourceAccountID})
MATCH (src:Account)-[tx:Transaction]->(dst:Account)
WHERE exists ( (origin)-[:ChildOf*0..]->(src) )
RETURN src.id, src.displayName, dst.id, dst.displayName, tx.id, tx.date, tx.referenceID, tx.amount, tx.description`

		result, err := tx.Run(ctx, getTxQuery, map[string]any{"sourceAccountID": obj.ID})
		if err != nil {
			return nil, errors.New("failed to execute query: %w\n%s", err, getTxQuery)
		}

		records, err := result.Collect(ctx)
		if err != nil {
			return nil, errors.New("failed to collect records: %w", err)
		}

		transactions := make([]*model.Transaction, 0)
		for _, rec := range records {
			src := &model.Account{Tenant: obj.Tenant, ID: rec.Values[0].(string), DisplayName: rec.Values[1].(string)}
			dst := &model.Account{Tenant: obj.Tenant, ID: rec.Values[2].(string), DisplayName: rec.Values[3].(string)}
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

func (s *AccountsService) IncomingTransactions(ctx context.Context, obj *model.Account) ([]*model.Transaction, error) {
	if !web.GetToken(ctx).IsPermittedPerTenant(obj.Tenant.ID, "Read accounts") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	} else if !web.GetToken(ctx).IsPermittedPerTenant(obj.Tenant.ID, "Read transactions") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	session := s.getNeo4jSessionForTenant(ctx, neo4j.AccessModeRead, obj.Tenant.ID)
	defer session.Close(ctx)

	v, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		const getTxQuery = `// Get databases representing tenants
MATCH (target:Account {id: $targetAccountID})
MATCH (src:Account)-[tx:Transaction]->(dst:Account)
WHERE exists ( (target)-[:ChildOf*0..]->(dst) )
RETURN src.id, src.displayName, dst.id, dst.displayName, tx.id, tx.date, tx.referenceID, tx.amount, tx.description`

		result, err := tx.Run(ctx, getTxQuery, map[string]any{"targetAccountID": obj.ID})
		if err != nil {
			return nil, errors.New("failed to execute query: %w\n%s", err, getTxQuery)
		}

		records, err := result.Collect(ctx)
		if err != nil {
			return nil, errors.New("failed to collect records: %w", err)
		}

		transactions := make([]*model.Transaction, 0)
		for _, rec := range records {
			src := &model.Account{Tenant: obj.Tenant, ID: rec.Values[0].(string), DisplayName: rec.Values[1].(string)}
			dst := &model.Account{Tenant: obj.Tenant, ID: rec.Values[2].(string), DisplayName: rec.Values[3].(string)}
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
