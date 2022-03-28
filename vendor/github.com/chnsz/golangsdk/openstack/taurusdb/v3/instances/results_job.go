package instances

import (
	"fmt"
	"time"

	"github.com/chnsz/golangsdk"
)

type JobResponse struct {
	JobID   string `json:"job_id"`
	OrderID string `json:"order_id"`
}

type JobStatus struct {
	Job Job `json:"job"`
}

type Job struct {
	Status     string `json:"status"`
	JobID      string `json:"id"`
	FailReason string `json:"fail_reason"`
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
		job_status := new(JobStatus)
		url := jobURL(client)
		url += fmt.Sprintf("?id=%s", jobID)
		_, err := client.Get(url, &job_status, &golangsdk.RequestOpts{
			MoreHeaders: requestOpts.MoreHeaders,
		})
		time.Sleep(15 * time.Second)
		if err != nil {
			return false, err
		}
		job := job_status.Job

		if job.Status == "Completed" {
			return true, nil
		}
		if job.Status == "Failed" {
			err = fmt.Errorf("Job failed: %s.", job.FailReason)
			return false, err
		}

		return false, nil
	})
}
