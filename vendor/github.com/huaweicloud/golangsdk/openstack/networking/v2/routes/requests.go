package routes

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the attributes you want to see returned.
type ListOpts struct {
	// Specifies the route type.
	Type string `q:"type"`

	// Specifies the next hop. If the route type is peering, enter the VPC peering connection ID.
	//NextHop string `q:"nexthop"`

	//Specifies the destination IP address or CIDR block.
	Destination string `q:"destination"`

	// Specifies the VPC for which a route is to be added.
	VPC_ID string `q:"vpc_id"`

	//Specifies the tenant ID. Only the administrator can specify the tenant ID of other tenants.
	Tenant_Id string `q:"tenant_id"`

	//Specifies the route ID.
	RouteID string `q:"id"`
}
type ListOptsBuilder interface {
	ToRouteListQuery() (string, error)
}

// ToRouteListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToRouteListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List returns a Pager which allows you to iterate over a collection of
// vpc routes  resources. It accepts a ListOpts struct, which allows you to
// filter  the returned collection for greater efficiency.
func List(c *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)

	if opts != nil {
		query, err := opts.ToRouteListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return RoutePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToRouteCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new routes. There are
// no required values.
type CreateOpts struct {
	Type        string `json:"type,omitempty" required:"true"`
	NextHop     string `json:"nexthop,omitempty" required:"true"`
	Destination string `json:"destination,omitempty" required:"true"`
	Tenant_Id   string `json:"tenant_id,omitempty"`
	VPC_ID      string `json:"vpc_id,omitempty" required:"true"`
}

// ToRouteCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToRouteCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "route")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// logical routes. When it is created, the routes does not have an internal
// interface - it is not associated to any routes.
//
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToRouteCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{201}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular route based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// Delete will permanently delete a particular route based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}
