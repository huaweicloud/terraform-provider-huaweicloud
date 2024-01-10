package cloudvolumes

import (
	"github.com/chnsz/golangsdk"
)

// JobResult contains the response body and error from UpdateQoS response
type JobResult struct {
	golangsdk.Result
}

// JobResponse contains all the information from UpdateQoS response
type JobResponse struct {
	JobID string `json:"job_id"`
}

// Extract will get the JobResponse object out of the JobResult
func (r JobResult) Extract() (*JobResponse, error) {
	job := new(JobResponse)
	err := r.ExtractInto(job)
	return job, err
}
