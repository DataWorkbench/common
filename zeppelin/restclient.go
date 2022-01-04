package zeppelin

import (
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/DataWorkbench/common/qerror"

	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
)

type ClientConfig struct {
	ZeppelinRestUrl string
	Timeout         time.Duration
	RetryCount      int
	QueryInterval   time.Duration
}

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
	return jsonparser.GetString(body, "body")
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
	return checkBodyStatus(body)
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
	return jsonparser.GetString(body, "body")
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

func (c *Client) submitParagraphWithSessionId(noteId string, paragraphId string, sessionId string) (*ParagraphResult, error) {
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
	return c.QueryParagraphResult(noteId, paragraphId)
}

func (c *Client) submitParagraph(noteId string, paragraphId string) (*ParagraphResult, error) {
	return c.submitParagraphWithSessionId(noteId, paragraphId, "")
}

func (c *Client) executeParagraphWithSessionId(noteId string, paragraphId string, sessionId string) (*ParagraphResult, error) {
	_, err := c.submitParagraphWithSessionId(noteId, paragraphId, sessionId)
	if err != nil {
		return nil, err
	}
	return c.waitUtilParagraphFinish(noteId, paragraphId)
}

func (c *Client) executeParagraph(noteId string, paragraphId string) (*ParagraphResult, error) {
	return c.executeParagraphWithSessionId(noteId, paragraphId, "")
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
	return checkBodyStatus(body)
}

func (c *Client) waitUtilParagraphRunning(noteId string, paragraphId string) (*ParagraphResult, error) {
	for {
		paragraphResult, err := c.QueryParagraphResult(noteId, paragraphId)
		if err != nil {
			return nil, err
		}
		if paragraphResult.Status.IsRunning() || paragraphResult.Status.IsFinished() {
			return paragraphResult, nil
		}
		if paragraphResult.Status.IsFailed() || paragraphResult.Status.IsUnknown() {
			return paragraphResult, nil
		}
		time.Sleep(time.Millisecond * c.ClientConfig.QueryInterval)
	}
}

func (c *Client) waitUtilParagraphFinish(noteId string, paragraphId string) (*ParagraphResult, error) {
	for {
		paragraphResult, err := c.QueryParagraphResult(noteId, paragraphId)
		if err != nil {
			return nil, err
		}
		if paragraphResult.Status.IsCompleted() {
			return paragraphResult, nil
		}
		time.Sleep(time.Millisecond * c.ClientConfig.QueryInterval)
	}
}

func (c *Client) QueryParagraphResult(noteId string, paragraphId string) (*ParagraphResult, error) {
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
	return NewParagraphResult(body)
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

func checkResponse(response *http.Response) ([]byte, error) {
	if response.StatusCode == 302 {
		return nil, qerror.InvalidateZeppelinUser
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
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
