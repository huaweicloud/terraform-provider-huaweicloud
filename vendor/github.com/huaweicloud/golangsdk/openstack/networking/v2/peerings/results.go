package peerings

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type VpcInfo struct {
	VpcId    string `json:"vpc_id" required:"true"`
	TenantId string `json:"tenant_id,omitempty"`
}

// Peering represents a Neutron VPC peering connection.
//Manage and perform other operations on VPC peering connections,
// including querying VPC peering connections as well as
// creating, querying, deleting, and updating a VPC peering connection.
type Peering struct {
	// ID is the unique identifier for the vpc_peering_connection.
	ID string `json:"id"`

	// Name is the human readable name for the vpc_peering_connection. It does not have to be
	// unique.
	Name string `json:"name"`

	// Status indicates whether or not a vpc_peering_connections is currently operational.
	Status string `json:"status"`

	// RequestVpcInfo indicates information about the local VPC
	RequestVpcInfo VpcInfo `json:"request_vpc_info"`

	// AcceptVpcInfo indicates information about the peer  VPC
	AcceptVpcInfo VpcInfo `json:"accept_vpc_info"`
}

// PeeringConnectionPage is the page returned by a pager when traversing over a
// collection of vpc_peering_connections.
type PeeringConnectionPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of vpc_peering_connections has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r PeeringConnectionPage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"peerings_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a PeeringConnectionPage struct is empty.
func (r PeeringConnectionPage) IsEmpty() (bool, error) {
	is, err := ExtractPeerings(r)
	return len(is) == 0, err
}

// ExtractPeerings accepts a Page struct, specifically a PeeringConnectionPage struct,
// and extracts the elements into a slice of Peering structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractPeerings(r pagination.Page) ([]Peering, error) {
	var s struct {
		Peerings []Peering `json:"peerings"`
	}
	err := (r.(PeeringConnectionPage)).ExtractInto(&s)
	return s.Peerings, err
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a Peering.
func (r commonResult) Extract() (*Peering, error) {
	var s struct {
		Peering *Peering `json:"peering"`
	}
	err := r.ExtractInto(&s)
	return s.Peering, err
}

// ExtractResult is a function that accepts a result and extracts a Peering.
func (r commonResult) ExtractResult() (Peering, error) {
	var s struct {
		// ID is the unique identifier for the vpc.
		ID string `json:"id"`
		// Name is the human readable name for the vpc. It does not have to be
		// unique.
		Name string `json:"name"`

		// Status indicates whether or not a vpc is currently operational.
		Status string `json:"status"`

		// Status indicates whether or not a vpc is currently operational.
		RequestVpcInfo VpcInfo `json:"request_vpc_info"`

		//Provides informaion about shared snat
		AcceptVpcInfo VpcInfo `json:"accept_vpc_info"`
	}
	err1 := r.ExtractInto(&s)
	return s, err1
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Vpc Peering Connection.
type GetResult struct {
	commonResult
}

// AcceptResult represents the result of a get operation. Call its Extract
// method to interpret it as a Vpc Peering Connection.
type AcceptResult struct {
	commonResult
}

// RejectResult represents the result of a get operation. Call its Extract
// method to interpret it as a Vpc Peering Connection.
type RejectResult struct {
	commonResult
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a vpc peering connection.
type CreateResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a vpc peering connection.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
