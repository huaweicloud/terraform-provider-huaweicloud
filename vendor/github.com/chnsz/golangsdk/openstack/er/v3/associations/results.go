package associations

import (
	"github.com/chnsz/golangsdk/pagination"
)

// createResp is the structure that represents the API response of the 'Create' method, which contains association
// details and the request information.
type createResp struct {
	// The response detail of the association.
	Association Association `json:"association"`
	// The request ID.
	RequestId string `json:"request_id"`
}

// Association is the structure that represents the details of the association under route table.
type Association struct {
	// The ID of the association.
	ID string `json:"id"`
	// The ID of the route table to which the association belongs.
	RouteTableId string `json:"route_table_id"`
	// The ID of the corresponding attachment.
	AttachmentId string `json:"attachment_id"`
	// The resource type for the corresponding attachment.
	ResourceType string `json:"resource_type"`
	// The resource ID for the corresponding attachment.
	ResourceId string `json:"resource_id"`
	// The configuration of the export routing policy.
	RoutePolicy ExportRoutePolicy `json:"route_policy"`
	// The current status of the association.
	Status string `json:"state"`
	// The creation time of the association.
	CreatedAt string `json:"created_at"`
	// The last update time of the association.
	UpdatedAt string `json:"updated_at"`
}

// listResp is the structure that represents the API response of the 'List' method, which contains association list,
// page details and the request information.
type listResp struct {
	// The list of the associations.
	Associations []Association `json:"associations"`
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

// AssociationPage represents the response pages of the List method.
type AssociationPage struct {
	pagination.MarkerPageBase
}

// IsEmpty returns true if a ListResult no association.
func (r AssociationPage) IsEmpty() (bool, error) {
	resp, err := extractAssociations(r)
	return len(resp) == 0, err
}

// LastMarker returns the last marker index in a ListResult.
func (r AssociationPage) LastMarker() (string, error) {
	resp, err := extractPageInfo(r)
	if err != nil {
		return "", err
	}

	return resp.NextMarker, nil
}

// extractPageInfo is a method which to extract the response of the page information.
func extractPageInfo(r pagination.Page) (*pageInfo, error) {
	var s listResp
	err := r.(AssociationPage).Result.ExtractInto(&s)
	return &s.PageInfo, err
}

// ExtractAssociations is a method which to extract the response to a association list.
func extractAssociations(r pagination.Page) ([]Association, error) {
	var s listResp
	err := r.(AssociationPage).Result.ExtractInto(&s)
	return s.Associations, err
}
