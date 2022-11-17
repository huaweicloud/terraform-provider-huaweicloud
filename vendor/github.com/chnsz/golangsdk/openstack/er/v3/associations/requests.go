package associations

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure used to create an association to a specified route table.
type CreateOpts struct {
	// The ID of the VPC attachment.
	AttachmentId string `json:"attachment_id,omitempty"`
	// The export routing policy.
	RoutePolicy ExportRoutePolicy `json:"route_policy,omitempty"`
}

// ExportRoutePolicy is an object that represents the configuration of the export routing policy.
type ExportRoutePolicy struct {
	// The export routing policy ID.
	ExportPoilicyId string `json:"export_policy_id,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to create a new association under a specified route table.
func Create(client *golangsdk.ServiceClient, instanceId, routeTableId string, opts CreateOpts) (*Association, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r createResp
	_, err = client.Post(enableURL(client, instanceId, routeTableId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Association, err
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// Number of records to be queried.
	// The valid value is range from 0 to 2000.
	Limit int `q:"limit"`
	// The ID of the association of the last record on the previous page.
	// If it is empty, it is the first page of the query.
	// This parameter must be used together with limit.
	// The valid value is range from 1 to 128.
	Marker string `q:"marker"`
	// The list of attachment IDs, support for querying multiple associations.
	AttachmentIds []string `q:"attachment_id"`
	// The list of attachment resource types, support for querying multiple associations.
	ResourceTypes []string `q:"resource_type"`
	// The list of current status of the associations, support for querying multiple associations.
	Statuses []string `q:"state"`
	// The list of keyword to sort the associations result, sort by ID by default.
	// The optional values are as follow:
	// + id
	// + name
	// + state
	SortKey []string `q:"sort_key"`
	// The returned results are arranged in ascending or descending order, the default is asc.
	SortDir []string `q:"sort_dir"`
}

// List is a method to query the list of the associations under specified route table using given parameters.
func List(client *golangsdk.ServiceClient, instanceId, routeTableId string, opts ListOpts) ([]Association, error) {
	url := queryURL(client, instanceId, routeTableId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := AssociationPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return extractAssociations(pages)
}

// DeleteOpts is the structure used to remove an association from a specified route table.
type DeleteOpts struct {
	// The ID of the VPC attachment.
	AttachmentId string `json:"attachment_id,omitempty"`
	// The export routing policy.
	RoutePolicy ExportRoutePolicy `json:"route_policy,omitempty"`
}

// Delete is a method to remove an existing association from a specified route table.
func Delete(client *golangsdk.ServiceClient, instanceId, routeTableId string, opts DeleteOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = client.Post(disableURL(client, instanceId, routeTableId), b, nil,
		&golangsdk.RequestOpts{
			MoreHeaders: requestOpts.MoreHeaders,
		})
	return err
}
