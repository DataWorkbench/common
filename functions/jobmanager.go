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
	"github.com/DataWorkbench/gproto/pkg/model"
	"github.com/DataWorkbench/gproto/pkg/request"
	"github.com/DataWorkbench/gproto/pkg/response"
	"gorm.io/gorm"
)

const (
	ParagraphUnknown = "UNKNOWN"
	ParagraphFinish  = "FINISHED"
	ParagraphRunning = "RUNNING"
	ParagraphReady   = "READY"
	ParagraphError   = "ERROR"
	ParagraphPending = "PENDING"
	ParagraphAbort   = "ABORT"

	MaxStatusFailedNum = 30

	JobTableName = "jobmanager"
)

type JobQueueType struct {
	Watch           JobWatchInfo
	StatusFailedNum int32
	HttpClient      HttpClient
}

type JobWatchInfo struct {
	JobID             string                        `json:"id"`
	NoteID            string                        `json:"noteid"`
	ServerAddr        string                        `json:"serveraddr"`
	FlinkParagraphIDs constants.FlinkParagraphsInfo `json:"flinkparagraphids"`
	FlinkResources    model.JobResources            `json:"jobresources"`
	JobState          response.JobState
}

type JobdevClient struct {
	Client jobdevpb.JobdeveloperClient
}

type JobmanagerInfo struct {
	JobID          string `gorm:"column:jobid;primaryKey"`
	SpaceID        string `gorm:"column:spaceid;"`
	NoteID         string `gorm:"column:noteid;"`
	Status         string `gorm:"column:status;"`
	Message        string `gorm:"column:message;"`
	Paragraph      string `gorm:"column:paragraph;"`
	CreateTime     string `gorm:"column:createtime;"`
	UpdateTime     string `gorm:"column:updatetime;"`
	Resources      string `gorm:"column:resources;"`
	ZeppelinServer string `gorm:"column:zeppelinserver;"`
}

func (smi JobmanagerInfo) TableName() string {
	return JobTableName
}

func NewJobdevClient(conn *grpcwrap.ClientConn) (c JobdevClient, err error) {
	c.Client = jobdevpb.NewJobdeveloperClient(conn)
	return c, nil
}

type HttpClient struct {
	ZeppelinServer string
	Client         *http.Client
}

