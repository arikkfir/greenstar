package tenant

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/auth"
	"github.com/arikkfir/greenstar/backend/internal/server/util"
	"github.com/arikkfir/greenstar/backend/internal/util/db"
	strings2 "github.com/arikkfir/greenstar/backend/internal/util/strings"
	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/client"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
	"github.com/jackc/pgx/v5"
	"github.com/pganalyze/pg_query_go/v5"
	"log/slog"
	"slices"
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

type HandlerImpl struct {
	Descope *client.DescopeClient
}

func (h *HandlerImpl) List(ctx context.Context, req ListRequest) (*ListResponse, error) {
	tx := db.TxFromContext(ctx)

	res := ListResponse{}
	res.Items = make([]Tenant, 0)

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
			var t Tenant
			if err := rows.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt, &t.DisplayName); err != nil {
				return nil, fmt.Errorf("failed scanning row: %w", err)
			}
			res.Items = append(res.Items, t)
		}
	}

	token := auth.GetToken(ctx)
	permittedTenants := slices.DeleteFunc(res.Items, func(t Tenant) bool {
		return !token.IsPermittedForTenant(t.ID, "tenants:read")
	})
	res.Items = permittedTenants
	res.TotalCount = uint(len(permittedTenants))

	return &res, nil
}

func (h *HandlerImpl) buildListQuery(_ context.Context, req ListRequest) (string, []any, error) {
	var args []any

	result, err := pg_query.Parse(listSQL)
	if err != nil {
		return "", nil, fmt.Errorf("failed parsing query: %w", err)
	}

	if req.HasID() {
		args = append(args, "%"+*req.ID+"%")
		util.AddWhereNode(
			result.Stmts[0].Stmt.GetSelectStmt(),
			util.MakeLikeExprNode(
				util.MakeColumnRefNode("t", "id"),
				util.MakeParamNode(len(args)),
			),
		)
	}

	if req.HasDisplayName() {
		args = append(args, "%"+*req.DisplayName+"%")
		util.AddWhereNode(
			result.Stmts[0].Stmt.GetSelectStmt(),
			util.MakeLikeExprNode(
				util.MakeColumnRefNode("t", "display_name"),
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
			orderByColumns = append(orderByColumns, fmt.Sprintf("t.%s %s", strcase.ToSnake(sortTokens[0]), sortTokens[1]))
		}
	} else {
		orderByColumns = []string{"t.display_name ASC"}
	}
	if err := util.SetOrderBy(result.Stmts[0].Stmt.GetSelectStmt(), orderByColumns); err != nil {
		return "", nil, fmt.Errorf("failed setting order by: %w", err)
	}

	return util.GetSQL(result), args, nil
}

func (h *HandlerImpl) buildCountQuery(ctx context.Context, req ListRequest) (string, []any, error) {
	sql, args, err := h.buildListQuery(ctx, req)
	if err != nil {
		return "", nil, fmt.Errorf("failed building list query: %w", err)
	}

	result, err := pg_query.Parse(sql)
	if err != nil {
		return "", nil, fmt.Errorf("failed parsing query: %w", err)
	}

	util.ReplaceSelectTargetsWithCountRows(result.Stmts[0].Stmt.GetSelectStmt())
	util.ClearOrderByClause(result.Stmts[0].Stmt.GetSelectStmt())

	return util.GetSQL(result), args, nil
}

func (h *HandlerImpl) Create(ctx context.Context, req CreateRequest) (*CreateResponse, error) {
	tx := db.TxFromContext(ctx)

	if !req.HasID() {
		if req.ID == "" {
			req.ID = uuid.NewString()
		}
	} else if req.ID == "" {
		return nil, fmt.Errorf("%w: id cannot be empty", util.ErrBadRequest)
	}
	if slug := strings2.Slugify(req.ID); slug != req.ID {
		return nil, fmt.Errorf("%w: id must be a valid slug (try '%s' instead)", util.ErrBadRequest, slug)
	}

	if req.DisplayName == "" {
		return nil, fmt.Errorf("%w: display name is required", util.ErrBadRequest)
	}

	if result, err := tx.Exec(ctx, createSQL, req.ID, req.DisplayName); err != nil {
		return nil, fmt.Errorf("failed inserting row: %w", err)
	} else if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return nil, fmt.Errorf("failed inserting row: no rows affected")
	} else if rowsAffected > 1 {
		return nil, fmt.Errorf("failed inserting row: unexpected number of rows affected (%d)", rowsAffected)
	}

	if err := h.Descope.Management.Tenant().Update(ctx, req.ID, &descope.TenantRequest{Name: req.DisplayName}); err != nil {
		slog.ErrorContext(ctx, "Failed updating Descope tenant - trying to create it", "err", err, "id", req.ID, "name", req.DisplayName)
		if err := h.Descope.Management.Tenant().CreateWithID(ctx, req.ID, &descope.TenantRequest{Name: req.DisplayName}); err != nil {
			return nil, fmt.Errorf("failed creating Descope tenant: %w", err)
		}
	}

	getResp, err := h.Get(ctx, GetRequest{ID: req.ID})
	if err != nil {
		return nil, fmt.Errorf("failed fetching saved row: %w", err)
	}
	resp := CreateResponse(*getResp)
	return &resp, err
}

