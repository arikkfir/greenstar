package account

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/server/middleware"
	"github.com/arikkfir/greenstar/backend/internal/server/util"
	"github.com/arikkfir/greenstar/backend/internal/util/db"
	"github.com/arikkfir/greenstar/backend/internal/util/lang"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pganalyze/pg_query_go/v5"
	"strings"
)

var (
	//go:embed sql/list.sql
	listSQL string

	//go:embed sql/insert.sql
	createSQL string

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

type HandlerImpl struct{}

func (h *HandlerImpl) List(ctx context.Context, req ListRequest) (*ListResponse, error) {
	tx := db.TxFromContext(ctx)

	res := ListResponse{}
	res.Items = make([]Account, 0)

	if sql, args, err := h.buildCountQuery(ctx, req); err != nil {
		return nil, fmt.Errorf("failed building count query: %w", err)
	} else if err := tx.QueryRow(ctx, sql, args...).Scan(&res.TotalCount); err != nil {
		return nil, fmt.Errorf("failed fetching rows count: %w\nSQL:\n%s\nArgs:\n%+v", err, sql, args)
	}

	if sql, args, err := h.buildListQuery(ctx, req); err != nil {
		return nil, fmt.Errorf("failed building list query: %w", err)
	} else if rows, err := tx.Query(ctx, sql, args...); err != nil {
		return nil, fmt.Errorf("failed fetching rows count: %w", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var a Account
			if err := rows.Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt, &a.DisplayName, &a.Icon, &a.ParentID, &a.TotalIncomingAmount, &a.TotalOutgoingAmount, &a.Balance); err != nil {
				return nil, fmt.Errorf("failed scanning row: %w", err)
			}
			res.Items = append(res.Items, a)
		}
	}

	return &res, nil
}

func (h *HandlerImpl) buildListQuery(ctx context.Context, req ListRequest) (string, []any, error) {
	args := []any{middleware.GetTenantID(ctx), req.Currency}

	result, err := pg_query.Parse(listSQL)
	if err != nil {
		return "", nil, fmt.Errorf("failed parsing query: %w", err)
	}

	if req.HasDisplayName() {
		pattern := "%" + *req.DisplayName + "%"
		args = append(args, pattern)
		util.AddWhereNode(
			result.Stmts[0].Stmt.GetSelectStmt(),
			util.MakeLikeExprNode(
				util.MakeColumnRefNode("a", "display_name"),
				util.MakeParamNode(len(args)),
			),
		)
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
		orderByColumns = []string{"a.display_name ASC"}
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
		return "", nil, fmt.Errorf("failed parsing list query for count query: %w", err)
	}

	if req.HasDisplayName() {
		pattern := "%" + *req.DisplayName + "%"
		args = append(args, pattern)
		util.AddWhereNode(
			result.Stmts[0].Stmt.GetSelectStmt(),
			util.MakeLikeExprNode(
				util.MakeColumnRefNode("a", "display_name"),
				util.MakeParamNode(len(args)),
			),
		)
	}

	util.ReplaceSelectTargetsWithCountRows(result.Stmts[0].Stmt.GetSelectStmt())
	util.ClearOrderByClause(result.Stmts[0].Stmt.GetSelectStmt())

	return util.GetSQL(result), args, nil
}

func (h *HandlerImpl) Create(ctx context.Context, req CreateRequest) (*CreateResponse, error) {
	tx := db.TxFromContext(ctx)

	if req.DisplayName == "" {
		return nil, fmt.Errorf("%w: display name is required", util.ErrBadRequest)
	}

	if req.ParentID != nil {
		_, err := h.Get(ctx, GetRequest{ID: *req.ParentID})
		if err != nil {
			if errors.Is(err, util.ErrNotFound) {
				return nil, fmt.Errorf("%w: parent account %s not found", util.ErrBadRequest, *req.ParentID)
			} else {
				return nil, fmt.Errorf("failed fetching parent account '%s': %w", *req.ParentID, err)
			}
		}
	}

	var id string
	args := []any{req.DisplayName, req.Icon, req.ParentID, middleware.GetTenantID(ctx)}
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
	tx := db.TxFromContext(ctx)

	if req.Currency == nil {
		req.Currency = lang.PtrOf(util.DefaultCurrency)
	}

	var res GetResponse
	if err := tx.QueryRow(ctx, getSQL, middleware.GetTenantID(ctx), req.ID, *req.Currency).Scan(&res.ID, &res.CreatedAt, &res.UpdatedAt, &res.DisplayName, &res.Icon, &res.ParentID, &res.TotalIncomingAmount, &res.TotalOutgoingAmount, &res.Balance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, util.ErrNotFound
		}

		var pge *pgconn.PgError
		if errors.As(err, &pge) {
			switch pge.Code {
			case pgerrcode.InvalidTextRepresentation:
				return nil, fmt.Errorf("%w: %s", util.ErrUnprocessableEntity, pge.Message)
			}
		}

		return nil, fmt.Errorf("failed fetching row: %w", err)
	}

	return &res, nil
}

