package logtanks

import (
	"github.com/chnsz/golangsdk"
)

type LogTank struct {
	// The unique ID for the LogTank.
	ID string `json:"id"`

	// The Loadbalancer on which the log associated with.
	LoadbalancerID string `json:"loadbalancer_id"`

	// The log group on which the log associated with.
	LogGroupId string `json:"log_group_id"`

	// The topic on which the log subscribe.
	LogTopicId string `json:"log_topic_id"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a logtank.
func (r commonResult) Extract() (*LogTank, error) {
	var s struct {
		LogTank *LogTank `json:"logtank"`
	}
	err := r.ExtractInto(&s)
	return s.LogTank, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a LogTank.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a LogTank.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a LogTank.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
