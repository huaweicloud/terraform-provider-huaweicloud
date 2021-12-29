package jobs

import (
	"github.com/chnsz/golangsdk"
)

// Job object contains the response to a Get request
type Job struct {
	Status     string    `json:"status"`
	Entities   JobEntity `json:"entities"`
	JobID      string    `json:"job_id"`
	JobType    string    `json:"job_type"`
	ErrorCode  string    `json:"error_code"`
	FailReason string    `json:"fail_reason"`
	Error      ErrorInfo `json:"error"`
	BeginTime  string    `json:"begin_time"`
	EndTime    string    `json:"end_time"`
}

// JobEntity contains the response to the job task
type JobEntity struct {
	Name       string `json:"name"`
	Size       int    `json:"size"`
	VolumeID   string `json:"volume_id"`
	VolumeType string `json:"volume_type"`
	SubJobs    []Job  `json:"sub_jobs"`
}

// ErrorInfo contains the error message returned when an error occurs
type ErrorInfo struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// GetResult contains the response body and error from a Get request
type GetResult struct {
	golangsdk.Result
}

// ExtractJob will get the *Job object out of the GetResult
func (r GetResult) ExtractJob() (*Job, error) {
	var job Job
	err := r.ExtractInto(&job)
	return &job, err
}
