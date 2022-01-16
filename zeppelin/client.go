package zeppelin

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/common/web/ghttp"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	*ghttp.Client
	zeppelinUrl string
}

func NewClient(ctx context.Context, cfg *ghttp.ClientConfig, zeppelinUrl string) *Client {
	client := ghttp.NewClient(ctx, cfg)
	return &Client{client, zeppelinUrl}
}

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

func (s Status) IsUnknown() bool {
	return UNKNOWN == s
}

func (s Status) IsReady() bool {
	return READY == s
}

func (s Status) IsRunning() bool {
	return RUNNING == s
}

func (s Status) IsPending() bool {
	return PENDING == s
}

func (s Status) IsCompleted() bool {
	return FINISHED == s || ERROR == s || ABORT == s
}

func (s Status) IsFinished() bool {
	return FINISHED == s
}

func (s Status) IsFailed() bool {
	return ERROR == s || ABORT == s
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

type ParagraphResult struct {
	NoteId      string    `json:"noteId"`
	ParagraphId string    `json:"paragraphId"`
	Status      Status    `json:"status"`
	Progress    int64     `json:"progress"`
	Results     []*Result `json:"results"`
	JobUrls     []string  `json:"jobUrl"`
	JobId       string    `json:"jobId"`
}

type Result struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func NewParagraphResult(value []byte) (*ParagraphResult, error) {
	var err error
	var paragraphResult ParagraphResult
	if paragraphResult.ParagraphId, err = jsonparser.GetString(value, "body", "id"); err != nil {
		return nil, err
	}
	status, err := jsonparser.GetString(value, "body", "status")
	if err != nil {
		return nil, err
	}
	paragraphResult.Status = valueOf(status)
	if paragraphResult.Progress, err = jsonparser.GetInt(value, "body", "progress"); err != nil {
		return nil, err
	}

	_, _ = jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			return
		}
		var result Result
		if err = json.Unmarshal(value, &result); err != nil {
			return
		}
		paragraphResult.Results = append(paragraphResult.Results, &result)
	}, "body", "results", "msg")

	_, _ = jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			return
		}
		jobUrl, err := jsonparser.GetString(value, "jobUrl")
		if err != nil {
			return
		}
		paragraphResult.JobUrls = append(paragraphResult.JobUrls, jobUrl)
	}, "body", "runtimeInfos", "jobUrl", "values")

	return &paragraphResult, nil
}

func (c *Client) getBaseUrl() string {
	return "http://" + c.zeppelinUrl + "/api"
}

func (c *Client) createNoteWithGroup(ctx context.Context, notePath string, defaultIntpGroup string) (string, error) {
	var response *http.Response
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()
	reqParam, err := json.Marshal(map[string]string{"name": notePath, "defaultInterpreterGroup": defaultIntpGroup})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, c.getBaseUrl()+"/notebook", strings.NewReader(string(reqParam)))
	if err != nil {
		return "", err
	}
	response, err = c.Send(ctx, req)
	if err != nil {
		return "", err
	}
	res, err := checkResponse(response)
	if err != nil {
		return "", err
	}
	if err = checkBodyStatus(res); err != nil {
		return "", err
	}
	return jsonparser.GetString(res, "body")
}

func (c *Client) CreateNote(ctx context.Context, notePath string) (string, error) {
	return c.createNoteWithGroup(ctx, notePath, "")
}

func (c *Client) ListNotes(ctx context.Context) (map[string]string, error) {
	var response *http.Response
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()
	req, err := http.NewRequest(http.MethodGet, c.getBaseUrl()+"/notebook", strings.NewReader(""))
	if err != nil {
		return nil, nil
	}
	response, err = c.Send(ctx, req)
	if err != nil {
		return nil, err
	}
	res, err := checkResponse(response)
	if err != nil {
		return nil, err
	}
	if err = checkBodyStatus(res); err != nil {
		return nil, err
	}
	result := map[string]string{}
	_, err = jsonparser.ArrayEach(res, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			return
		}
		id, err := jsonparser.GetString(value, "id")
		if err != nil {
			return
		}
		path, err := jsonparser.GetString(value, "path")
		if err != nil {
			return
		}
		if len(id) > 0 && len(path) > 0 {
			result[path] = id
		}
	}, "body")
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) DeleteNote(ctx context.Context, noteId string) error {
	var response *http.Response
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/notebook/%s", c.getBaseUrl(), noteId), strings.NewReader(""))
	if err != nil {
		return err
	}
	response, err = c.Send(ctx, req)
	if err != nil {
		return err
	}
	res, err := checkResponse(response)
	if err != nil {
		return err
	}
	return checkBodyStatus(res)
}

