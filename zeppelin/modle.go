package zeppelin

import (
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
	ParagraphId string
	Status      Status
	Progress    int
	Results     []*Result
	JobUrls     []string
}

type Result struct {
	rType string
	rData string
}

func NewResult(jsonObj *fastjson.Value) *Result {
	rType := string(jsonObj.GetStringBytes("type"))
	rData := string(jsonObj.GetStringBytes("data"))
	return &Result{
		rType: rType,
		rData: rData,
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

type SessionInfo struct {
	SessionId   string
	NoteId      string
	Interpreter string
	State       string
	WebUrl      string
	StartTime   string
}

func NewSessionInfo(sessionStr string) (*SessionInfo, error) {
	sessionInfo := SessionInfo{}
	sessionObj, err := fastjson.Parse(sessionStr)
	body := sessionObj.Get("body")
	if err != nil {
		return nil, err
	}
	if strings.Contains(sessionStr, "sessionId") {
		sessionInfo.SessionId = string(body.GetStringBytes("sessionId"))
	}
	if strings.Contains(sessionStr, "noteId") {
		sessionInfo.NoteId = string(body.GetStringBytes("noteId"))
	}
	if strings.Contains(sessionStr, "interpreter") {
		sessionInfo.Interpreter = string(body.GetStringBytes("interpreter"))
	}
	if strings.Contains(sessionStr, "state") {
		sessionInfo.State = string(body.GetStringBytes("state"))
	}
	if strings.Contains(sessionStr, "weburl") {
		sessionInfo.WebUrl = string(body.GetStringBytes("weburl"))
	}
	if strings.Contains(sessionStr, "startTime") {
		sessionInfo.StartTime = string(body.GetStringBytes("startTime"))
	}
	return &sessionInfo, nil
}
