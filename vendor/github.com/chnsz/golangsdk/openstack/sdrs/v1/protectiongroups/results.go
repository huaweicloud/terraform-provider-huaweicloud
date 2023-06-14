package protectiongroups

import (
	"github.com/chnsz/golangsdk"
)

type Group struct {
	//Group ID
	Id string `json:"id"`
	//Group Name
	Name string `json:"name"`
	//Group Description
	Description string `json:"description"`
	//The source AZ of a protection group
	SourceAZ string `json:"source_availability_zone"`
	//The target AZ of a protection group
	TargetAZ string `json:"target_availability_zone"`
	//An active-active domain
	DomainID string `json:"domain_id"`
	//ID of the source VPC
	SourceVpcID string `json:"source_vpc_id"`
	//Deployment model
	DrType string `json:"dr_type"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a group.
func (r commonResult) Extract() (*Group, error) {
	var response Group
	err := r.ExtractInto(&response)
	return &response, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "server_group")
}

// UpdateResult represents the result of a update operation. Call its Extract
// method to interpret it as a Group.
type UpdateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Group.
type GetResult struct {
	commonResult
}
