package routetables

import (
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/pagination"
)

// createResp is the structure that represents the API response of the 'Create' method, which contains route table
// details and the request information.
type createResp struct {
	// The response detail of the route table.
	RouteTable RouteTable `json:"route_table"`
	// The request ID.
	RequestId string `json:"request_id"`
}

// RouteTable is the structure that represents the details of the route table for ER service.
type RouteTable struct {
	// The ID of the route table.
	ID string `json:"id"`
	// The name of the route table.
	// The value can contain 1 to 64 characters, only english and chinese letters, digits, underscore (_), hyphens (-)
	// and dots (.) are allowed.
	Name string `json:"name"`
	// The description of the route table.
	// The value contain a maximum of 255 characters, and the angle brackets (< and >) are not allowed.
	Description string `json:"description"`
	// The configuration of the BGP route selection.
	BgpOptions BgpOptions `json:"bgp_options"`
	// The key/value pairs to associate with the route table.
	Tags []tags.ResourceTag `json:"tags"`
	// Whether this route table is the default association route table.
	IsDefaultAssociation bool `json:"is_default_association"`
	// Whether this route table is the default propagation route table.
	IsDefaultPropagation bool `json:"is_default_propagation"`
	// The current status of the route table.
	Status string `json:"state"`
	// The creation time of the route table.
	CreatedAt string `json:"created_at"`
	// The last update time of the route table.
	UpdatedAt string `json:"updated_at"`
}

// getResp is the structure that represents the API response of the 'Get' method, which contains route table details and
// the request information.
type getResp struct {
	// The response detail of the route table.
	RouteTable RouteTable `json:"route_table"`
	// The request ID.
	RequestId string `json:"request_id"`
}

// listResp is the structure that represents the API response of the 'List' method, which contains route table list,
// page details and the request information.
type listResp struct {
	// The list of the route tables.
	RouteTables []RouteTable `json:"route_tables"`
	// The request ID.
	RequestId string `json:"request_id"`
	// The page information.
	PageInfo pageInfo `json:"page_info"`
}

// pageInfo is the structure that represents the page information.
type pageInfo struct {
	// The next marker information.
	NextMarker string `json:"next_marker"`
	// The number of the route table in current page.
	CurrentCount int `json:"current_count"`
}

// RouteTablePage represents the response pages of the List method.
type RouteTablePage struct {
	pagination.MarkerPageBase
}

// IsEmpty returns true if a ListResult no route table.
func (r RouteTablePage) IsEmpty() (bool, error) {
	resp, err := extractRouteTables(r)
	return len(resp) == 0, err
}

// LastMarker returns the last marker index in a ListResult.
func (r RouteTablePage) LastMarker() (string, error) {
	resp, err := extractPageInfo(r)
	if err != nil {
		return "", err
	}
	if resp.NextMarker != "" {
		return "", nil
	}
	return resp.NextMarker, nil
}

// extractPageInfo is a method which to extract the response of the page information.
func extractPageInfo(r pagination.Page) (*pageInfo, error) {
	var s listResp
	err := r.(RouteTablePage).Result.ExtractInto(&s)
	return &s.PageInfo, err
}

// extractRouteTables is a method which to extract the response to a route table list.
func extractRouteTables(r pagination.Page) ([]RouteTable, error) {
	var s listResp
	err := r.(RouteTablePage).Result.ExtractInto(&s)
	return s.RouteTables, err
}

// updateResp is the structure that represents the API response of the 'Update' method, which contains route table
// details and the request information.
type updateResp struct {
	// The response detail of the route table.
	RouteTable RouteTable `json:"route_table"`
	// The request ID.
	RequestId string `json:"request_id"`
}
