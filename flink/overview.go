package flink

// Overview represents the response of '/overview'
type Overview struct {
	TaskManagers   int    `json:"taskmanagers"`
	SlotsTotal     int    `json:"slots-total"`
	SlotSavailable int    `json:"slots-available"`
	JobsRunning    int    `json:"jobs-running"`
	JobsFinished   int    `json:"jobs-finished"`
	JobsCancelled  int    `json:"jobs-cancelled"`
	JobsFailed     int    `json:"jobs-failed"`
	FlinkVersion   string `json:"flink-version"`
	FlinkCommit    string `json:"flink-commit"`
}
