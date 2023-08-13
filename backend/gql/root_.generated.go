// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gql

import (
	"bytes"
	"context"
	"errors"
	"sync/atomic"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/introspection"
	"github.com/arikkfir/greenstar/backend/model"
	gqlparser "github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

// NewExecutableSchema creates an ExecutableSchema from the ResolverRoot interface.
func NewExecutableSchema(cfg Config) graphql.ExecutableSchema {
	return &executableSchema{
		resolvers:  cfg.Resolvers,
		directives: cfg.Directives,
		complexity: cfg.Complexity,
	}
}

type Config struct {
	Resolvers  ResolverRoot
	Directives DirectiveRoot
	Complexity ComplexityRoot
}

type ResolverRoot interface {
	Account() AccountResolver
	Mutation() MutationResolver
	Query() QueryResolver
	Tenant() TenantResolver
}

type DirectiveRoot struct {
}

type ComplexityRoot struct {
	Account struct {
		ChildCount           func(childComplexity int) int
		Children             func(childComplexity int) int
		DisplayName          func(childComplexity int) int
		ID                   func(childComplexity int) int
		IncomingTransactions func(childComplexity int) int
		Labels               func(childComplexity int) int
		OutgoingTransactions func(childComplexity int) int
		Parent               func(childComplexity int) int
		Tenant               func(childComplexity int) int
	}

	KeyAndValue struct {
		Key   func(childComplexity int) int
		Value func(childComplexity int) int
	}

	Mutation struct {
		CreateAccount         func(childComplexity int, tenantID string, accountID *string, account model.AccountChanges) int
		CreateTenant          func(childComplexity int, tenantID *string, tenant model.TenantChanges) int
		CreateTransaction     func(childComplexity int, tenantID string, transaction model.TransactionChanges) int
		CreateTransactions    func(childComplexity int, tenantID string, transactions []*model.TransactionChanges) int
		DeleteAccount         func(childComplexity int, tenantID string, accountID string) int
		DeleteTenant          func(childComplexity int, tenantID string) int
		ScrapeIsraelBankYahav func(childComplexity int, tenantID string, username string, id string, password string) int
		UpdateAccount         func(childComplexity int, tenantID string, accountID string, account model.AccountChanges) int
		UpdateOperation       func(childComplexity int, id string, op model.OperationChanges) int
		UpdateTenant          func(childComplexity int, tenantID string, tenant model.TenantChanges) int
	}

	Operation struct {
		CreatedAt   func(childComplexity int) int
		Description func(childComplexity int) int
		ID          func(childComplexity int) int
		Name        func(childComplexity int) int
		Result      func(childComplexity int) int
		Status      func(childComplexity int) int
		UpdatedAt   func(childComplexity int) int
	}

	Query struct {
		Operation func(childComplexity int, id string) int
		Tenant    func(childComplexity int, id string) int
		Tenants   func(childComplexity int) int
	}

	Tenant struct {
		Account      func(childComplexity int, id string) int
		Accounts     func(childComplexity int) int
		DisplayName  func(childComplexity int) int
		ID           func(childComplexity int) int
		Transactions func(childComplexity int) int
	}

	Transaction struct {
		Amount        func(childComplexity int) int
		Date          func(childComplexity int) int
		Description   func(childComplexity int) int
		ID            func(childComplexity int) int
		ReferenceID   func(childComplexity int) int
		SourceAccount func(childComplexity int) int
		TargetAccount func(childComplexity int) int
	}
}

type executableSchema struct {
	resolvers  ResolverRoot
	directives DirectiveRoot
	complexity ComplexityRoot
}

func (e *executableSchema) Schema() *ast.Schema {
	return parsedSchema
}

func (e *executableSchema) Complexity(typeName, field string, childComplexity int, rawArgs map[string]interface{}) (int, bool) {
	ec := executionContext{nil, e, 0, 0, nil}
	_ = ec
	switch typeName + "." + field {

	case "Account.childCount":
		if e.complexity.Account.ChildCount == nil {
			break
		}

		return e.complexity.Account.ChildCount(childComplexity), true

	case "Account.children":
		if e.complexity.Account.Children == nil {
			break
		}

		return e.complexity.Account.Children(childComplexity), true

	case "Account.displayName":
		if e.complexity.Account.DisplayName == nil {
			break
		}

		return e.complexity.Account.DisplayName(childComplexity), true

	case "Account.id":
		if e.complexity.Account.ID == nil {
			break
		}

		return e.complexity.Account.ID(childComplexity), true

	case "Account.incomingTransactions":
		if e.complexity.Account.IncomingTransactions == nil {
			break
		}

		return e.complexity.Account.IncomingTransactions(childComplexity), true

	case "Account.labels":
		if e.complexity.Account.Labels == nil {
			break
		}

		return e.complexity.Account.Labels(childComplexity), true

	case "Account.outgoingTransactions":
		if e.complexity.Account.OutgoingTransactions == nil {
			break
		}

		return e.complexity.Account.OutgoingTransactions(childComplexity), true

	case "Account.parent":
		if e.complexity.Account.Parent == nil {
			break
		}

		return e.complexity.Account.Parent(childComplexity), true

	case "Account.tenant":
		if e.complexity.Account.Tenant == nil {
			break
		}

		return e.complexity.Account.Tenant(childComplexity), true

	case "KeyAndValue.key":
		if e.complexity.KeyAndValue.Key == nil {
			break
		}

		return e.complexity.KeyAndValue.Key(childComplexity), true

	case "KeyAndValue.value":
		if e.complexity.KeyAndValue.Value == nil {
			break
		}

		return e.complexity.KeyAndValue.Value(childComplexity), true

	case "Mutation.createAccount":
		if e.complexity.Mutation.CreateAccount == nil {
			break
		}

		args, err := ec.field_Mutation_createAccount_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Mutation.CreateAccount(childComplexity, args["tenantID"].(string), args["accountID"].(*string), args["account"].(model.AccountChanges)), true

	case "Mutation.createTenant":
		if e.complexity.Mutation.CreateTenant == nil {
			break
		}

		args, err := ec.field_Mutation_createTenant_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Mutation.CreateTenant(childComplexity, args["tenantID"].(*string), args["tenant"].(model.TenantChanges)), true

	case "Mutation.createTransaction":
		if e.complexity.Mutation.CreateTransaction == nil {
			break
		}

		args, err := ec.field_Mutation_createTransaction_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Mutation.CreateTransaction(childComplexity, args["tenantID"].(string), args["transaction"].(model.TransactionChanges)), true

	case "Mutation.createTransactions":
		if e.complexity.Mutation.CreateTransactions == nil {
			break
		}

		args, err := ec.field_Mutation_createTransactions_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Mutation.CreateTransactions(childComplexity, args["tenantID"].(string), args["transactions"].([]*model.TransactionChanges)), true

	case "Mutation.deleteAccount":
		if e.complexity.Mutation.DeleteAccount == nil {
			break
		}

		args, err := ec.field_Mutation_deleteAccount_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Mutation.DeleteAccount(childComplexity, args["tenantID"].(string), args["accountID"].(string)), true

	case "Mutation.deleteTenant":
		if e.complexity.Mutation.DeleteTenant == nil {
			break
		}

		args, err := ec.field_Mutation_deleteTenant_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Mutation.DeleteTenant(childComplexity, args["tenantID"].(string)), true

	case "Mutation.scrapeIsraelBankYahav":
		if e.complexity.Mutation.ScrapeIsraelBankYahav == nil {
			break
		}

		args, err := ec.field_Mutation_scrapeIsraelBankYahav_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Mutation.ScrapeIsraelBankYahav(childComplexity, args["tenantID"].(string), args["username"].(string), args["id"].(string), args["password"].(string)), true

	case "Mutation.updateAccount":
		if e.complexity.Mutation.UpdateAccount == nil {
			break
		}

		args, err := ec.field_Mutation_updateAccount_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Mutation.UpdateAccount(childComplexity, args["tenantID"].(string), args["accountID"].(string), args["account"].(model.AccountChanges)), true

	case "Mutation.updateOperation":
		if e.complexity.Mutation.UpdateOperation == nil {
			break
		}

		args, err := ec.field_Mutation_updateOperation_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Mutation.UpdateOperation(childComplexity, args["id"].(string), args["op"].(model.OperationChanges)), true

	case "Mutation.updateTenant":
		if e.complexity.Mutation.UpdateTenant == nil {
			break
		}

		args, err := ec.field_Mutation_updateTenant_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Mutation.UpdateTenant(childComplexity, args["tenantID"].(string), args["tenant"].(model.TenantChanges)), true

	case "Operation.createdAt":
		if e.complexity.Operation.CreatedAt == nil {
			break
		}

		return e.complexity.Operation.CreatedAt(childComplexity), true

	case "Operation.description":
		if e.complexity.Operation.Description == nil {
			break
		}

		return e.complexity.Operation.Description(childComplexity), true

	case "Operation.id":
		if e.complexity.Operation.ID == nil {
			break
		}

		return e.complexity.Operation.ID(childComplexity), true

	case "Operation.name":
		if e.complexity.Operation.Name == nil {
			break
		}

		return e.complexity.Operation.Name(childComplexity), true

	case "Operation.result":
		if e.complexity.Operation.Result == nil {
			break
		}

		return e.complexity.Operation.Result(childComplexity), true

	case "Operation.status":
		if e.complexity.Operation.Status == nil {
			break
		}

		return e.complexity.Operation.Status(childComplexity), true

	case "Operation.updatedAt":
		if e.complexity.Operation.UpdatedAt == nil {
			break
		}

		return e.complexity.Operation.UpdatedAt(childComplexity), true

	case "Query.operation":
		if e.complexity.Query.Operation == nil {
			break
		}

		args, err := ec.field_Query_operation_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Query.Operation(childComplexity, args["id"].(string)), true

	case "Query.tenant":
		if e.complexity.Query.Tenant == nil {
			break
		}

		args, err := ec.field_Query_tenant_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Query.Tenant(childComplexity, args["id"].(string)), true

	case "Query.tenants":
		if e.complexity.Query.Tenants == nil {
			break
		}

		return e.complexity.Query.Tenants(childComplexity), true

	case "Tenant.account":
		if e.complexity.Tenant.Account == nil {
			break
		}

		args, err := ec.field_Tenant_account_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Tenant.Account(childComplexity, args["id"].(string)), true

	case "Tenant.accounts":
		if e.complexity.Tenant.Accounts == nil {
			break
		}

		return e.complexity.Tenant.Accounts(childComplexity), true

	case "Tenant.displayName":
		if e.complexity.Tenant.DisplayName == nil {
			break
		}

		return e.complexity.Tenant.DisplayName(childComplexity), true

	case "Tenant.id":
		if e.complexity.Tenant.ID == nil {
			break
		}

		return e.complexity.Tenant.ID(childComplexity), true

	case "Tenant.transactions":
		if e.complexity.Tenant.Transactions == nil {
			break
		}

		return e.complexity.Tenant.Transactions(childComplexity), true

	case "Transaction.amount":
		if e.complexity.Transaction.Amount == nil {
			break
		}

		return e.complexity.Transaction.Amount(childComplexity), true

	case "Transaction.Date":
		if e.complexity.Transaction.Date == nil {
			break
		}

		return e.complexity.Transaction.Date(childComplexity), true

	case "Transaction.description":
		if e.complexity.Transaction.Description == nil {
			break
		}

		return e.complexity.Transaction.Description(childComplexity), true

	case "Transaction.id":
		if e.complexity.Transaction.ID == nil {
			break
		}

		return e.complexity.Transaction.ID(childComplexity), true

	case "Transaction.referenceID":
		if e.complexity.Transaction.ReferenceID == nil {
			break
		}

		return e.complexity.Transaction.ReferenceID(childComplexity), true

	case "Transaction.sourceAccount":
		if e.complexity.Transaction.SourceAccount == nil {
			break
		}

		return e.complexity.Transaction.SourceAccount(childComplexity), true

	case "Transaction.targetAccount":
		if e.complexity.Transaction.TargetAccount == nil {
			break
		}

		return e.complexity.Transaction.TargetAccount(childComplexity), true

	}
	return 0, false
}

func (e *executableSchema) Exec(ctx context.Context) graphql.ResponseHandler {
	rc := graphql.GetOperationContext(ctx)
	ec := executionContext{rc, e, 0, 0, make(chan graphql.DeferredResult)}
	inputUnmarshalMap := graphql.BuildUnmarshalerMap(
		ec.unmarshalInputAccountChanges,
		ec.unmarshalInputKeyAndValueInput,
		ec.unmarshalInputOperationChanges,
		ec.unmarshalInputTenantChanges,
		ec.unmarshalInputTransactionChanges,
	)
	first := true

	switch rc.Operation.Operation {
	case ast.Query:
		return func(ctx context.Context) *graphql.Response {
			var response graphql.Response
			var data graphql.Marshaler
			if first {
				first = false
				ctx = graphql.WithUnmarshalerMap(ctx, inputUnmarshalMap)
				data = ec._Query(ctx, rc.Operation.SelectionSet)
			} else {
				if atomic.LoadInt32(&ec.pendingDeferred) > 0 {
					result := <-ec.deferredResults
					atomic.AddInt32(&ec.pendingDeferred, -1)
					data = result.Result
					response.Path = result.Path
					response.Label = result.Label
					response.Errors = result.Errors
				} else {
					return nil
				}
			}
			var buf bytes.Buffer
			data.MarshalGQL(&buf)
			response.Data = buf.Bytes()
			if atomic.LoadInt32(&ec.deferred) > 0 {
				hasNext := atomic.LoadInt32(&ec.pendingDeferred) > 0
				response.HasNext = &hasNext
			}

			return &response
		}
	case ast.Mutation:
		return func(ctx context.Context) *graphql.Response {
			if !first {
				return nil
			}
			first = false
			ctx = graphql.WithUnmarshalerMap(ctx, inputUnmarshalMap)
			data := ec._Mutation(ctx, rc.Operation.SelectionSet)
			var buf bytes.Buffer
			data.MarshalGQL(&buf)

			return &graphql.Response{
				Data: buf.Bytes(),
			}
		}

	default:
		return graphql.OneShot(graphql.ErrorResponse(ctx, "unsupported GraphQL operation"))
	}
}

type executionContext struct {
	*graphql.OperationContext
	*executableSchema
	deferred        int32
	pendingDeferred int32
	deferredResults chan graphql.DeferredResult
}

func (ec *executionContext) processDeferredGroup(dg graphql.DeferredGroup) {
	atomic.AddInt32(&ec.pendingDeferred, 1)
	go func() {
		ctx := graphql.WithFreshResponseContext(dg.Context)
		dg.FieldSet.Dispatch(ctx)
		ds := graphql.DeferredResult{
			Path:   dg.Path,
			Label:  dg.Label,
			Result: dg.FieldSet,
			Errors: graphql.GetErrors(ctx),
		}
		// null fields should bubble up
		if dg.FieldSet.Invalids > 0 {
			ds.Result = graphql.Null
		}
		ec.deferredResults <- ds
	}()
}

func (ec *executionContext) introspectSchema() (*introspection.Schema, error) {
	if ec.DisableIntrospection {
		return nil, errors.New("introspection disabled")
	}
	return introspection.WrapSchema(parsedSchema), nil
}

func (ec *executionContext) introspectType(name string) (*introspection.Type, error) {
	if ec.DisableIntrospection {
		return nil, errors.New("introspection disabled")
	}
	return introspection.WrapTypeFromDef(parsedSchema, parsedSchema.Types[name]), nil
}

var sources = []*ast.Source{
	{Name: "../../schema/accounts.graphql", Input: `type Account {
    tenant: Tenant!
    id: ID!
    displayName: String!
    labels: [KeyAndValue!]!
    childCount: Int!
    children: [Account!]!
    parent: Account

    outgoingTransactions: [Transaction!]!
    incomingTransactions: [Transaction!]!
}

input AccountChanges {
    displayName: String
    labels: [KeyAndValueInput!]
    parentID: ID
}
`, BuiltIn: false},
	{Name: "../../schema/common.graphql", Input: `type KeyAndValue {
    key: String!
    value: String!
}

input KeyAndValueInput {
    key: String!
    value: String!
}
`, BuiltIn: false},
	{Name: "../../schema/mutations.graphql", Input: `type Mutation {
    createTenant(tenantID: String, tenant: TenantChanges!): Tenant!
    updateTenant(tenantID: String!, tenant: TenantChanges!): Tenant!
    deleteTenant(tenantID: String!): ID!

    createAccount(tenantID: ID!, accountID: ID, account: AccountChanges!): Account!
    updateAccount(tenantID: ID!, accountID: ID!, account: AccountChanges!): Account!
    deleteAccount(tenantID: ID!, accountID: ID!): ID!

    createTransaction(tenantID: ID!, transaction: TransactionChanges!): Transaction!
    createTransactions(tenantID: ID!, transactions: [TransactionChanges!]!): Int!

    scrapeIsraelBankYahav(tenantID: ID!, username: String!, id: String!, password: String!): String!

    updateOperation(id: ID!, op: OperationChanges!): Operation!
}
`, BuiltIn: false},
	{Name: "../../schema/operations.graphql", Input: `scalar DateTime

enum OperationStatus {
    PENDING
    ACCEPTED
    REJECTED
    STARTED
    COMPLETED
}

enum OperationResult {
    SUCCEEDED
    FAILED
}

type Operation {
    id: ID!
    name: String!
    description: String
    status: OperationStatus!
    result: OperationResult!
    createdAt: DateTime!
    updatedAt: DateTime!
}

input OperationChanges {
    name: String!
    description: String
    status: OperationStatus!
    result: OperationResult!
}
`, BuiltIn: false},
	{Name: "../../schema/query.graphql", Input: `type Query {
    tenants: [Tenant!]!
    tenant(id: ID!): Tenant
    operation(id: ID!): Operation
}
`, BuiltIn: false},
	{Name: "../../schema/tenants.graphql", Input: `type Tenant {
    id: ID!
    displayName: String!

    accounts: [Account!]!
    account(id: ID!): Account
    transactions: [Transaction!]!
}

input TenantChanges {
    displayName: String!
}
`, BuiltIn: false},
	{Name: "../../schema/transactions.graphql", Input: `scalar Time
scalar Money

type Transaction {
    id: ID!
    Date: Time!
    targetAccount: Account!
    sourceAccount: Account!
    referenceID: String!
    amount: Money!
    description: String!
}

input TransactionChanges {
    Date: Time!
    targetAccountID: ID!
    sourceAccountID: ID!
    referenceID: String!
    amount: Money!
    description: String!
}
`, BuiltIn: false},
}
var parsedSchema = gqlparser.MustLoadSchema(sources...)