func (c *Client) AddParagraph(ctx context.Context, noteId string, title string, text string) (string, error) {
	var response *http.Response
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()
	reqParam, err := json.Marshal(map[string]string{"title": title, "text": text})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/notebook/%s/paragraph", c.getBaseUrl(), noteId), strings.NewReader(string(reqParam)))
	if err != nil {
		return "", err
	}
	response, err = c.Send(ctx, req)
	if err != nil {
		return "", err
	}
	res, err := checkResponse(response)
	if err != nil {
		return "", err
	}
	if err = checkBodyStatus(res); err != nil {
		return "", err
	}
	return jsonparser.GetString(res, "body")
}

func (c *Client) SubmitParagraph(ctx context.Context, noteId string, paragraphId string) (*ParagraphResult, error) {
	var response *http.Response
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/notebook/job/%s/%s", c.getBaseUrl(), noteId, paragraphId), strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	response, err = c.Send(ctx, req)
	if err != nil {
		return nil, err
	}
	res, err := checkResponse(response)
	if err != nil {
		return nil, err
	}
	if err = checkBodyStatus(res); err != nil {
		return nil, err
	}
	return c.QueryParagraphResult(ctx, noteId, paragraphId)
}

func (c *Client) SubmitWithProperties(ctx context.Context, intp string, secondIntp string, noteId string, code string, properties map[string]string) (*ParagraphResult, error) {
	sb := strings.Builder{}
	sb.WriteString("%" + intp)
	if len(secondIntp) > 0 {
		sb.WriteString("." + secondIntp)
	}
	if properties != nil && len(properties) > 0 {
		sb.WriteString("(")
		var propStr []string
		for k, v := range properties {
			propStr = append(propStr, fmt.Sprintf("\"%s\"=\"%s\"", k, v))
		}
		sb.WriteString(strings.Join(propStr, ","))
		sb.WriteString(")")
	}
	sb.WriteString(" " + code)
	paragraphId, err := c.AddParagraph(ctx, noteId, "code", sb.String())
	if err != nil {
		return nil, err
	}
	return c.SubmitParagraph(ctx, noteId, paragraphId)
}

func (c *Client) Submit(ctx context.Context, intp string, secondIntp string, noteId string, code string) (*ParagraphResult, error) {
	return c.SubmitWithProperties(ctx, intp, secondIntp, noteId, code, map[string]string{})
}

func (c *Client) WaitUntilFinish(ctx context.Context, noteId string, paragraphId string, intervalMs time.Duration) (*ParagraphResult, error) {
	for {
		paragraphResult, err := c.QueryParagraphResult(ctx, noteId, paragraphId)
		if err != nil {
			return nil, err
		}
		if paragraphResult.Status.IsCompleted() {
			return paragraphResult, nil
		}
		time.Sleep(time.Millisecond * intervalMs)
	}
}

func (c *Client) ExecuteParagraph(ctx context.Context, noteId string, paragraphId string) (*ParagraphResult, error) {
	_, err := c.SubmitParagraph(ctx, noteId, paragraphId)
	if err != nil {
		return nil, err
	}
	return c.WaitUntilFinish(ctx, noteId, paragraphId, 10000)
}

func (c *Client) Execute(ctx context.Context, intp string, secondIntp string, noteId string, code string) (*ParagraphResult, error) {
	sb := strings.Builder{}
	sb.WriteString("%" + intp)
	if len(secondIntp) > 0 {
		sb.WriteString("." + secondIntp)
	}
	sb.WriteString(" " + code)
	paragraphId, err := c.AddParagraph(ctx, noteId, "code", sb.String())
	if err != nil {
		return nil, err
	}
	return c.ExecuteParagraph(ctx, noteId, paragraphId)
}

func (c *Client) QueryParagraphResult(ctx context.Context, noteId string, paragraphId string) (*ParagraphResult, error) {
	var response *http.Response
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/notebook/%s/paragraph/%s", c.getBaseUrl(), noteId, paragraphId), strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	response, err = c.Send(ctx, req)
	if err != nil {
		return nil, err
	}
	res, err := checkResponse(response)
	if err != nil {
		return nil, err
	}
	if err = checkBodyStatus(res); err != nil {
		return nil, err
	}
	return NewParagraphResult(res)
}

func checkResponse(response *http.Response) ([]byte, error) {
	if response.StatusCode == 302 {
		return nil, qerror.InvalidateZeppelinUser
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		if body != nil && strings.Contains(string(body), "org.apache.zeppelin.notebook.exception.NotePathAlreadyExistsException") {
			return nil, qerror.ZeppelinNoteAlreadyExists
		}
		return nil, qerror.CallZeppelinRestApiFailed.Format(response.StatusCode, response.Status, string(body))
	}
	return body, nil
}

func checkBodyStatus(resBody []byte) error {
	status, err := jsonparser.GetString(resBody, "status")
	if err != nil {
		return err
	}
	if !strings.EqualFold("OK", status) {
		message, err := jsonparser.GetString(resBody, "status")
		if err != nil {
			return err
		}
		return qerror.ZeppelinReturnStatusError.Format(message)
	}
	return nil
}
