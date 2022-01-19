package flink

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
)

type Client struct {
	*ghttp.Client
}

type Job struct {
	Jid              string         `json:"jid"`
	Name             string         `json:"name"`
	State            string         `json:"state"`
	StartTime        int64          `json:"start-time"`
	EndTime          int64          `json:"end-time"`
	Duration         int64          `json:"duration"`
	LastModification int64          `json:"last-modification"`
	Exceptions       *JobExceptions `json:"exceptions"`
}

type JobExceptions struct {
	RootException string           `json:"root-exception"`
	Timestamp     int64            `json:"timestamp"`
	AllExceptions []*AllExceptions `json:"all-exceptions"`
	Truncated     bool             `json:"truncated"`
}

type AllExceptions struct {
	Exception string `json:"exception"`
	Task      string `json:"task"`
	Location  string `json:"location"`
	Timestamp int64  `json:"timestamp"`
}

func NewClient(ctx context.Context, cfg *ghttp.ClientConfig) *Client {
	httpclient := ghttp.NewClient(ctx, cfg)
	return &Client{httpclient}
}

func (c *Client) ListJobs(ctx context.Context, flinkUrl string) ([]*Job, error) {
	var response *http.Response
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/jobs/overview", flinkUrl), strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	response, err = c.Send(ctx, req)
	if err != nil {
		return nil, err
	}
	body, err := checkResponse(response)
	if err != nil {
		return nil, err
	}
	var jobs []*Job
	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			return
		}
		job := Job{}
		if err = json.Unmarshal(value, &job); err != nil {
			return
		}
		jobs = append(jobs, &job)
	}, "jobs")
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (c *Client) GetInfo(ctx context.Context, flinkUrl string, flinkId string) (*Job, error) {
	var (
		response *http.Response
		job      *Job
	)
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/jobs/%s", flinkUrl, flinkId), strings.NewReader(""))
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
	if err = json.Unmarshal(res, &job); err != nil {
		return nil, err
	}
	return job, nil
}

func (c *Client) CancelJob(ctx context.Context, flinkUrl string, flinkId string) error {
	var response *http.Response
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()
	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("http://%s/jobs/%s", flinkUrl, flinkId), strings.NewReader(""))
	if err != nil {
		return err
	}
	response, err = c.Send(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetExceptions(ctx context.Context, flinkUrl string, flinkId string) (*JobExceptions, error) {
	var response *http.Response
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/jobs/%s/exceptions", flinkUrl, flinkId), strings.NewReader(""))
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
	var exceptions JobExceptions
	if err = json.Unmarshal(res, &exceptions); err != nil {
		return nil, err
	}
	return &exceptions, nil
}

func checkResponse(response *http.Response) (res []byte, err error) {
	res, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, qerror.FlinkRestError.Format(response.StatusCode, response.Status, string(res))
	}
	return
}
