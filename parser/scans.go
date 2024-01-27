package parser

import (
	"vitess.io/vitess/go/vt/sqlparser"
)

type InternalComparisonType string

const (
	defaultType   InternalComparisonType = "default"
	colNameType   InternalComparisonType = "colName"
	subQueryType  InternalComparisonType = "subQuery"
	andComparison InternalComparisonType = "and"
	orComparison  InternalComparisonType = "or"
)

func (t InternalComparisonType) isColName() bool {
	return t == colNameType
}
func (t InternalComparisonType) isSubQuery() bool {
	return t == subQueryType
}
func (t InternalComparisonType) isAndOr() bool {
	return t == andComparison || t == orComparison
}

func (p *Parser) extractScans(ast sqlparser.Statement) {
	switch stmt := ast.(type) {
	case *sqlparser.Select:
		p.extractScansFromSelectStatement(stmt)
	}
}

func (p *Parser) extractScansFromSelectStatement(stmt *sqlparser.Select) {
	scanResult := []Scan{}
	if stmt.Where != nil {
		switch expr := stmt.Where.Expr.(type) {
		case *sqlparser.ComparisonExpr:
			switch left := expr.Left.(type) {
			case *sqlparser.ColName:
				scanResult = append(scanResult, Scan{
					Column: left.Name.String(),
					Table:  left.Qualifier.Name.String(),
				})
			}
		}
	}
	p.scans = scanResult
}

func getInternalComparisonType(expr sqlparser.Expr) InternalComparisonType {
	switch expr.(type) {
	case *sqlparser.ColName:
		return colNameType
	case *sqlparser.Subquery:
		return subQueryType
	case *sqlparser.AndExpr:
		return andComparison
	case *sqlparser.OrExpr:
		return orComparison
	default:
		return defaultType
	}
}

func (p *Parser) extractFromAndExpr(expr *sqlparser.AndExpr) {
	switch leftExpr := expr.Left.(type) {
	case *sqlparser.ComparisonExpr:
		p.extractFromComparisonExpr(leftExpr)
	case *sqlparser.AndExpr:
		p.extractFromAndExpr(leftExpr)
	case *sqlparser.OrExpr:
		p.extractFromOrExpr(leftExpr)
	}

	switch rightExpr := expr.Right.(type) {
	case *sqlparser.ComparisonExpr:
		p.extractFromComparisonExpr(rightExpr)
	case *sqlparser.AndExpr:
		p.extractFromAndExpr(rightExpr)
	case *sqlparser.OrExpr:
		p.extractFromOrExpr(rightExpr)
	}
}

func (p *Parser) extractFromOrExpr(expr *sqlparser.OrExpr) {
	switch leftExpr := expr.Left.(type) {
	case *sqlparser.ComparisonExpr:
		p.extractFromComparisonExpr(leftExpr)
	case *sqlparser.AndExpr:
		p.extractFromAndExpr(leftExpr)
	case *sqlparser.OrExpr:
		p.extractFromOrExpr(leftExpr)
	}

	switch rightExpr := expr.Right.(type) {
	case *sqlparser.ComparisonExpr:
		p.extractFromComparisonExpr(rightExpr)
	case *sqlparser.AndExpr:
		p.extractFromAndExpr(rightExpr)
	case *sqlparser.OrExpr:
		p.extractFromOrExpr(rightExpr)
	}
}

func (p *Parser) extractFromComparisonExpr(expr *sqlparser.ComparisonExpr) {
	// extract the join!
	if expr.Operator == sqlparser.EqualOp &&
		getInternalComparisonType(expr.Left).isColName() &&
		getInternalComparisonType(expr.Right).isColName() {

		var innerCol, outerCol, innerTable, outerTable string
		switch colName := expr.Left.(type) {
		case *sqlparser.ColName:
			innerCol = colName.Name.String()
			innerTable = colName.Qualifier.Name.String()
		}

		switch colName := expr.Right.(type) {
		case *sqlparser.ColName:
			outerCol = colName.Name.String()
			outerTable = colName.Qualifier.Name.String()
		}

		p.joins = append(p.joins, Joins{
			Inner: JoinColMetadata{
				Column:    innerCol,
				Qualifier: innerTable,
			},
			Outer: JoinColMetadata{
				Column:    outerCol,
				Qualifier: outerTable,
			},
		})
	}

	if getInternalComparisonType(expr.Left).isSubQuery() {
		p.extractFromSubQuery(expr.Left)
	}
	if getInternalComparisonType(expr.Right).isSubQuery() {
		p.extractFromSubQuery(expr.Right)
	}
}

func (p *Parser) extractFromSubQuery(expr sqlparser.Expr) {
	switch inner := expr.(type) {
	case *sqlparser.Subquery:
		switch selectStmt := inner.Select.(type) {
		case *sqlparser.Select:
			p.extractFromSelectAST(selectStmt)
		}
	}
}
