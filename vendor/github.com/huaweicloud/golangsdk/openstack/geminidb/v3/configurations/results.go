package configurations

import (
	"github.com/huaweicloud/golangsdk"
)

type ApplyResponse struct {
	JobId   string `json:"job_id"`
	Success bool   `json:"success"`
}

type Parameter struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	ReadOnly bool   `json:"readonly"`
}

type GetResponse struct {
	Parameters []Parameter `json:"configuration_parameters"`
}

type commonResult struct {
	golangsdk.Result
}

type ApplyResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

func (r ApplyResult) Extract() (*ApplyResponse, error) {
	var response ApplyResponse
	err := r.ExtractInto(&response)
	return &response, err
}

func (r GetResult) Extract() (*GetResponse, error) {
	var response GetResponse
	err := r.ExtractInto(&response)
	return &response, err
}