func (h *HandlerImpl) Patch(ctx context.Context, req PatchRequest) (*PatchResponse, error) {
	tx := db.TxFromContext(ctx)

	args := []any{middleware.GetTenantID(ctx), req.ID}

	q, err := pg_query.Parse(patchSQL)
	if err != nil {
		return nil, fmt.Errorf("failed parsing query: %w", err)
	}

	if req.HasDisplayName() {
		if req.DisplayName == nil || *req.DisplayName == "" {
			return nil, fmt.Errorf("%w: display name must not be empty", util.ErrBadRequest)
		}
		args = append(args, *req.DisplayName)
		util.AddSetClause(q.Stmts[0].Stmt.GetUpdateStmt(), "display_name", len(args))
	}

	if req.HasIcon() {
		args = append(args, req.Icon)
		util.AddSetClause(q.Stmts[0].Stmt.GetUpdateStmt(), "icon", len(args))
	}

	if req.HasParentID() {
		if req.ParentID != nil {
			_, err := h.Get(ctx, GetRequest{ID: *req.ParentID})
			if err != nil {
				if errors.Is(err, util.ErrNotFound) {
					return nil, fmt.Errorf("%w: parent account %s not found", util.ErrBadRequest, *req.ParentID)
				} else {
					return nil, fmt.Errorf("failed fetching parent account '%s': %w", *req.ParentID, err)
				}
			}
		}
		args = append(args, req.ParentID)
		util.AddSetClause(q.Stmts[0].Stmt.GetUpdateStmt(), "parent_id", len(args))
	}

	if result, err := tx.Exec(ctx, util.GetSQL(q), args...); err != nil {
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
	tx := db.TxFromContext(ctx)

	if req.DisplayName == "" {
		return nil, fmt.Errorf("%w: display name must not be empty", util.ErrBadRequest)
	}

	if req.ParentID != nil {
		_, err := h.Get(ctx, GetRequest{ID: *req.ParentID})
		if err != nil {
			if errors.Is(err, util.ErrNotFound) {
				return nil, fmt.Errorf("%w: parent account %s not found", util.ErrBadRequest, *req.ParentID)
			} else {
				return nil, fmt.Errorf("failed fetching parent account '%s': %w", *req.ParentID, err)
			}
		}
	}

	args := []any{middleware.GetTenantID(ctx), req.ID, req.DisplayName, req.Icon, req.ParentID}
	if result, err := tx.Exec(ctx, updateSQL, args...); err != nil {
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
	tx := db.TxFromContext(ctx)

	if result, err := tx.Exec(ctx, deleteSQL, middleware.GetTenantID(ctx), req.ID); err != nil {
		return fmt.Errorf("failed deleting row: %w", err)
	} else if result.RowsAffected() == 0 {
		return fmt.Errorf("failed deleting row: no rows affected: %w", util.ErrNotFound)
	} else if rowsAffected := result.RowsAffected(); rowsAffected > 1 {
		return fmt.Errorf("failed deleting row: more than one row was affected (%d)", rowsAffected)
	}

	return nil
}

func (h *HandlerImpl) DeleteAll(ctx context.Context, _ DeleteAllRequest) error {
	tx := db.TxFromContext(ctx)

	_, err := tx.Exec(ctx, deleteAllSQL, middleware.GetTenantID(ctx))
	if err != nil {
		return fmt.Errorf("failed deleting rows: %w", err)
	}

	return nil
}
