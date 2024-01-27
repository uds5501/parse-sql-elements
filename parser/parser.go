package parser

import (
	"log"
	"vitess.io/vitess/go/vt/sqlparser"
)

type Parser struct {
	joins      []Joins
	aggregates []Aggregations
	scans      []Scan
}

func NewParser() Parser {
	return Parser{
		joins:      []Joins{},
		aggregates: []Aggregations{},
		scans:      []Scan{},
	}
}

func (p *Parser) GetJoins() []Joins {
	return p.joins
}

func (p *Parser) GetScans() []Scan {
	return p.scans
}

func (p *Parser) GetAggregates() []Aggregations {
	return p.aggregates
}

func (p *Parser) ParseQuery(sql string) error {
	ast, _, err := sqlparser.Parse2(sql)
	if err != nil {
		log.Println("error while parsing ->", err)
		return err
	}
	p.extract(ast)

	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) extractFromWhere(where *sqlparser.Where) {
	switch expr := where.Expr.(type) {
	case *sqlparser.ComparisonExpr:
		p.extractFromComparisonExpr(expr)
	case *sqlparser.AndExpr:
		p.extractFromAndExpr(expr)
	case *sqlparser.OrExpr:
		p.extractFromOrExpr(expr)
	}

}
