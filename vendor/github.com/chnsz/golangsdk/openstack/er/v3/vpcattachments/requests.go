package vpcattachments

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure required by the Create method to create a VPC attachment under the ER instance.
type CreateOpts struct {
	// The VPC ID corresponding to the VPC attachment.
	VpcId string `json:"vpc_id" required:"true"`
	// The VPC subnet ID corresponding to the VPC attachment.
	SubnetId string `json:"virsubnet_id" required:"true"`
	// The name of the VPC attachment.
	// The value can contain 1 to 64 characters, only english and chinese letters, digits, underscore (_), hyphens (-)
	// and dots (.) are allowed.
	Name string `json:"name" required:"true"`
	// The description of the VPC attachment.
	// The value contain a maximum of 255 characters, and the angle brackets (< and >) are not allowed.
	Description string `json:"description,omitempty"`
	// Whether automatically configure a route pointing to the ER instance for the VPC, defaults to false.
	AutoCreateVpcRoutes *bool `json:"auto_create_vpc_routes,omitempty"`
	// The key/value pairs to associate with the VPC attachment.
	Tags []tags.ResourceTag `json:"tags,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to create a VPC attachment under the ER instance using create option.
func Create(client *golangsdk.ServiceClient, instanceId string, opts CreateOpts) (*Attachment, error) {
	b, err := golangsdk.BuildRequestBody(opts, "vpc_attachment")
	if err != nil {
		return nil, err
	}

	var r SingleResp
	_, err = client.Post(rootURL(client, instanceId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Attachment, err
}

// Get is a method to obtain the details of a specified VPC attachment under ER instance using its ID.
func Get(client *golangsdk.ServiceClient, instanceId, attachmentId string) (*Attachment, error) {
	var r SingleResp
	_, err := client.Get(resourceURL(client, instanceId, attachmentId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Attachment, err
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// Number of records to be queried.
	// The valid value is range from 0 to 2000.
	Limit int `q:"limit"`
	// The ID of the VPC attachment of the last record on the previous page.
	// If it is empty, it is the first page of the query.
	// This parameter must be used together with limit.
	// The valid value is range from 1 to 128.
	Marker string `q:"marker"`
	// The list of VPC IDs corresponding to the VPC attachments, support for querying multiple VPC attachments.
	VpcId []string `q:"vpc_id"`
	// The list of VPC attachment IDs, support for querying multiple VPC attachments.
	AttachmentId []string `q:"id"`
	// The list of current status of the VPC attachments, support for querying multiple VPC attachments.
	Status []string `q:"state"`
	// The list of keyword to sort the VPC attachments result, sort by ID by default.
	// The optional values are as follow:
	// + id
	// + name
	// + state
	SortKey []string `q:"sort_key"`
	// The returned results are arranged in ascending or descending order, the default is asc.
	SortDir []string `q:"sort_dir"`
}

// List is a method to query the list of the VPC attachments under specified ER instance using given opts.
func List(client *golangsdk.ServiceClient, instanceId string, opts ListOpts) ([]Attachment, error) {
	url := rootURL(client, instanceId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := AttachmentPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	}).AllPages()
	if err != nil {
		return nil, err
	}

	return ExtractAttachments(pages)
}

// UpdateOpts is the structure required by the Update method to update the VPC attachment configuration.
type UpdateOpts struct {
	// The name of the VPC attachment.
	// The value can contain 1 to 64 characters, only english and chinese letters, digits, underscore (_), hyphens (-)
	// and dots (.) are allowed.
	Name string `json:"name,omitempty"`
	// The description of the VPC attachment.
	// The value contain a maximum of 255 characters, and the angle brackets (< and >) are not allowed.
	Description *string `json:"description,omitempty"`
}

// Update is a method to update the VPC attachment under the ER instance using update option.
func Update(client *golangsdk.ServiceClient, instanceId, attachmentId string, opts UpdateOpts) (*Attachment, error) {
	b, err := golangsdk.BuildRequestBody(opts, "vpc_attachment")
	if err != nil {
		return nil, err
	}

	var r SingleResp
	_, err = client.Put(resourceURL(client, instanceId, attachmentId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Attachment, err
}

// Delete is a method to remove an existing VPC attachment under specified ER instance.
func Delete(client *golangsdk.ServiceClient, instanceId, attachmentId string) error {
	_, err := client.Delete(resourceURL(client, instanceId, attachmentId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
