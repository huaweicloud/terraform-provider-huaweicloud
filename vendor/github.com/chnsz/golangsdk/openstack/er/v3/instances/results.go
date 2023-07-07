package instances

import (
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/pagination"
)

// Instance is the structure that represents the details of the ER instance.
type Instance struct {
	// The ID of the instance.
	ID string `json:"id"`
	// The name of the instance.
	Name string `json:"name"`
	// The destination of the instance.
	Description string `json:"description"`
	// The current status of the instance.
	Status string `json:"state"`
	// The key/value pairs to associate with the instance.
	Tags []tags.ResourceTag `json:"tags"`
	// The creation time of the instance.
	CreatedAt string `json:"created_at"`
	// The last update time of the instance.
	UpdatedAt string `json:"updated_at"`
	// The ID of enterprise project to which the instance belongs.
	EnterpriseProjectId string `json:"enterprise_project_id"`
	// The project ID.
	ProjectId string `json:"project_id"`
	// The BGP AS number of the ER instance.
	ASN int `json:"asn"`
	// Whether to enable the propagation of the default route table.
	EnableDefaultPropagation bool `json:"enable_default_propagation"`
	// Whether to enable the association of the default route table.
	EnableDefaultAssociation bool `json:"enable_default_association"`
	// Whether to automatically accept the creation of shared attachment.
	AutoAcceptSharedAttachments bool `json:"auto_accept_shared_attachments"`
	// The ID of the default propagation route table.
	DefaultPropagationRouteTableId string `json:"default_propagation_route_table_id"`
	// The ID of the default association route table.
	DefaultAssociationRouteTableId string `json:"default_association_route_table_id"`
	// The availability zone list where the ER instance is located.
	AvailabilityZoneIds []string `json:"availability_zone_ids"`
}

// listResp is the structure that represents the API response of the 'List' method, which contains instance list, page
// details and the request information.
type listResp struct {
	// The list of the instances.
	Instances []Instance `json:"instances"`
	// The request ID.
	RequestId string `json:"request_id"`
	// The page information.
	PageInfo pageInfo `json:"page_info"`
}

// pageInfo is the structure that represents the page information.
type pageInfo struct {
	// The next marker information.
	NextMarker string `json:"next_marker"`
	// The number of the instances in current page.
	CurrentCount int `json:"current_count"`
}

// InstancePage represents the response pages of the List method.
type InstancePage struct {
	pagination.MarkerPageBase
}

// IsEmpty returns true if a ListResult no instance.
func (r InstancePage) IsEmpty() (bool, error) {
	resp, err := extractInstances(r)
	return len(resp) == 0, err
}

// LastMarker returns the last marker index in a ListResult.
func (r InstancePage) LastMarker() (string, error) {
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
	err := r.(InstancePage).Result.ExtractInto(&s)
	return &s.PageInfo, err
}

// extractInstances is a method which to extract the response to an instance list.
func extractInstances(r pagination.Page) ([]Instance, error) {
	var s listResp
	err := r.(InstancePage).Result.ExtractInto(&s)
	return s.Instances, err
}
