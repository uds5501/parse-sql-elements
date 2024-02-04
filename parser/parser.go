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

func (p *Parser) Reset() {
	p.aggregates = []Aggregations{}
	p.scans = []Scan{}
	p.joins = []Joins{}
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
