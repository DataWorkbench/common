package flink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/common/web/ghttp"
	"io/ioutil"
	"net/http"
)

type Client struct {
	*ghttp.Client
}

func New(ctx context.Context, cfg *ghttp.ClientConfig) *Client {
	httpclient := ghttp.NewClient(ctx, cfg)
	return &Client{httpclient}
}

func (c *Client) get(ctx context.Context, url string, data interface{}) (err error) {
	var resp *http.Response
	var req *http.Request
	var body []byte

	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	if resp, err = c.Send(ctx, req); err != nil {
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		return qerror.FlinkRestError.Format(resp.StatusCode, resp.Status, string(body))
	}

	if err = json.Unmarshal(body, data); err != nil {
		return
	}
	return nil
}

func (c *Client) Overview(ctx context.Context, flinkUrl string) (*Overview, error) {
	var response *http.Response
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()

	url := fmt.Sprintf("http://%s/overview", flinkUrl)
	data := new(Overview)
	if err := c.get(ctx, url, data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) ListJobs(ctx context.Context, flinkUrl string) (*JobsOverview, error) {
	var response *http.Response
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()
	url := fmt.Sprintf("http://%s/jobs/overview", flinkUrl)
	data := new(JobsOverview)
	if err := c.get(ctx, url, data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) ListTaskManagers(ctx context.Context, flinkUrl string) (*TaskManagers, error) {
	var response *http.Response
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()
	url := fmt.Sprintf("http://%s/taskmanagers", flinkUrl)
	data := new(TaskManagers)
	if err := c.get(ctx, url, data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) DescribeJob(ctx context.Context, flinkUrl string, flinkId string) (*JobInfo, error) {
	var response *http.Response
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()
	url := fmt.Sprintf("http://%s/jobs/%s", flinkUrl, flinkId)
	data := new(JobInfo)
	if err := c.get(ctx, url, data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) DescribeJobPlan(ctx context.Context, flinkUrl string, flinkId string) (*JobPlan, error) {
	var response *http.Response
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()
	url := fmt.Sprintf("http://%s/jobs/%s/plan", flinkUrl, flinkId)
	data := new(JobPlan)
	if err := c.get(ctx, url, data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) DescribeJobExceptions(ctx context.Context, flinkUrl string, flinkId string) (*JobExceptions, error) {
	var response *http.Response
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()
	url := fmt.Sprintf("http://%s/jobs/%s/exceptions", flinkUrl, flinkId)
	data := new(JobExceptions)
	if err := c.get(ctx, url, data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) CancelJob(ctx context.Context, flinkUrl string, flinkId string) error {
	var response *http.Response
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()
	url := fmt.Sprintf("http://%s/jobs/%s", flinkUrl, flinkId)
	req, err := http.NewRequest(http.MethodPatch, url, nil)
	if err != nil {
		return err
	}
	response, err = c.Send(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) TriggerSavepoint(ctx context.Context, flinkUrl string, flinkId string, targetDirectory string) (string, error) {
	return c.SavePoint(ctx, flinkUrl, flinkId, false, targetDirectory)
}

func (c *Client) CancelWithSavepoint(ctx context.Context, flinkUrl string, flinkId string, targetDirectory string) (string, error) {
	return c.SavePoint(ctx, flinkUrl, flinkId, true, targetDirectory)
}

func (c *Client) SavePoint(ctx context.Context, flinkUrl string, flinkId string, cancelJob bool, targetDirectory string) (requestId string, err error) {
	var req *http.Request
	var response *http.Response
	var body []byte

	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()
	url := fmt.Sprintf("http://%s/jobs/%s/savepoints", flinkUrl, flinkId)
	body, err = json.Marshal(&map[string]interface{}{
		"cancel-job":       cancelJob,
		"target-directory": targetDirectory,
	})
	if err != nil {
		return
	}
	req, err = http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return
	}
	response, err = c.Send(ctx, req)
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	requestId = string(body)
	return
}
