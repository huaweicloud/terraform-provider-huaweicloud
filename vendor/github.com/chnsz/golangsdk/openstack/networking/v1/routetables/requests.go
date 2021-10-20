package routetables

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// ListOptsBuilder is an interface by which can build the request body of listing route tables
type ListOptsBuilder interface {
	ToRouteTableListQuery() (string, error)
}

// ListOpts allows to query all route tables or filter collections by parameters
// Marker and Limit are used for pagination.
type ListOpts struct {
	// ID is the unique identifier for the route table
	ID string `q:"id"`
	// VpcID is the unique identifier for the vpc
	VpcID string `q:"vpc_id"`
	// SubnetID the unique identifier for the subnet
	SubnetID string `q:"subnet_id"`
	// Limit is the number of records returned for each page query, the value range is 0~intmax
	Limit int `q:"limit"`
	// Marker is the starting resource ID of the paging query,
	// which means that the query starts from the next record of the specified resource
	Marker string `q:"marker"`
}

// ToRouteTableListQuery formats a ListOpts into a query string
func (opts ListOpts) ToRouteTableListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List returns a Pager which allows you to iterate over a collection of
// vpc route tables. It accepts a ListOpts struct, which allows you to
// filter  the returned collection for greater efficiency.
func List(c *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)

	if opts != nil {
		query, err := opts.ToRouteTableListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := RouteTablePage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	})
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToRouteTableCreateMap() (map[string]interface{}, error)
}

// RouteOpts contains all the values needed to manage a vpc route
type RouteOpts struct {
	// The destination CIDR block. The destination of each route must be unique.
	// The destination cannot overlap with any subnet CIDR block in the VPC.
	Destination string `json:"destination" required:"true"`
	// the type of the next hop. value range:
	// ecs, eni, vip, nat, peering, vpn, dc, cc, egw
	Type string `json:"type" required:"true"`
	// the instance ID of next hop
	NextHop string `json:"nexthop" required:"true"`
	// The supplementary information about the route. The description can contain
	// a maximum of 255 characters and cannot contain angle brackets (< or >).
	Description *string `json:"description,omitempty"`
}

// CreateOpts contains all the values needed to create a new route table
type CreateOpts struct {
	// The VPC ID that the route table belongs to
	VpcID string `json:"vpc_id" required:"true"`
	// The name of the route table. The name can contain a maximum of 64 characters,
	// which may consist of letters, digits, underscores (_), hyphens (-), and periods (.).
	// The name cannot contain spaces.
	Name string `json:"name" required:"true"`
	// The supplementary information about the route table. The description can contain
	// a maximum of 255 characters and cannot contain angle brackets (< or >).
	Description string `json:"description,omitempty"`
	// The route information
	Routes []RouteOpts `json:"routes,omitempty"`
}

// ToRouteTableCreateMap builds a create request body from CreateOpts
func (opts CreateOpts) ToRouteTableCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "routetable")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// route table
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToRouteTableCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}

// Get retrieves a particular route table based on its unique ID
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request
type UpdateOptsBuilder interface {
	ToRouteTableUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains the values used when updating a route table
type UpdateOpts struct {
	Name        string                 `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	Routes      map[string][]RouteOpts `json:"routes,omitempty"`
}

// ToRouteTableUpdateMap builds an update body based on UpdateOpts
func (opts UpdateOpts) ToRouteTableUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "routetable")
}

// Update allows route tables to be updated
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToRouteTableUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will permanently delete a particular route table based on its unique ID
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}

// ActionOptsBuilder allows extensions to add additional parameters to the
// Action request: associate or disassociate subnets with a route table
type ActionOptsBuilder interface {
	ToRouteTableActionMap() (map[string]interface{}, error)
}

// ActionSubnetsOpts contains the subnets list that associate or disassociate with a route tabl
type ActionSubnetsOpts struct {
	Associate    []string `json:"associate,omitempty"`
	Disassociate []string `json:"disassociate,omitempty"`
}

// ActionOpts contains the values used when associating or disassociating subnets with a route table
type ActionOpts struct {
	Subnets ActionSubnetsOpts `json:"subnets" required:"true"`
}

// ToRouteTableActionMap builds an update body based on UpdateOpts.
func (opts ActionOpts) ToRouteTableActionMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "routetable")
}

// Action will associate or disassociate subnets with a particular route table based on its unique ID
func Action(c *golangsdk.ServiceClient, id string, opts ActionOptsBuilder) (r ActionResult) {
	b, err := opts.ToRouteTableActionMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Post(actionURL(c, id), b, &r.Body, nil)
	return
}
