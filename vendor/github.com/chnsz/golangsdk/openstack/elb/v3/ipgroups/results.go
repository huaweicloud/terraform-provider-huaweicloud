package ipgroups

import (
	"github.com/chnsz/golangsdk"
)

type ListenerID struct {
	ID string `json:"id"`
}

type IpGroup struct {
	// The unique ID for the IpGroup.
	ID string `json:"id"`

	// Human-readable name for the IpGroup.
	Name string `json:"name"`

	// Human-readable description for the IpGroup.
	Description string `json:"description"`

	// whether to use HTTP2.
	Http2Enable bool `json:"http2_enable"`

	// A list of listener IDs.
	Listeners []ListenerID `json:"listeners"`

	// A list of IP addresses.
	IpList []IpListOpt `json:"ip_list"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a ipgroup.
func (r commonResult) Extract() (*IpGroup, error) {
	var s struct {
		IpGroup *IpGroup `json:"ipgroup"`
	}
	err := r.ExtractInto(&s)
	return s.IpGroup, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a IpGroup.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a IpGroup.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a IpGroup.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
