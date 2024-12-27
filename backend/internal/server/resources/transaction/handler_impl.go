package transaction

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/account"
	util2 "github.com/arikkfir/greenstar/backend/internal/server/resources/util"
	"github.com/arikkfir/greenstar/backend/internal/server/util"
	"github.com/arikkfir/greenstar/backend/internal/util/db"
	"github.com/arikkfir/greenstar/backend/internal/util/lang"
	"github.com/arikkfir/greenstar/backend/internal/util/tenant"
	"github.com/jackc/pgx/v5"
	"github.com/pganalyze/pg_query_go/v5"
	"strings"
)

var (
	//go:embed sql/list.sql
	listSQL string

	//go:embed sql/insert.sql
	createSQL string

	//go:embed sql/insert_mock_rates.sql
	insertMockRatesSQL string

	//go:embed sql/get.sql
	getSQL string

	//go:embed sql/patch.sql
	patchSQL string

	//go:embed sql/update.sql
	updateSQL string

	//go:embed sql/delete.sql
	deleteSQL string

	//go:embed sql/delete_all.sql
	deleteAllSQL string
)

type HandlerImpl struct {
	AccountsHandler account.Handler
}

func (h *HandlerImpl) List(ctx context.Context, req ListRequest) (*ListResponse, error) {
	tx := db.GetTransaction(ctx)

	res := ListResponse{}
	res.Items = make([]Transaction, 0)

	if sql, args, err := h.buildCountQuery(ctx, req); err != nil {
		return nil, fmt.Errorf("failed building count query: %w", err)
	} else if err := tx.QueryRow(ctx, sql, args...).Scan(&res.TotalCount); err != nil {
		return nil, fmt.Errorf("failed fetching rows count: %w", err)
	}

	if sql, args, err := h.buildListQuery(ctx, req); err != nil {
		return nil, fmt.Errorf("failed building list query: %w", err)
	} else if rows, err := tx.Query(ctx, sql, args...); err != nil {
		return nil, fmt.Errorf("failed fetching rows count: %w", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var t Transaction
			if err := rows.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt, &t.Date, &t.ReferenceID, &t.Amount, &t.Currency, &t.ConvertedAmount, &t.Description, &t.SourceAccountID, &t.TargetAccountID); err != nil {
				return nil, fmt.Errorf("failed scanning row: %w", err)
			}
			res.Items = append(res.Items, t)
		}
	}

	return &res, nil
}

