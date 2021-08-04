package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/DataWorkbench/common/constants"
	"github.com/DataWorkbench/common/grpcwrap"
	"github.com/DataWorkbench/glog"
	"github.com/DataWorkbench/gproto/pkg/jobdevpb"
	"gorm.io/gorm"
)

const (
	ParagraphUnknown    = "UNKNOWN"
	ParagraphFinish     = "FINISHED"
	ParagraphRunning    = "RUNNING"
	ParagraphReady      = "READY"
	ParagraphError      = "ERROR"
	ParagraphPending    = "PENDING"
	ParagraphAbort      = "ABORT"
	JobmanagerTableName = "jobmanager"
	MaxStatusFailedNum  = 100
)

type JobmanagerInfo struct {
	ID             string `gorm:"column:id;primaryKey"`
	NoteID         string `gorm:"column:noteid;"`
	Status         string `gorm:"column:status;"`
	Message        string `gorm:"column:message;"`
	Paragraph      string `gorm:"column:paragraph;"`
	CreateTime     string `gorm:"column:createtime;"`
	UpdateTime     string `gorm:"column:updatetime;"`
	Resources      string `gorm:"column:resources;"`
	SpaceID        string `gorm:"column:spaceid;"`
	EngineType     string `gorm:"column:enginetype;"`
	ZeppelinServer string `gorm:"column:zeppelinserver;"`
}

func (smi JobmanagerInfo) TableName() string {
	return JobmanagerTableName
}

type JobWatchInfo struct {
	ID                string                        `json:"id"`
	EngineType        string                        `json:"enginetype"`
	NoteID            string                        `json:"noteid"`
	ServerAddr        string                        `json:"serveraddr"`
	FlinkParagraphIDs constants.FlinkParagraphsInfo `json:"flinkparagraphids"`
	FlinkResources    constants.JobResources        `json:"jobresources"`
}

func StringStatusToInt32(s string) (r int32) {
	if s == constants.StatusRunningString {
		r = constants.StatusRunning
	} else if s == constants.StatusFinishString {
		r = constants.StatusFinish
	} else if s == constants.StatusFailedString {
		r = constants.StatusFailed
	}
	return r
}

func Int32StatusToString(i int32) (r string) {
	if i == constants.StatusRunning {
		r = constants.StatusRunningString
	} else if i == constants.StatusFinish {
		r = constants.StatusFinishString
	} else if i == constants.StatusFailed {
		r = constants.StatusFailedString
	}
	return r
}

type HttpClient struct {
	ZeppelinServer string
	Client         *http.Client
}

func NewHttpClient(serverAddr string) HttpClient {
	return HttpClient{ZeppelinServer: "http://" + serverAddr, Client: &http.Client{Timeout: time.Second * 60}}
}

func doRequest(client *http.Client, method string, status int, api string, body string, retJson bool) (repJson map[string]string, repString string, err error) {
	var (
		req     *http.Request
		rep     *http.Response
		reqBody io.Reader
	)

	if body == "" {
		reqBody = nil
	} else {
		reqBody = strings.NewReader(body)
	}

	req, err = http.NewRequest(method, api, reqBody)
	if err != nil {
		return
	}

	rep, err = client.Do(req)
	if err != nil {
		//rep.Body.Close()
		return
	}

	repBody, _ := ioutil.ReadAll(rep.Body)
	rep.Body.Close()

	repString = string(repBody)
	if retJson {
		err = json.Unmarshal(repBody, &repJson)
		if err != nil {
			rep.Body.Close()
			return
		}
	}

	if rep.StatusCode != status {
		err = fmt.Errorf("%s request failed, http status code %d, message %s", api, rep.StatusCode, repString)
		rep.Body.Close()
		return
	}

	return
}

func (ex *HttpClient) CreateNote(ID string) (noteID string, err error) {
	var repJson map[string]string

	repJson, _, err = doRequest(ex.Client, http.MethodPost, http.StatusOK, ex.ZeppelinServer+"/api/notebook", fmt.Sprintf("{\"name\": \"%s\"}", ID), true)
	if err != nil {
		return
	}
	noteID = repJson["body"]

	return noteID, nil
}

func (ex *HttpClient) DeleteNote(ID string) (err error) {
	_, _, err = doRequest(ex.Client, http.MethodDelete, http.StatusOK, ex.ZeppelinServer+"/api/notebook/"+ID, "", false)
	return
}

func (ex *HttpClient) StopAllParagraphs(noteID string) (err error) {
	_, _, err = doRequest(ex.Client, http.MethodDelete, http.StatusOK, ex.ZeppelinServer+"/api/notebook/job/"+noteID, "", false)
	return
}

func (ex *HttpClient) CreateParagraph(noteID string, index int32, name string, text string) (paragraphID string, err error) {
	var repJson map[string]string

	repJson, _, err = doRequest(ex.Client, http.MethodPost, http.StatusOK, ex.ZeppelinServer+"/api/notebook/"+noteID+"/paragraph", fmt.Sprintf("{\"title\": \"%s\", \"text\": %s, \"index\": %d}", name, strconv.Quote(text), index), true)
	if err != nil {
		return
	}
	paragraphID = repJson["body"]

	return paragraphID, nil
}

