package zeppelin

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/valyala/fastjson"
)

type Status string

const (
	UNKNOWN  Status = "UNKNOWN"
	READY    Status = "READY"
	PENDING  Status = "PENDING"
	RUNNING  Status = "RUNNING"
	FINISHED Status = "FINISHED"
	ERROR    Status = "ERROR"
	ABORT    Status = "ABORT"
)

func (s Status) isReady() bool {
	return READY == s
}

func (s Status) isRunning() bool {
	return RUNNING == s
}

func (s Status) isPending() bool {
	return PENDING == s
}

func (s Status) isCompleted() bool {
	return FINISHED == s || ERROR == s || ABORT == s
}

func (s Status) isFinished() bool {
	return FINISHED == s
}

func valueOf(value string) Status {
	switch value {
	case "UNKNOWN":
		return UNKNOWN
	case "READY":
		return READY
	case "PENDING":
		return PENDING
	case "RUNNING":
		return RUNNING
	case "FINISHED":
		return FINISHED
	case "ERROR":
		return ERROR
	case "ABORT":
		return ABORT
	}
	return UNKNOWN
}

type ClientConfig struct {
	ZeppelinRestUrl string
	Timeout         time.Duration
	RetryCount      int
	QueryInterval   time.Duration
}

type ParagraphResult struct {
	ParagraphId string    `json:"paragraphId"`
	Status      Status    `json:"status"`
	Progress    int       `json:"progress"`
	Results     []*Result `json:"results"`
	JobUrls     []string  `json:"jobUrl"`
}

type Result struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func NewResult(jsonObj *fastjson.Value) *Result {
	rType := string(jsonObj.GetStringBytes("type"))
	rData := string(jsonObj.GetStringBytes("data"))
	return &Result{
		Type: rType,
		Data: rData,
	}
}

func NewParagraphResult(paragraphJson *fastjson.Value) *ParagraphResult {
	result := &ParagraphResult{}
	result.ParagraphId = string(paragraphJson.GetStringBytes("id"))
	result.Status = valueOf(string(paragraphJson.GetStringBytes("status")))
	result.Progress = paragraphJson.GetInt("progress")
	if strings.Contains(paragraphJson.String(), "results") {
		resultJson := paragraphJson.Get("results")
		msgArray := resultJson.GetArray("msg")
		for _, resultObj := range msgArray {
			result.Results = append(result.Results, NewResult(resultObj))
		}
	}

	if strings.Contains(paragraphJson.String(), "runtimeInfos") {
		runtimeInfosJson := paragraphJson.Get("runtimeInfos")
		if strings.Contains(runtimeInfosJson.String(), "jobUrl") {
			jobUrlJson := runtimeInfosJson.Get("jobUrl")
			if strings.Contains(jobUrlJson.String(), "values") {
				valuesArray := jobUrlJson.GetArray("values")
				for _, value := range valuesArray {
					result.JobUrls = append(result.JobUrls, string(value.GetStringBytes("jobUrl")))
				}
			}
		}
	}
	return result
}

type ExecuteResult struct {
	statementId string
	status      Status
	results     []*Result
	jobUrls     []string
	progress    int
}

func NewExecuteResult(paragraphResult *ParagraphResult) *ExecuteResult {
	return &ExecuteResult{
		statementId: paragraphResult.ParagraphId,
		status:      paragraphResult.Status,
		results:     paragraphResult.Results,
		jobUrls:     paragraphResult.JobUrls,
		progress:    paragraphResult.Progress,
	}
}

type SessionInfo struct {
	SessionId   string `json:"sessionId"`
	NoteId      string `json:"noteId"`
	Interpreter string `json:"interpreter"`
	State       string `json:"state"`
	WebUrl      string `json:"webUrl"`
	StartTime   string `json:"starTime"`
}

func NewSessionInfo(sessionStr string) (*SessionInfo, error) {
	sessionInfo := SessionInfo{}
	sessionObj, err := fastjson.Parse(sessionStr)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(sessionObj.Get("body").String()), &sessionInfo); err != nil {
		return nil, err
	}
	return &sessionInfo, nil
}