func (h *HandlerImpl) buildListQuery(ctx context.Context, req ListRequest) (string, []any, error) {
	args := []any{nil, nil, tenant.GetTenantID(ctx), req.Currency}

	result, err := pg_query.Parse(listSQL)
	if err != nil {
		return "", nil, fmt.Errorf("failed parsing query: %w", err)
	}

	if req.HasMinDate() {
		args = append(args, req.MinDate)
		util.AddWhereNode(
			result.Stmts[0].Stmt.GetSelectStmt(),
			util.MakeOpExprNode(
				">=",
				util.MakeColumnRefNode("t", "date"),
				util.MakeParamNode(len(args)),
			),
		)
	}
	if req.HasMaxDate() {
		args = append(args, req.MaxDate)
		util.AddWhereNode(
			result.Stmts[0].Stmt.GetSelectStmt(),
			util.MakeOpExprNode(
				"<=",
				util.MakeColumnRefNode("t", "date"),
				util.MakeParamNode(len(args)),
			),
		)
	}
	if req.HasReferenceID() {
		args = append(args, req.MaxDate)
		util.AddWhereNode(
			result.Stmts[0].Stmt.GetSelectStmt(),
			util.MakeLikeExprNode(
				util.MakeColumnRefNode("t", "date"),
				util.MakeParamNode(len(args)),
			),
		)
	}
	if req.HasMinAmount() {
		args = append(args, req.MinAmount)
		util.AddWhereNode(
			result.Stmts[0].Stmt.GetSelectStmt(),
			util.MakeOpExprNode(
				">=",
				util.MakeColumnRefNode("t", "amount"),
				util.MakeParamNode(len(args)),
			),
		)
	}
	if req.HasMaxAmount() {
		args = append(args, req.MaxAmount)
		util.AddWhereNode(
			result.Stmts[0].Stmt.GetSelectStmt(),
			util.MakeOpExprNode(
				"<=",
				util.MakeColumnRefNode("t", "amount"),
				util.MakeParamNode(len(args)),
			),
		)
	}
	if req.HasDescription() {
		if req.Description == nil {
			util.AddWhereNode(
				result.Stmts[0].Stmt.GetSelectStmt(),
				util.MakeIsNullNode(util.MakeColumnRefNode("t", "description")),
			)
		} else {
			args = append(args, req.Description)
			util.AddWhereNode(
				result.Stmts[0].Stmt.GetSelectStmt(),
				util.MakeLikeExprNode(
					util.MakeColumnRefNode("t", "description"),
					util.MakeParamNode(len(args)),
				),
			)
		}
	}
	if req.HasSourceAccountID() || req.HasTargetAccountID() {
		if req.HasSourceAccountID() {
			util.AddJoin(
				result.Stmts[0].Stmt.GetSelectStmt(),
				pg_query.JoinType_JOIN_LEFT,
				pg_query.MakeFullRangeVarNode("", "source_accounts", "wsa", -1),
				pg_query.MakeAExprNode(
					pg_query.A_Expr_Kind_AEXPR_OP,
					[]*pg_query.Node{pg_query.MakeStrNode("=")},
					pg_query.MakeColumnRefNode(
						[]*pg_query.Node{pg_query.MakeStrNode("t"), pg_query.MakeStrNode("source_account_id")},
						-1,
					),
					pg_query.MakeColumnRefNode(
						[]*pg_query.Node{pg_query.MakeStrNode("wsa"), pg_query.MakeStrNode("id")},
						-1,
					),
					-1,
				),
			)
			args[0] = req.SourceAccountID
		}

		if req.HasTargetAccountID() {
			util.AddJoin(
				result.Stmts[0].Stmt.GetSelectStmt(),
				pg_query.JoinType_JOIN_LEFT,
				pg_query.MakeFullRangeVarNode("", "target_accounts", "wta", -1),
				pg_query.MakeAExprNode(
					pg_query.A_Expr_Kind_AEXPR_OP,
					[]*pg_query.Node{pg_query.MakeStrNode("=")},
					pg_query.MakeColumnRefNode(
						[]*pg_query.Node{pg_query.MakeStrNode("t"), pg_query.MakeStrNode("target_account_id")},
						-1,
					),
					pg_query.MakeColumnRefNode(
						[]*pg_query.Node{pg_query.MakeStrNode("wta"), pg_query.MakeStrNode("id")},
						-1,
					),
					-1,
				),
			)
			args[1] = req.TargetAccountID
		}

		// Now add a WHERE clause that ensures that the transaction complies with source, target or both
		if req.HasSourceAccountID() && req.HasTargetAccountID() {
			util.AddWhereNode(
				result.Stmts[0].Stmt.GetSelectStmt(),
				util.MakeOrExprNode(
					util.MakeIsNotNullNode(util.MakeColumnRefNode("wsa", "id")),
					util.MakeIsNotNullNode(util.MakeColumnRefNode("wta", "id")),
				),
			)
		} else if req.HasSourceAccountID() && !req.HasTargetAccountID() {
			util.AddWhereNode(
				result.Stmts[0].Stmt.GetSelectStmt(),
				util.MakeIsNotNullNode(util.MakeColumnRefNode("wsa", "id")),
			)
		} else if !req.HasSourceAccountID() && req.HasTargetAccountID() {
			util.AddWhereNode(
				result.Stmts[0].Stmt.GetSelectStmt(),
				util.MakeIsNotNullNode(util.MakeColumnRefNode("wta", "id")),
			)
		}
	}

	if req.Offset != nil {
		util.SetOffset(result.Stmts[0].Stmt.GetSelectStmt(), int64(*req.Offset))
	}

	if req.Count != nil {
		util.SetLimit(result.Stmts[0].Stmt.GetSelectStmt(), int64(*req.Count))
	}

	var orderByColumns []string
	if len(req.Sort) > 0 {
		for _, sort := range req.Sort {
			sortTokens := strings.Split(sort, ":")
			orderByColumns = append(orderByColumns, "a."+sortTokens[0]+" "+sortTokens[1])
		}
	} else {
		orderByColumns = []string{"t.date DESC"}
	}
	if err := util.SetOrderBy(result.Stmts[0].Stmt.GetSelectStmt(), orderByColumns); err != nil {
		return "", nil, fmt.Errorf("failed setting order by: %w", err)
	}

	return util.GetSQL(result), args, nil
}

