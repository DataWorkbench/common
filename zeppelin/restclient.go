package zeppelin
//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/buger/jsonparser"
//	"io/ioutil"
//	"net/http"
//	"strings"
//	"time"
//
//	"github.com/DataWorkbench/common/qerror"
//
//	"github.com/gojek/heimdall/v7"
//	"github.com/gojek/heimdall/v7/httpclient"
//)
//
//type ClientConfig struct {
//	ZeppelinRestUrl string
//	Timeout         time.Duration
//	RetryCount      int
//	QueryInterval   time.Duration
//}
//
//type Client struct {
//	*httpclient.Client
//	ClientConfig ClientConfig
//}
//
//func NewZeppelinClient(config ClientConfig) *Client {
//	client := httpclient.NewClient(
//		httpclient.WithHTTPTimeout(config.Timeout),
//		httpclient.WithRetryCount(config.RetryCount),
//		httpclient.WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(time.Millisecond*10, time.Millisecond*50))),
//	)
//	return &Client{client, config}
//}
//
//func (c *Client) getBaseUrl() string {
//	return "http://" + c.ClientConfig.ZeppelinRestUrl + "/api"
//}
//
//func (c *Client) CreateNoteWithGroup(notePath string, defaultInterpreterGroup string) (string, error) {
//	var response *http.Response
//	var err error
//	defer func() {
//		if response != nil && response.Body != nil {
//			_ = response.Body.Close()
//		}
//	}()
//	reqObj := map[string]interface{}{}
//	reqObj["name"] = notePath
//	reqObj["defaultInterpreterGroup"] = defaultInterpreterGroup
//	reqBytes, err := json.Marshal(&reqObj)
//	if err != nil {
//		return "", err
//	}
//	response, err = c.Post(c.getBaseUrl()+"/notebook", strings.NewReader(string(reqBytes)), http.Header{})
//	if err != nil {
//		return "", err
//	}
//	body, err := checkResponse(response)
//	if err != nil {
//		return "", err
//	}
//	if err = checkBodyStatus(body); err != nil {
//		return "", err
//	}
//	return jsonparser.GetString(body, "body")
//}
//
//func (c *Client) CreateNote(notePath string) (string, error) {
//	return c.CreateNoteWithGroup(notePath, "")
//}
//
//func (c *Client) ListNotes() (map[string]string, error) {
//	var response *http.Response
//	var err error
//	defer func() {
//		if response != nil && response.Body != nil {
//			_ = response.Body.Close()
//		}
//	}()
//	response, err = c.Get(c.getBaseUrl()+"/notebook", http.Header{})
//	body, err := checkResponse(response)
//	if err != nil {
//		return nil, err
//	}
//	if err = checkBodyStatus(body); err != nil {
//		return nil, err
//	}
//	result := map[string]string{}
//	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
//		if err != nil {
//			return
//		}
//		id, err := jsonparser.GetString(value, "id")
//		if err != nil {
//			return
//		}
//		path, err := jsonparser.GetString(value, "path")
//		if err != nil {
//			return
//		}
//		if len(id) > 0 && len(path) > 0 {
//			result[path] = id
//		}
//	}, "body")
//	if err != nil {
//		return nil, err
//	}
//	return result, nil
//}
//
//func (c *Client) DeleteNote(noteId string) error {
//	var response *http.Response
//	var err error
//	defer func() {
//		if response != nil && response.Body != nil {
//			_ = response.Body.Close()
//		}
//	}()
//	response, err = c.Delete(c.getBaseUrl()+fmt.Sprintf("/notebook/%s", noteId), http.Header{})
//	if err != nil {
//		return err
//	}
//	body, err := checkResponse(response)
//	if err != nil {
//		return err
//	}
//	return checkBodyStatus(body)
//}
//
//func (c *Client) AddParagraph(noteId string, title string, text string) (string, error) {
//	var response *http.Response
//	var err error
//	defer func() {
//		if response != nil && response.Body != nil {
//			_ = response.Body.Close()
//		}
//	}()
//	reqObj := map[string]interface{}{}
//	reqObj["title"] = title
//	reqObj["text"] = text
//	reqBytes, err := json.Marshal(&reqObj)
//	if err != nil {
//		return "", err
//	}
//	response, err = c.Post(c.getBaseUrl()+fmt.Sprintf("/notebook/%s/paragraph", noteId), strings.NewReader(string(reqBytes)), http.Header{})
//	if err != nil {
//		return "", err
//	}
//	body, err := checkResponse(response)
//	if err != nil {
//		return "", err
//	}
//	if err = checkBodyStatus(body); err != nil {
//		return "", err
//	}
//	return jsonparser.GetString(body, "body")
//}
//
////func (c *Client) updateParagraph(noteId string, paragraphId string, title string, text string) error {
////	reqObj := map[string]string{}
////	reqObj["title"] = title
////	reqObj["text"] = text
////	reqBytes, err := json.Marshal(reqObj)
////	if err != nil {
////		return err
////	}
////	response, err := c.Put(c.getBaseUrl()+fmt.Sprintf("/notebook/%s/paragraph/%s", noteId, paragraphId), strings.NewReader(string(reqBytes)), http.Header{})
////	if err != nil {
////		return err
////	}
////	body, err := checkResponse(response)
////	if err != nil {
////		return err
////	}
////	return checkBodyStatus(body)
////}
//
////func (c *Client) SubmitParagraphWithSessionId(noteId string, paragraphId string, sessionId string) (*ParagraphResult, error) {
////	response, err := c.Post(c.getBaseUrl()+fmt.Sprintf("/notebook/job/%s/%s", noteId, paragraphId),
////		strings.NewReader(""), http.Header{})
////	if err != nil {
////		return nil, err
////	}
////	body, err := checkResponse(response)
////	if err != nil {
////		return nil, err
////	}
////	if err = checkBodyStatus(body); err != nil {
////		return nil, err
////	}
////	return c.QueryParagraphResult(noteId, paragraphId)
////}
//
//func (c *Client) SubmitParagraph(noteId string, paragraphId string) (*ParagraphResult, error) {
//	var response *http.Response
//	var err error
//	defer func() {
//		if response != nil && response.Body != nil {
//			_ = response.Body.Close()
//		}
//	}()
//	response, err = c.Post(c.getBaseUrl()+fmt.Sprintf("/notebook/job/%s/%s", noteId, paragraphId),
//		strings.NewReader(""), http.Header{})
//	if err != nil {
//		return nil, err
//	}
//	body, err := checkResponse(response)
//	if err != nil {
//		return nil, err
//	}
//	if err = checkBodyStatus(body); err != nil {
//		return nil, err
//	}
//	return c.QueryParagraphResult(noteId, paragraphId)
//}
//
//func (c *Client) SubmitWithProperties(interceptor string, secondIntp string,
//	noteId string, code string, properties map[string]string) (*ParagraphResult, error) {
//	builder := strings.Builder{}
//	builder.WriteString("%" + interceptor)
//	if len(secondIntp) > 0 {
//		builder.WriteString("." + secondIntp)
//	}
//	if properties != nil && len(properties) > 0 {
//		builder.WriteString("(")
//		var propertyStr []string
//		for k, v := range properties {
//			propertyStr = append(propertyStr, fmt.Sprintf("\"%s\"=\"%s\"", k, v))
//		}
//		builder.WriteString(strings.Join(propertyStr, ","))
//		builder.WriteString(")")
//	}
//	builder.WriteString(" " + code)
//	paragraphId, err := c.AddParagraph(noteId, "code", builder.String())
//	if err != nil {
//		return nil, err
//	}
//	paragraphResult, err := c.SubmitParagraph(noteId, paragraphId)
//	if err != nil {
//		return nil, err
//	}
//	return paragraphResult, nil
//}
//
//func (c *Client) Submit(interceptor string, secondIntp string, noteId string, code string) (*ParagraphResult, error) {
//	return c.SubmitWithProperties(interceptor, secondIntp, noteId, code, map[string]string{})
//}
//
////func (c *Client) executeParagraphWithSessionId(noteId string, paragraphId string, sessionId string) (*ParagraphResult, error) {
////	_, err := c.SubmitParagraphWithSessionId(noteId, paragraphId, sessionId)
////	if err != nil {
////		return nil, err
////	}
////	return c.waitUtilParagraphFinish(noteId, paragraphId)
////}
//
//func (c *Client) ExecuteParagraph(noteId string, paragraphId string) (*ParagraphResult, error) {
//	_, err := c.SubmitParagraph(noteId, paragraphId)
//	if err != nil {
//		return nil, err
//	}
//	return c.waitUtilParagraphFinish(noteId, paragraphId)
//}
//
//func (c *Client) Execute(interceptor string, secondIntp string, noteId string, code string) (*ParagraphResult, error) {
//	builder := strings.Builder{}
//	builder.WriteString("%" + interceptor)
//	if len(secondIntp) > 0 {
//		builder.WriteString("." + secondIntp)
//	}
//	builder.WriteString(" " + code)
//	paragraphId, err := c.AddParagraph(noteId, "code", builder.String())
//	if err != nil {
//		return nil, err
//	}
//	paragraphResult, err := c.ExecuteParagraph(noteId, paragraphId)
//	if err != nil {
//		return nil, err
//	}
//	return paragraphResult, nil
//}
//
//func (c *Client) CancelParagraph(noteId string, paragraphId string) error {
//	var response *http.Response
//	var err error
//	defer func() {
//		if response != nil && response.Body != nil {
//			_ = response.Body.Close()
//		}
//	}()
//	response, err = c.Delete(c.getBaseUrl()+fmt.Sprintf("/notebook/job/%s/%s", noteId, paragraphId), http.Header{})
//	if err != nil {
//		return err
//	}
//	body, err := checkResponse(response)
//	if err != nil {
//		return err
//	}
//	return checkBodyStatus(body)
//}
//
//func (c *Client) waitUtilParagraphRunning(noteId string, paragraphId string) (*ParagraphResult, error) {
//	for {
//		paragraphResult, err := c.QueryParagraphResult(noteId, paragraphId)
//		if err != nil {
//			return nil, err
//		}
//		if paragraphResult.Status.IsRunning() || paragraphResult.Status.IsFinished() {
//			return paragraphResult, nil
//		}
//		if paragraphResult.Status.IsFailed() || paragraphResult.Status.IsUnknown() {
//			return paragraphResult, nil
//		}
//		time.Sleep(time.Millisecond * c.ClientConfig.QueryInterval)
//	}
//}
//
//func (c *Client) waitUtilParagraphFinish(noteId string, paragraphId string) (*ParagraphResult, error) {
//	for {
//		paragraphResult, err := c.QueryParagraphResult(noteId, paragraphId)
//		if err != nil {
//			return nil, err
//		}
//		if paragraphResult.Status.IsCompleted() {
//			return paragraphResult, nil
//		}
//		time.Sleep(time.Millisecond * c.ClientConfig.QueryInterval)
//	}
//}
//
//func (c *Client) QueryParagraphResult(noteId string, paragraphId string) (*ParagraphResult, error) {
//	var response *http.Response
//	var err error
//	defer func() {
//		if response != nil && response.Body != nil {
//			_ = response.Body.Close()
//		}
//	}()
//	response, err = c.Get(c.getBaseUrl()+fmt.Sprintf("/notebook/%s/paragraph/%s", noteId, paragraphId), http.Header{})
//	if err != nil {
//		return nil, err
//	}
//	body, err := checkResponse(response)
//	if err != nil {
//		return nil, err
//	}
//	if err = checkBodyStatus(body); err != nil {
//		return nil, err
//	}
//	return NewParagraphResult(body)
//}
//
////func (c *Client) newSession(interpreter string) (*SessionInfo, error) {
////	response, err := c.Post(c.getBaseUrl()+fmt.Sprintf("/session%s", queryString("interpreter", interpreter)), strings.NewReader(""), http.Header{})
////	if err != nil {
////		return nil, err
////	}
////	body, err := checkResponse(response)
////	if err != nil {
////		return nil, err
////	}
////	if err = checkBodyStatus(body); err != nil {
////		return nil, err
////	}
////	return NewSessionInfo(body)
////}
//
////func (c *Client) stopSession(sessionId string) error {
////	response, err := c.Delete(c.getBaseUrl()+fmt.Sprintf("/session/%s", sessionId), http.Header{})
////	if err != nil {
////		return err
////	}
////	body, err := checkResponse(response)
////	if err != nil {
////		return err
////	}
////	return checkBodyStatus(body)
////}
//
////func (c *Client) getSession(sessionId string) (*SessionInfo, error) {
////	response, err := c.Get(c.getBaseUrl()+fmt.Sprintf("/session/%s", sessionId), http.Header{})
////	if err != nil {
////		return nil, err
////	}
////	if response.StatusCode == 404 {
////		body, err := ioutil.ReadAll(response.Body)
////		if err != nil {
////			return nil, err
////		}
////		if strings.Contains(string(body), "No such session") {
////			return nil, nil
////		}
////	}
////	body, err := checkResponse(response)
////	if err != nil {
////		return nil, err
////	}
////	if err = checkBodyStatus(body); err != nil {
////		return nil, err
////	}
////	return NewSessionInfo(body)
////}
//
////func queryString(name string, value string) (queryStr string) {
////	queryStr = "?" + name
////	if value != "" && len(value) > 0 {
////		queryStr = queryStr + "=" + value
////	}
////	return queryStr
////}
//
//func checkResponse(response *http.Response) ([]byte, error) {
//	if response.StatusCode == 302 {
//		return nil, qerror.InvalidateZeppelinUser
//	}
//	body, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		return nil, err
//	}
//	if response.StatusCode != 200 {
//		if body != nil && strings.Contains(string(body), "org.apache.zeppelin.notebook.exception.NotePathAlreadyExistsException") {
//			return nil, qerror.ZeppelinNoteAlreadyExists
//		}
//		return nil, qerror.CallZeppelinRestApiFailed.Format(response.StatusCode, response.Status, string(body))
//	}
//	return body, nil
//}
//
//func checkBodyStatus(resBody []byte) error {
//	status, err := jsonparser.GetString(resBody, "status")
//	if err != nil {
//		return err
//	}
//	if !strings.EqualFold("OK", status) {
//		message, err := jsonparser.GetString(resBody, "status")
//		if err != nil {
//			return err
//		}
//		return qerror.ZeppelinReturnStatusError.Format(message)
//	}
//	return nil
//}
