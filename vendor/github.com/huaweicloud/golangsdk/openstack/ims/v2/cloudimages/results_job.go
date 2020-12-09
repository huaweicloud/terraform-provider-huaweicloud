package cloudimages

import (
	"fmt"

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
	ImageID     string `json:"image_id"`
	DataImageID string `json:"__data_images"`
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
	jobClient := *client
	// v1/{project_id}/jobs/{job_id}
	jobClient.ResourceBase = jobClient.Endpoint + "v1/" + jobClient.ProjectID + "/"
	return golangsdk.WaitFor(secs, func() (bool, error) {
		job := new(JobStatus)
		_, err := jobClient.Get(jobClient.ServiceURL("jobs", jobID), &job, nil)
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

func GetJobEntity(client *golangsdk.ServiceClient, jobId string, label string) (interface{}, error) {
	if label != "image_id" && label != "__data_images" {
		return nil, fmt.Errorf("Unsupported label %s in GetJobEntity.", label)
	}

	jobClient := *client
	// v1/{project_id}/jobs/{job_id}
	jobClient.ResourceBase = jobClient.Endpoint + "v1/" + jobClient.ProjectID + "/"
	job := new(JobStatus)
	_, err := jobClient.Get(jobClient.ServiceURL("jobs", jobId), &job, nil)
	if err != nil {
		return nil, err
	}

	if job.Status == "SUCCESS" {
		if e := job.Entities.ImageID; e != "" {
			return e, nil
		}
		if e := job.Entities.DataImageID; e != "" {
			return e, nil
		}
	}

	return nil, fmt.Errorf("Unexpected conversion error in GetJobEntity.")
}
