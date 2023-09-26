package util

import (
	"fmt"
	pg_query "github.com/pganalyze/pg_query_go/v5"
	"regexp"
	"strings"
)

var (
	orderByColRefRE = regexp.MustCompile(`^([a-zA-Z0-9_.]+)(\s+(ASC|DESC)?(NULLS FIRST|LAST)?)?$`)
)

func ReplaceSelectTargetsWithCountRows(selectStmt *pg_query.SelectStmt) {
	funcCallNode := pg_query.MakeFuncCallNode(
		[]*pg_query.Node{pg_query.MakeStrNode("count")},
		nil,
		-1,
	)
	funcCallNode.GetFuncCall().AggStar = true
	selectStmt.TargetList = []*pg_query.Node{pg_query.MakeResTargetNodeWithVal(funcCallNode, -1)}
}

func MakeOpExprNode(op string, left, right *pg_query.Node) *pg_query.Node {
	return pg_query.MakeAExprNode(
		pg_query.A_Expr_Kind_AEXPR_OP,
		[]*pg_query.Node{pg_query.MakeStrNode(op)},
		left,
		right,
		-1,
	)
}

func MakeIsNullNode(node *pg_query.Node) *pg_query.Node {
	return &pg_query.Node{
		Node: &pg_query.Node_NullTest{
			NullTest: &pg_query.NullTest{
				Arg:          node,
				Nulltesttype: pg_query.NullTestType_IS_NULL,
				Location:     -1,
			},
		},
	}
}

func MakeIsNotNullNode(node *pg_query.Node) *pg_query.Node {
	return &pg_query.Node{
		Node: &pg_query.Node_NullTest{
			NullTest: &pg_query.NullTest{
				Arg:          node,
				Nulltesttype: pg_query.NullTestType_IS_NOT_NULL,
				Location:     -1,
			},
		},
	}
}

func MakeLikeExprNode(left, right *pg_query.Node) *pg_query.Node {
	return pg_query.MakeAExprNode(
		pg_query.A_Expr_Kind_AEXPR_LIKE,
		[]*pg_query.Node{pg_query.MakeStrNode("~~")},
		left,
		right,
		-1,
	)
}

func MakeAndExprNode(left, right *pg_query.Node) *pg_query.Node {
	return pg_query.MakeBoolExprNode(pg_query.BoolExprType_AND_EXPR, []*pg_query.Node{left, right}, -1)
}

func MakeOrExprNode(left, right *pg_query.Node) *pg_query.Node {
	return pg_query.MakeBoolExprNode(pg_query.BoolExprType_OR_EXPR, []*pg_query.Node{left, right}, -1)
}

func MakeColumnRefNode(tableAlias, columnName string) *pg_query.Node {
	return pg_query.MakeColumnRefNode(
		[]*pg_query.Node{pg_query.MakeStrNode(tableAlias), pg_query.MakeStrNode(columnName)},
		-1,
	)
}

func MakeParamNode(index int) *pg_query.Node {
	return pg_query.MakeParamRefNode(int32(index), -1)
}

func AddWhereNode(selectStmt *pg_query.SelectStmt, node *pg_query.Node) {
	whereClause := selectStmt.WhereClause

	if whereClause == nil {
		selectStmt.WhereClause = node
	} else if _, ok := whereClause.Node.(*pg_query.Node_AExpr); ok {
		selectStmt.WhereClause = pg_query.MakeBoolExprNode(
			pg_query.BoolExprType_AND_EXPR,
			[]*pg_query.Node{whereClause, node},
			-1,
		)
	} else if boolExpr, ok := whereClause.Node.(*pg_query.Node_BoolExpr); ok {
		boolExpr.BoolExpr.Args = append(boolExpr.BoolExpr.Args, node)
	} else {
		panic("Unknown where clause type")
	}
}

