package parser

import (
	"log"
	"vitess.io/vitess/go/vt/sqlparser"
)

func (p *Parser) extract(ast sqlparser.Statement) {
	switch stmt := ast.(type) {
	case *sqlparser.Select:
		p.extractFromSelectAST(stmt)
	default:
		log.Println("SQL Query is not SELECT type")
	}
}

func (p *Parser) extractFromSelectAST(stmt *sqlparser.Select) {
	if stmt.From != nil {
		for _, tableExpr := range stmt.From {
			switch expr := tableExpr.(type) {
			case *sqlparser.JoinTableExpr:
				p.extractFromJoinTableExpr(expr)
			}
		}
	}
	// TODO: Extract other elements as well here.
	if stmt.Where != nil {
		p.extractFromWhere(stmt.Where)
	}
}

func (p *Parser) extractFromJoinTableExpr(expr *sqlparser.JoinTableExpr) {
	_, traversedJoinConditions := traverseAST(expr)
	tableCache := map[string]string{}
	for _, c := range traversedJoinConditions {
		p.joins = append(p.joins, Joins{
			Inner: JoinColMetadata{
				Column:    c.InnerColumn,
				Table:     tableCache[c.InnerTable],
				Qualifier: c.InnerTable,
			},
			Outer: JoinColMetadata{
				Column:    c.OuterColumn,
				Table:     tableCache[c.OuterTable],
				Qualifier: c.OuterTable,
			},
		})
	}
}

func traverseAST(expr *sqlparser.JoinTableExpr) ([]ASTTraversalJoinTableMetadata, []ASTTraversalJoinConditionMetadata) {
	tableMetadata := []ASTTraversalJoinTableMetadata{}
	conditionMetadata := []ASTTraversalJoinConditionMetadata{}

	// definitely there is a condition here, extract the condition here :D
	conditionMetadata = append(conditionMetadata, extractASTTableConditionFromExpression(expr))

	switch leftExpr := expr.LeftExpr.(type) {
	case *sqlparser.AliasedTableExpr:
		tableMetadata = append(tableMetadata, extractASTTableMetadataFromExpression(leftExpr))
	case *sqlparser.JoinTableExpr:
		// if child itself is a join expression, go deeper :)
		childTables, childConditions := traverseAST(leftExpr)
		tableMetadata = append(tableMetadata, childTables...)
		conditionMetadata = append(conditionMetadata, childConditions...)
	}

	switch rightExpr := expr.RightExpr.(type) {
	case *sqlparser.AliasedTableExpr:
		tableMetadata = append(tableMetadata, extractASTTableMetadataFromExpression(rightExpr))
	case *sqlparser.JoinTableExpr:
		childTables, childConditions := traverseAST(rightExpr)
		tableMetadata = append(tableMetadata, childTables...)
		conditionMetadata = append(conditionMetadata, childConditions...)
	}
	return tableMetadata, conditionMetadata
}

func extractASTTableMetadataFromExpression(aliasedExpr *sqlparser.AliasedTableExpr) ASTTraversalJoinTableMetadata {
	result := ASTTraversalJoinTableMetadata{}

	switch expr := aliasedExpr.Expr.(type) {
	case sqlparser.TableName:
		result.Table = expr.Name.String()
	}
	result.Alias = aliasedExpr.As.String()
	return result
}

func extractASTTableConditionFromExpression(joinExpr *sqlparser.JoinTableExpr) ASTTraversalJoinConditionMetadata {
	result := ASTTraversalJoinConditionMetadata{}
	switch condition := joinExpr.Condition.On.(type) {
	case *sqlparser.ComparisonExpr:
		switch inner := condition.Left.(type) {
		case *sqlparser.ColName:
			result.InnerColumn = inner.Name.String()
			result.InnerTable = inner.Qualifier.Name.String()
		}

		switch outer := condition.Right.(type) {
		case *sqlparser.ColName:
			result.OuterColumn = outer.Name.String()
			result.OuterTable = outer.Qualifier.Name.String()
		}
	}
	return result
}
