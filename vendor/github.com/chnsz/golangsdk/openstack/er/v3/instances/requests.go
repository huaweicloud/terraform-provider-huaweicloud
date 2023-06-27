package instances

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// Number of records to be queried.
	// The valid value is range from 0 to 2000.
	Limit int `q:"limit"`
	// The ID of the instance of the last record on the previous page.
	// If it is empty, it is the first page of the query.
	// This parameter must be used together with limit.
	// The valid value is range from 1 to 128.
	Marker string `q:"marker"`
	// The enterprise project IDs of the instance to be queried.
	EnterpriseProjectIds []string `q:"enterprise_project_id"`
	// The status list of the instance to be queried.
	Statuses []string `q:"state"`
	// The instance IDs to be queried.
	IDs []string `q:"id"`
	// The resource IDs corresponding to the connection.
	ResourceIds []string `q:"resource_id"`
	// Whether resources belong to the current renant. If this parameter is set to true, only resources belonging to the
	// current tenant are queried, excluding shared resources. If the value is false, the current tenant and resources
	// shared with the tenant are queried.
	OwnedBySelf bool `q:"owned_by_self"`
	// The list of the destinations, support for querying multiple instances.
	SortKey []string `q:"sort_key"`
	// The returned results are arranged in ascending or descending order, the default is asc.
	SortDir []string `q:"sort_dir"`
}

// List is a method to query the list of the instances using given parameters.
func List(client *golangsdk.ServiceClient, opts ListOpts) ([]Instance, error) {
	url := rootURL(client)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := InstancePage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return extractInstances(pages)
}
