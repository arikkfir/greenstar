// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gql

import (
	"context"
	"errors"
	"strconv"
	"sync/atomic"

	"github.com/99designs/gqlgen/graphql"
	"github.com/arikkfir/greenstar/backend/model"
	"github.com/vektah/gqlparser/v2/ast"
)

// region    ************************** generated!.gotpl **************************

// endregion ************************** generated!.gotpl **************************

// region    ***************************** args.gotpl *****************************

// endregion ***************************** args.gotpl *****************************

// region    ************************** directives.gotpl **************************

// endregion ************************** directives.gotpl **************************

// region    **************************** field.gotpl *****************************

func (ec *executionContext) _Operation_id(ctx context.Context, field graphql.CollectedField, obj *model.Operation) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Operation_id(ctx, field)
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

func (ec *executionContext) fieldContext_Operation_id(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Operation",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type ID does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Operation_name(ctx context.Context, field graphql.CollectedField, obj *model.Operation) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Operation_name(ctx, field)
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
		return obj.Name, nil
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

func (ec *executionContext) fieldContext_Operation_name(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Operation",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Operation_description(ctx context.Context, field graphql.CollectedField, obj *model.Operation) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Operation_description(ctx, field)
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
		return obj.Description, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.(*string)
	fc.Result = res
	return ec.marshalOString2ᚖstring(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Operation_description(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Operation",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Operation_status(ctx context.Context, field graphql.CollectedField, obj *model.Operation) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Operation_status(ctx, field)
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
		return obj.Status, nil
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
	res := resTmp.(model.OperationStatus)
	fc.Result = res
	return ec.marshalNOperationStatus2githubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐOperationStatus(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Operation_status(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Operation",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type OperationStatus does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Operation_result(ctx context.Context, field graphql.CollectedField, obj *model.Operation) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Operation_result(ctx, field)
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
		return obj.Result, nil
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
	res := resTmp.(model.OperationResult)
	fc.Result = res
	return ec.marshalNOperationResult2githubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐOperationResult(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Operation_result(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Operation",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type OperationResult does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Operation_createdAt(ctx context.Context, field graphql.CollectedField, obj *model.Operation) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Operation_createdAt(ctx, field)
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
		return obj.CreatedAt, nil
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
	return ec.marshalNDateTime2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Operation_createdAt(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Operation",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type DateTime does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Operation_updatedAt(ctx context.Context, field graphql.CollectedField, obj *model.Operation) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Operation_updatedAt(ctx, field)
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
		return obj.UpdatedAt, nil
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
	return ec.marshalNDateTime2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Operation_updatedAt(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Operation",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type DateTime does not have child fields")
		},
	}
	return fc, nil
}

// endregion **************************** field.gotpl *****************************

// region    **************************** input.gotpl *****************************

func (ec *executionContext) unmarshalInputOperationChanges(ctx context.Context, obj interface{}) (model.OperationChanges, error) {
	var it model.OperationChanges
	asMap := map[string]interface{}{}
	for k, v := range obj.(map[string]interface{}) {
		asMap[k] = v
	}

	fieldsInOrder := [...]string{"name", "description", "status", "result"}
	for _, k := range fieldsInOrder {
		v, ok := asMap[k]
		if !ok {
			continue
		}
		switch k {
		case "name":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("name"))
			data, err := ec.unmarshalNString2string(ctx, v)
			if err != nil {
				return it, err
			}
			it.Name = data
		case "description":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("description"))
			data, err := ec.unmarshalOString2ᚖstring(ctx, v)
			if err != nil {
				return it, err
			}
			it.Description = data
		case "status":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("status"))
			data, err := ec.unmarshalNOperationStatus2githubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐOperationStatus(ctx, v)
			if err != nil {
				return it, err
			}
			it.Status = data
		case "result":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("result"))
			data, err := ec.unmarshalNOperationResult2githubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐOperationResult(ctx, v)
			if err != nil {
				return it, err
			}
			it.Result = data
		}
	}

	return it, nil
}

// endregion **************************** input.gotpl *****************************

// region    ************************** interface.gotpl ***************************

// endregion ************************** interface.gotpl ***************************

// region    **************************** object.gotpl ****************************

var operationImplementors = []string{"Operation"}

func (ec *executionContext) _Operation(ctx context.Context, sel ast.SelectionSet, obj *model.Operation) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, operationImplementors)

	out := graphql.NewFieldSet(fields)
	deferred := make(map[string]*graphql.FieldSet)
	for i, field := range fields {
		switch field.Name {
		case "__typename":
			out.Values[i] = graphql.MarshalString("Operation")
		case "id":
			out.Values[i] = ec._Operation_id(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "name":
			out.Values[i] = ec._Operation_name(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "description":
			out.Values[i] = ec._Operation_description(ctx, field, obj)
		case "status":
			out.Values[i] = ec._Operation_status(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "result":
			out.Values[i] = ec._Operation_result(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "createdAt":
			out.Values[i] = ec._Operation_createdAt(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "updatedAt":
			out.Values[i] = ec._Operation_updatedAt(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
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

func (ec *executionContext) unmarshalNDateTime2string(ctx context.Context, v interface{}) (string, error) {
	res, err := graphql.UnmarshalString(v)
	return res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalNDateTime2string(ctx context.Context, sel ast.SelectionSet, v string) graphql.Marshaler {
	res := graphql.MarshalString(v)
	if res == graphql.Null {
		if !graphql.HasFieldError(ctx, graphql.GetFieldContext(ctx)) {
			ec.Errorf(ctx, "the requested element is null which the schema does not allow")
		}
	}
	return res
}

func (ec *executionContext) marshalNOperation2githubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐOperation(ctx context.Context, sel ast.SelectionSet, v model.Operation) graphql.Marshaler {
	return ec._Operation(ctx, sel, &v)
}

func (ec *executionContext) marshalNOperation2ᚖgithubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐOperation(ctx context.Context, sel ast.SelectionSet, v *model.Operation) graphql.Marshaler {
	if v == nil {
		if !graphql.HasFieldError(ctx, graphql.GetFieldContext(ctx)) {
			ec.Errorf(ctx, "the requested element is null which the schema does not allow")
		}
		return graphql.Null
	}
	return ec._Operation(ctx, sel, v)
}

func (ec *executionContext) unmarshalNOperationChanges2githubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐOperationChanges(ctx context.Context, v interface{}) (model.OperationChanges, error) {
	res, err := ec.unmarshalInputOperationChanges(ctx, v)
	return res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) unmarshalNOperationResult2githubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐOperationResult(ctx context.Context, v interface{}) (model.OperationResult, error) {
	var res model.OperationResult
	err := res.UnmarshalGQL(v)
	return res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalNOperationResult2githubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐOperationResult(ctx context.Context, sel ast.SelectionSet, v model.OperationResult) graphql.Marshaler {
	return v
}

func (ec *executionContext) unmarshalNOperationStatus2githubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐOperationStatus(ctx context.Context, v interface{}) (model.OperationStatus, error) {
	var res model.OperationStatus
	err := res.UnmarshalGQL(v)
	return res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalNOperationStatus2githubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐOperationStatus(ctx context.Context, sel ast.SelectionSet, v model.OperationStatus) graphql.Marshaler {
	return v
}

func (ec *executionContext) marshalOOperation2ᚖgithubᚗcomᚋarikkfirᚋgreenstarᚋbackendᚋmodelᚐOperation(ctx context.Context, sel ast.SelectionSet, v *model.Operation) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	return ec._Operation(ctx, sel, v)
}

// endregion ***************************** type.gotpl *****************************
