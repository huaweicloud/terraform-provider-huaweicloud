package tags

import (
	"github.com/huaweicloud/golangsdk"
)

type commonResult struct {
	golangsdk.Result
}

// Tags model
type Tags struct {
	// Tags is a list of any tags. Tags are arbitrarily defined strings
	// attached to a resource.
	Tags map[string]string `json:"tags"`
}

// Extract interprets any commonResult as a Tags.
func (r commonResult) Extract() (*Tags, error) {
	var s *Tags
	err := r.ExtractInto(&s)
	return s, err
}

// CreateResult represents the result of a Create operation
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a Get operation
type GetResult struct {
	commonResult
}

//DeleteResult model
type DeleteResult struct {
	golangsdk.ErrResult
}
