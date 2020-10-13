package vpcs

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// Route is a possible route in a vpc.
type Route struct {
	NextHop         string `json:"nexthop"`
	DestinationCIDR string `json:"destination"`
}

// Vpc represents a Neutron vpc. A vpc is a logical entity that
// forwards packets across internal subnets and NATs (network address
// translation) them on external networks through an appropriate gateway.
//
// A vpc has an interface for each subnet with which it is associated. By
// default, the IP address of such interface is the subnet's gateway IP. Also,
// whenever a vpc is associated with a subnet, a port for that vpc
// interface is added to the subnet's network.
type Vpc struct {
	// ID is the unique identifier for the vpc.
	ID string `json:"id"`

	// Name is the human readable name for the vpc. It does not have to be
	// unique.
	Name string `json:"name"`

	//Specifies the range of available subnets in the VPC.
	CIDR string `json:"cidr"`

	//Enterprise Project ID.
	EnterpriseProjectID string `json:"enterprise_project_id"`

	// Status indicates whether or not a vpc is currently operational.
	Status string `json:"status"`

	// Routes are a collection of static routes that the vpc will host.
	Routes []Route `json:"routes"`

	//Provides informaion about shared snat
	EnableSharedSnat bool `json:"enable_shared_snat"`
}

// VpcPage is the page returned by a pager when traversing over a
// collection of vpcs.
type VpcPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of vpcs has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r VpcPage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"vpcs_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a VpcPage struct is empty.
func (r VpcPage) IsEmpty() (bool, error) {
	is, err := ExtractVpcs(r)
	return len(is) == 0, err
}

// ExtractVpcs accepts a Page struct, specifically a VpcPage struct,
// and extracts the elements into a slice of Vpc structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractVpcs(r pagination.Page) ([]Vpc, error) {
	var s struct {
		Vpcs []Vpc `json:"vpcs"`
	}
	err := (r.(VpcPage)).ExtractInto(&s)
	return s.Vpcs, err
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a vpc.
func (r commonResult) Extract() (*Vpc, error) {
	var s struct {
		Vpc *Vpc `json:"vpc"`
	}
	err := r.ExtractInto(&s)
	return s.Vpc, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Vpc.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Vpc.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Vpc.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
