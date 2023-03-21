package imagecopy

import "github.com/chnsz/golangsdk"

type JobResponse struct {
	JobID string `json:"job_id"`
}

type JobResult struct {
	golangsdk.Result
}

func (r JobResult) ExtractJobResponse() (*JobResponse, error) {
	job := new(JobResponse)
	err := r.ExtractInto(job)
	return job, err
}
