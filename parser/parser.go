package parser

import (
	"log"
	"fmt"
	"errors"
	"vitess.io/vitess/go/vt/sqlparser"
	"github.com/kr/pretty"
)

type Parser struct {
	alias          map[string]string
	joins          []Joins
	aggregates     []Aggregations
	scans          []Scan
	existingTables []Table
}

func NewParser() Parser {
	return Parser{
		alias:      map[string]string{},
		joins:      []Joins{},
		aggregates: []Aggregations{},
		scans:      []Scan{},
	}
}

func (p *Parser) UpdateTables(t []Table) {
	p.existingTables = t
}

func (p *Parser) ParseQuery(sql string) error {
	ast, _, err := sqlparser.Parse2(sql)
	if err != nil {
		log.Println("error while parsing ->", err)
		return err
	}
	fmt.Println(pretty.Sprint(ast))
	p.extract(ast)

	fmt.Println("joins: ", p.joins)
	fmt.Println("aggregates: ", p.aggregates)
	fmt.Println("scans: ", p.scans)
	fmt.Println("alias: ", p.alias)
	err = p.verifyConfiguration()
	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) verifyConfiguration() error {
	// now let the parser cross verify whether the aliases match?
	for _, table := range p.alias {
		found := false
		for _, actualTable := range p.existingTables {
			if actualTable.Name == p.alias[table] {
				found = true
				break
			}
		}
		if !found {
			return errors.New(fmt.Sprintf("Table %s not found", table))
		}
	}

	// whether the scan columns match?
	for _, scan := range p.scans {
		foundAlias := false
		var existingTableName string
		// find actual table name
		for _, alias := range p.alias {
			if alias == scan.Table {
				foundAlias = true
				existingTableName = p.alias[alias]
				break
			}
		}
		if !foundAlias {
			for _, alias := range p.alias {
				if p.alias[alias] == scan.Table {
					foundAlias = true
					existingTableName = scan.Table
					break
				}
			}
		}
		if !foundAlias {
			return errors.New(fmt.Sprintf("Alias %s not found", scan.Table))
		}

		foundCol := false
		for _, table := range p.existingTables {
			if table.Name == existingTableName {
				for _, col := range table.Columns {
					if col.Name == scan.Column {
						foundCol = true
						break
					}
				}
			}
			if foundCol {
				break
			}
		}

		if !foundCol {
			return errors.New(fmt.Sprintf("Alias %s not found", scan.Table))
		}

	}
	//TODO: implement validations whether the aggregates match?

	return nil
}

func (p *Parser) extractFromWhere(where *sqlparser.Where) {
	fmt.Println("in where...")
	switch expr := where.Expr.(type) {
	case *sqlparser.ComparisonExpr:
		p.extractFromComparisonExpr(expr)
	case *sqlparser.AndExpr:
		p.extractFromAndExpr(expr)
	case *sqlparser.OrExpr:
		p.extractFromOrExpr(expr)
	}

}
