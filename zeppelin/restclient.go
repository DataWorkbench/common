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

func (c *Client) getBaseUrl() string {
	return c.ClientConfig.ZeppelinRestUrl + "/api"
}

func (c *Client) createNoteWithGroup(notePath string, defaultInterpreterGroup string) (string, error) {
	reqObj := map[string]interface{}{}
	reqObj["name"] = notePath
	reqObj["defaultInterpreterGroup"] = defaultInterpreterGroup
	reqBytes, err := json.Marshal(&reqObj)
	if err != nil {
		return "", err
	}
	response, err := c.Post(c.getBaseUrl()+"/notebook", strings.NewReader(string(reqBytes)), http.Header{})
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

func (c *Client) createNote(notePath string) (string, error) {
	return c.createNoteWithGroup(notePath, "")
}

func (c *Client) deleteNote(noteId string) error {
	response, err := c.Delete(c.getBaseUrl()+fmt.Sprintf("/notebook/%s", noteId), http.Header{})
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

func (c *Client) addParagraph(noteId string, title string, text string) (string, error) {
	reqObj := map[string]interface{}{}
	reqObj["title"] = title
	reqObj["text"] = text
	reqBytes, err := json.Marshal(&reqObj)
	if err != nil {
		return "", err
	}
	response, err := c.Post(c.getBaseUrl()+fmt.Sprintf("/notebook/%s/paragraph", noteId), strings.NewReader(string(reqBytes)), http.Header{})
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

func (c *Client) updateParagraph(noteId string, paragraphId string, title string, text string) error {
	reqObj := map[string]string{}
	reqObj["title"] = title
	reqObj["text"] = text
	reqBytes, err := json.Marshal(reqObj)
	if err != nil {
		return err
	}
	response, err := c.Put(c.getBaseUrl()+fmt.Sprintf("/notebook/%s/paragraph/%s", noteId, paragraphId), strings.NewReader(string(reqBytes)), http.Header{})
	if err != nil {
		return err
	}
	body, err := checkResponse(response)
	if err != nil {
		return err
	}
	return checkBodyStatus(body)
}

func (c *Client) submitParagraphWithAll(noteId string, paragraphId string, sessionId string, parameters map[string]string) (*ParagraphResult, error) {
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
	url := c.getBaseUrl() + fmt.Sprintf("/notebook/job/%s/%s%s", noteId, paragraphId, queryString("sessionId", sessionId))
	fmt.Println(url)
	response, err := c.Post(c.getBaseUrl()+fmt.Sprintf("/notebook/job/%s/%s%s", noteId, paragraphId, queryString("sessionId", sessionId)),
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
	return c.queryParagraphResult(noteId, paragraphId)
}

func (c *Client) submitParagraphWithSessionId(noteId string, paragraphId string, sessionId string) (*ParagraphResult, error) {
	return c.submitParagraphWithAll(noteId, paragraphId, sessionId, make(map[string]string))
}

func (c *Client) submitParagraphWithParameters(noteId string, paragraphId string, parameters map[string]string) (*ParagraphResult, error) {
	return c.submitParagraphWithAll(noteId, paragraphId, "", parameters)
}

func (c *Client) submitParagraph(noteId string, paragraphId string) (*ParagraphResult, error) {
	return c.submitParagraphWithAll(noteId, paragraphId, "", make(map[string]string))
}

func (c *Client) executeParagraphWithAll(noteId string, paragraphId string, sessionId string, parameters map[string]string) (*ParagraphResult, error) {
	_, err := c.submitParagraphWithAll(noteId, paragraphId, sessionId, parameters)
	if err != nil {
		return nil, err
	}
	return c.waitUtilParagraphFinish(noteId, paragraphId)
}

func (c *Client) executeParagraphWithSessionId(noteId string, paragraphId string, sessionId string) (*ParagraphResult, error) {
	return c.executeParagraphWithAll(noteId, paragraphId, sessionId, make(map[string]string))
}

func (c *Client) executeParagraph(noteId string, paragraphId string) (*ParagraphResult, error) {
	return c.executeParagraphWithAll(noteId, paragraphId, "", make(map[string]string))
}

func (c *Client) cancelParagraph(noteId string, paragraphId string) error {
	response, err := c.Delete(c.getBaseUrl()+fmt.Sprintf("/notebook/job/%s/%s", noteId, paragraphId), http.Header{})
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

func (c *Client) waitUtilParagraphRunning(noteId string, paragraphId string) (*ParagraphResult, error) {
	for {
		paragraphResult, err := c.queryParagraphResult(noteId, paragraphId)
		if err != nil {
			return nil, err
		}
		if paragraphResult.Status.isRunning() {
			return paragraphResult, nil
		}
		time.Sleep(time.Millisecond * c.ClientConfig.QueryInterval)
	}
}

func (c *Client) waitUtilParagraphFinish(noteId string, paragraphId string) (*ParagraphResult, error) {
	for {
		paragraphResult, err := c.queryParagraphResult(noteId, paragraphId)
		if err != nil {
			return nil, err
		}
		if paragraphResult.Status.isCompleted() {
			return paragraphResult, nil
		}
		time.Sleep(time.Millisecond * c.ClientConfig.QueryInterval)
	}
}

func (c *Client) waitUtilParagraphFinishWithTimeout(noteId string, paragraphId string, timeoutInSec int) (*ParagraphResult, error) {
	start := time.Now().Second()
	for (time.Now().Second() - start) < timeoutInSec {
		paragraphResult, err := c.queryParagraphResult(noteId, paragraphId)
		if err != nil {
			return nil, err
		}
		if paragraphResult.Status.isCompleted() {
			return paragraphResult, nil
		}
		time.Sleep(time.Millisecond * c.ClientConfig.QueryInterval)
	}
	return nil, qerror.ZeppelinRunParagraphTimeout.Format(timeoutInSec)
}

func (c *Client) queryParagraphResult(noteId string, paragraphId string) (*ParagraphResult, error) {
	response, err := c.Get(c.getBaseUrl()+fmt.Sprintf("/notebook/%s/paragraph/%s", noteId, paragraphId), http.Header{})
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

func (c *Client) newSession(interpreter string) (*SessionInfo, error) {
	response, err := c.Post(c.getBaseUrl()+fmt.Sprintf("/session%s", queryString("interpreter", interpreter)), strings.NewReader(""), http.Header{})
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

func (c *Client) stopSession(sessionId string) error {
	response, err := c.Delete(c.getBaseUrl()+fmt.Sprintf("/session/%s", sessionId), http.Header{})
	if err != nil {
		return err
	}
	body, err := checkResponse(response)
	if err != nil {
		return err
	}
	return checkBodyStatus(body)
}

func (c *Client) getSession(sessionId string) (*SessionInfo, error) {
	response, err := c.Get(c.getBaseUrl()+fmt.Sprintf("/session/%s", sessionId), http.Header{})
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
