package parser

type AggregationType string
type DataType string

const (
	AggregationTypeCount AggregationType = "COUNT"
	AggregationTypeAvg   AggregationType = "AVG"
)

const (
	DataTypeText     DataType = "TEXT"
	DataTypeInt      DataType = "INT"
	DataTypeSmallInt DataType = "SMALLINT"
	DataTypeBoolean  DataType = "BOOLEAN"
	DataTypeFloat    DataType = "FLOAT"
)

type Column struct {
	Name      string
	IsIndexed bool
	DataType  DataType
}

type Table struct {
	Name    string
	Columns []Column
	Rows    int
}

type JoinColMetadata struct {
	Column    string
	Table     string
	Qualifier string
}

type Joins struct {
	Inner JoinColMetadata
	Outer JoinColMetadata
}

type ASTTraversalJoinTableMetadata struct {
	Table string
	Alias string
}

type ASTTraversalJoinConditionMetadata struct {
	InnerTable  string
	InnerColumn string
	OuterTable  string
	OuterColumn string
}

type Aggregations struct {
	Table           string
	Column          string
	AggregationType AggregationType
}

// Scan TODO: In documents, specify why we don't take operation type.
type Scan struct {
	Table     string
	Column    string
	Qualifier string
}