func AddSetClause(updateStmt *pg_query.UpdateStmt, columnName string, paramIndex int) {
	updateStmt.TargetList = append(updateStmt.TargetList,
		pg_query.MakeResTargetNodeWithNameAndVal(columnName, pg_query.MakeParamRefNode(int32(paramIndex), -1), -1),
	)
}

func SetOffset(selectStmt *pg_query.SelectStmt, offset int64) {
	selectStmt.LimitOffset = pg_query.MakeAConstIntNode(offset, -1)
}

func RemoveOffset(selectStmt *pg_query.SelectStmt) {
	selectStmt.LimitOffset = nil
}

func SetLimit(selectStmt *pg_query.SelectStmt, count int64) {
	selectStmt.LimitCount = pg_query.MakeAConstIntNode(count, -1)
	selectStmt.LimitOption = pg_query.LimitOption_LIMIT_OPTION_COUNT
}

func RemoveLimit(selectStmt *pg_query.SelectStmt) {
	selectStmt.LimitCount = nil
	selectStmt.LimitOption = pg_query.LimitOption_LIMIT_OPTION_UNDEFINED
}

func ClearWithClause(selectStmt *pg_query.SelectStmt) {
	selectStmt.WithClause = nil
}

func ClearOrderByClause(selectStmt *pg_query.SelectStmt) {
	selectStmt.SortClause = nil
}

func AddJoin(selectStmt *pg_query.SelectStmt, joinType pg_query.JoinType, rightTable *pg_query.Node, condition *pg_query.Node) {
	if len(selectStmt.FromClause) != 1 {
		panic(fmt.Errorf("expected single expression in FROM clause: %s", selectStmt.FromClause[0].String()))
	}

	selectStmt.FromClause = []*pg_query.Node{
		pg_query.MakeJoinExprNode(
			joinType,
			selectStmt.FromClause[0],
			rightTable,
			condition,
		),
	}
}

func SetOrderBy(selectStmt *pg_query.SelectStmt, elements []string) error {
	if len(elements) == 0 {
		return nil
	}

	var sortByNodes []*pg_query.Node
	for _, elt := range elements {

		matches := orderByColRefRE.FindStringSubmatch(elt)
		if matches == nil {
			return fmt.Errorf("invalid order by element: %s", elt)
		}

		var columnRefStrNodes []*pg_query.Node
		direction := pg_query.SortByDir_SORTBY_DEFAULT
		nullsStrategy := pg_query.SortByNulls_SORTBY_NULLS_DEFAULT
		if len(matches) > 0 {
			for _, token := range strings.Split(matches[1], ".") {
				columnRefStrNodes = append(columnRefStrNodes, pg_query.MakeStrNode(token))
			}
		} else {
			return fmt.Errorf("invalid order by element: %s", elt)
		}

		if len(matches) >= 3 && matches[3] != "" {
			switch matches[3] {
			case "ASC":
				direction = pg_query.SortByDir_SORTBY_ASC
			case "DESC":
				direction = pg_query.SortByDir_SORTBY_DESC
			default:
				return fmt.Errorf("invalid order by direction element: %s", elt)
			}
		}

		if len(matches) >= 4 && matches[4] != "" {
			switch matches[4] {
			case "NULLS FIRST":
				nullsStrategy = pg_query.SortByNulls_SORTBY_NULLS_FIRST
			case "NULLS LAST":
				nullsStrategy = pg_query.SortByNulls_SORTBY_NULLS_LAST
			default:
				return fmt.Errorf("invalid order by null strategy element: %s", elt)
			}
		}

		sortByNodes = append(sortByNodes,
			pg_query.MakeSortByNode(
				pg_query.MakeColumnRefNode(columnRefStrNodes, -1),
				direction,
				nullsStrategy,
				-1,
			),
		)
	}

	selectStmt.SortClause = sortByNodes
	return nil
}

func GetSQL(result *pg_query.ParseResult) string {
	sql, err := pg_query.Deparse(result)
	if err != nil {
		panic(fmt.Errorf("error deparsing SQL: %s", err))
	}
	return sql
}
