package parser

import (
	"log"
	"vitess.io/vitess/go/vt/sqlparser"
)

func (p *Parser) extractJoins(ast sqlparser.Statement) {
	switch stmt := ast.(type) {
	case *sqlparser.Select:
		p.extractJoinFromSelectStatement(stmt)
	default:
		log.Println("SQL Query is not SELECT type")
	}
}

func (p *Parser) extractJoinFromSelectStatement(stmt *sqlparser.Select) {
	if stmt.From != nil {
		for _, tableExpr := range stmt.From {
			switch expr := tableExpr.(type) {
			case *sqlparser.JoinTableExpr:
				p.extractJoinFromJoinTableExpr(expr)
			}
		}
	}
}

func (p *Parser) extractJoinFromJoinTableExpr(expr *sqlparser.JoinTableExpr) {
	traversedJoinTables, traversedJoinConditions := traverseAST(expr)
	tableCache := map[string]string{}
	for _, t := range traversedJoinTables {
		p.alias[t.Alias] = t.Table
	}
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
	defaultTableMetadata := []ASTTraversalJoinTableMetadata{}
	defaultConditionMetadata := []ASTTraversalJoinConditionMetadata{}

	// definitely there is a condition here, extract the condition here :D
	defaultConditionMetadata = append(defaultConditionMetadata, extractASTTableConditionFromExpression(expr))

	switch leftExpr := expr.LeftExpr.(type) {
	case *sqlparser.AliasedTableExpr:
		defaultTableMetadata = append(defaultTableMetadata, extractASTTableMetadataFromExpression(leftExpr))
	case *sqlparser.JoinTableExpr:
		// if child itself is a join expression, go deeper :)
		childTables, childConditions := traverseAST(leftExpr)
		defaultTableMetadata = append(defaultTableMetadata, childTables...)
		defaultConditionMetadata = append(defaultConditionMetadata, childConditions...)
	}

	switch rightExpr := expr.RightExpr.(type) {
	case *sqlparser.AliasedTableExpr:
		defaultTableMetadata = append(defaultTableMetadata, extractASTTableMetadataFromExpression(rightExpr))
	case *sqlparser.JoinTableExpr:
		childTables, childConditions := traverseAST(rightExpr)
		defaultTableMetadata = append(defaultTableMetadata, childTables...)
		defaultConditionMetadata = append(defaultConditionMetadata, childConditions...)
	}
	return defaultTableMetadata, defaultConditionMetadata
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
