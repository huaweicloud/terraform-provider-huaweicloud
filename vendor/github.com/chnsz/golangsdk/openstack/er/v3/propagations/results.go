package propagations

import (
	"github.com/chnsz/golangsdk/pagination"
)

// createResp is the structure that represents the API response of the 'Create' method, which contains propagation
// details and the request information.
type createResp struct {
	// The response detail of the propagation.
	Propagation Propagation `json:"propagation"`
	// The request ID.
	RequestId string `json:"request_id"`
}

// Propagation is the structure that represents the details of the propagation under route table.
type Propagation struct {
	// The ID of the propagation.
	ID string `json:"id"`
	// The ID of the project where the propagation is located.
	ProjectId string `json:"project_id"`
	// The ID of the ER instance to which the propagation belongs.
	InstanceId string `json:"er_id"`
	// The ID of the route table to which the association belongs.
	RouteTableId string `json:"route_table_id"`
	// The ID of the corresponding attachment.
	AttachmentId string `json:"attachment_id"`
	// The resource type for the corresponding attachment.
	ResourceType string `json:"resource_type"`
	// The resource ID for the corresponding attachment.
	ResourceId string `json:"resource_id"`
	// The configuration of the import routing policy.
	RoutePolicy ImportRoutePolicy `json:"route_policy"`
	// The current status of the propagation.
	Status string `json:"state"`
	// The creation time of the propagation.
	CreatedAt string `json:"created_at"`
	// The last update time of the propagation.
	UpdatedAt string `json:"updated_at"`
}

// listResp is the structure that represents the API response of the 'List' method, which contains propagation list,
// page details and the request information.
type listResp struct {
	// The list of the propagations.
	Propagations []Propagation `json:"propagations"`
	// The request ID.
	RequestId string `json:"request_id"`
	// The page information.
	PageInfo pageInfo `json:"page_info"`
}

// pageInfo is the structure that represents the page information.
type pageInfo struct {
	// The next marker information.
	NextMarker string `json:"next_marker"`
	// The number of the propagations in current page.
	CurrentCount int `json:"current_count"`
}

// PropagationPage represents the response pages of the List method.
type PropagationPage struct {
	pagination.MarkerPageBase
}

// IsEmpty returns true if a ListResult no propagation.
func (r PropagationPage) IsEmpty() (bool, error) {
	resp, err := extractPropagations(r)
	return len(resp) == 0, err
}

// LastMarker returns the last marker index in a ListResult.
func (r PropagationPage) LastMarker() (string, error) {
	resp, err := extractPageInfo(r)
	if err != nil {
		return "", err
	}

	return resp.NextMarker, nil
}

// extractPageInfo is a method which to extract the response of the page information.
func extractPageInfo(r pagination.Page) (*pageInfo, error) {
	var s listResp
	err := r.(PropagationPage).Result.ExtractInto(&s)
	return &s.PageInfo, err
}

// ExtractPropagations is a method which to extract the response to a propagation list.
func extractPropagations(r pagination.Page) ([]Propagation, error) {
	var s listResp
	err := r.(PropagationPage).Result.ExtractInto(&s)
	return s.Propagations, err
}
