package parser

type AggregationType string

const (
	AggregationTypeCount AggregationType = "COUNT"
	AggregationTypeAvg   AggregationType = "AVG"
	AggregationTypeSum   AggregationType = "SUM"
	AggregationTypeMin   AggregationType = "MIN"
	AggregationTypeMax   AggregationType = "MAX"
)

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
	Qualifier       string
	AggregationType AggregationType
}

// Scan TODO: In documents, specify why we don't take operation type.
type Scan struct {
	Table     string
	Column    string
	Qualifier string
}
