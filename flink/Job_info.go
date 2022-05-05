package flink

type VertexTasks struct {
	DEPLOYING   int `json:"DEPLOYING"`
	SCHEDULED   int `json:"SCHEDULED"`
	FINISHED    int `json:"FINISHED"`
	CANCELING   int `json:"CANCELING"`
	CREATED     int `json:"CREATED"`
	CANCELED    int `json:"CANCELED"`
	RECONCILING int `json:"RECONCILING"`
	RUNNING     int `json:"RUNNING"`
	FAILED      int `json:"FAILED"`
}

type VertexMetrics struct {
	ReadBytes            int  `json:"read-bytes"`
	ReadBytesComplete    bool `json:"read-bytes-complete"`
	WriteBytes           int  `json:"write-bytes"`
	WriteBytesComplete   bool `json:"write-bytes-complete"`
	ReadRecords          int  `json:"read-records"`
	ReadRecordsComplete  bool `json:"read-records-complete"`
	WriteRecords         int  `json:"write-records"`
	WriteRecordsComplete bool `json:"write-records-complete"`
}

type Vertex struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Parallelism int            `json:"parallelism"`
	Status      string         `json:"status"`
	StartTime   int64          `json:"start-time"`
	EndTime     int            `json:"end-time"`
	Duration    int            `json:"duration"`
	Tasks       *VertexTasks   `json:"tasks"`
	Metrics     *VertexMetrics `json:"metrics"`
}

// Timestamps represents the timestamp of job state changed.
type Timestamps struct {
	CANCELLING   int   `json:"CANCELLING"`
	CANCELED     int   `json:"CANCELED"`
	RESTARTING   int   `json:"RESTARTING"`
	INITIALIZING int64 `json:"INITIALIZING"`
	CREATED      int64 `json:"CREATED"`
	RECONCILING  int   `json:"RECONCILING"`
	RUNNING      int64 `json:"RUNNING"`
	FAILING      int   `json:"FAILING"`
	FAILED       int   `json:"FAILED"`
	FINISHED     int   `json:"FINISHED"`
	SUSPENDED    int   `json:"SUSPENDED"`
}

// StatusCounts represents the task status count.
type StatusCounts struct {
	DEPLOYING   int `json:"DEPLOYING"`
	SCHEDULED   int `json:"SCHEDULED"`
	FINISHED    int `json:"FINISHED"`
	CANCELING   int `json:"CANCELING"`
	CREATED     int `json:"CREATED"`
	CANCELED    int `json:"CANCELED"`
	RECONCILING int `json:"RECONCILING"`
	RUNNING     int `json:"RUNNING"`
	FAILED      int `json:"FAILED"`
}

type NodeInput struct {
	Num          int    `json:"num"`
	Id           string `json:"id"`
	ShipStrategy string `json:"ship_strategy"`
	Exchange     string `json:"exchange"`
}

type PlanNode struct {
	Id               string       `json:"id"`
	Parallelism      int          `json:"parallelism"`
	Operator         string       `json:"operator"`
	OperatorStrategy string       `json:"operator_strategy"`
	Description      string       `json:"description"`
	Inputs           []*NodeInput `json:"inputs"`
}

type Plan struct {
	Jid   string      `json:"jid"`
	Name  string      `json:"name"`
	Nodes []*PlanNode `json:"nodes"`
}

// JobInfo represents the response of '/jobs/<jod_id>'
type JobInfo struct {
	Jid          string        `json:"jid"`
	Name         string        `json:"name"`
	IsStoppable  bool          `json:"isStoppable"`
	State        string        `json:"state"` // The job state. RUNNING...
	StartTime    int64         `json:"start-time"`
	EndTime      int           `json:"end-time"`
	Duration     int           `json:"duration"`
	Now          int64         `json:"now"`
	Timestamps   *Timestamps   `json:"timestamps"`
	Vertices     []*Vertex     `json:"vertices"`
	StatusCounts *StatusCounts `json:"status-counts"`
	Plan         *Plan         `json:"plan"`
}
