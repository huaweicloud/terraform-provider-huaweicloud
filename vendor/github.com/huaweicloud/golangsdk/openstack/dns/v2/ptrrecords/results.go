package ptrrecords

import (
	"github.com/huaweicloud/golangsdk"
)

type commonResult struct {
	golangsdk.Result
}

// Extract interprets a GetResult, CreateResult as a Ptr.
// An error is returned if the original call or the extraction failed.
func (r commonResult) Extract() (*Ptr, error) {
	var s *Ptr
	err := r.ExtractInto(&s)
	return s, err
}

// CreateResult is the result of a Create request. Call its Extract method
// to interpret the result as a Ptr.
type CreateResult struct {
	commonResult
}

// GetResult is the result of a Get request. Call its Extract method
// to interpret the result as a Ptr.
type GetResult struct {
	commonResult
}

// DeleteResult is the result of a Delete request. Call its ExtractErr method
// to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}

// Ptr represents a ptr record.
type Ptr struct {
	// ID uniquely identifies this ptr amongst all other ptr records.
	ID string `json:"id"`

	// Name for this ptr.
	PtrName string `json:"ptrdname"`

	// Description for this ptr.
	Description string `json:"description"`

	// TTL is the Time to Live for the ptr.
	TTL int `json:"ttl"`

	// Address of the floating ip.
	Address string `json:"address"`

	// Status of the PTR.
	Status string `json:"status"`
}
