package policies

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type Policy struct {
	//Backup policy ID
	ID string `json:"backup_policy_id"`
	//Backup policy name
	Name string `json:"backup_policy_name"`
	//Details about the scheduling policy
	ScheduledPolicy ScheduledPolicy `json:"scheduled_policy"`
	//Number of volumes associated with the backup policy
	ResourceCount int `json:"policy_resource_count"`
}

type ResultResources struct {
	//List of successfully associated/disassociated resources
	SuccessResources []Resource `json:"success_resources"`
	//List of the resources that fail to be associated/disassociated
	FailResources []Resource `json:"fail_resources"`
}

type Resource struct {
	//Resource ID
	ResourceID string `json:"resource_id"`
	//Resource type
	ResourceType string `json:"resource_type"`
	//Availability zone to which the resource belongs
	AvailabilityZone string `json:"availability_zone"`
	//POD to which the resource belongs
	Pod string `json:"os_vol_host_attr"`
}

type commonResult struct {
	golangsdk.Result
}

// CreateResult represents the result of a create operation.
type CreateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	commonResult
}

// UpdateResult represents the result of a update operation.
type UpdateResult struct {
	commonResult
}

// ResourceResult represents the result of a associate/diassociate operation.
type ResourceResult struct {
	commonResult
}

// Extract will get the Policy object from the commonResult
func (r commonResult) Extract() (*Policy, error) {
	var response Policy
	err := r.ExtractInto(&response)
	return &response, err
}

type PolicyPage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if a page contains no Policies results.
func (r PolicyPage) IsEmpty() (bool, error) {
	s, err := ExtractPolicies(r)
	return len(s) == 0, err
}

// NextPageURL uses the response's embedded link reference to navigate to the
// next page of results.
func (r PolicyPage) NextPageURL() (string, error) {
	var s struct {
		Policies []golangsdk.Link `json:"backup_policies_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Policies)
}

// ExtractPolicies accepts a Page struct, specifically a PolicyPage struct,
// and extracts the elements into a slice of Policy structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractPolicies(r pagination.Page) ([]Policy, error) {
	var s struct {
		Policies []Policy `json:"backup_policies"`
	}
	err := (r.(PolicyPage)).ExtractInto(&s)
	return s.Policies, err
}

// ExtractResource will get the Association/disassociation of resources from the ResourceResult
func (r ResourceResult) ExtractResource() (*ResultResources, error) {
	var response ResultResources
	err := r.ExtractInto(&response)
	return &response, err
}
