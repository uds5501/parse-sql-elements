package parser

import "vitess.io/vitess/go/vt/sqlparser"

func (p *Parser) extractAggregates(ast sqlparser.Statement) {
	switch stmt := ast.(type) {
	case *sqlparser.Select:
		p.extractAggregationsFromSelectStatement(stmt)
	}
}
func (p *Parser) extractAggregationsFromSelectStatement(stmt *sqlparser.Select) {
	aggregations := []Aggregations{}
	if stmt.SelectExprs != nil {
		for _, selectExpr := range stmt.SelectExprs {
			switch aliasedExpr := selectExpr.(type) {
			case *sqlparser.AliasedExpr:
				switch expr := aliasedExpr.Expr.(type) {
				case *sqlparser.ColName:
					continue
				case *sqlparser.Count:
					aggregations = append(aggregations, extractCountFromExpr(expr))
				case *sqlparser.Avg:
					aggregations = append(aggregations, extractAvgFromExpr(expr))
				}
			}
		}
	}
	p.aggregates = aggregations
}

func extractCountFromExpr(expr *sqlparser.Count) Aggregations {
	aggregation := Aggregations{
		AggregationType: AggregationTypeCount,
	}
	for _, e := range expr.Args {
		switch col := e.(type) {
		case *sqlparser.ColName:
			aggregation.Table = col.Qualifier.Name.String()
			aggregation.Column = col.Name.String()
		}
	}
	return aggregation
}

func extractAvgFromExpr(expr *sqlparser.Avg) Aggregations {
	aggregation := Aggregations{
		AggregationType: AggregationTypeAvg,
	}
	switch col := expr.Arg.(type) {
	case *sqlparser.ColName:
		aggregation.Table = col.Qualifier.Name.String()
		aggregation.Column = col.Name.String()
	}
	return aggregation
}
