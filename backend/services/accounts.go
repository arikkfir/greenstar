package services

import (
	"context"
	"fmt"
	"github.com/arikkfir/greenstar/backend/model"
	"github.com/arikkfir/greenstar/backend/util"
	"github.com/arikkfir/greenstar/backend/web"
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

	createAccountQuery := `CREATE (account:Account {accountID: $accountID, displayName: $displayName, icon: $icon})`
	createAccountParams["accountID"] = id
	createAccountParams["displayName"] = account.DisplayName
	createAccountParams["icon"] = account.Icon

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

	createAccountQuery += "RETURN account.accountID, account.displayName"

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

	updateAccountQuery := `MATCH (account:Account {accountID: $accountID}) `
	updateAccountParams["accountID"] = accountID

	if account.DisplayName != nil {
		updateAccountQuery += "SET account.displayName = $displayName\n"
		updateAccountParams["displayName"] = *account.DisplayName
	}

	if account.Icon != nil {
		updateAccountQuery += "SET account.icon = $icon\n"
		updateAccountParams["icon"] = *account.Icon
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

	updateAccountQuery += "RETURN account.accountID, account.displayName"

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
	if v == nil {
		return nil, err
	} else {
		return v.(*model.Account), err
	}
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
		const deleteAccountQuery = `MATCH (account:Account {accountID: $accountID}) DELETE account`

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
	if v == nil {
		return "", err
	} else {
		return v.(string), err
	}
}

func (s *AccountsService) Accounts(ctx context.Context, tenant *model.Tenant) ([]*model.Account, error) {
	if tenant.ID == GlobalTenantID {
		return nil, errors.New(util.ErrBadRequest, util.UserFacingTag)
	} else if !web.GetToken(ctx).IsPermittedPerTenant(tenant.ID, "Read accounts") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	session := s.getNeo4jSessionForTenant(ctx, neo4j.AccessModeRead, tenant.ID)
	defer session.Close(ctx)

	v, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		const query = `MATCH (a:Account)  WHERE NOT exists ((a)-[:ChildOf]->(:Account))  RETURN a.accountID, a.displayName, a.icon`

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
				Icon:        rec.Values[2].(string),
			})
		}
		return accounts, nil
	})
	if v == nil {
		return nil, err
	} else {
		return v.([]*model.Account), err
	}
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
		const getAccountCypher = `MATCH (account:Account {accountID: $accountID}) RETURN account.accountID, account.displayName, account.icon`

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
			Icon:        record.Values[2].(string),
		}, nil
	})
	if err != nil {
		return nil, err
	} else {
		return v.(*model.Account), err
	}
}

func (s *AccountsService) Labels(ctx context.Context, obj *model.Account) ([]*model.KeyAndValue, error) {
	if !web.GetToken(ctx).IsPermittedPerTenant(obj.Tenant.ID, "Read accounts") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	session := s.getNeo4jSessionForTenant(ctx, neo4j.AccessModeRead, obj.Tenant.ID)
	defer session.Close(ctx)

	v, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		const getLabelsCypher = `MATCH (acc:Account {accountID: $accountID})-[r:HasLabel]->(l:Label) RETURN l.name, r.value`

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
	if v == nil {
		return nil, err
	} else {
		return v.([]*model.KeyAndValue), err
	}
}

func (s *AccountsService) ChildCount(ctx context.Context, obj *model.Account) (int, error) {
	if !web.GetToken(ctx).IsPermittedPerTenant(obj.Tenant.ID, "Read accounts") {
		return 0, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	session := s.getNeo4jSessionForTenant(ctx, neo4j.AccessModeRead, obj.Tenant.ID)
	defer session.Close(ctx)

	v, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		const getChildCountCypher = `MATCH (acc:Account)-[:ChildOf]->(parent:Account {accountID: $accountID}) RETURN count(acc)`

		result, err := tx.Run(ctx, getChildCountCypher, map[string]any{"accountID": obj.ID})
		if err != nil {
			return 0, errors.New("failed to execute query: %w\n%s", err, getChildCountCypher)
		}

		rec, err := result.Single(ctx)
		if err != nil {
			return 0, errors.New("failed to collect record: %w", err)
		}

		if len(rec.Values) != 1 {
			return 0, errors.New("incorrect record returned: %d values", len(rec.Values))
		} else if rec.Values[0] == nil {
			return 0, nil
		} else if cnt, ok := rec.Values[0].(int64); !ok {
			return 0, errors.New("incorrect record returned: %T(%+v)", rec.Values[0], rec.Values[0])
		} else {
			return int(cnt), nil
		}
	})
	if v == nil {
		return 0, err
	} else {
		return v.(int), err
	}
}

