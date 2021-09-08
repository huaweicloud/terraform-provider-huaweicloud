package routetables
import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type Route struct {
	Type  string `json:"type"`
	Destination string `json:"destination,omitempty"`
	Nexthop string `json:"nexthop"`
	System string `json:"system,omitempty"`
}

type Subnet struct {
	Id  string `json:"id"`
}

type RouteTable struct {
	// Specifies the route name.
	Name string `json:"name"`

	// Specifies the routes. 
	Routes []Route `json:"routes"`

	// Specifies the subnets. 
	Subnets []Subnet `json:"subnets"`

	//Specifies the destination .
	Destination string `json:"destination"`

	// Specifies the VPC .
	VPC_ID string `json:"vpc_id"`

	//Specifies the tenant ID. Only the administrator can specify the tenant ID of other tenants.
	Tenant_Id string `json:"tenant_id"`

	//Specifies the route ID.
	RouteID string `json:"id"`
}

// RoutePage is the page returned by a pager when traversing over a
// collection of routes.
type RouteTablePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of routes has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r RouteTablePage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"routetables_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a RoutePage struct is empty.
func (r RouteTablePage) IsEmpty() (bool, error) {
	is, err := ExtractRouteTables(r)
	return len(is) == 0, err
}

// ExtractRoutes accepts a Page struct, specifically a RoutePage struct,
// and extracts the elements into a slice of Roue structs. In other words,
// a generic collection is mapped into a relevant slice.
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

// Extract is a function that accepts a result and extracts a Route.
func (r commonResult) Extract() (*RouteTable, error) {
	var s struct {
		RouteTable *RouteTable `json:"routetable"`
	}
	err := r.ExtractInto(&s)
	return s.RouteTable, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Route.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Route.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}

