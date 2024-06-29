package attachments

import (
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/pagination"
)

// listResp is the structure that represents the ER attachment list, page detail and the request information.
type listResp struct {
	// The list of the ER attachment.
	Attachments []Attachment `json:"attachments"`
	// The request ID.
	RequestId string `json:"request_id"`
	// The page information.
	PageInfo PageInfo `json:"page_info"`
}

// Attachment is the structure that represents the ER attachment (type of VPC, VPN, VGW and PEERING) details.
type Attachment struct {
	// The attachment ID.
	ID string `json:"id"`
	// The attachment name.
	Name string `json:"name"`
	// The attachment description.
	Description string `json:"description"`
	// The current status of the attachment.
	// + pending
	// + available
	// + modifying
	// + deleting
	// + failed
	// + pending_acceptance
	// + rejected
	// + initiating_request
	// + freezed
	Status string `json:"state"`
	// The creation time of the attachment.
	CreatedAt string `json:"created_at"`
	// The latest update time of the attachment.
	UpdatedAt string `json:"updated_at"`
	// The tag list associated with the attachment.
	Tags []tags.ResourceTag `json:"tags"`
	// The project ID where the attachment is located.
	ProjectId string `json:"project_id"`
	// The ER instance ID to which the attachment is belongs.
	InstanceId string `json:"er_id"`
	// The resource ID associated with the internal connection.
	ResourceId string `json:"resource_id"`
	// The resource type.
	// + vpc: virtual private cloud.
	// + vpn: vpn gateway.
	// + vgw: virtual gateway of cloud private line.
	// + peering: Peering connection, through the cloud connection (CC) to load enterprise routers in different regions
	//   to create a peering connection.
	ResourceType string `json:"resource_type"`
	// The project ID to which the attachment is belongs.
	ResourceProjectId string `json:"resource_project_id"`
	// Whether this attachment is associated.
	Associated bool `json:"associated"`
	// The associated route table ID.
	RouteTableId string `json:"route_table_id"`
}

// PageInfo is the structure that represents the page information.
type PageInfo struct {
	// The next marker (next page index) information.
	NextMarker string `json:"next_marker"`
	// The number of the ER attahcment in current page.
	CurrentCount int `json:"current_count"`
}

// AttachmentPage represents the response pages of the List method.
type AttachmentPage struct {
	pagination.MarkerPageBase
}

// IsEmpty returns true if current page no ER attachment.
func (r AttachmentPage) IsEmpty() (bool, error) {
	resp, err := extractAttachments(r)
	return len(resp) == 0, err
}

// LastMarker returns the last marker index during current page.
func (r AttachmentPage) LastMarker() (string, error) {
	resp, err := extractPageInfo(r)
	if err != nil {
		return "", err
	}
	return resp.NextMarker, nil
}

// NextPageURL generates the URL for the page of results after this one.
func (r AttachmentPage) NextPageURL() (string, error) {
	currentURL := r.URL

	mark, err := r.Owner.LastMarker()
	if err != nil {
		return "", err
	}
	if mark == "" {
		return "", nil
	}

	q := currentURL.Query()
	q.Set("marker", mark)
	currentURL.RawQuery = q.Encode()

	return currentURL.String(), nil
}

// extractPageInfo is a method which to extract the response of the page information.
func extractPageInfo(r pagination.Page) (*PageInfo, error) {
	var s listResp
	err := r.(AttachmentPage).Result.ExtractInto(&s)
	return &s.PageInfo, err
}

// extractAttachments is a method which to extract the response to an attachment list.
func extractAttachments(r pagination.Page) ([]Attachment, error) {
	var s listResp
	err := r.(AttachmentPage).Result.ExtractInto(&s)
	return s.Attachments, err
}
