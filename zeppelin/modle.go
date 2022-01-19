package zeppelin

//
//import (
//	"encoding/json"
//	"github.com/buger/jsonparser"
//)
//
//type Status string
//
//const (
//	UNKNOWN  Status = "UNKNOWN"
//	READY    Status = "READY"
//	PENDING  Status = "PENDING"
//	RUNNING  Status = "RUNNING"
//	FINISHED Status = "FINISHED"
//	ERROR    Status = "ERROR"
//	ABORT    Status = "ABORT"
//)
//
//func (s Status) IsUnknown() bool {
//	return UNKNOWN == s
//}
//
//func (s Status) IsReady() bool {
//	return READY == s
//}
//
//func (s Status) IsRunning() bool {
//	return RUNNING == s
//}
//
//func (s Status) IsPending() bool {
//	return PENDING == s
//}
//
//func (s Status) IsCompleted() bool {
//	return FINISHED == s || ERROR == s || ABORT == s
//}
//
//func (s Status) IsFinished() bool {
//	return FINISHED == s
//}
//
//func (s Status) IsFailed() bool {
//	return ERROR == s || ABORT == s
//}
//
//func valueOf(value string) Status {
//	switch value {
//	case "UNKNOWN":
//		return UNKNOWN
//	case "READY":
//		return READY
//	case "PENDING":
//		return PENDING
//	case "RUNNING":
//		return RUNNING
//	case "FINISHED":
//		return FINISHED
//	case "ERROR":
//		return ERROR
//	case "ABORT":
//		return ABORT
//	}
//	return UNKNOWN
//}
//
//type ParagraphResult struct {
//	NoteId      string    `json:"noteId"`
//	ParagraphId string    `json:"paragraphId"`
//	Status      Status    `json:"status"`
//	Progress    int64     `json:"progress"`
//	Results     []*Result `json:"results"`
//	JobUrls     []string  `json:"jobUrl"`
//	JobId       string    `json:"jobId"`
//}
//
//type Result struct {
//	Type string `json:"type"`
//	Data string `json:"data"`
//}
//
//func NewParagraphResult(value []byte) (*ParagraphResult, error) {
//	var err error
//	var paragraphResult ParagraphResult
//	if paragraphResult.ParagraphId, err = jsonparser.GetString(value, "body", "id"); err != nil {
//		return nil, err
//	}
//	status, err := jsonparser.GetString(value, "body", "status")
//	if err != nil {
//		return nil, err
//	}
//	paragraphResult.Status = valueOf(status)
//	if paragraphResult.Progress, err = jsonparser.GetInt(value, "body", "progress"); err != nil {
//		return nil, err
//	}
//
//	_, _ = jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
//		if err != nil {
//			return
//		}
//		var result Result
//		if err = json.Unmarshal(value, &result); err != nil {
//			return
//		}
//		paragraphResult.Results = append(paragraphResult.Results, &result)
//	}, "body", "results", "msg")
//
//	_, _ = jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
//		if err != nil {
//			return
//		}
//		jobUrl, err := jsonparser.GetString(value, "jobUrl")
//		if err != nil {
//			return
//		}
//		paragraphResult.JobUrls = append(paragraphResult.JobUrls, jobUrl)
//	}, "body", "runtimeInfos", "jobUrl", "values")
//
//	return &paragraphResult, nil
//}
//
////type ExecuteResult struct {
////	StatementId string
////	Status      Status
////	Results     []*Result
////	JobUrls     []string
////	JobId       string
////	Progress    int64
////	FlinkUrl    string
////	SessionInfo *SessionInfo
////}
////
////func NewExecuteResult(paragraphResult *ParagraphResult, sessionInfo *SessionInfo) *ExecuteResult {
////	return &ExecuteResult{
////		StatementId: paragraphResult.ParagraphId,
////		Status:      paragraphResult.Status,
////		Results:     paragraphResult.Results,
////		JobUrls:     paragraphResult.JobUrls,
////		Progress:    paragraphResult.Progress,
////		SessionInfo: sessionInfo,
////	}
////}
////
////type SessionInfo struct {
////	SessionId   string `json:"sessionId"`
////	NoteId      string `json:"noteId"`
////	Interpreter string `json:"interpreter"`
////	State       string `json:"state"`
////	WebUrl      string `json:"webUrl"`
////	StartTime   string `json:"starTime"`
////}
//
////func NewSessionInfo(session []byte) (*SessionInfo, error) {
////	sessionInfo := SessionInfo{}
////	body, _, _, err := jsonparser.Get(session, "body")
////	if err != nil {
////		return nil, err
////	}
////	if err = json.Unmarshal(body, &sessionInfo); err != nil {
////		return nil, err
////	}
////	return &sessionInfo, nil
////}
