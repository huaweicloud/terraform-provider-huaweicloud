package attachments

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// Number of records to be queried.
	// The valid value is range from 0 to 2000.
	Limit int `q:"limit"`
	// The next page index for the attachments query.
	// If it is empty, it is the first page of the query.
	// This parameter must be used together with limit.
	Marker string `q:"marker"`
	// The list of current status of the attachments (type of VPC, VPN, VGW and PEERING), support for querying multiple
	// attachments.
	Statuses []string `q:"state"`
	// The resource types to be filtered.
	// + vpc: virtual private cloud
	// + vpn: vpn gateway
	// + vgw: virtual gateway of cloud private line
	// + peering: Peering connection, through the cloud connection (CC) to load enterprise routers in different regions
	//   to create a peering connection
	ResourceTypes []string `q:"resource_type"`
	// The resource IDs corresponding to the attachment to be filtered.
	ResourceIds []string `q:"resource_id"`
	// The list of keyword to sort the attachments result, sort by ID by default.
	// The optional values are as follow:
	// + id
	// + name
	// + state
	SortKey []string `q:"sort_key"`
	// The returned results are arranged in ascending or descending order, the default is asc.
	// The valid values are as follows:
	// + asc
	// + desc
	SortDir []string `q:"sort_dir"`
}

// List is a method to query the list of the propagations under specified route table using given opts.
func List(client *golangsdk.ServiceClient, instanceId string, opts ListOpts) ([]Attachment, error) {
	url := queryURL(client, instanceId)
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
	return extractAttachments(pages)
}
