package routes

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type Route struct {
	// Specifies the route type.
	Type string `json:"type"`

	// Specifies the next hop. If the route type is peering, enter the VPC peering connection ID.
	NextHop string `json:"nexthop"`

	//Specifies the destination IP address or CIDR block.
	Destination string `json:"destination"`

	// Specifies the VPC for which a route is to be added.
	VPC_ID string `json:"vpc_id"`

	//Specifies the tenant ID. Only the administrator can specify the tenant ID of other tenants.
	Tenant_Id string `json:"tenant_id"`

	//Specifies the route ID.
	RouteID string `json:"id"`
}

// RoutePage is the page returned by a pager when traversing over a
// collection of routes.
type RoutePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of routes has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r RoutePage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"routes_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a RoutePage struct is empty.
func (r RoutePage) IsEmpty() (bool, error) {
	is, err := ExtractRoutes(r)
	return len(is) == 0, err
}

// ExtractRoutes accepts a Page struct, specifically a RoutePage struct,
// and extracts the elements into a slice of Roue structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractRoutes(r pagination.Page) ([]Route, error) {
	var s struct {
		Routes []Route `json:"routes"`
	}
	err := (r.(RoutePage)).ExtractInto(&s)
	return s.Routes, err
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a Route.
func (r commonResult) Extract() (*Route, error) {
	var s struct {
		Route *Route `json:"route"`
	}
	err := r.ExtractInto(&s)
	return s.Route, err
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
