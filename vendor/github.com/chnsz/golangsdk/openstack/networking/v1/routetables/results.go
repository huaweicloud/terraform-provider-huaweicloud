package routetables

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// Route represents a route object in a route table
type Route struct {
	Type            string `json:"type"`
	DestinationCIDR string `json:"destination"`
	NextHop         string `json:"nexthop"`
	Description     string `json:"description"`
}

// Subnet represents a subnet object associated with a route table
type Subnet struct {
	ID string `json:"id"`
}

// RouteTable represents a route table
type RouteTable struct {
	// Name is the human readable name for the route table
	Name string `json:"name"`

	// Description is the supplementary information about the route table
	Description string `json:"description"`

	// ID is the unique identifier for the route table
	ID string `json:"id"`

	// the VPC ID that the route table belongs to.
	VpcID string `json:"vpc_id"`

	// project id
	TenantID string `json:"tenant_id"`

	// Default indicates whether it is a default route table
	Default bool `json:"default"`

	// Routes is an array of static routes that the route table will host
	Routes []Route `json:"routes"`

	// Subnets is an array of subnets that associated with the route table
	Subnets []Subnet `json:"subnets"`
}

// RouteTablePage is the page returned by a pager when traversing over a
// collection of route tables
type RouteTablePage struct {
	pagination.MarkerPageBase
}

// LastMarker returns the last route table ID in a ListResult
func (r RouteTablePage) LastMarker() (string, error) {
	tables, err := ExtractRouteTables(r)
	if err != nil {
		return "", err
	}
	if len(tables) == 0 {
		return "", nil
	}
	return tables[len(tables)-1].ID, nil
}

// IsEmpty checks whether a RouteTablePage struct is empty.
func (r RouteTablePage) IsEmpty() (bool, error) {
	tables, err := ExtractRouteTables(r)
	return len(tables) == 0, err
}

// ExtractRouteTables accepts a Page struct, specifically a RouteTablePage struct,
// and extracts the elements into a slice of RouteTable structs.
func ExtractRouteTables(r pagination.Page) ([]RouteTable, error) {
	var s struct {
		RouteTables []RouteTable `json:"routetables"`
	}
	err := (r.(RouteTablePage)).ExtractInto(&s)
	return s.RouteTables, err
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a route table
func (r commonResult) Extract() (*RouteTable, error) {
	var s struct {
		RouteTable *RouteTable `json:"routetable"`
	}
	err := r.ExtractInto(&s)
	return s.RouteTable, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a RouteTable
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a RouteTable
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a RouteTable
type UpdateResult struct {
	commonResult
}

// ActionResult represents the result of an action operation. Call its Extract
// method to interpret it as a RouteTable
type ActionResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
