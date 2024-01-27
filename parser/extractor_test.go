package parser

import (
	"testing"

	"vitess.io/vitess/go/vt/sqlparser"
	"github.com/stretchr/testify/assert"
	"fmt"
	"github.com/kr/pretty"
)

func TestJoinExtract(t *testing.T) {
	sql := `
SELECT
    t1.columnA,
    t2.columnB,
    t3.columnC
FROM
    table1 t1
    INNER JOIN table2 t2 ON t1.id = t2.table1_id
    INNER JOIN table3 t3 ON t2.id = t3.table2_id
    INNER JOIN table4 t4 ON t3.id = t4.table3_id
    INNER JOIN table5 t5 ON t4.id = t5.table4_id
ORDER BY
    t3.columnC DESC, t2.columnB ASC
LIMIT 10;
`
	ast, _, _ := sqlparser.Parse2(sql)
	expectedJoins := []Joins{
		{
			Inner: JoinColMetadata{
				Column:    "id",
				Table:     "",
				Qualifier: "t4",
			},
			Outer: JoinColMetadata{
				Column:    "table4_id",
				Table:     "",
				Qualifier: "t5",
			},
		},
		{
			Inner: JoinColMetadata{
				Column:    "id",
				Table:     "",
				Qualifier: "t3",
			},
			Outer: JoinColMetadata{
				Column:    "table3_id",
				Table:     "",
				Qualifier: "t4",
			},
		},
		{
			Inner: JoinColMetadata{
				Column:    "id",
				Table:     "",
				Qualifier: "t2",
			},
			Outer: JoinColMetadata{
				Column:    "table2_id",
				Table:     "",
				Qualifier: "t3",
			},
		},
		{
			Inner: JoinColMetadata{
				Column:    "id",
				Table:     "",
				Qualifier: "t1",
			},
			Outer: JoinColMetadata{
				Column:    "table1_id",
				Table:     "",
				Qualifier: "t2",
			},
		},
	}
	p := NewParser()
	p.extract(ast)
	assert.ElementsMatch(t, p.joins, expectedJoins)
}

func TestJoinExtractWithinSelect(t *testing.T) {
	sql := `
SELECT *
FROM
    Parcels P
	INNER JOIN Users U on P.user_id = U.id
WHERE id IN 
(SELECT parcel_id FROM ParcelList PL 
WHERE PL.parcel_id = P.id);
`
	ast, _, _ := sqlparser.Parse2(sql)
	expectedJoins := []Joins{
		{
			Inner: JoinColMetadata{
				Column:    "parcel_id",
				Table:     "",
				Qualifier: "PL",
			},
			Outer: JoinColMetadata{
				Column:    "id",
				Table:     "",
				Qualifier: "P",
			},
		},
		{
			Inner: JoinColMetadata{
				Column:    "user_id",
				Table:     "",
				Qualifier: "P",
			},
			Outer: JoinColMetadata{
				Column:    "id",
				Table:     "",
				Qualifier: "U",
			},
		},
	}
	p := NewParser()
	p.extract(ast)
	fmt.Println(pretty.Sprint(ast))
	assert.ElementsMatch(t, p.joins, expectedJoins)
}

func TestJoinExtractWithOrClause(t *testing.T) {
	sql := `
SELECT *
FROM
    Parcels P
	INNER JOIN Users U on P.user_id = U.id
WHERE id = 100 OR  
(SELECT parcel_id FROM ParcelList PL 
WHERE PL.parcel_id = P.id) = check_id;
`
	ast, _, _ := sqlparser.Parse2(sql)
	expectedJoins := []Joins{
		{
			Inner: JoinColMetadata{
				Column:    "parcel_id",
				Table:     "",
				Qualifier: "PL",
			},
			Outer: JoinColMetadata{
				Column:    "id",
				Table:     "",
				Qualifier: "P",
			},
		},
		{
			Inner: JoinColMetadata{
				Column:    "user_id",
				Table:     "",
				Qualifier: "P",
			},
			Outer: JoinColMetadata{
				Column:    "id",
				Table:     "",
				Qualifier: "U",
			},
		},
	}
	p := NewParser()
	p.extract(ast)
	assert.ElementsMatch(t, p.joins, expectedJoins)
}