func (h *HandlerImpl) Get(ctx context.Context, req GetRequest) (*GetResponse, error) {
	tx := db.TxFromContext(ctx)

	var res GetResponse
	if err := tx.QueryRow(ctx, getSQL, req.ID).Scan(&res.ID, &res.CreatedAt, &res.UpdatedAt, &res.DisplayName); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, util.ErrNotFound
		} else {
			return nil, fmt.Errorf("failed fetching row: %w", err)
		}
	}

	return &res, nil
}

func (h *HandlerImpl) Patch(ctx context.Context, req PatchRequest) (*PatchResponse, error) {
	tx := db.TxFromContext(ctx)

	args := []any{req.ID}

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

		if err := h.Descope.Management.Tenant().Update(ctx, req.ID, &descope.TenantRequest{Name: *req.DisplayName}); err != nil {
			return nil, fmt.Errorf("failed updating Descope tenant: %w", err)
		}
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

	if result, err := tx.Exec(ctx, updateSQL, req.ID, req.DisplayName); err != nil {
		return nil, fmt.Errorf("failed updating row: %w", err)
	} else if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return nil, fmt.Errorf("failed updating row: no rows affected: %w", util.ErrNotFound)
	} else if rowsAffected > 1 {
		return nil, fmt.Errorf("failed updating row: more than one row was affected (%d)", rowsAffected)
	}

	if err := h.Descope.Management.Tenant().Update(ctx, req.ID, &descope.TenantRequest{Name: req.DisplayName}); err != nil {
		return nil, fmt.Errorf("failed updating Descope tenant: %w", err)
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

	if result, err := tx.Exec(ctx, deleteSQL, req.ID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("failed deleting row: no rows affected: %w", util.ErrNotFound)
		} else {
			return fmt.Errorf("failed deleting row: %w", err)
		}
	} else if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("failed updating row: no rows affected: %w", util.ErrNotFound)
	} else if rowsAffected > 1 {
		return fmt.Errorf("failed updating row: more than one row was affected (%d)", rowsAffected)
	}

	if err := h.Descope.Management.Tenant().Delete(ctx, req.ID, true); err != nil {
		return fmt.Errorf("failed deleting Descope tenant: %w", err)
	}

	return nil
}

func (h *HandlerImpl) DeleteAll(ctx context.Context, _ DeleteAllRequest) error {
	tx := db.TxFromContext(ctx)

	_, err := tx.Exec(ctx, deleteAllSQL)
	if err != nil {
		return fmt.Errorf("failed deleting rows: %w", err)
	}

	tenants, err := h.Descope.Management.Tenant().LoadAll(ctx)
	if err != nil {
		return fmt.Errorf("failed loading Descope tenants: %w", err)
	}

	for _, tenant := range tenants {
		if err := h.Descope.Management.Tenant().Delete(ctx, tenant.ID, true); err != nil {
			return fmt.Errorf("failed deleting Descope tenant: %w", err)
		}
	}

	return nil
}
