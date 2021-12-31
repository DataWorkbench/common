package flink

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/DataWorkbench/common/qerror"

	"github.com/buger/jsonparser"
	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"

)

type Client struct {
	*httpclient.Client
	ClientConfig ClientConfig
}

func NewFlinkClient(config ClientConfig) *Client {
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(config.Timeout),
		httpclient.WithRetryCount(config.RetryCount),
		httpclient.WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(time.Millisecond*10, time.Millisecond*50))),
	)
	return &Client{client, config}
}

func (c *Client) listJobs() ([]*Job, error) {
	var jobs []*Job
	response, err := c.Get(c.getBaseUrl("jobs")+"/overview", http.Header{})
	if err != nil {
		return nil, err
	}
	body, err := c.checkResponse(response)
	if err != nil {
		return nil, err
	}
	_, _ = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			return
		}
		job := Job{}
		if err = json.Unmarshal(value, &job); err != nil {
			return
		}
		jobs = append(jobs, &job)
	}, "jobs")
	return jobs, nil
}

func (c *Client) getJobInfoByJobId(jobId string) (*Job, error) {
	var job *Job
	response, err := c.Get(c.getBaseUrl("jobs")+"/"+jobId, http.Header{})
	if err != nil {
		return nil, err
	}
	bytes, err := c.checkResponse(response)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bytes, &job); err != nil {
		return nil, err
	}
	return job, nil
}

func (c *Client) getJobInfo(jobId string, jobName string) (*Job, error) {
	job, err := c.getJobInfoByJobId(jobId)
	if err != nil {
		jobs, err := c.listJobs()
		if err != nil {
			return nil, err
		}
		for _, j := range jobs {
			if strings.EqualFold(j.Name, jobName) {
				return j, nil
			}
		}
		return nil, qerror.FlinkJobNotExists.Format(jobName)
	}
	return job, nil
}

func (c *Client) cancelJob(jobId string) error {
	response, err := c.Patch(c.getBaseUrl("jobs")+"/"+jobId, strings.NewReader(""), http.Header{})
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return qerror.FlinkRestError.Format(response.StatusCode, response.Status, string(body))
	}
	return nil
}

func (c *Client) savepoint(jobId string, target string, cancel bool) (*Savepoint, error) {
	reqObj := map[string]interface{}{}
	reqObj["cancel-job"] = cancel
	reqObj["target-directory"] = target
	reqBytes, err := json.Marshal(reqObj)
	if err != nil {
		return nil, err
	}
	response, err := c.Post(c.getBaseUrl("jobs")+fmt.Sprintf("/%s/savepoints", jobId), strings.NewReader(string(reqBytes)), http.Header{})
	bytes, err := c.checkResponse(response)
	if err != nil {
		return nil, err
	}
	var savepoint *Savepoint
	if err = json.Unmarshal(bytes, savepoint); err != nil {
		return nil, err
	}
	return savepoint, err
}

func (c *Client) triggerSavepoint(jobId string, target string) (string, error) {
	savepoint, err := c.savepoint(jobId, target, false)
	if err != nil {
		return "", err
	}
	return savepoint.RequestId, nil
}

func (c *Client) cancelWithSavepoint(jobId string, target string) (string, error) {
	savepoint, err := c.savepoint(jobId, target, true)
	if err != nil {
		return "", err
	}
	return savepoint.RequestId, nil
}

func (c *Client) getSavepoint(jobId string, requestId string) (*Savepoint, error) {
	response, err := c.Get(c.getBaseUrl("jobs")+fmt.Sprintf("/%s/savepoints/%s", jobId, requestId), http.Header{})
	if err != nil {
		return nil, err
	}
	bytes, err := c.checkResponse(response)
	if err != nil {
		return nil, err
	}
	var savepoint *Savepoint
	if err = json.Unmarshal(bytes, savepoint); err != nil {
		return nil, err
	}
	return savepoint, nil
}

func (c *Client) getBaseUrl(baseType string) string {
	return c.ClientConfig.FlinkRestUrl + "/" + baseType
}

func (c *Client) checkResponse(response *http.Response) (res []byte, err error) {
	res, err = ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		if err != nil {
			return nil, err
		}
		return nil, qerror.FlinkRestError.Format(response.StatusCode, response.Status, string(res))
	}
	return res, nil
}
