package flink

type JobOverviewTasks struct {
	Total       int `json:"total"`
	Created     int `json:"created"`
	Scheduled   int `json:"scheduled"`
	Deploying   int `json:"deploying"`
	Running     int `json:"running"`
	Finished    int `json:"finished"`
	Canceling   int `json:"canceling"`
	Canceled    int `json:"canceled"`
	Failed      int `json:"failed"`
	Reconciling int `json:"reconciling"`
}

type JobOverview struct {
	Jid              string            `json:"jid"`
	Name             string            `json:"name"`
	State            string            `json:"state"` // The job state. RUNNING...
	StartTime        int64             `json:"start-time"`
	EndTime          int64             `json:"end-time"`
	Duration         int64             `json:"duration"`
	LastModification int64             `json:"last-modification"`
	Tasks            *JobOverviewTasks `json:"tasks"`
}

// JobsOverview represents the response of '/jobs/overview'
type JobsOverview struct {
	Jobs []*JobOverview
}
