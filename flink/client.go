package flink

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/common/web/ghttp"
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
