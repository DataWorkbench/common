package constants

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/DataWorkbench/common/grpcwrap"
	"github.com/DataWorkbench/common/utils/idgenerator"
	"github.com/DataWorkbench/glog"
	"github.com/DataWorkbench/gproto/pkg/jobdevpb"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type EngineRequestOptions struct {
	JobID       string  `json:"jobid"`
	EngineID    string  `json:"engineid"`
	WorkspaceID string  `json:"workspaceid"`
	Parallelism int32   `json:"parallelism"`
	JobMem      int32   `json:"job_mem"` // in MB
	JobCpu      float32 `json:"job_cpu"`
	TaskCpu     float32 `json:"task_cpu"`
	TaskMem     int32   `json:"task_mem"` // in MB
	TaskNum     int32   `json:"task_num"`
	AccessKey   string  `json:"accesskey"`
	SecretKey   string  `json:"secretkey"`
	EndPoint    string  `json:"endpoint"`
}

type EngineResponseOptions struct {
	EngineType      string `json:"enginetype"`
	EngineHost      string `json:"enginehost"`
	EnginePort      string `json:"engineport"`
	EngineExtension string `json:"engineextension"`
}

const (
	StatusFailed        = InstanceStateFailed
	StatusFailedString  = "failed"
	StatusFinish        = InstanceStateSucceed
	StatusFinishString  = "finish"
	StatusRunning       = InstanceStateRunning
	StatusRunningString = "running"
	JobSuccess          = "success"
	JobAbort            = "job abort"
	JobRunning          = "job running"
	jobError            = "error happend"
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
	ID                string              `json:"id"`
	EngineType        string              `json:"enginetype"`
	NoteID            string              `json:"noteid"`
	ServerAddr        string              `json:"serveraddr"`
	FlinkParagraphIDs FlinkParagraphsInfo `json:"flinkparagraphids"`
	FlinkResources    JobResources        `json:"jobresources"`
}

func StringStatusToInt32(s string) (r int32) {
	if s == StatusRunningString {
		r = StatusRunning
	} else if s == StatusFinishString {
		r = StatusFinish
	} else if s == StatusFailedString {
		r = StatusFailed
	}
	return r
}

func Int32StatusToString(i int32) (r string) {
	if i == StatusRunning {
		r = StatusRunningString
	} else if i == StatusFinish {
		r = StatusFinishString
	} else if i == StatusFailed {
		r = StatusFailedString
	}
	return r
}

type HttpClient struct {
	ZeppelinServer string
	Client         *http.Client
}

const (
	ParagraphUnknown = "UNKNOWN"
	ParagraphFinish  = "FINISHED"
	ParagraphRunning = "RUNNING"
	ParagraphReady   = "READY"
	ParagraphError   = "ERROR"
	ParagraphPending = "PENDING"
	ParagraphAbort   = "ABORT"
)

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

	repJson, _, err = doRequest(ex.Client, http.MethodPost, http.StatusOK, ex.ZeppelinServer+"/api/notebook/"+noteID+"/paragraph", fmt.Sprintf("{\"title\": \"%s\", \"text\": \"%s\", \"index\": %d}", name, text, index), true)
	if err != nil {
		return
	}
	paragraphID = repJson["body"]

	return paragraphID, nil
}

func (ex *HttpClient) RunParagraphSync(noteID string, paragraphID string) (err error) {
	var status string
	_, _, err = doRequest(ex.Client, http.MethodPost, http.StatusOK, ex.ZeppelinServer+"/api/notebook/run/"+noteID+"/"+paragraphID, "", false)
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
	Ctx    context.Context
}

func NewJobdevClient(serverAddr string) (c JobdevClient, err error) {
	var conn *grpc.ClientConn

	ctx := glog.WithContext(context.Background(), glog.NewDefault())
	conn, err = grpcwrap.NewConn(ctx, &grpcwrap.ClientConfig{
		Address:      serverAddr,
		LogLevel:     2,
		LogVerbosity: 99,
	})
	if err != nil {
		return
	}

	c.Client = jobdevpb.NewJobdeveloperClient(conn)

	ln := glog.NewDefault().Clone()
	reqId, _ := idgenerator.New("").Take()
	ln.WithFields().AddString("rid", reqId)

	c.Ctx = grpcwrap.ContextWithRequest(context.Background(), ln, reqId)

	return c, nil
}

func FreeJobResources(resources JobResources, EngineType string, logger *glog.Logger, httpClient HttpClient, jobdevClient JobdevClient) (err error) {
	if EngineType == ServerTypeFlink {
		var (
			req          jobdevpb.JobFreeRequest
			zeppelinFree JobFreeActionFlink
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

		resp, tmperr := jobdevClient.Client.JobFree(jobdevClient.Ctx, &req)
		if tmperr != nil {
			err = tmperr
			return
		}
		fmt.Println(resp)

		//if err = json.Unmarshal([]byte(resp.GetJobResources()), &zeppelinFree); err != nil {
		//	fmt.Println("UNKNOWN where resp is empty---------/lzzhang")
		//	return
		//}

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

func ModifyStatus(ctx context.Context, ID string, status int32, message string, resources JobResources, EngineType string, db *gorm.DB, logger *glog.Logger, httpClient HttpClient, jobdevClient JobdevClient) (err error) {
	var info JobmanagerInfo

	info.ID = ID
	info.Status = Int32StatusToString(status)
	info.Message = message
	info.UpdateTime = time.Now().Format("2006-01-02 15:04:05")

	edb := db.WithContext(ctx)
	if err = edb.Select("status", "message", "updatetime").Where("id = ? ", info.ID).Updates(info).Error; err != nil {
		return
	}

	if status == StatusFinish || status == StatusFailed {
		err = FreeJobResources(resources, EngineType, logger, httpClient, jobdevClient)
		if err != nil {
			logger.Warn().String("can't delete jar", resources.Jar).String("can't FreeEngine", resources.JobID).Fire()
		}
	}

	return
}
