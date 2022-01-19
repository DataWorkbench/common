package flink

//
//import (
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"net/http"
//	"strings"
//	"time"
//
//	"github.com/DataWorkbench/common/qerror"
//
//	"github.com/buger/jsonparser"
//	"github.com/gojek/heimdall/v7"
//	"github.com/gojek/heimdall/v7/httpclient"
//)
//
//type Client struct {
//	*httpclient.Client
//	ClientConfig ClientConfig
//}
//
//func NewFlinkClient(config ClientConfig) *Client {
//	client := httpclient.NewClient(
//		httpclient.WithHTTPTimeout(config.Timeout),
//		httpclient.WithRetryCount(config.RetryCount),
//		httpclient.WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(time.Millisecond*10, time.Millisecond*50))),
//	)
//	return &Client{client, config}
//}
//
//func (c *Client) ListJobs(flinkUrl string) ([]*Job, error) {
//	var jobs []*Job
//	response, err := c.Get(fmt.Sprintf("http://%s/jobs/overview", flinkUrl), http.Header{})
//	if err != nil {
//		return nil, err
//	}
//	body, err := c.checkResponse(response)
//	if err != nil {
//		return nil, err
//	}
//	_, _ = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
//		if err != nil {
//			return
//		}
//		job := Job{}
//		if err = json.Unmarshal(value, &job); err != nil {
//			return
//		}
//		jobs = append(jobs, &job)
//	}, "jobs")
//	return jobs, nil
//}
//
//func (c *Client) GetJobInfoByJobId(flinkUrl string, jobId string) (*Job, error) {
//	var job *Job
//	response, err := c.Get(fmt.Sprintf("http://%s/jobs/%s", flinkUrl, jobId), http.Header{})
//	if err != nil {
//		return nil, err
//	}
//	bytes, err := c.checkResponse(response)
//	if err != nil {
//		return nil, err
//	}
//	if err = json.Unmarshal(bytes, &job); err != nil {
//		return nil, err
//	}
//	return job, nil
//}
//
//func (c *Client) GetJobInfo(flinkUrl string, jobId string, jobName string) (*Job, error) {
//	job, err := c.GetJobInfoByJobId(flinkUrl, jobId)
//	if err != nil {
//		jobs, err := c.ListJobs(flinkUrl)
//		if err != nil {
//			return nil, err
//		}
//		for _, j := range jobs {
//			if strings.EqualFold(j.Name, jobName) {
//				return j, nil
//			}
//		}
//		return nil, qerror.FlinkJobNotExists.Format(jobName)
//	}
//	return job, nil
//}
//
//func (c *Client) CancelJob(flinkUrl string, jobId string) error {
//	_, err := c.Patch(fmt.Sprintf("http://%s/jobs/%s", flinkUrl, jobId), strings.NewReader(""), http.Header{})
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (c *Client) savepoint(flinkUrl string, jobId string, target string, cancel bool) (*Savepoint, error) {
//	reqObj := map[string]interface{}{}
//	reqObj["cancel-job"] = cancel
//	reqObj["target-directory"] = target
//	reqBytes, err := json.Marshal(reqObj)
//	if err != nil {
//		return nil, err
//	}
//	response, err := c.Post(fmt.Sprintf("http://%s/jobs/%s/savepoints", flinkUrl, jobId), strings.NewReader(string(reqBytes)), http.Header{})
//	if err != nil {
//		return nil, err
//	}
//	bytes, err := c.checkResponse(response)
//	if err != nil {
//		return nil, err
//	}
//	var savepoint *Savepoint
//	if bytes != nil {
//		if err = json.Unmarshal(bytes, savepoint); err != nil {
//			return nil, err
//		}
//	}
//	return savepoint, err
//}
//
//func (c *Client) TriggerSavepoint(flinkUrl string, jobId string, target string) (string, error) {
//	savepoint, err := c.savepoint(flinkUrl, jobId, target, false)
//	if err != nil {
//		return "", err
//	}
//	return savepoint.RequestId, nil
//}
//
//func (c *Client) CancelWithSavepoint(flinkUrl string, jobId string, target string) (string, error) {
//	savepoint, err := c.savepoint(flinkUrl, jobId, target, true)
//	if err != nil {
//		return "", err
//	}
//	return savepoint.RequestId, nil
//}
//
//func (c *Client) GetSavepoint(flinkUrl string, jobId string, requestId string) (*Savepoint, error) {
//	response, err := c.Get(fmt.Sprintf("http://%s/jobs/%s/savepoints/%s", flinkUrl, jobId, requestId), http.Header{})
//	if err != nil {
//		return nil, err
//	}
//	bytes, err := c.checkResponse(response)
//	if err != nil {
//		return nil, err
//	}
//	var savepoint Savepoint
//	if err = json.Unmarshal(bytes, &savepoint); err != nil {
//		return nil, err
//	}
//	return &savepoint, nil
//}
//
//func (c *Client) checkResponse(response *http.Response) (res []byte, err error) {
//	res, err = ioutil.ReadAll(response.Body)
//	if response.StatusCode != 200 {
//		if err != nil {
//			return nil, err
//		}
//		return nil, qerror.FlinkRestError.Format(response.StatusCode, response.Status, string(res))
//	}
//	return res, nil
//}