func (ex *HttpClient) RunParagraphSync(noteID string, paragraphID string) (err error) {
	var status string
	_, _, err = doRequest(ex.Client, http.MethodPost, http.StatusOK, ex.ZeppelinServer+"/api/notebook/run/"+noteID+"/"+paragraphID, "", false)
	if err != nil {
		return
	}
	status, err = ex.GetParagraphStatus(noteID, paragraphID)
	if status != "OK" && status != "FINISHED" {
		err = fmt.Errorf("run failed. status is " + status)
	}
	return
}

func (ex *HttpClient) RunParagraphAsync(noteID string, paragraphID string) (err error) {
	_, _, err = doRequest(ex.Client, http.MethodPost, http.StatusOK, ex.ZeppelinServer+"/api/notebook/job/"+noteID+"/"+paragraphID, "", false)
	return
}

func (ex *HttpClient) GetParagraphStatus(noteID string, paragraphID string) (status string, err error) {
	var repString string
	var repJsonLevel1 map[string]json.RawMessage
	var repJsonLevel2 map[string]string

	_, repString, err = doRequest(ex.Client, http.MethodGet, http.StatusOK, ex.ZeppelinServer+"/api/notebook/job/"+noteID+"/"+paragraphID, "", false)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(repString), &repJsonLevel1)
	if err != nil {
		return
	}
	err = json.Unmarshal(repJsonLevel1["body"], &repJsonLevel2)
	if err != nil {
		return
	}
	status = repJsonLevel2["status"]

	return
}

func (ex *HttpClient) GetParagraphResultOutput(noteID string, paragraphID string) (msg string, err error) {
	var repString string
	var repJsonLevel1 map[string]json.RawMessage
	var repJsonLevel2 map[string]json.RawMessage
	var repJsonLevel3 map[string]json.RawMessage

	_, repString, err = doRequest(ex.Client, http.MethodGet, http.StatusOK, ex.ZeppelinServer+"/api/notebook/"+noteID+"/paragraph/"+paragraphID, "", false)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(repString), &repJsonLevel1)
	if err != nil {
		return
	}
	err = json.Unmarshal(repJsonLevel1["body"], &repJsonLevel2)
	if err != nil {
		return
	}
	err = json.Unmarshal(repJsonLevel2["results"], &repJsonLevel3)
	if err != nil {
		return
	}

	msg = string(repJsonLevel3["msg"])

	return
}

type JobdevClient struct {
	Client jobdevpb.JobdeveloperClient
}

func NewJobdevClient(conn *grpcwrap.ClientConn) (c JobdevClient, err error) {
	c.Client = jobdevpb.NewJobdeveloperClient(conn)
	return c, nil
}

func FreeJobResources(ctx context.Context, resources constants.JobResources, EngineType string, logger *glog.Logger, httpClient HttpClient, jobdevClient JobdevClient) (err error) {
	if EngineType == constants.ServerTypeFlink {
		var (
			req          jobdevpb.JobFreeRequest
			zeppelinFree constants.JobFreeActionFlink
			resp         *jobdevpb.JobFreeAction
			noteID       string
			paragraphID  string
		)

		defer func() {
			if err != nil {
				logger.Warn().String("can't delete jar", resources.Jar).String("FreeEngine", resources.JobID).Error("message", err).Fire()
			}
			if noteID != "" {
				_ = httpClient.DeleteNote(noteID)
			}
		}()

		req.EngineType = EngineType
		resourcesByte, _ := json.Marshal(resources)
		req.JobResources = string(resourcesByte)

		resp, err = jobdevClient.Client.JobFree(ctx, &req)
		if err != nil {
			return
		}

		respString := resp.GetJobResources()
		if respString != "" {
			if err = json.Unmarshal([]byte(respString), &zeppelinFree); err != nil {
				return
			}
		}

		if zeppelinFree.ZeppelinDeleteJar != "" {
			noteID, err = httpClient.CreateNote(resources.JobID + "_delete_resources")
			if err != nil {
				return
			}

			paragraphID, err = httpClient.CreateParagraph(noteID, 0, "delete resources", zeppelinFree.ZeppelinDeleteJar)
			if err != nil {
				return
			}

			if err = httpClient.RunParagraphSync(noteID, paragraphID); err != nil {
				return
			}
		}

		logger.Info().String("delete jar", resources.Jar).String("FreeEngine", resources.JobID).Fire()
	}

	return
}

func ModifyStatus(ctx context.Context, ID string, status int32, message string, resources constants.JobResources, EngineType string, db *gorm.DB, logger *glog.Logger, httpClient HttpClient, jobdevClient JobdevClient) (err error) {
	var info JobmanagerInfo

	info.ID = ID
	info.Status = Int32StatusToString(status)
	info.Message = message
	info.UpdateTime = time.Now().Format("2006-01-02 15:04:05")

	edb := db.WithContext(ctx)
	if err = edb.Select("status", "message", "updatetime").Where("id = ? ", info.ID).Updates(info).Error; err != nil {
		return
	}

	if status == constants.StatusFinish || status == constants.StatusFailed {
		err = FreeJobResources(ctx, resources, EngineType, logger, httpClient, jobdevClient)
		if err != nil {
			logger.Warn().String("can't delete jar", resources.Jar).String("can't FreeEngine", resources.JobID).Fire()
		}
	}

	return
}