func (h *HandlerImpl) buildCountQuery(ctx context.Context, req ListRequest) (string, []any, error) {
	sql, args, err := h.buildListQuery(ctx, req)
	if err != nil {
		return "", nil, fmt.Errorf("failed building count query: %w", err)
	}

	result, err := pg_query.Parse(sql)
	if err != nil {
		return "", nil, fmt.Errorf("failed parsing list query for count query: %w\nSQL: %s", err, sql)
	}

	util.RemoveOffset(result.Stmts[0].Stmt.GetSelectStmt())
	util.RemoveLimit(result.Stmts[0].Stmt.GetSelectStmt())
	util.ReplaceSelectTargetsWithCountRows(result.Stmts[0].Stmt.GetSelectStmt())
	util.ClearOrderByClause(result.Stmts[0].Stmt.GetSelectStmt())

	return util.GetSQL(result), args, nil
}

func (h *HandlerImpl) Create(ctx context.Context, req CreateRequest) (*CreateResponse, error) {
	tx := db.GetTransaction(ctx)

	if req.Amount.IsZero() {
		return nil, fmt.Errorf("%w: amount is required", util.ErrBadRequest)
	} else if req.Currency == "" {
		return nil, fmt.Errorf("%w: currency is required", util.ErrBadRequest)
	} else if req.Date.IsZero() {
		return nil, fmt.Errorf("%w: date is required", util.ErrBadRequest)
	} else if req.ReferenceID == "" {
		return nil, fmt.Errorf("%w: reference ID is required", util.ErrBadRequest)
	} else if req.SourceAccountID == "" {
		return nil, fmt.Errorf("%w: source account is required", util.ErrBadRequest)
	} else if req.TargetAccountID == "" {
		return nil, fmt.Errorf("%w: target account is required", util.ErrBadRequest)
	} else if _, err := h.AccountsHandler.Get(ctx, account.GetRequest{ID: req.SourceAccountID}); err != nil {
		if errors.Is(err, util.ErrNotFound) {
			return nil, fmt.Errorf("%w: source account %s not found", util.ErrBadRequest, req.SourceAccountID)
		} else {
			return nil, fmt.Errorf("failed fetching source account '%s': %w", req.SourceAccountID, err)
		}
	} else if _, err := h.AccountsHandler.Get(ctx, account.GetRequest{ID: req.TargetAccountID}); err != nil {
		if errors.Is(err, util.ErrNotFound) {
			return nil, fmt.Errorf("%w: source target %s not found", util.ErrBadRequest, req.TargetAccountID)
		} else {
			return nil, fmt.Errorf("failed fetching target account '%s': %w", req.TargetAccountID, err)
		}
	}

	if _, err := tx.Exec(ctx, insertMockRatesSQL, req.Date, req.Currency); err != nil {
		return nil, fmt.Errorf("failed inserting mock rate row: %w", err)
	}

	var id string
	args := []any{req.Amount, req.Currency, req.Date, req.Description, req.ReferenceID, req.SourceAccountID, req.TargetAccountID}
	if err := tx.QueryRow(ctx, createSQL, args...).Scan(&id); err != nil {
		return nil, fmt.Errorf("failed inserting row: %w", err)
	}

	getResp, err := h.Get(ctx, GetRequest{ID: id})
	if err != nil {
		return nil, fmt.Errorf("failed fetching saved row: %w", err)
	}
	resp := CreateResponse(*getResp)
	return &resp, err
}

