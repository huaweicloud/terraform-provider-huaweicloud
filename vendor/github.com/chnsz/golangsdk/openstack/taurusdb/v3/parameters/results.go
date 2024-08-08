package parameters

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type JobResponse struct {
	JobID           string `json:"job_id"`
	RestartRequired bool   `json:"restart_required"`
}

type JobResult struct {
	golangsdk.Result
}

func (r JobResult) ExtractJobResponse() (*JobResponse, error) {
	var job JobResponse
	err := r.ExtractInto(&job)
	return &job, err
}

type ParameterPage struct {
	pagination.OffsetPageBase
}

func (r ParameterPage) IsEmpty() (bool, error) {
	data, err := ExtractParameters(r)
	if err != nil {
		return false, err
	}
	return len(data.ParameterValues) == 0, err
}

func ExtractParameters(r pagination.Page) (ListParametersResponse, error) {
	var s ListParametersResponse
	err := (r.(ParameterPage)).ExtractInto(&s)
	return s, err
}

type ListParametersResponse struct {
	Configurations  Configurations   `json:"configurations"`
	ParameterValues []ParameterValue `json:"parameter_values"`
	TotalCount      int              `json:"total_count"`
}

type Configurations struct {
	// the version of the datastore
	DatastoreVersionName string `json:"datastore_version_name"`
	// the name of the datastore
	DatastoreName string `json:"datastore_name"`
	// the created time
	Created string `json:"created"`
	// the update time
	Updated string `json:"updated"`
}

type ParameterValue struct {
	// the parameter name
	Name string `json:"name"`
	// the parameter value
	Value string `json:"value"`
	// whether the instance need restart
	RestartRequired bool `json:"restart_required"`
	// whether the parameter is readonly
	Readonly bool `json:"readonly"`
	// the parameter value range
	ValueRange string `json:"value_range"`
	// the parameter type
	Type string `json:"type"`
	// the description
	Description string `json:"description"`
}
