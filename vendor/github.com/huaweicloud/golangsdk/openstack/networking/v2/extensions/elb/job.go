package elb

import (
	"github.com/huaweicloud/golangsdk"
)

type Job struct {
	Uri   string `json:"uri"`
	JobId string `json:"job_id"`
}

type JobResult struct {
	golangsdk.Result
}

func (r JobResult) Extract() (*Job, error) {
	j := &Job{}
	err := r.ExtractInto(j)
	return j, err
}

type JobInfo struct {
	Status     string                 `json:"status"`
	Entities   map[string]interface{} `json:"entities"`
	JobId      string                 `json:"job_id"`
	JobType    string                 `json:"job_type"`
	BeginTime  string                 `json:"begin_time"`
	EndTime    string                 `json:"end_time"`
	ErrorCode  string                 `json:"error_code"`
	FailReason string                 `json:"fail_reason"`
	Message    string                 `json:"message"`
	Code       string                 `json:"code"`
	SubJobs    []JobInfo              `json:"sub_jobs"`
}

type JobInfoResult struct {
	golangsdk.Result
}

func (r JobInfoResult) Extract() (*JobInfo, error) {
	j := &JobInfo{}
	err := r.ExtractInto(j)
	return j, err
}

func QueryJobInfo(c *golangsdk.ServiceClient, jobId string) (r JobInfoResult) {
	e := c.ResourceBaseURL()
	e = e + c.ProjectID + "/jobs/" + jobId
	_, r.Err = c.Get(e, &r.Body, nil)
	return
}
