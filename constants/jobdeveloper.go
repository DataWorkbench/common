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
	UDTFNode      = "UDTF"  //don't print Upstream node
	UDTTFNode     = "UDTTF" //don't print Upstream node
	JarNode       = "Jar"
	ScalaNode     = "Scala"
	PythonNode    = "Python"
)

type DagNode struct {
	NodeType      string     `json:"nodetype"`
	NodeID        string     `json:"nodeid"`
	NodeName      string     `json:"nodename"`
	UpStream      string     `json:"upstream"`
	UpStreamRight string     `json:"upstreamright"`
	DownStream    string     `json:"downstream"`
	PointX        string     `json:"pointx"` //computer have different pixel, so this is the percent of max pixel
	PointY        string     `json:"pointy"`
	Property      JSONString `json:"property"`
}

type DestNodeProperty struct {
	Column []string `json:"column"`
	ID     string   `json:"id"`
}

type ValuesType struct {
	Values []string `json:"values"`
}

type ValuesNodeProperty struct {
	Rows []ValuesType `json:"rows"`
}

type ColumnAs struct {
	Field       string `json:"field"`
	Func        string `json:"func"`
	WindowsName string `json:"windowsname"`
	Type        string `json:"type"`
	As          string `json:"as"`
}

type ConstNodeProperty struct {
	Table  string     `json:"table"`
	Column []ColumnAs `json:"column"`
}

type SourceNodeProperty struct {
	ID           string     `json:"id"`
	TableAS      string     `json:"table"`
	Distinct     string     `json:"distinct"`
	Column       []ColumnAs `json:"column"`
	CustomColumn []ColumnAs `json:"customcolumn"`
}

type DimensionNodeProperty struct {
	ID           string     `json:"id"`
	TableAS      string     `json:"table"`    //not null.
	Distinct     string     `json:"distinct"` //Distinct or ALL
	Column       []ColumnAs `json:"column"`
	CustomColumn []ColumnAs `json:"customcolumn"`
	TimeColumu   []ColumnAs `json:"timecolumn"`
}

type OrderByColumn struct {
	Field string `json:"field"`
	Order string `json:"order"` //asc or desc.
}

type OrderByNodeProperty struct {
	Column []OrderByColumn `json:"column"` // order column must have serial number. source table selected column.
}

type LimitNodeProperty struct {
	Limit int32 `json:"limit"`
}

type OffsetNodeProperty struct {
	Offset int32 `json:"offset"`
}

type FetchNodeProperty struct {
	Fetch int32 `json:"fetch"`
}

type FilterNodeProperty struct {
	Where  string `json:"where"`  //left/right node
	In     string `json:"in"`     //left column
	Exists string `json:"exists"` //left column
}

type UnionNodeProperty struct {
	All bool `json:"all"`
}

type GroupByNodeProperty struct {
	Groupby []string `json:"groupby"`
}

type HavingNodeProperty struct {
	Having string `json:"having"`
}

type WindowNodeItem struct {
	Name string `json:"name"`
	Spec string `json:"spec"`
}

type WindowNodeProperty struct {
	Window []WindowNodeItem `json:"window"`
}

type SqlNodeProperty struct {
	Sql string `json:"sql"`
}

type ScalaNodeProperty struct {
	Code string `json:"code"`
}

type PythonNodeProperty struct {
	Code string `json:"code"`
}

type UDTFNodeProperty struct {
	ID           string     `json:"id"`
	Args         string     `json:"args"`
	TableAs      string     `json:"tableas"`
	SelectColumn []ColumnAs `json:"selectcolumn"`
	Column       []ColumnAs `json:"column"`
}

type JoinNodeProperty struct {
	Join           string     `json:"join"`
	Expression     string     `json:"expression"`
	Column         []ColumnAs `json:"column"`
	TableAs        string     `json:"tableas"`
	TableAsRight   string     `json:"tableasright"`
	Args           string     `json:"args"`
	GenerateColumn []ColumnAs `json:"generatecolumn"`
}

type UDTTFNodeProperty struct {
	ID       string     `json:"id"`
	FuncName string     `json:"funcname"`
	Args     string     `json:"args"`
	Column   []ColumnAs `json:"column"`
}

type JarNodeProperty struct {
	JarId     string `json:"jar_id"`
	JarArgs   string `json:"jar_args"`  // allow regex `^[a-zA-Z0-9_/. ]+$`
	JarEntry  string `json:"jar_entry"` // allow regex `^[a-zA-Z0-9_/. ]+$`
	AccessKey string `json:"accesskey"`
	SecretKey string `json:"secretkey"`
	EndPoint  string `json:"endpoint"`
	//TODO HbaseHosts []HostType `json:"hbasehosts"`
}

const (
	JOIN              = "JOIN"
	LEFT_JOIN         = "LEFT JOIN"
	RIGHT_JOIN        = "RIGHT JOIN"
	FULL_OUT_JOIN     = "FULL OUTER JOIN"
	CROSS_JOIN        = "CROSS JOIN"
	INTERVAL_JOIN     = "WHERE"
	DISTINCT_ALL      = ""
	DISTINCT_DISTINCT = "DISTINCT"
)

// these resources will delete when job finish
type JobResources struct {
	Jar      string `json:"jar"`
	JobID    string `json:"jobid"`
	EngineID string `json:"engineID"`
}

type JobElementFlink struct {
	ZeppelinConf      string `json:"conf"`
	ZeppelinDepends   string `json:"depends"`
	ZeppelinScalaUDF  string `json:"scalaudf"`
	ZeppelinPythonUDF string `json:"pythonudf"`
	ZeppelinMainRun   string `json:"mainrun"`
	//TODOS3info            SourceS3Params `json:"s3"`
	//TODO HbaseHosts        []HostType     `json:"hbasehosts"`
	Resources JobResources `json:"resource"`
}

type FlinkParagraphsInfo struct {
	Conf      string `json:"conf"`
	Depends   string `json:"depends"`
	ScalaUDF  string `json:"scalaudf"`  //jar in zeppelin flink.conf
	PythonUDF string `json:"pythonudf"` //jar in zeppelin flink.conf
	MainRun   string `json:"mainrun"`
}

type JobFreeActionFlink struct {
	ZeppelinDeleteJar string `json:"zeppelindeletejar"`
}
