package vpcattachments

import (
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/pagination"
)

// SingleResp is the structure that represents the VPC attachment detail and the request information of the API request.
type SingleResp struct {
	// The response detail of the VPC attachment.
	Attachment Attachment `json:"vpc_attachment"`
	// The request ID.
	RequestId string `json:"request_id"`
}

// Attachment is the structure that represents the details of the VPC attachment for ER service.
type Attachment struct {
	// The project ID.
	ProjectId string `json:"project_id"`
	// The ID of the project where the VPC is located.
	VpcProjectId string `json:"vpc_project_id"`
	// The ID of the VPC attachment.
	ID string `json:"id"`
	// The name of the VPC attachment.
	Name string `json:"name"`
	// The description of the VPC attachment.
	Description string `json:"description"`
	// The VPC ID corresponding to the VPC attachment.
	VpcId string `json:"vpc_id"`
	// The VPC subnet ID corresponding to the VPC attachment.
	SubnetId string `json:"virsubnet_id"`
	// Whether automatically configure a route pointing to the ER instance for the VPC.
	AutoCreateVpcRoutes bool `json:"auto_create_vpc_routes"`
	// The current status of the VPC attachment.
	Status string `json:"state"`
	// The creation time of the VPC attachment.
	CreatedAt string `json:"created_at"`
	// The last update time of the VPC attachment.
	UpdatedAt string `json:"updated_at"`
	// The key/value pairs to associate with the VPC attachment.
	Tags []tags.ResourceTag `json:"tags"`
}

// MultipleResp is the structure that represents the VPC attachment list, page detail and the request information.
type MultipleResp struct {
	// The list of the VPC attachment.
	Attachments []Attachment `json:"vpc_attachments"`
	// The request ID.
	RequestId string `json:"request_id"`
	// The page information.
	PageInfo PageInfo `json:"page_info"`
}

// PageInfo is the structure that represents the page information.
type PageInfo struct {
	// The next marker information.
	NextMarker string `json:"next_marker"`
	// The number of the VPC attahcment in current page.
	CurrentCount int `json:"current_count"`
}

// AttachmentPage represents the response pages of the List method.
type AttachmentPage struct {
	pagination.MarkerPageBase
}

// IsEmpty returns true if current page no VPC attachment.
func (r AttachmentPage) IsEmpty() (bool, error) {
	resp, err := ExtractAttachments(r)
	return len(resp) == 0, err
}

// LastMarker returns the last marker index during current page.
func (r AttachmentPage) LastMarker() (string, error) {
	resp, err := ExtractPageInfo(r)
	if err != nil {
		return "", err
	}
	if resp.NextMarker != "" {
		return "", nil
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

// ExtractPageInfo is a method which to extract the response of the page information.
func ExtractPageInfo(r pagination.Page) (*PageInfo, error) {
	var s MultipleResp
	err := r.(AttachmentPage).Result.ExtractInto(&s)
	return &s.PageInfo, err
}

// ExtractAttachments is a method which to extract the response to an attachment list.
func ExtractAttachments(r pagination.Page) ([]Attachment, error) {
	var s MultipleResp
	err := r.(AttachmentPage).Result.ExtractInto(&s)
	return s.Attachments, err
}
