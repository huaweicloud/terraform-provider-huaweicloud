package configurations

import "github.com/chnsz/golangsdk"

type Configuration struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	DataStoreVer  string `json:"datastore_version_name"`
	DataStoreName string `json:"datastore_name"`
}

type ListResult struct {
	golangsdk.Result
}

func (lr ListResult) Extract() ([]Configuration, error) {
	var a struct {
		Configurations []Configuration `json:"configurations"`
	}
	err := lr.Result.ExtractInto(&a)
	return a.Configurations, err
}

type JobResponse struct {
	JobID string `json:"job_id"`
}

type JobResult struct {
	golangsdk.Result
}

func (r JobResult) ExtractJobResponse() (*JobResponse, error) {
	var job JobResponse
	err := r.ExtractInto(&job)
	return &job, err
}
