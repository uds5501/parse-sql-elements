package parser

import "vitess.io/vitess/go/vt/sqlparser"

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
