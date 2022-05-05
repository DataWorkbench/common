package flink

type AllExceptions struct {
	Exception string `json:"exception"`
	Task      string `json:"task"`
	Location  string `json:"location"`
	Timestamp int64  `json:"timestamp"`
}

type JobExceptions struct {
	RootException string           `json:"root-exception"`
	Timestamp     int64            `json:"timestamp"`
	AllExceptions []*AllExceptions `json:"all-exceptions"`
	Truncated     bool             `json:"truncated"`
}
