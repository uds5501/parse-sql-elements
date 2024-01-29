# parse-sql-elements

Welcome to a small package leveraging [vitess](https://pkg.go.dev/vitess.io/vitess) that helps you list all the Joins,
Scans and Aggregates within your SQL query.

## Why should you use this package?

It's an easy API which helps you utilize all the different elements of your queries without having to do the hardwork of
parsing the query yourself and doing the hardwork of traversing the AST by yourself.

## How to use?

You can install this package using -

```shell
go install github.com/uds5501/parse-sql-elements
```

Demo program -

```go
package main

import (
	"fmt"

	"github.com/uds5501/parse-sql-elements/parser"
)

func main() {
	parser := parser.NewParser()

	sql := `SELECT departments.department_name, AVG(employees.salary) AS average_salary
        FROM employees JOIN departments ON employees.department_id = departments.department_id;`

	parser.ParseQuery(sql)
	fmt.Println(parser.GetJoins())
	fmt.Println(parser.GetAggregates())
}
```

This should give you the following output (as per release >= `v0.0.5`)

```shell
[{{department_id  employees} {department_id  departments}}]
[{ salary employees AVG}]
```

Note: Currently, I'd recommend you to create a new parser for each query you want to parse.
I'll be adding the support to reset the parser metadata in future releases.

### Query Support

Currently, the parser should be able to identify the joins, selects and aggregates from the following query styles:

| Query Type     | Aggregates Extracted      | Joins extracted              | Scans extracted |
|----------------|---------------------------|------------------------------|-----------------|
| Select query   | SUM, AVG, MIN, MAX, COUNT | INNER JOINS                  | ALL             |
| Subqueries     | SUM, AVG, MIN, MAX, COUNT | INNER JOINS                  | ALL             |
| Derived tables | To support                | To support                   | To support      |
| Update         | N/A                       | N/A                          | N/A             |
| Insert         | N/A                       | N/A                          | N/A             |

### Current stable version:
Current stable published version is `v0.0.5`, I'll be updating this as I keep supporting newer queries.

### Have a feature request?
If you have a feature request and there is something I like to implement in future versions, I'd recommend you to
* [email me](mailto:singhuddeshyaofficial@gmail.com) the requirements and why do you need the said
  feature.
* Or you can raise an issue and send a PR to implement the same!