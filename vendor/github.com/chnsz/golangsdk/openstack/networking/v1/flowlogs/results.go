package flowlogs

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// VPC flow log struct
type FlowLog struct {
	// Specifies the VPC flow log UUID.
	ID string `json:"id"`

	// Specifies the VPC flow log name.
	Name string `json:"name"`

	// Provides supplementary information about the VPC flow log.
	Description string `json:"description"`

	// Specifies the type of resource on which to create the VPC flow log.
	ResourceType string `json:"resource_type"`

	// Specifies the unique resource ID.
	ResourceID string `json:"resource_id"`

	// Specifies the type of traffic to log.
	TrafficType string `json:"traffic_type"`

	// Specifies the log group ID..
	LogGroupID string `json:"log_group_id"`

	// Specifies the log topic ID.
	LogTopicID string `json:"log_topic_id"`

	// Specifies the VPC flow log status, the value can be ACTIVE, DOWN or ERROR.
	Status string `json:"status"`

	// Specifies the project ID.
	TenantID string `json:"tenant_id"`

	// Specifies whether to enable the VPC flow log function.
	AdminState bool `json:"admin_state"`

	// Specifies the time when the VPC flow log was created.
	CreatedAt string `json:"created_at"`

	// Specifies the time when the VPC flow log was updated.
	UpdatedAt string `json:"updated_at"`
}

// FlowLogPage is the page returned by a pager when traversing over a collection
// of flow logs.
type FlowLogPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of flow logs has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r FlowLogPage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"flowlogs_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a FlowLogPage struct is empty.
func (r FlowLogPage) IsEmpty() (bool, error) {
	is, err := ExtractFlowLogs(r)
	return len(is) == 0, err
}

// ExtractFlowLogs accepts a Page struct, specifically a FlowLogPage struct,
// and extracts the elements into a slice of FlowLog structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractFlowLogs(r pagination.Page) ([]FlowLog, error) {
	var s struct {
		FlowLogs []FlowLog `json:"flow_logs"`
	}
	err := (r.(FlowLogPage)).ExtractInto(&s)
	return s.FlowLogs, err
}

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

func (r CreateResult) Extract() (*FlowLog, error) {
	var entity FlowLog
	err := r.ExtractIntoStructPtr(&entity, "flow_log")
	return &entity, err
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type GetResult struct {
	commonResult
}

func (r GetResult) Extract() (*FlowLog, error) {
	var entity FlowLog
	err := r.ExtractIntoStructPtr(&entity, "flow_log")
	return &entity, err
}

type UpdateResult struct {
	commonResult
}

func (r UpdateResult) Extract() (*FlowLog, error) {
	var entity FlowLog
	err := r.ExtractIntoStructPtr(&entity, "flow_log")
	return &entity, err
}
