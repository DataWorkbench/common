package constants

const (
	EmptyNode     = "Empty"
	ValuesNode    = "Values"
	ConstNode     = "Const"
	SourceNode    = "Source"
	DimensionNode = "Dimension"
	DestNode      = "Dest"
	OrderByNode   = "OrderBy"
	LimitNode     = "Limit"
	OffsetNode    = "Offset"
	FetchNode     = "Fetch"
	FilterNode    = "Filter"
	UnionNode     = "Union"
	ExceptNode    = "Except"
	IntersectNode = "Intersect"
	GroupByNode   = "GroupBy"
	HavingNode    = "Having"
	WindowNode    = "Window"
	JoinNode      = "Join"
	SqlNode       = "Sql"
	UDTFNode      = "UDTF"
	UDTTFNode     = "UDTTF"
	JarNode       = "Jar"
	ScalaNode     = "Scala"
	PythonNode    = "Python"

	JOIN          = "JOIN"
	LEFT_JOIN     = "LEFT JOIN"
	RIGHT_JOIN    = "RIGHT JOIN"
	FULL_OUT_JOIN = "FULL OUTER JOIN"
	CROSS_JOIN    = "CROSS JOIN"
	INTERVAL_JOIN = "INTERVAL JOIN"

	DISTINCT_ALL      = ""
	DISTINCT_DISTINCT = "DISTINCT"
)

type FlinkParagraphsInfo struct {
	Conf      string `json:"conf"`
	Depends   string `json:"depends"`
	ScalaUDF  string `json:"scalaudf"`
	PythonUDF string `json:"pythonudf"`
	MainRun   string `json:"mainrun"`
}
