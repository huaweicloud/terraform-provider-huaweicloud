package routes

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure used to create the route under a specified route table.
type CreateOpts struct {
	// The destination of the route.
	Destination string `json:"destination" required:"true"`
	// The ID of the corresponding attachment.
	AttachmentId string `json:"attachment_id,omitempty"`
	// Whether route is the black hole route.
	IsBlackHole *bool `json:"is_blackhole,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to create a new route under a specified route table.
func Create(client *golangsdk.ServiceClient, routeTableId string, opts CreateOpts) (*Route, error) {
	b, err := golangsdk.BuildRequestBody(opts, "route")
	if err != nil {
		return nil, err
	}

	var r createResp
	_, err = client.Post(rootURL(client, routeTableId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Route, err
}

// Get is a method to obtain the route details using given parameters.
func Get(client *golangsdk.ServiceClient, routeTableId, routeId string) (*Route, error) {
	var r getResp
	_, err := client.Get(resourceURL(client, routeId, routeTableId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Route, err
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// Number of records to be queried.
	// The valid value is range from 0 to 2000.
	Limit int `q:"limit"`
	// The ID of the route of the last record on the previous page.
	// If it is empty, it is the first page of the query.
	// This parameter must be used together with limit.
	// The valid value is range from 1 to 128.
	Marker string `q:"marker"`
	// The list of the destinations, support for querying multiple routes.
	Destination []string `q:"destination"`
	// The list of attachment IDs, support for querying multiple routes.
	AttachmentIds []string `json:"attachment_id"`
	// The list of attachment resource types, support for querying multiple routes.
	ResourceType []string `json:"resource_type"`
	// The list of keyword to sort the associations result, sort by ID by default.
	// The optional values are as follow:
	// + id
	// + name
	// + state
	SortKey []string `q:"sort_key"`
	// The returned results are arranged in ascending or descending order, the default is asc.
	SortDir []string `q:"sort_dir"`
}

// List is a method to query the list of the routes under a specified route table using given parameters.
func List(client *golangsdk.ServiceClient, routeTableId string, opts ListOpts) ([]Route, error) {
	url := rootURL(client, routeTableId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := RoutePage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return extractRoutes(pages)
}

// UpdateOpts is the structure used to update the route configuration.
type UpdateOpts struct {
	// The ID of the corresponding attachment.
	AttachmentId string `json:"attachment_id,omitempty"`
	// Whether route is the black hole route.
	IsBlackHole *bool `json:"is_blackhole,omitempty"`
}

// Update is a method to update route configuration using update option.
func Update(client *golangsdk.ServiceClient, routeTableId, routeId string, opts UpdateOpts) (*Route, error) {
	b, err := golangsdk.BuildRequestBody(opts, "route_table")
	if err != nil {
		return nil, err
	}

	var r updateResp
	_, err = client.Put(resourceURL(client, routeTableId, routeId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Route, err
}

// Delete is a method to remove an existing route from a specified route table.
func Delete(client *golangsdk.ServiceClient, routeTableId, routeId string) error {
	_, err := client.Delete(resourceURL(client, routeTableId, routeId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
