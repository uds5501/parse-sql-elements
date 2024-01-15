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

func newParser() Parser {
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
	p.extractJoins(ast)
	p.extractAggregates(ast)
	p.extractScans(ast)

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

//TODO: Take this responsibility away from parser and model engine accordingly.
func (p *Parser) GeneratePlans() {
	// Generate plans for I/O related configs only.
	// start with SCANS. (imagine every scan is binary tree.)
	// then go to JOINS. (try doing different combo of joins there :))

}
