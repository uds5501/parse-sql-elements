package parser

import "vitess.io/vitess/go/vt/sqlparser"

func (p *Parser) extractAggregates(expr *sqlparser.AliasedExpr) {
	switch expr := expr.Expr.(type) {
	case *sqlparser.CountStar:
		p.aggregates = append(p.aggregates, Aggregations{AggregationType: AggregationTypeCount})
	case *sqlparser.Count:
		p.aggregates = append(p.aggregates, extractCountFromExpr(expr))
	case *sqlparser.Sum:
		p.aggregates = append(p.aggregates, extractSumFromExpr(expr))
	case *sqlparser.Avg:
		p.aggregates = append(p.aggregates, extractAvgFromExpr(expr))
	case *sqlparser.Max:
		p.aggregates = append(p.aggregates, extractMaxFromExpr(expr))
	case *sqlparser.Min:
		p.aggregates = append(p.aggregates, extractMinFromExpr(expr))
	}
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

func extractSumFromExpr(expr *sqlparser.Sum) Aggregations {
	aggregation := Aggregations{
		AggregationType: AggregationTypeSum,
	}
	switch col := expr.Arg.(type) {
	case *sqlparser.ColName:
		aggregation.Qualifier = col.Qualifier.Name.String()
		aggregation.Column = col.Name.String()
	}

	return aggregation
}

func extractAvgFromExpr(expr *sqlparser.Avg) Aggregations {
	aggregation := Aggregations{
		AggregationType: AggregationTypeAvg,
	}
	switch col := expr.Arg.(type) {
	case *sqlparser.ColName:
		aggregation.Qualifier = col.Qualifier.Name.String()
		aggregation.Column = col.Name.String()
	}
	return aggregation
}

func extractMinFromExpr(expr *sqlparser.Min) Aggregations {
	aggregation := Aggregations{
		AggregationType: AggregationTypeMin,
	}
	switch col := expr.Arg.(type) {
	case *sqlparser.ColName:
		aggregation.Qualifier = col.Qualifier.Name.String()
		aggregation.Column = col.Name.String()
	}
	return aggregation
}

func extractMaxFromExpr(expr *sqlparser.Max) Aggregations {
	aggregation := Aggregations{
		AggregationType: AggregationTypeMax,
	}
	switch col := expr.Arg.(type) {
	case *sqlparser.ColName:
		aggregation.Qualifier = col.Qualifier.Name.String()
		aggregation.Column = col.Name.String()
	}
	return aggregation
}
