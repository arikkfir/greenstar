// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gql

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/99designs/gqlgen/graphql"
	"github.com/arikkfir/greenstar/backend/model"
	"github.com/vektah/gqlparser/v2/ast"
)

// region    ************************** generated!.gotpl **************************

type AccountResolver interface {
	Labels(ctx context.Context, obj *model.Account) ([]*model.KeyAndValue, error)
	ChildCount(ctx context.Context, obj *model.Account) (int, error)
	Children(ctx context.Context, obj *model.Account) ([]*model.Account, error)
	Parent(ctx context.Context, obj *model.Account) (*model.Account, error)
	OutgoingTransactions(ctx context.Context, obj *model.Account) ([]*model.Transaction, error)
	IncomingTransactions(ctx context.Context, obj *model.Account) ([]*model.Transaction, error)
}

// endregion ************************** generated!.gotpl **************************

// region    ***************************** args.gotpl *****************************

// endregion ***************************** args.gotpl *****************************

// region    ************************** directives.gotpl **************************

// endregion ************************** directives.gotpl **************************

// region    **************************** field.gotpl *****************************

func (ec *executionContext) _Account_tenant(ctx context.Context, field graphql.CollectedField, obj *model.Account) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Account_tenant(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Tenant, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(*model.Tenant)
	fc.Result = res
	return ec.marshalNTenant2ᚖgithubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐTenant(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Account_tenant(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Account",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "id":
				return ec.fieldContext_Tenant_id(ctx, field)
			case "displayName":
				return ec.fieldContext_Tenant_displayName(ctx, field)
			case "accounts":
				return ec.fieldContext_Tenant_accounts(ctx, field)
			case "account":
				return ec.fieldContext_Tenant_account(ctx, field)
			case "transactions":
				return ec.fieldContext_Tenant_transactions(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type Tenant", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) _Account_id(ctx context.Context, field graphql.CollectedField, obj *model.Account) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Account_id(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.ID, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(string)
	fc.Result = res
	return ec.marshalNID2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Account_id(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Account",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type ID does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Account_displayName(ctx context.Context, field graphql.CollectedField, obj *model.Account) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Account_displayName(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.DisplayName, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(string)
	fc.Result = res
	return ec.marshalNString2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Account_displayName(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Account",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Account_icon(ctx context.Context, field graphql.CollectedField, obj *model.Account) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Account_icon(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Icon, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(string)
	fc.Result = res
	return ec.marshalNString2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Account_icon(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Account",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Account_labels(ctx context.Context, field graphql.CollectedField, obj *model.Account) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Account_labels(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Account().Labels(rctx, obj)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.([]*model.KeyAndValue)
	fc.Result = res
	return ec.marshalNKeyAndValue2ᚕᚖgithubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐKeyAndValueᚄ(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Account_labels(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Account",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "key":
				return ec.fieldContext_KeyAndValue_key(ctx, field)
			case "value":
				return ec.fieldContext_KeyAndValue_value(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type KeyAndValue", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) _Account_childCount(ctx context.Context, field graphql.CollectedField, obj *model.Account) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Account_childCount(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Account().ChildCount(rctx, obj)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(int)
	fc.Result = res
	return ec.marshalNInt2int(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Account_childCount(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Account",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type Int does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Account_children(ctx context.Context, field graphql.CollectedField, obj *model.Account) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Account_children(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Account().Children(rctx, obj)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.([]*model.Account)
	fc.Result = res
	return ec.marshalNAccount2ᚕᚖgithubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐAccountᚄ(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Account_children(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Account",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "tenant":
				return ec.fieldContext_Account_tenant(ctx, field)
			case "id":
				return ec.fieldContext_Account_id(ctx, field)
			case "displayName":
				return ec.fieldContext_Account_displayName(ctx, field)
			case "icon":
				return ec.fieldContext_Account_icon(ctx, field)
			case "labels":
				return ec.fieldContext_Account_labels(ctx, field)
			case "childCount":
				return ec.fieldContext_Account_childCount(ctx, field)
			case "children":
				return ec.fieldContext_Account_children(ctx, field)
			case "parent":
				return ec.fieldContext_Account_parent(ctx, field)
			case "outgoingTransactions":
				return ec.fieldContext_Account_outgoingTransactions(ctx, field)
			case "incomingTransactions":
				return ec.fieldContext_Account_incomingTransactions(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type Account", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) _Account_parent(ctx context.Context, field graphql.CollectedField, obj *model.Account) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Account_parent(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Account().Parent(rctx, obj)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.(*model.Account)
	fc.Result = res
	return ec.marshalOAccount2ᚖgithubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐAccount(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Account_parent(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Account",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "tenant":
				return ec.fieldContext_Account_tenant(ctx, field)
			case "id":
				return ec.fieldContext_Account_id(ctx, field)
			case "displayName":
				return ec.fieldContext_Account_displayName(ctx, field)
			case "icon":
				return ec.fieldContext_Account_icon(ctx, field)
			case "labels":
				return ec.fieldContext_Account_labels(ctx, field)
			case "childCount":
				return ec.fieldContext_Account_childCount(ctx, field)
			case "children":
				return ec.fieldContext_Account_children(ctx, field)
			case "parent":
				return ec.fieldContext_Account_parent(ctx, field)
			case "outgoingTransactions":
				return ec.fieldContext_Account_outgoingTransactions(ctx, field)
			case "incomingTransactions":
				return ec.fieldContext_Account_incomingTransactions(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type Account", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) _Account_outgoingTransactions(ctx context.Context, field graphql.CollectedField, obj *model.Account) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Account_outgoingTransactions(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Account().OutgoingTransactions(rctx, obj)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.([]*model.Transaction)
	fc.Result = res
	return ec.marshalNTransaction2ᚕᚖgithubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐTransactionᚄ(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Account_outgoingTransactions(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Account",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "id":
				return ec.fieldContext_Transaction_id(ctx, field)
			case "Date":
				return ec.fieldContext_Transaction_Date(ctx, field)
			case "targetAccount":
				return ec.fieldContext_Transaction_targetAccount(ctx, field)
			case "sourceAccount":
				return ec.fieldContext_Transaction_sourceAccount(ctx, field)
			case "referenceID":
				return ec.fieldContext_Transaction_referenceID(ctx, field)
			case "amount":
				return ec.fieldContext_Transaction_amount(ctx, field)
			case "description":
				return ec.fieldContext_Transaction_description(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type Transaction", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) _Account_incomingTransactions(ctx context.Context, field graphql.CollectedField, obj *model.Account) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Account_incomingTransactions(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Account().IncomingTransactions(rctx, obj)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.([]*model.Transaction)
	fc.Result = res
	return ec.marshalNTransaction2ᚕᚖgithubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐTransactionᚄ(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Account_incomingTransactions(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Account",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "id":
				return ec.fieldContext_Transaction_id(ctx, field)
			case "Date":
				return ec.fieldContext_Transaction_Date(ctx, field)
			case "targetAccount":
				return ec.fieldContext_Transaction_targetAccount(ctx, field)
			case "sourceAccount":
				return ec.fieldContext_Transaction_sourceAccount(ctx, field)
			case "referenceID":
				return ec.fieldContext_Transaction_referenceID(ctx, field)
			case "amount":
				return ec.fieldContext_Transaction_amount(ctx, field)
			case "description":
				return ec.fieldContext_Transaction_description(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type Transaction", field.Name)
		},
	}
	return fc, nil
}

// endregion **************************** field.gotpl *****************************

// region    **************************** input.gotpl *****************************

func (ec *executionContext) unmarshalInputAccountChanges(ctx context.Context, obj interface{}) (model.AccountChanges, error) {
	var it model.AccountChanges
	asMap := map[string]interface{}{}
	for k, v := range obj.(map[string]interface{}) {
		asMap[k] = v
	}

	fieldsInOrder := [...]string{"displayName", "icon", "labels", "parentID"}
	for _, k := range fieldsInOrder {
		v, ok := asMap[k]
		if !ok {
			continue
		}
		switch k {
		case "displayName":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("displayName"))
			data, err := ec.unmarshalOString2ᚖstring(ctx, v)
			if err != nil {
				return it, err
			}
			it.DisplayName = data
		case "icon":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("icon"))
			data, err := ec.unmarshalOString2ᚖstring(ctx, v)
			if err != nil {
				return it, err
			}
			it.Icon = data
		case "labels":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("labels"))
			data, err := ec.unmarshalOKeyAndValueInput2ᚕᚖgithubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐKeyAndValueInputᚄ(ctx, v)
			if err != nil {
				return it, err
			}
			it.Labels = data
		case "parentID":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("parentID"))
			data, err := ec.unmarshalOID2ᚖstring(ctx, v)
			if err != nil {
				return it, err
			}
			it.ParentID = data
		}
	}

	return it, nil
}

// endregion **************************** input.gotpl *****************************

// region    ************************** interface.gotpl ***************************

// endregion ************************** interface.gotpl ***************************

// region    **************************** object.gotpl ****************************

var accountImplementors = []string{"Account"}

func (ec *executionContext) _Account(ctx context.Context, sel ast.SelectionSet, obj *model.Account) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, accountImplementors)

	out := graphql.NewFieldSet(fields)
	deferred := make(map[string]*graphql.FieldSet)
	for i, field := range fields {
		switch field.Name {
		case "__typename":
			out.Values[i] = graphql.MarshalString("Account")
		case "tenant":
			out.Values[i] = ec._Account_tenant(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				atomic.AddUint32(&out.Invalids, 1)
			}
		case "id":
			out.Values[i] = ec._Account_id(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				atomic.AddUint32(&out.Invalids, 1)
			}
		case "displayName":
			out.Values[i] = ec._Account_displayName(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				atomic.AddUint32(&out.Invalids, 1)
			}
		case "icon":
			out.Values[i] = ec._Account_icon(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				atomic.AddUint32(&out.Invalids, 1)
			}
		case "labels":
			field := field

			innerFunc := func(ctx context.Context, fs *graphql.FieldSet) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Account_labels(ctx, field, obj)
				if res == graphql.Null {
					atomic.AddUint32(&fs.Invalids, 1)
				}
				return res
			}

			if field.Deferrable != nil {
				dfs, ok := deferred[field.Deferrable.Label]
				di := 0
				if ok {
					dfs.AddField(field)
					di = len(dfs.Values) - 1
				} else {
					dfs = graphql.NewFieldSet([]graphql.CollectedField{field})
					deferred[field.Deferrable.Label] = dfs
				}
				dfs.Concurrently(di, func(ctx context.Context) graphql.Marshaler {
					return innerFunc(ctx, dfs)
				})

				// don't run the out.Concurrently() call below
				out.Values[i] = graphql.Null
				continue
			}

			out.Concurrently(i, func(ctx context.Context) graphql.Marshaler { return innerFunc(ctx, out) })
		case "childCount":
			field := field

			innerFunc := func(ctx context.Context, fs *graphql.FieldSet) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Account_childCount(ctx, field, obj)
				if res == graphql.Null {
					atomic.AddUint32(&fs.Invalids, 1)
				}
				return res
			}

			if field.Deferrable != nil {
				dfs, ok := deferred[field.Deferrable.Label]
				di := 0
				if ok {
					dfs.AddField(field)
					di = len(dfs.Values) - 1
				} else {
					dfs = graphql.NewFieldSet([]graphql.CollectedField{field})
					deferred[field.Deferrable.Label] = dfs
				}
				dfs.Concurrently(di, func(ctx context.Context) graphql.Marshaler {
					return innerFunc(ctx, dfs)
				})

				// don't run the out.Concurrently() call below
				out.Values[i] = graphql.Null
				continue
			}

			out.Concurrently(i, func(ctx context.Context) graphql.Marshaler { return innerFunc(ctx, out) })
		case "children":
			field := field

			innerFunc := func(ctx context.Context, fs *graphql.FieldSet) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Account_children(ctx, field, obj)
				if res == graphql.Null {
					atomic.AddUint32(&fs.Invalids, 1)
				}
				return res
			}

			if field.Deferrable != nil {
				dfs, ok := deferred[field.Deferrable.Label]
				di := 0
				if ok {
					dfs.AddField(field)
					di = len(dfs.Values) - 1
				} else {
					dfs = graphql.NewFieldSet([]graphql.CollectedField{field})
					deferred[field.Deferrable.Label] = dfs
				}
				dfs.Concurrently(di, func(ctx context.Context) graphql.Marshaler {
					return innerFunc(ctx, dfs)
				})

				// don't run the out.Concurrently() call below
				out.Values[i] = graphql.Null
				continue
			}

			out.Concurrently(i, func(ctx context.Context) graphql.Marshaler { return innerFunc(ctx, out) })
		case "parent":
			field := field

			innerFunc := func(ctx context.Context, fs *graphql.FieldSet) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Account_parent(ctx, field, obj)
				return res
			}

			if field.Deferrable != nil {
				dfs, ok := deferred[field.Deferrable.Label]
				di := 0
				if ok {
					dfs.AddField(field)
					di = len(dfs.Values) - 1
				} else {
					dfs = graphql.NewFieldSet([]graphql.CollectedField{field})
					deferred[field.Deferrable.Label] = dfs
				}
				dfs.Concurrently(di, func(ctx context.Context) graphql.Marshaler {
					return innerFunc(ctx, dfs)
				})

				// don't run the out.Concurrently() call below
				out.Values[i] = graphql.Null
				continue
			}

			out.Concurrently(i, func(ctx context.Context) graphql.Marshaler { return innerFunc(ctx, out) })
		case "outgoingTransactions":
			field := field

			innerFunc := func(ctx context.Context, fs *graphql.FieldSet) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Account_outgoingTransactions(ctx, field, obj)
				if res == graphql.Null {
					atomic.AddUint32(&fs.Invalids, 1)
				}
				return res
			}

			if field.Deferrable != nil {
				dfs, ok := deferred[field.Deferrable.Label]
				di := 0
				if ok {
					dfs.AddField(field)
					di = len(dfs.Values) - 1
				} else {
					dfs = graphql.NewFieldSet([]graphql.CollectedField{field})
					deferred[field.Deferrable.Label] = dfs
				}
				dfs.Concurrently(di, func(ctx context.Context) graphql.Marshaler {
					return innerFunc(ctx, dfs)
				})

				// don't run the out.Concurrently() call below
				out.Values[i] = graphql.Null
				continue
			}

			out.Concurrently(i, func(ctx context.Context) graphql.Marshaler { return innerFunc(ctx, out) })
		case "incomingTransactions":
			field := field

			innerFunc := func(ctx context.Context, fs *graphql.FieldSet) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Account_incomingTransactions(ctx, field, obj)
				if res == graphql.Null {
					atomic.AddUint32(&fs.Invalids, 1)
				}
				return res
			}

			if field.Deferrable != nil {
				dfs, ok := deferred[field.Deferrable.Label]
				di := 0
				if ok {
					dfs.AddField(field)
					di = len(dfs.Values) - 1
				} else {
					dfs = graphql.NewFieldSet([]graphql.CollectedField{field})
					deferred[field.Deferrable.Label] = dfs
				}
				dfs.Concurrently(di, func(ctx context.Context) graphql.Marshaler {
					return innerFunc(ctx, dfs)
				})

				// don't run the out.Concurrently() call below
				out.Values[i] = graphql.Null
				continue
			}

			out.Concurrently(i, func(ctx context.Context) graphql.Marshaler { return innerFunc(ctx, out) })
		default:
			panic("unknown field " + strconv.Quote(field.Name))
		}
	}
	out.Dispatch(ctx)
	if out.Invalids > 0 {
		return graphql.Null
	}

	atomic.AddInt32(&ec.deferred, int32(len(deferred)))

	for label, dfs := range deferred {
		ec.processDeferredGroup(graphql.DeferredGroup{
			Label:    label,
			Path:     graphql.GetPath(ctx),
			FieldSet: dfs,
			Context:  ctx,
		})
	}

	return out
}

// endregion **************************** object.gotpl ****************************

// region    ***************************** type.gotpl *****************************

func (ec *executionContext) marshalNAccount2githubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐAccount(ctx context.Context, sel ast.SelectionSet, v model.Account) graphql.Marshaler {
	return ec._Account(ctx, sel, &v)
}

func (ec *executionContext) marshalNAccount2ᚕᚖgithubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐAccountᚄ(ctx context.Context, sel ast.SelectionSet, v []*model.Account) graphql.Marshaler {
	ret := make(graphql.Array, len(v))
	var wg sync.WaitGroup
	isLen1 := len(v) == 1
	if !isLen1 {
		wg.Add(len(v))
	}
	for i := range v {
		i := i
		fc := &graphql.FieldContext{
			Index:  &i,
			Result: &v[i],
		}
		ctx := graphql.WithFieldContext(ctx, fc)
		f := func(i int) {
			defer func() {
				if r := recover(); r != nil {
					ec.Error(ctx, ec.Recover(ctx, r))
					ret = nil
				}
			}()
			if !isLen1 {
				defer wg.Done()
			}
			ret[i] = ec.marshalNAccount2ᚖgithubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐAccount(ctx, sel, v[i])
		}
		if isLen1 {
			f(i)
		} else {
			go f(i)
		}

	}
	wg.Wait()

	for _, e := range ret {
		if e == graphql.Null {
			return graphql.Null
		}
	}

	return ret
}

func (ec *executionContext) marshalNAccount2ᚖgithubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐAccount(ctx context.Context, sel ast.SelectionSet, v *model.Account) graphql.Marshaler {
	if v == nil {
		if !graphql.HasFieldError(ctx, graphql.GetFieldContext(ctx)) {
			ec.Errorf(ctx, "the requested element is null which the schema does not allow")
		}
		return graphql.Null
	}
	return ec._Account(ctx, sel, v)
}

func (ec *executionContext) unmarshalNAccountChanges2githubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐAccountChanges(ctx context.Context, v interface{}) (model.AccountChanges, error) {
	res, err := ec.unmarshalInputAccountChanges(ctx, v)
	return res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalOAccount2ᚖgithubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐAccount(ctx context.Context, sel ast.SelectionSet, v *model.Account) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	return ec._Account(ctx, sel, v)
}

// endregion ***************************** type.gotpl *****************************
