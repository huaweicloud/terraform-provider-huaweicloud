package routetables

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure required by the 'Create' method to create a route table under a specified ER instance.
type CreateOpts struct {
	// The name of the route table.
	// The value can contain 1 to 64 characters, only english and chinese letters, digits, underscore (_), hyphens (-)
	// and dots (.) are allowed.
	Name string `json:"name" required:"true"`
	// The description of the route table.
	// The value contain a maximum of 255 characters, and the angle brackets (< and >) are not allowed.
	Description string `json:"description,omitempty"`
	// The configuration of the BGP route selection.
	BgpOptions BgpOptions `json:"bgp_options,omitempty"`
	// The key/value pairs to associate with the route table.
	Tags []tags.ResourceTag `json:"tags,omitempty"`
}

// BgpOptions is an object that represents the BGP configuration for routing.
type BgpOptions struct {
	// Whether the AS path attributes of the routes are not compared during load balancing.
	LoadBalancingAsPathIgnore *bool `json:"load_balancing_as_path_ignore,omitempty"`
	// Whether the AS path attributes of the same length are not compared during load balancing.
	LoadBalancingAsPathRelax *bool `json:"load_balancing_as_path_relax,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to create a new route table under a specified ER instance using given parameters.
func Create(client *golangsdk.ServiceClient, instanceId string, opts CreateOpts) (*RouteTable, error) {
	b, err := golangsdk.BuildRequestBody(opts, "route_table")
	if err != nil {
		return nil, err
	}

	var r createResp
	_, err = client.Post(rootURL(client, instanceId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.RouteTable, err
}

// Get is a method to obtain the route table details under a specified ER instance.
func Get(client *golangsdk.ServiceClient, instanceId, routeTableId string) (*RouteTable, error) {
	var r getResp
	_, err := client.Get(resourceURL(client, instanceId, routeTableId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.RouteTable, err
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// Number of records to be queried.
	// The valid value is range from 0 to 2000.
	Limit int `q:"limit"`
	// The ID of the route table of the last record on the previous page.
	// If it is empty, it is the first page of the query.
	// This parameter must be used together with limit.
	// The valid value is range from 1 to 128.
	Marker string `q:"marker"`
	// The list of current status of the route tables, support for querying multiple route tables.
	Status []string `q:"state"`
	// Whether this route table is the default association route table.
	IsDefaultAssociation bool `q:"is_default_association"`
	// Whether this route table is the default propagation route table.
	IsDefaultPropagation bool `q:"is_default_propagation"`
	// The list of keyword to sort the route tables result, sort by ID by default.
	// The optional values are as follow:
	// + id
	// + name
	// + state
	SortKey []string `q:"sort_key"`
	// The returned results are arranged in ascending or descending order, the default is asc.
	SortDir []string `q:"sort_dir"`
}

// List is a method to query the list of the route tables under a specified ER instance using given parameters.
func List(client *golangsdk.ServiceClient, instanceId string, opts ListOpts) ([]RouteTable, error) {
	url := rootURL(client, instanceId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := RouteTablePage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return extractRouteTables(pages)
}

// UpdateOpts is the structure required by the 'Update' method to update the route table configuration.
type UpdateOpts struct {
	// The name of the route table.
	// The value can contain 1 to 64 characters, only english and chinese letters, digits, underscore (_), hyphens (-)
	// and dots (.) are allowed.
	Name string `json:"name,omitempty"`
	// The description of the route table.
	// The value contain a maximum of 255 characters, and the angle brackets (< and >) are not allowed.
	Description *string `json:"description,omitempty"`
}

// Update is a method to update the route table under a specified ER instance using parameters.
func Update(client *golangsdk.ServiceClient, instanceId, routeTableId string, opts UpdateOpts) (*RouteTable, error) {
	b, err := golangsdk.BuildRequestBody(opts, "route_table")
	if err != nil {
		return nil, err
	}

	var r updateResp
	_, err = client.Put(resourceURL(client, instanceId, routeTableId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.RouteTable, err
}

// Delete is a method to remove an existing route table under a specified ER instance.
func Delete(client *golangsdk.ServiceClient, instanceId, routeTableId string) error {
	_, err := client.Delete(resourceURL(client, instanceId, routeTableId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
