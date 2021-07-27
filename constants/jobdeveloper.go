package constants

const (
	EmptyNode     = "Empty"
	UDFNode       = "UDF"
	ValuesNode    = "Values"
	ConstNode     = "Const"
	SourceNode    = "Source"
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
	UDTFNode      = "UDTF"   //don't print Upstream node
	UDTTFNode     = "UDTTF"  //don't print Upstream node
	ArraysNode    = "Arrays" //don't print Upstream node
	JarNode       = "Jar"
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
	Table  string   `json:"table"` //sourceID, parser as $qc$sot-0123456789012363$qc$
	Column []string `json:"column"`
	ID     string   `json:"id"`
}

type ValuesNodeProperty struct {
	Row []string `json:"row"`
}

type UDFNodeProperty struct {
	ID       string `json:"id"`
	FuncName string `json:"funcname"`
}

type ColumnAs struct {
	Field string `json:"field"`
	As    string `json:"as"`
}

type ConstNodeProperty struct {
	Column []ColumnAs `json:"column"`
}

type SourceNodeProperty struct {
	ID       string     `json:"id"`
	Table    string     `json:"table"`    //sourceid, $qc$som-xxxx$qc$
	Distinct string     `json:"distinct"` //Distinct or ALL
	Column   []ColumnAs `json:"column"`
}

type OrderByColumn struct {
	Field string `json:"field"`
	Order string `json:"order"` //asc or desc
}

type OrderByNodeProperty struct {
	Column []OrderByColumn `json:"column"`
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
	Filter string `json:"filter"`
	In     string `json:"in"`
	Exists string `json:"exists"`
}

type UnionNodeProperty struct {
	All string `json:"all"`
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

type UDTFNodeProperty struct {
	ID           string     `json:"id"`
	FuncName     string     `json:"funcname"`
	Args         string     `json:"args"`
	As           string     `json:"as"`
	SelectColumn []ColumnAs `json:"selectcolumn"`
	Column       []ColumnAs `json:"column"`
}

type JoinNodeProperty struct {
	Join       string     `json:"join"`
	Expression string     `json:"expression"`
	Column     []ColumnAs `json:"column"`
}

type ArraysNodeProperty struct {
	Args         string     `json:"args"`
	As           string     `json:"as"`
	SelectColumn []ColumnAs `json:"selectcolumn"`
	Column       []ColumnAs `json:"column"`
}

type UDTTFNodeProperty struct {
	ID       string     `json:"id"`
	FuncName string     `json:"funcname"`
	Args     string     `json:"args"`
	Column   []ColumnAs `json:"column"`
}

type JarNodeProperty struct {
	JarArgs    string `json:"jar_args"`  // allow regex `^[a-zA-Z0-9_/. ]+$`
	JarEntry   string `json:"jar_entry"` // allow regex `^[a-zA-Z0-9_/. ]+$`
	JarPath    string `json:"jar_path"`
	AccessKey  string `json:"accesskey"`
	SecretKey  string `json:"secretkey"`
	EndPoint   string `json:"endpoint"`
	HbaseHosts string `json:"hbasehosts"`
}

const (
	JOIN              = "JOIN"
	LEFT_JOIN         = "LEFT JOIN"
	RIGHT_JOIN        = "RIGHT JOIN"
	FULL_OUT_JOIN     = "FULL OUTER JOIN"
	CROSS_JOIN        = "CROSS JOIN"
	DISTINCT_ALL      = "ALL"
	DISTINCT_DISTINCT = "DISTINCT"
)

// these resources will delete when job finish
type JobResources struct {
	Jar      string `json:"jar"`
	JobID    string `json:"jobid"`
	EngineID string `json:"engineID"`
}

type JobElementFlink struct {
	ZeppelinConf      string         `json:"conf"`
	ZeppelinDepends   string         `json:"depends"`
	ZeppelinFuncScala string         `json:"funcscala"`
	ZeppelinMainRun   string         `json:"mainrun"`
	S3info            SourceS3Params `json:"s3"`
	HbaseHosts        string         `json:"hbasehosts"`
	Resources         JobResources   `json:"resource"`
}

type FlinkParagraphsInfo struct {
	Conf      string `json:"conf"`
	Depends   string `json:"depends"`
	FuncScala string `json:"funcscala"` //jar in zeppelin flink.conf
	MainRun   string `json:"mainrun"`
}

type JobFreeActionFlink struct {
	ZeppelinDeleteJar string `json:"zeppelindeletejar"`
}