func (h *HandlerImpl) Get(ctx context.Context, req GetRequest) (*GetResponse, error) {
	tx := db.GetTransaction(ctx)

	if req.Currency == nil {
		req.Currency = lang.PtrOf(util2.DefaultCurrency)
	}

	var res GetResponse
	if err := tx.QueryRow(ctx, getSQL, req.ID, tenant.GetTenantID(ctx), req.Currency).Scan(&res.ID, &res.CreatedAt, &res.UpdatedAt, &res.Amount, &res.Currency, &res.ConvertedAmount, &res.Date, &res.Description, &res.ReferenceID, &res.SourceAccountID, &res.TargetAccountID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, util.ErrNotFound
		} else {
			return nil, fmt.Errorf("failed fetching row: %w", err)
		}
	}

	return &res, nil
}

func (h *HandlerImpl) Patch(ctx context.Context, req PatchRequest) (*PatchResponse, error) {
	tx := db.GetTransaction(ctx)

	args := []any{req.ID}

	q, err := pg_query.Parse(patchSQL)
	if err != nil {
		return nil, fmt.Errorf("failed parsing query: %w", err)
	}

	if req.HasAmount() {
		if req.Amount.IsZero() {
			return nil, fmt.Errorf("%w: amount is required", util.ErrBadRequest)
		}
		args = append(args, *req.Amount)
		util.AddSetClause(q.Stmts[0].Stmt.GetUpdateStmt(), "amount", len(args))
	}
	if req.HasCurrency() {
		if *req.Currency == "" {
			return nil, fmt.Errorf("%w: currency is required", util.ErrBadRequest)
		}
		args = append(args, req.Currency)
		util.AddSetClause(q.Stmts[0].Stmt.GetUpdateStmt(), "currency", len(args))
	}
	if req.HasDate() {
		if req.Date.IsZero() {
			return nil, fmt.Errorf("%w: date is required", util.ErrBadRequest)
		}
		args = append(args, req.Date)
		util.AddSetClause(q.Stmts[0].Stmt.GetUpdateStmt(), "date", len(args))
	}
	if req.HasDescription() {
		args = append(args, req.Description)
		util.AddSetClause(q.Stmts[0].Stmt.GetUpdateStmt(), "description", len(args))
	}
	if req.HasReferenceID() {
		if *req.ReferenceID == "" {
			return nil, fmt.Errorf("%w: reference ID is required", util.ErrBadRequest)
		}
		args = append(args, req.ReferenceID)
		util.AddSetClause(q.Stmts[0].Stmt.GetUpdateStmt(), "reference_id", len(args))
	}
	if req.HasSourceAccountID() {
		if *req.SourceAccountID == "" {
			return nil, fmt.Errorf("%w: source account is required", util.ErrBadRequest)
		} else if _, err := h.AccountsHandler.Get(ctx, account.GetRequest{ID: *req.SourceAccountID}); err != nil {
			if errors.Is(err, util.ErrNotFound) {
				return nil, fmt.Errorf("%w: source account %s not found", util.ErrBadRequest, *req.SourceAccountID)
			} else {
				return nil, fmt.Errorf("failed fetching source account '%s': %w", *req.SourceAccountID, err)
			}
		}
		args = append(args, req.SourceAccountID)
		util.AddSetClause(q.Stmts[0].Stmt.GetUpdateStmt(), "source_account_id", len(args))
	}
	if req.HasTargetAccountID() {
		if *req.TargetAccountID == "" {
			return nil, fmt.Errorf("%w: target account is required", util.ErrBadRequest)
		} else if _, err := h.AccountsHandler.Get(ctx, account.GetRequest{ID: *req.SourceAccountID}); err != nil {
			if errors.Is(err, util.ErrNotFound) {
				return nil, fmt.Errorf("%w: target account %s not found", util.ErrBadRequest, *req.TargetAccountID)
			} else {
				return nil, fmt.Errorf("failed fetching target account '%s': %w", *req.TargetAccountID, err)
			}
		}
		args = append(args, req.TargetAccountID)
		util.AddSetClause(q.Stmts[0].Stmt.GetUpdateStmt(), "target_account_id", len(args))
	}

	result, err := tx.Exec(ctx, util.GetSQL(q), args...)
	if err != nil {
		return nil, fmt.Errorf("failed saving row: %w", err)
	} else if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return nil, fmt.Errorf("failed saving row: no rows affected: %w", util.ErrNotFound)
	} else if rowsAffected > 1 {
		return nil, fmt.Errorf("failed saving row: unexpected number of rows affected (%d)", rowsAffected)
	}

	getResp, err := h.Get(ctx, GetRequest{ID: req.ID})
	if err != nil {
		return nil, fmt.Errorf("failed fetching updated row: %w", err)
	}
	resp := PatchResponse(*getResp)
	return &resp, err
}

