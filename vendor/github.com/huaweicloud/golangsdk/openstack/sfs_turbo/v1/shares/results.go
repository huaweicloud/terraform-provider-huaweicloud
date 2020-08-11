package shares

import (
	"encoding/json"
	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// TurboResponse contains the information of creating response
type TurboResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// Turbo contains all information associated with an SFS Turbo file system
type Turbo struct {
	// The UUID of the SFS Turbo file system
	ID string `json:"id"`
	// The name of the SFS Turbo file system
	Name string `json:"name"`
	// Size of the share in GB
	Size string `json:"size"`
	// The statue of the SFS Turbo file system
	Status string `json:"status"`
	// The sub-statue of the SFS Turbo file system
	SubStatus string `json:"sub_status"`
	// The version ID of the SFS Turbo file system
	Version string `json:"version"`
	// The mount location
	ExportLocation string `json:"export_location"`
	// The creation progress of the SFS Turbo file system
	Actions []string `json:"actions"`
	// The protocol type of the SFS Turbo file system
	ShareProto string `json:"share_proto"`
	// The type of the SFS Turbo file system, STANDARD or PERFORMANCE.
	ShareType string `json:"share_type"`
	// The region of the SFS Turbo file system
	Region string `json:"region"`
	// The code of the availability zone
	AvailabilityZone string `json:"availability_zone"`
	// The name of the availability zone
	AZName string `json:"az_name"`
	// The VPC ID
	VpcID string `json:"vpc_id"`
	// The subnet ID
	SubnetID string `json:"subnet_id"`
	// The security group ID
	SecurityGroupID string `json:"security_group_id"`
	// The avaliable capacity if the SFS Turbo file system
	AvailCapacity string `json:"avail_capacity"`
	// bandwidth is returned for an enhanced file system
	ExpandType string `json:"expand_type"`
	// The ID of the encryption key
	CryptKeyID string `json:"crypt_key_id"`
	// The billing mode, 0 indicates pay-per-use, 1 indicates yearly/monthly subscription
	PayModel string `json:"pay_model"`
	// Timestamp when the share was created
	CreatedAt time.Time `json:"-"`
}

func (r *Turbo) UnmarshalJSON(b []byte) error {
	type tmp Turbo
	var s struct {
		tmp
		CreatedAt golangsdk.JSONRFC3339MilliNoZ `json:"created_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Turbo(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)

	return nil
}

type commonResult struct {
	golangsdk.Result
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

// GetResult contains the response body and error from a Get request.
type GetResult struct {
	commonResult
}

// DeleteResult contains the response body and error from a Delete request.
type DeleteResult struct {
	golangsdk.ErrResult
}

// ExpandResult contains the response body and error from a Expand request.
type ExpandResult struct {
	golangsdk.ErrResult
}

// Extract will get the Turbo response object from the CreateResult
func (r CreateResult) Extract() (*TurboResponse, error) {
	var resp TurboResponse
	err := r.ExtractInto(&resp)
	return &resp, err
}

// Extract will get the Turbo object from the GetResult
func (r GetResult) Extract() (*Turbo, error) {
	var object Turbo
	err := r.ExtractInto(&object)
	return &object, err
}

// TurboPage is the page returned by a pager when traversing over a
// collection of Shares.
type TurboPage struct {
	pagination.LinkedPageBase
}

// ExtractTurbos accepts a Page struct, specifically a TurboPage struct,
// and extracts the elements into a slice of share structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractTurbos(r pagination.Page) ([]Turbo, error) {
	var s struct {
		ListedShares []Turbo `json:"shares"`
	}
	err := (r.(TurboPage)).ExtractInto(&s)
	return s.ListedShares, err
}

// IsEmpty returns true if a ListResult contains no Shares.
func (r TurboPage) IsEmpty() (bool, error) {
	shares, err := ExtractTurbos(r)
	return len(shares) == 0, err
}

// NextPageURL is invoked when a paginated collection of shares has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r TurboPage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"shares_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}
