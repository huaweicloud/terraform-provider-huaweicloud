package certificates

import (
	"github.com/huaweicloud/golangsdk"
)

type Certificate struct {
	// Certificate ID
	Id string `json:"id"`
	// Certificate Name
	Name string `json:"name"`
	// the time when the certificate expires in unix timestamp
	ExpireTime int `json:"expireTime"`
	// the time when the certificate is uploaded in unix timestamp
	TimeStamp int `json:"timestamp"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a certificate.
func (r commonResult) Extract() (*Certificate, error) {
	var response Certificate
	err := r.ExtractInto(&response)
	return &response, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Certificate.
type CreateResult struct {
	commonResult
}

// UpdateResult represents the result of a update operation. Call its Extract
// method to interpret it as a Certificate.
type UpdateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Certificate.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
