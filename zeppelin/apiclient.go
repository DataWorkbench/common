package zeppelin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/DataWorkbench/common/qerror"

	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/valyala/fastjson"
)

type Client struct {
	*httpclient.Client
	ClientConfig ClientConfig
}

func NewZeppelinClient(config ClientConfig) *Client {
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(config.Timeout),
		httpclient.WithRetryCount(config.RetryCount),
		httpclient.WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(time.Millisecond*10, time.Millisecond*50))),
	)
	return &Client{client, config}
}

func (z *Client) getBaseUrl() string {
	return z.ClientConfig.ZeppelinRestUrl + "/api"
}

func (z *Client) createNoteWithGroup(notePath string, defaultInterpreterGroup string) (string, error) {
	reqObj := map[string]interface{}{}
	reqObj["name"] = notePath
	reqObj["defaultInterpreterGroup"] = defaultInterpreterGroup
	reqBytes, err := json.Marshal(&reqObj)
	if err != nil {
		return "", err
	}
	response, err := z.Post(z.getBaseUrl()+"/notebook", strings.NewReader(string(reqBytes)), http.Header{})
	if err != nil {
		return "", err
	}
	body, err := checkResponse(response)
	if err != nil {
		return "", err
	}
	if err = checkBodyStatus(body); err != nil {
		return "", err
	}
	bodyJson, err := fastjson.Parse(body)
	if err != nil {
		return "", err
	}
	return string(bodyJson.GetStringBytes("body")), nil
}

func (z *Client) createNote(notePath string) (string, error) {
	return z.createNoteWithGroup(notePath, "")
}

func (z *Client) deleteNote(noteId string) error {
	response, err := z.Delete(z.getBaseUrl()+fmt.Sprintf("/notebook/%s", noteId), http.Header{})
	if err != nil {
		return err
	}
	body, err := checkResponse(response)
	if err != nil {
		return err
	}
	if err = checkBodyStatus(body); err != nil {
		return err
	}
	return nil
}

func (z *Client) addParagraph(noteId string, title string, text string) (string, error) {
	reqObj := map[string]interface{}{}
	reqObj["title"] = title
	reqObj["text"] = text
	reqBytes, err := json.Marshal(&reqObj)
	if err != nil {
		return "", err
	}
	response, err := z.Post(z.getBaseUrl()+fmt.Sprintf("/notebook/%s/paragraph", noteId), strings.NewReader(string(reqBytes)), http.Header{})
	if err != nil {
		return "", err
	}
	body, err := checkResponse(response)
	if err != nil {
		return "", err
	}
	if err = checkBodyStatus(body); err != nil {
		return "", err
	}
	bodyJson, err := fastjson.Parse(body)
	if err != nil {
		return "", err
	}
	return string(bodyJson.GetStringBytes("body")), nil
}

func (z *Client) submitParagraphWithAll(noteId string, paragraphId string, sessionId string, parameters map[string]string) (*ParagraphResult, error) {
	reqObj := map[string]interface{}{}
	parametersJson, err := json.Marshal(parameters)
	if err != nil {
		return nil, err
	}
	// TODO the params is zeppelin web params,
	reqObj["params"] = string(parametersJson)
	//reqBytes, err := json.Marshal(reqObj)
	if err != nil {
		return nil, err
	}
	url := z.getBaseUrl() + fmt.Sprintf("/notebook/job/%s/%s%s", noteId, paragraphId, queryString("sessionId", sessionId))
	fmt.Println(url)
	response, err := z.Post(z.getBaseUrl()+fmt.Sprintf("/notebook/job/%s/%s%s", noteId, paragraphId, queryString("sessionId", sessionId)),
		strings.NewReader(""), http.Header{})
	if err != nil {
		return nil, err
	}
	body, err := checkResponse(response)
	if err != nil {
		return nil, err
	}
	if err = checkBodyStatus(body); err != nil {
		return nil, err
	}
	return z.queryParagraphResult(noteId, paragraphId)
}

func (z *Client) submitParagraphWithSessionId(noteId string, paragraphId string, sessionId string) (*ParagraphResult, error) {
	return z.submitParagraphWithAll(noteId, paragraphId, sessionId, make(map[string]string))
}

func (z *Client) submitParagraphWithParameters(noteId string, paragraphId string, parameters map[string]string) (*ParagraphResult, error) {
	return z.submitParagraphWithAll(noteId, paragraphId, "", parameters)
}

func (z *Client) submitParagraph(noteId string, paragraphId string) (*ParagraphResult, error) {
	return z.submitParagraphWithAll(noteId, paragraphId, "", make(map[string]string))
}

func (z *Client) executeParagraphWithAll(noteId string, paragraphId string, sessionId string, parameters map[string]string) (*ParagraphResult, error) {
	_, err := z.submitParagraphWithAll(noteId, paragraphId, sessionId, parameters)
	if err != nil {
		return nil, err
	}
	return z.waitUtilParagraphFinish(noteId, paragraphId)
}

func (z *Client) executeParagraphWithSessionId(noteId string, paragraphId string, sessionId string) (*ParagraphResult, error) {
	return z.executeParagraphWithAll(noteId, paragraphId, sessionId, make(map[string]string))
}

func (z *Client) executeParagraph(noteId string, paragraphId string) (*ParagraphResult, error) {
	return z.executeParagraphWithAll(noteId, paragraphId, "", make(map[string]string))
}