func (h *HandlerImpl) Update(ctx context.Context, req UpdateRequest) (*UpdateResponse, error) {
	tx := db.GetTransaction(ctx)

	if req.Amount.IsZero() {
		return nil, fmt.Errorf("%w: amount is required", util.ErrBadRequest)
	} else if req.Currency == "" {
		return nil, fmt.Errorf("%w: currency is required", util.ErrBadRequest)
	} else if req.Date.IsZero() {
		return nil, fmt.Errorf("%w: date is required", util.ErrBadRequest)
	} else if req.ReferenceID == "" {
		return nil, fmt.Errorf("%w: reference ID is required", util.ErrBadRequest)
	} else if req.SourceAccountID == "" {
		return nil, fmt.Errorf("%w: source account is required", util.ErrBadRequest)
	} else if _, err := h.Get(ctx, GetRequest{ID: req.SourceAccountID}); err != nil {
		if errors.Is(err, util.ErrNotFound) {
			return nil, fmt.Errorf("%w: source account %s not found", util.ErrBadRequest, req.SourceAccountID)
		} else {
			return nil, fmt.Errorf("failed fetching source account '%s': %w", req.SourceAccountID, err)
		}
	} else if req.TargetAccountID == "" {
		return nil, fmt.Errorf("%w: target account is required", util.ErrBadRequest)
	} else if _, err := h.Get(ctx, GetRequest{ID: req.TargetAccountID}); err != nil {
		if errors.Is(err, util.ErrNotFound) {
			return nil, fmt.Errorf("%w: target target %s not found", util.ErrBadRequest, req.TargetAccountID)
		} else {
			return nil, fmt.Errorf("failed fetching target account '%s': %w", req.TargetAccountID, err)
		}
	}

	args := []any{req.Amount, req.Currency, req.Date, req.Description, req.ReferenceID, req.SourceAccountID, req.TargetAccountID, tenant.GetTenantID(ctx), req.ID}
	result, err := tx.Exec(ctx, updateSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("failed updating row: %w", err)
	} else if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return nil, fmt.Errorf("failed updating row: no rows affected: %w", util.ErrNotFound)
	} else if rowsAffected > 1 {
		return nil, fmt.Errorf("failed updating row: more than one row was affected (%d)", rowsAffected)
	}

	getResp, err := h.Get(ctx, GetRequest{ID: req.ID})
	if err != nil {
		return nil, fmt.Errorf("failed fetching updated row: %w", err)
	}
	resp := UpdateResponse(*getResp)
	return &resp, err
}

func (h *HandlerImpl) Delete(ctx context.Context, req DeleteRequest) error {
	tx := db.GetTransaction(ctx)

	result, err := tx.Exec(ctx, deleteSQL, tenant.GetTenantID(ctx), req.ID)
	if err != nil {
		return fmt.Errorf("failed deleting row: %w", err)
	} else if result.RowsAffected() == 0 {
		return fmt.Errorf("failed deleting row: no rows affected: %w", util.ErrNotFound)
	} else if rowsAffected := result.RowsAffected(); rowsAffected > 1 {
		return fmt.Errorf("failed deleting row: more than one row was affected (%d)", rowsAffected)
	}

	return nil
}

func (h *HandlerImpl) DeleteAll(ctx context.Context, _ DeleteAllRequest) error {
	tx := db.GetTransaction(ctx)

	_, err := tx.Exec(ctx, deleteAllSQL, tenant.GetTenantID(ctx))
	if err != nil {
		return fmt.Errorf("failed deleting rows: %w\nSQL: %s", err, deleteAllSQL)
	}

	return nil
}
