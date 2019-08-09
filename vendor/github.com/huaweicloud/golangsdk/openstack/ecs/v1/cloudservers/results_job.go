package cloudservers

import (
	"fmt"
	"time"

	"github.com/huaweicloud/golangsdk"
)

type JobResponse struct {
	JobID string `json:"job_id"`
}

type JobStatus struct {
	Status     string    `json:"status"`
	Entities   JobEntity `json:"entities"`
	JobID      string    `json:"job_id"`
	JobType    string    `json:"job_type"`
	BeginTime  string    `json:"begin_time"`
	EndTime    string    `json:"end_time"`
	ErrorCode  string    `json:"error_code"`
	FailReason string    `json:"fail_reason"`
}

type JobEntity struct {
	// Specifies the number of subtasks.
	// When no subtask exists, the value of this parameter is 0.
	SubJobsTotal int `json:"sub_jobs_total"`

	// Specifies the execution information of a subtask.
	// When no subtask exists, the value of this parameter is left blank.
	SubJobs []SubJob `json:"sub_jobs"`
}

type SubJob struct {
	// Specifies the task ID.
	Id string `json:"job_id"`

	// Task type.
	Type string `json:"job_type"`

	//Specifies the task status.
	//  SUCCESS: indicates the task is successfully executed.
	//  RUNNING: indicates that the task is in progress.
	//  FAIL: indicates that the task failed.
	//  INIT: indicates that the task is being initialized.
	Status string `json:"status"`

	// Specifies the time when the task started.
	BeginTime string `json:"begin_time"`

	// Specifies the time when the task finished.
	EndTime string `json:"end_time"`

	// Specifies the returned error code when the task execution fails.
	ErrorCode string `json:"error_code"`

	// Specifies the cause of the task execution failure.
	FailReason string `json:"fail_reason"`

	// Specifies the object of the task.
	Entities map[string]string `json:"entities"`
}

type JobResult struct {
	golangsdk.Result
}

func (r JobResult) ExtractJobResponse() (*JobResponse, error) {
	job := new(JobResponse)
	err := r.ExtractInto(job)
	return job, err
}

func (r JobResult) ExtractJobStatus() (*JobStatus, error) {
	job := new(JobStatus)
	err := r.ExtractInto(job)
	return job, err
}

func WaitForJobSuccess(client *golangsdk.ServiceClient, secs int, jobID string) error {
	return golangsdk.WaitFor(secs, func() (bool, error) {
		job := new(JobStatus)
		_, err := client.Get(jobURL(client, jobID), &job, nil)
		time.Sleep(5 * time.Second)
		if err != nil {
			return false, err
		}

		if job.Status == "SUCCESS" {
			return true, nil
		}
		if job.Status == "FAIL" {
			err = fmt.Errorf("Job failed with code %s: %s.\n", job.ErrorCode, job.FailReason)
			return false, err
		}

		return false, nil
	})
}

func GetJobEntity(client *golangsdk.ServiceClient, jobID string, label string) (interface{}, error) {
	job := new(JobStatus)
	_, err := client.Get(jobURL(client, jobID), &job, nil)
	if err != nil {
		return nil, err
	}

	if job.Status == "SUCCESS" {
		if e := job.Entities.SubJobs[0].Entities[label]; e != "" {
			return e, nil
		}
	}

	return nil, fmt.Errorf("Unexpected conversion error in GetJobEntity.")
}