func (z *Client) cancelParagraph(noteId string, paragraphId string) error {
	response, err := z.Delete(z.getBaseUrl()+fmt.Sprintf("/notebook/job/%s/%s", noteId, paragraphId), http.Header{})
	if err != nil {
		return err
	}
	body, err := checkResponse(response)
	if err != nil {
		return err
	}
	if err = checkBodyStatus(body); err != nil {
		return err
	}
	return nil
}

func (z *Client) waitUtilParagraphRunning(noteId string, paragraphId string) (*ParagraphResult, error) {
	for {
		paragraphResult, err := z.queryParagraphResult(noteId, paragraphId)
		if err != nil {
			return nil, err
		}
		if paragraphResult.Status.isRunning() {
			return paragraphResult, nil
		}
		time.Sleep(time.Millisecond * z.ClientConfig.QueryInterval)
	}
}

func (z *Client) waitUtilParagraphJobUrlReturn(noteId string, paragraphId string) (*ParagraphResult, error) {
	paragraphResult := &ParagraphResult{}
	for len(paragraphResult.JobUrls) == 0 {
		result, err := z.waitUtilParagraphRunning(noteId, paragraphId)
		if err != nil {
			return nil, err
		}
		paragraphResult = result
		time.Sleep(time.Millisecond * z.ClientConfig.QueryInterval)
	}
	return paragraphResult, nil
}

func (z *Client) waitUtilParagraphFinish(noteId string, paragraphId string) (*ParagraphResult, error) {
	for {
		paragraphResult, err := z.queryParagraphResult(noteId, paragraphId)
		if err != nil {
			return nil, err
		}
		if paragraphResult.Status.isCompleted() {
			return paragraphResult, nil
		}
		time.Sleep(time.Millisecond * z.ClientConfig.QueryInterval)
	}
}

func (z *Client) waitUtilParagraphFinishWithTimeout(noteId string, paragraphId string, timeoutInSec int) (*ParagraphResult, error) {
	start := time.Now().Second()
	for (time.Now().Second() - start) < timeoutInSec {
		paragraphResult, err := z.queryParagraphResult(noteId, paragraphId)
		if err != nil {
			return nil, err
		}
		if paragraphResult.Status.isCompleted() {
			return paragraphResult, nil
		}
		time.Sleep(time.Millisecond * z.ClientConfig.QueryInterval)
	}
	return nil, qerror.ZeppelinRunParagraphTimeout.Format(timeoutInSec)
}

func (z *Client) queryParagraphResult(noteId string, paragraphId string) (*ParagraphResult, error) {
	response, err := z.Get(z.getBaseUrl()+fmt.Sprintf("/notebook/%s/paragraph/%s", noteId, paragraphId), http.Header{})
	if err != nil {
		return nil, err
	}
	body, err := checkResponse(response)
	if err != nil {
		return nil, err
	}
	if err = checkBodyStatus(body); err != nil {
		return nil, err
	}
	jsonBody, err := fastjson.Parse(body)
	if err != nil {
		return nil, err
	}
	paragraphJson := jsonBody.Get("body")
	return NewParagraphResult(paragraphJson), nil
}

func (z *Client) newSession(interpreter string) (*SessionInfo, error) {
	response, err := z.Post(z.getBaseUrl()+fmt.Sprintf("/session%s", queryString("interpreter", interpreter)), strings.NewReader(""), http.Header{})
	if err != nil {
		return nil, err
	}
	body, err := checkResponse(response)
	if err != nil {
		return nil, err
	}
	if err = checkBodyStatus(body); err != nil {
		return nil, err
	}
	return NewSessionInfo(body)
}

func (z *Client) stopSession(sessionId string) error {
	response, err := z.Delete(z.getBaseUrl()+fmt.Sprintf("/session/%s", sessionId), http.Header{})
	if err != nil {
		return err
	}
	body, err := checkResponse(response)
	if err != nil {
		return err
	}
	return checkBodyStatus(body)
}

func (z *Client) getSession(sessionId string) (*SessionInfo, error) {
	response, err := z.Get(z.getBaseUrl()+fmt.Sprintf("/session/%s", sessionId), http.Header{})
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 404 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		if strings.Contains(string(body), "No such session") {
			return nil, nil
		}
	}
	body, err := checkResponse(response)
	if err != nil {
		return nil, err
	}
	if err = checkBodyStatus(body); err != nil {
		return nil, err
	}
	return NewSessionInfo(body)
}

func queryString(name string, value string) (queryStr string) {
	queryStr = "?" + name
	if value != "" && len(value) > 0 {
		queryStr = queryStr + "=" + value
	}
	return queryStr
}

func checkResponse(response *http.Response) (string, error) {
	if response.StatusCode == 302 {
		return "", qerror.InvalidateZeppelinUser
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if response.StatusCode != 200 {
		return "", qerror.CallZeppelinRestApiFailed.Format(response.StatusCode, response.Status, string(body))
	}
	return string(body), nil
}

func checkBodyStatus(resBody string) error {
	bodyJson, err := fastjson.Parse(resBody)
	if err != nil {
		return err
	}
	if !strings.EqualFold("OK", string(bodyJson.GetStringBytes("status"))) {
		return qerror.ZeppelinReturnStatusError.Format(bodyJson.Get("message").String())
	}
	return nil
}