func (s *AccountsService) Children(ctx context.Context, obj *model.Account) ([]*model.Account, error) {
	if !web.GetToken(ctx).IsPermittedPerTenant(obj.Tenant.ID, "Read accounts") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	session := s.getNeo4jSessionForTenant(ctx, neo4j.AccessModeRead, obj.Tenant.ID)
	defer session.Close(ctx)

	v, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		const getChildrenCypher = `MATCH (acc:Account)-[:ChildOf]->(parent:Account {accountID: $accountID}) RETURN acc.accountID, acc.displayName, acc.icon`

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
				Tenant:      obj.Tenant,
				ID:          rec.Values[0].(string),
				DisplayName: rec.Values[1].(string),
				Icon:        rec.Values[2].(string),
			})
		}
		return accounts, nil
	})
	if v == nil {
		return nil, err
	} else {
		return v.([]*model.Account), err
	}
}

func (s *AccountsService) Parent(ctx context.Context, obj *model.Account) (*model.Account, error) {
	if !web.GetToken(ctx).IsPermittedPerTenant(obj.Tenant.ID, "Read accounts") {
		return nil, errors.New(util.ErrPermissionDenied, util.UserFacingTag)
	}

	session := s.getNeo4jSessionForTenant(ctx, neo4j.AccessModeRead, obj.Tenant.ID)
	defer session.Close(ctx)

	v, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		const getParentCypher = `MATCH (child:Account {accountID: $accountID})-[:ChildOf]->(parent:Account) RETURN parent.accountID, parent.displayName, parent.icon`

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
			Tenant:      obj.Tenant,
			ID:          records[0].Values[0].(string),
			DisplayName: records[0].Values[1].(string),
			Icon:        records[0].Values[2].(string),
		}, nil
	})
	if v == nil {
		return nil, err
	} else {
		return v.(*model.Account), err
	}
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
MATCH (origin:Account {accountID: $sourceAccountID})
MATCH (src:Account)-[tx:Transaction]->(dst:Account)
WHERE exists ( (origin)<-[:ChildOf*0..]-(src) )
RETURN src.accountID, src.displayName, src.icon, dst.accountID, dst.displayName, dst.icon, tx.txID, tx.date, tx.referenceID, tx.amount, tx.description`

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
			src := &model.Account{Tenant: obj.Tenant, ID: rec.Values[0].(string), DisplayName: rec.Values[1].(string), Icon: rec.Values[2].(string)}
			dst := &model.Account{Tenant: obj.Tenant, ID: rec.Values[3].(string), DisplayName: rec.Values[4].(string), Icon: rec.Values[5].(string)}
			transactions = append(transactions, &model.Transaction{
				ID:            rec.Values[6].(string),
				Date:          rec.Values[7].(time.Time),
				TargetAccount: dst,
				SourceAccount: src,
				ReferenceID:   rec.Values[8].(string),
				Amount:        model.MustParseMoney(rec.Values[9].(string)),
				Description:   rec.Values[10].(string),
			})
		}
		return transactions, nil
	})
	if v == nil {
		return nil, err
	} else {
		return v.([]*model.Transaction), err
	}
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
MATCH (target:Account {accountID: $targetAccountID})
MATCH (src:Account)-[tx:Transaction]->(dst:Account)
WHERE exists ( (target)<-[:ChildOf*0..]-(dst) )
RETURN src.accountID, src.displayName, src.icon, dst.accountID, dst.displayName, dst.icon, tx.txID, tx.date, tx.referenceID, tx.amount, tx.description`

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
			src := &model.Account{Tenant: obj.Tenant, ID: rec.Values[0].(string), DisplayName: rec.Values[1].(string), Icon: rec.Values[2].(string)}
			dst := &model.Account{Tenant: obj.Tenant, ID: rec.Values[3].(string), DisplayName: rec.Values[4].(string), Icon: rec.Values[5].(string)}
			transactions = append(transactions, &model.Transaction{
				ID:            rec.Values[6].(string),
				Date:          rec.Values[7].(time.Time),
				TargetAccount: dst,
				SourceAccount: src,
				ReferenceID:   rec.Values[8].(string),
				Amount:        model.MustParseMoney(rec.Values[9].(string)),
				Description:   rec.Values[10].(string),
			})
		}
		return transactions, nil
	})
	if v == nil {
		return nil, err
	} else {
		return v.([]*model.Transaction), err
	}
}
