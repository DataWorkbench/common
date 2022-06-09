package flink

type Counts struct {
	Restored   int32 `json:"restored"`
	Total      int32 `json:"total"`
	InProgress int32 `json:"in_progress"`
	Completed  int32 `json:"completed"`
	Failed     int32 `json:"failed"`
}

type Latest struct {
}

type Completed struct {
	AlignmentBuffered int32  `json:"alignment_buffered"`
	CheckpointType    string `json:"checkpoint_type"`
	Discarded         bool   `json:"discarded"`
	EndToEndDuration  int32  `json:"end_to_end_duration"`
	ExternalPath      string `json:"external_path"`
	Id                int32  `json:"id"`
	IsSavepoint       bool   `json:"is_savepoint"`
}