func NewHttpClient(serverAddr string) HttpClient {
	return HttpClient{ZeppelinServer: "http://" + serverAddr, Client: &http.Client{Timeout: time.Second * 600}}
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

func (ex *HttpClient) CreateNote(jobID string) (noteID string, err error) {
	var repJson map[string]string

	repJson, _, err = doRequest(ex.Client, http.MethodPost, http.StatusOK, ex.ZeppelinServer+"/api/notebook", fmt.Sprintf("{\"name\": \"%s\"}", jobID), true)
	if err != nil {
		return
	}
	noteID = repJson["body"]

	return noteID, nil
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

func (ex *HttpClient) RunParagraphSync(noteID string, paragraphID string) (err error) {
	var status string
	_, _, err = doRequest(ex.Client, http.MethodPost, http.StatusOK, ex.ZeppelinServer+"/api/notebook/run/"+noteID+"/"+paragraphID, "", false)
	if err != nil {
		return
	}
	status, err = ex.GetParagraphStatus(noteID, paragraphID)
	if status != "OK" && status != "FINISHED" {
		msg, _ := ex.GetParagraphResultOutput(noteID, paragraphID)
		err = fmt.Errorf("run note " + noteID + " failed. status is " + status + ". the output message is " + msg)
	}
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

func (ex *HttpClient) CreateParagraph(noteID string, index int32, name string, text string) (paragraphID string, err error) {
	var repJson map[string]string

	repJson, _, err = doRequest(ex.Client, http.MethodPost, http.StatusOK, ex.ZeppelinServer+"/api/notebook/"+noteID+"/paragraph", fmt.Sprintf("{\"title\": \"%s\", \"text\": %s, \"index\": %d}", name, strconv.Quote(text), index), true)
	if err != nil {
		return
	}
	paragraphID = repJson["body"]

	return paragraphID, nil
}

func (ex *HttpClient) RunParagraphAsync(noteID string, paragraphID string) (err error) {
	_, _, err = doRequest(ex.Client, http.MethodPost, http.StatusOK, ex.ZeppelinServer+"/api/notebook/job/"+noteID+"/"+paragraphID, "", false)
	return
}

func StringStatusToInt32(s string) (r int32) {
	if s == constants.StatusRunningString {
		r = int32(constants.StatusRunning.Number())
	} else if s == constants.StatusFinishString {
		r = int32(constants.StatusFinish.Number())
	} else if s == constants.StatusFailedString {
		r = int32(constants.StatusFailed.Number())
	} else if s == constants.StatusTerminatedString {
		r = int32(constants.StatusTerminated.Number())
	}
	return r
}

func Int32StatusToString(i int32) (r string) {
	if i == int32(constants.StatusRunning.Number()) {
		r = constants.StatusRunningString
	} else if i == int32(constants.StatusFinish.Number()) {
		r = constants.StatusFinishString
	} else if i == int32(constants.StatusFailed) {
		r = constants.StatusFailedString
	} else if i == int32(constants.StatusTerminated) {
		r = constants.StatusTerminatedString
	}
	return r
}

func (ex *HttpClient) DeleteNote(ID string) (err error) {
	_, _, err = doRequest(ex.Client, http.MethodDelete, http.StatusOK, ex.ZeppelinServer+"/api/notebook/"+ID, "", false)
	return
}

func (ex *HttpClient) StopAllParagraphs(noteID string) (err error) {
	_, _, err = doRequest(ex.Client, http.MethodDelete, http.StatusOK, ex.ZeppelinServer+"/api/notebook/job/"+noteID, "", false)
	return
}

func FreeJobResources(ctx context.Context, resources model.JobResources, logger *glog.Logger, httpClient HttpClient, jobdevClient JobdevClient) (err error) {
	var (
		resp   *response.JobFree
		noteID string
	)

	defer func() {
		if err != nil {
			logger.Warn().String("can't free resources ", resources.JobID).Error("message", err).Fire()
			err = nil
		}
		if noteID != "" {
			_ = httpClient.DeleteNote(noteID)
		}
	}()

	resp, err = jobdevClient.Client.JobFree(ctx, &request.JobFree{Resources: &resources})
	if err != nil {
		return
	}

	if resp.ZeppelinDeleteJar != "" {
		var paragraphID string

		noteID, err = httpClient.CreateNote(resources.JobID + "_delete_resources")
		if err != nil {
			return
		}

		paragraphID, err = httpClient.CreateParagraph(noteID, 0, "delete_resources", resp.ZeppelinDeleteJar)
		if err != nil {
			return
		}

		if err = httpClient.RunParagraphSync(noteID, paragraphID); err != nil {
			return
		}
	}

	return
}

func ModifyState(ctx context.Context, jobID string, state int32, message string, db *gorm.DB) (err error) {
	var info JobmanagerInfo

	info.JobID = jobID
	info.Status = Int32StatusToString(state)
	info.Message = message
	info.UpdateTime = time.Now().Format("2006-01-02 15:04:05")

	edb := db.WithContext(ctx)
	err = edb.Select("status", "message", "updatetime").Where("jobid = ? ", info.JobID).Updates(info).Error

	return
}

func InitJobInfo(watchInfo JobWatchInfo) (job JobQueueType) {
	job.Watch = watchInfo
	job.StatusFailedNum = 0
	job.HttpClient = NewHttpClient(watchInfo.ServerAddr)

	return
}

func JobInfoToWatchInfo(jobinfo JobmanagerInfo) (watchInfo JobWatchInfo) {
	var Pa constants.FlinkParagraphsInfo
	var resource model.JobResources

	watchInfo.JobID = jobinfo.JobID
	watchInfo.NoteID = jobinfo.NoteID
	watchInfo.ServerAddr = jobinfo.ZeppelinServer
	_ = json.Unmarshal([]byte(jobinfo.Paragraph), &Pa)
	watchInfo.FlinkParagraphIDs = Pa
	if jobinfo.Resources != "" {
		_ = json.Unmarshal([]byte(jobinfo.Resources), &resource)
	}
	watchInfo.FlinkResources = resource
	watchInfo.JobState.State = StringStatusToInt32(jobinfo.Status)
	watchInfo.JobState.Message = jobinfo.Message

	return
}

func GetZeppelinJobState(ctx context.Context, jobInput JobQueueType, logger *glog.Logger, db *gorm.DB, jobdevClient JobdevClient) (job JobQueueType, err error) {
	var status string

	defer func() {
		if err != nil {
			job.StatusFailedNum += 1
		}
	}()

	job = jobInput
	if status, err = job.HttpClient.GetParagraphStatus(job.Watch.NoteID, job.Watch.FlinkParagraphIDs.MainRun); err != nil {
		logger.Error().Msg("can't get this paragraph status").String("noteid", job.Watch.NoteID).String("jobid", job.Watch.JobID).Int32("failed number", job.StatusFailedNum).Fire()
		job.StatusFailedNum += 1

		if job.StatusFailedNum <= MaxStatusFailedNum {
			err = nil
			return
		}

		return
	}

	if status == ParagraphFinish {
		var jobmsg string

		if jobmsg, err = job.HttpClient.GetParagraphResultOutput(job.Watch.NoteID, job.Watch.FlinkParagraphIDs.MainRun); err != nil {
			jobmsg = "job finish, but can't get the MainRun paragraph output"
			logger.Error().Msg(jobmsg).String("noteid", job.Watch.NoteID).String("jobid", job.Watch.JobID).String("error msg", err.Error()).Fire()
			err = nil
		}

		if err = ModifyState(ctx, job.Watch.JobID, int32(constants.StatusFinish), jobmsg, db); err != nil {
			logger.Error().Msg("can't change the job status to finish").String("jobid", job.Watch.JobID).Fire()
			return
		}

		if err = job.HttpClient.DeleteNote(job.Watch.NoteID); err != nil {
			logger.Error().Msg("can't delete the note").String("noteid", job.Watch.NoteID).String("jobid", job.Watch.JobID).String("error msg", err.Error()).Fire()
			err = nil
		}
		_ = FreeJobResources(ctx, job.Watch.FlinkResources, logger, job.HttpClient, jobdevClient)

		job.Watch.JobState.State = int32(constants.StatusFinish)
		job.Watch.JobState.Message = jobmsg
		return
	} else if status == ParagraphError {
		var jobmsg string

		if jobmsg, err = job.HttpClient.GetParagraphResultOutput(job.Watch.NoteID, job.Watch.FlinkParagraphIDs.MainRun); err != nil {
			jobmsg = "job error, but can't get the MainRun paragraph output"
			logger.Error().Msg(jobmsg).String("noteid", job.Watch.NoteID).String("jobid", job.Watch.JobID).String("error msg", err.Error()).Fire()
			err = nil
		}

		if err = ModifyState(ctx, job.Watch.JobID, int32(constants.StatusFailed), jobmsg, db); err != nil {
			logger.Error().Msg("can't change the job status to failed").String("jobid", job.Watch.JobID).Fire()
			return
		}

		if err = job.HttpClient.DeleteNote(job.Watch.NoteID); err != nil {
			logger.Error().Msg("can't delete the note").String("noteid", job.Watch.NoteID).String("jobid", job.Watch.JobID).String("error msg", err.Error()).Fire()
			err = nil
		}
		_ = FreeJobResources(ctx, job.Watch.FlinkResources, logger, job.HttpClient, jobdevClient)

		job.Watch.JobState.State = int32(constants.StatusFailed)
		job.Watch.JobState.Message = jobmsg
		return
	} else if status == ParagraphAbort {
		var jobmsg string

		if jobmsg, err = job.HttpClient.GetParagraphResultOutput(job.Watch.NoteID, job.Watch.FlinkParagraphIDs.MainRun); err != nil {
			jobmsg = "job terminated, but can't get the MainRun paragraph output"
			logger.Error().Msg(jobmsg).String("noteid", job.Watch.NoteID).String("jobid", job.Watch.JobID).String("error msg", err.Error()).Fire()
			err = nil
		}

		if err = ModifyState(ctx, job.Watch.JobID, int32(constants.StatusTerminated), jobmsg, db); err != nil {
			logger.Error().Msg("can't change the job status to terminated").String("jobid", job.Watch.JobID).Fire()
			return
		}

		if err = job.HttpClient.DeleteNote(job.Watch.NoteID); err != nil {
			logger.Error().Msg("can't delete the note").String("noteid", job.Watch.NoteID).String("jobid", job.Watch.JobID).String("error msg", err.Error()).Fire()
			err = nil
		}
		_ = FreeJobResources(ctx, job.Watch.FlinkResources, logger, job.HttpClient, jobdevClient)

		job.Watch.JobState.State = int32(constants.StatusTerminated)
		job.Watch.JobState.Message = jobmsg
		return
	} else {
		/* paragraph is running
		   ParagraphUnknown = "UNKNOWN"
		   ParagraphRunning = "RUNNING"
		   ParagraphReady = "READY"
		   ParagraphPending = "PENDING"
		*/
	}

	return
}