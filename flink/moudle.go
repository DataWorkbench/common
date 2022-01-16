package flink
//
//import "time"
//
//type ClientConfig struct {
//	Timeout       time.Duration
//	RetryCount    int
//	QueryInterval time.Duration
//}
//
//type Job struct {
//	Jid              string `json:"jid"`
//	Name             string `json:"name"`
//	State            string `json:"state"`
//	StartTime        int64  `json:"start-time"`
//	EndTime          int64  `json:"end-time"`
//	Duration         int64  `json:"duration"`
//	LastModification int64  `json:"last-modification"`
//}
//
//type Savepoint struct {
//	RequestId string             `json:"request-id"`
//	Status    SavepointStatus    `json:"status"`
//	Operation SavepointOperation `json:"operation"`
//}
//
//type SavepointStatus struct {
//	Id string `json:"id"`
//}
//
//type SavepointOperation struct {
//	Location     string `json:"location"`
//	FailureCause string `json:"failure-cause"`
//}