func TestJoinsWithSubqueries(t *testing.T) {

	testCases := []struct {
		Name          string
		Sql           string
		ExpectedJoins []Joins
	}{
		{
			"Subquery with AND clause",
			`SELECT U.id, P.user_id
				FROM
					Parcels P
					INNER JOIN Users U on P.user_id = U.id
				WHERE id IN
				(
				   SELECT parcel_id
				   FROM ParcelList PL
					   INNER JOIN sample_table T on T.pl_id = PL.id
				WHERE PL.parcel_id = P.id AND PL.t_id = T.id);`,
			[]Joins{
				{
					Inner: JoinColMetadata{
						Column:    "user_id",
						Table:     "",
						Qualifier: "P",
					},
					Outer: JoinColMetadata{
						Column:    "id",
						Table:     "",
						Qualifier: "U",
					},
				},
				{
					Inner: JoinColMetadata{
						Column:    "pl_id",
						Table:     "",
						Qualifier: "T",
					},
					Outer: JoinColMetadata{
						Column:    "id",
						Table:     "",
						Qualifier: "PL",
					},
				},
				{
					Inner: JoinColMetadata{
						Column:    "parcel_id",
						Table:     "",
						Qualifier: "PL",
					},
					Outer: JoinColMetadata{
						Column:    "id",
						Table:     "",
						Qualifier: "P",
					},
				},
				{
					Inner: JoinColMetadata{
						Column:    "t_id",
						Table:     "",
						Qualifier: "PL",
					},
					Outer: JoinColMetadata{
						Column:    "id",
						Table:     "",
						Qualifier: "T",
					},
				},
			},
		},
		{
			"Subquery with INNER JOIN only",
			`SELECT U.id, P.user_id
				FROM
					Parcels P
					INNER JOIN Users U on P.user_id = U.id
				WHERE id IN
				(
				   SELECT parcel_id
				   FROM ParcelList PL
					   INNER JOIN sample_table T on T.pl_id = PL.id)
				AND P.user_id = U.x_id;`,
			[]Joins{
				{
					Inner: JoinColMetadata{
						Column:    "user_id",
						Table:     "",
						Qualifier: "P",
					},
					Outer: JoinColMetadata{
						Column:    "id",
						Table:     "",
						Qualifier: "U",
					},
				},
				{
					Inner: JoinColMetadata{
						Column:    "pl_id",
						Table:     "",
						Qualifier: "T",
					},
					Outer: JoinColMetadata{
						Column:    "id",
						Table:     "",
						Qualifier: "PL",
					},
				},
				{
					Inner: JoinColMetadata{
						Column:    "user_id",
						Table:     "",
						Qualifier: "P",
					},
					Outer: JoinColMetadata{
						Column:    "x_id",
						Table:     "",
						Qualifier: "U",
					},
				},
			},
		},
		{
			"Subquery with simple join",
			`SELECT *
				FROM
					Parcels P
				WHERE id =  
				(SELECT parcel_id FROM ParcelList PL 
				WHERE PL.parcel_id = P.id);`,
			[]Joins{
				{
					Inner: JoinColMetadata{
						Column:    "parcel_id",
						Table:     "",
						Qualifier: "PL",
					},
					Outer: JoinColMetadata{
						Column:    "id",
						Table:     "",
						Qualifier: "P",
					},
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			fmt.Println("SQL: ", testCase.Sql)
			ast, _, err := sqlparser.Parse2(testCase.Sql)
			if err != nil {
				fmt.Println(err)
			}
			p := NewParser()
			p.extract(ast)
			assert.ElementsMatch(t, p.joins, testCase.ExpectedJoins)
		})
	}
}

func TestScans(t *testing.T) {
	testCases := []struct {
		Name          string
		Sql           string
		ExpectedScans []Scan
	}{
		{
			Name: "Basic Scan test",
			Sql: `SELECT *
					FROM
						Parcels P
						INNER JOIN Users U on P.user_id = U.id
					WHERE U.id = 100 OR  
					(SELECT parcel_id FROM ParcelList PL 
					WHERE PL.parcel_id = P.id) = U.check_id;`,
			ExpectedScans: []Scan{
				{
					Column:    "id",
					Qualifier: "U",
				},
				{
					Column:    "check_id",
					Qualifier: "U",
				},
			},
		},
		{
			Name: "Scan with sub queries",
			Sql: `SELECT *
					FROM
						Parcels P
						INNER JOIN Users U on P.user_id = U.id
					WHERE U.id = 100 OR  
					(SELECT parcel_id FROM ParcelList PL 
					WHERE PL.parcel_id = P.id AND P.x = 90) = U.check_id;`,
			ExpectedScans: []Scan{
				{
					Column:    "id",
					Qualifier: "U",
				},
				{
					Column:    "check_id",
					Qualifier: "U",
				},
				{
					Column:    "x",
					Qualifier: "P",
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			fmt.Println("SQL: ", testCase.Sql)
			ast, _, err := sqlparser.Parse2(testCase.Sql)
			if err != nil {
				fmt.Println(err)
			}
			p := NewParser()
			p.extract(ast)
			assert.ElementsMatch(t, p.scans, testCase.ExpectedScans)
		})
	}
}
