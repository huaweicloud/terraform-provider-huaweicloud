package protectedinstances

import (
	"github.com/chnsz/golangsdk"
)

type Instance struct {
	//Instance ID
	Id string `json:"id"`
	//Instance Name
	Name string `json:"name"`
	//Instance Description
	Description string `json:"description"`
	//Protection Group ID
	GroupID string `json:"server_group_id"`
	//Instance Status
	Status string `json:"status"`
	//Source Server
	SourceServer string `json:"source_server"`
	//Target Server
	TargetServer string `json:"target_server"`
	//Attachment
	Attachment []Attachment `json:"attachment"`
}

type Attachment struct {
	//Replication ID
	Replication string `json:"replication"`
	//Device Name
	Device string `json:"device"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a instance.
func (r commonResult) Extract() (*Instance, error) {
	var response Instance
	err := r.ExtractInto(&response)
	return &response, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "protected_instance")
}

// UpdateResult represents the result of a update operation. Call its Extract
// method to interpret it as a Instance.
type UpdateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Instance.
type GetResult struct {
	commonResult
}
