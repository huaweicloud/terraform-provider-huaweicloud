package routes

import (
	"github.com/chnsz/golangsdk/pagination"
)

// createResp is the structure that represents the API response of the 'Create' method, which contains route details and
// the request information.
type createResp struct {
	// The response detail of the route.
	Route Route `json:"route"`
	// The request ID.
	RequestId string `json:"request_id"`
}

// Route is the structure that represents the details of the route under a specified route table.
type Route struct {
	// The ID of the route.
	ID string `json:"id"`
	// The type of the route.
	Type string `json:"type"`
	// Whether route is the black hole route.
	IsBlackHole bool `json:"is_blackhole"`
	// The destination of the route.
	Destination string `json:"destination"`
	// The corresponding attachments.
	Attachments []Attachment `json:"attachments"`
	// The ID of the route table to which the route belongs.
	RouteTableId string `json:"route_table_id"`
	// The current status of the route.
	Status string `json:"state"`
	// The creation time of the association.
	CreatedAt string `json:"created_at"`
	// The last update time of the association.
	UpdatedAt string `json:"updated_at"`
}

// Attachment is an object that represents the details of the corresponding route.
type Attachment struct {
	// The resource ID for the corresponding attachment.
	ResourceId string `json:"resource_id"`
	// The resource type for the corresponding attachment.
	ResourceType string `json:"resource_type"`
	// The ID of the corresponding attachment.
	AttachmentId string `json:"attachment_id"`
}

// getResp is the structure that represents the API response of the 'Get' method, which contains route details and the
// request information.
type getResp struct {
	// The response detail of the route.
	Route Route `json:"route"`
	// The request ID.
	RequestId string `json:"request_id"`
}

// listResp is the structure that represents the API response of the 'List' method, which contains route list, page
// details and the request information.
type listResp struct {
	// The list of the routes.
	Routes []Route `json:"routes"`
	// The request ID.
	RequestId string `json:"request_id"`
	// The page information.
	PageInfo pageInfo `json:"page_info"`
}

// pageInfo is the structure that represents the page information.
type pageInfo struct {
	// The next marker information.
	NextMarker string `json:"next_marker"`
	// The number of the associations in current page.
	CurrentCount int `json:"current_count"`
}

// RoutePage represents the response pages of the List method.
type RoutePage struct {
	pagination.MarkerPageBase
}

// IsEmpty returns true if a ListResult no route.
func (r RoutePage) IsEmpty() (bool, error) {
	resp, err := extractRoutes(r)
	return len(resp) == 0, err
}

// LastMarker returns the last marker index in a ListResult.
func (r RoutePage) LastMarker() (string, error) {
	resp, err := extractPageInfo(r)
	if err != nil {
		return "", err
	}
	if resp.NextMarker != "" {
		return "", nil
	}
	return resp.NextMarker, nil
}

// extractRoutes is a method which to extract the response of the page information.
func extractPageInfo(r pagination.Page) (*pageInfo, error) {
	var s listResp
	err := r.(RoutePage).Result.ExtractInto(&s)
	return &s.PageInfo, err
}

// extractRoutes is a method which to extract the response to a route list.
func extractRoutes(r pagination.Page) ([]Route, error) {
	var s listResp
	err := r.(RoutePage).Result.ExtractInto(&s)
	return s.Routes, err
}

// updateResp is the structure that represents the API response of the 'Update' method, which contains route details and
// the request information.
type updateResp struct {
	// The response detail of the route.
	Route Route `json:"route"`
	// The request ID.
	RequestId string `json:"request_id"`
}
