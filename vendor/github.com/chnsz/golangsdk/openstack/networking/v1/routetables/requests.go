package routetables
import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the attributes you want to see returned.
type ListOpts struct {
	// Specifies the route name.
	Name string `q:"name"`

	//Specifies the destination route.
	Destination string `q:"destination"`

	// Specifies the VPC.
	VPC_ID string `q:"vpc_id"`

	//Specifies the tenant ID. Only the administrator can specify the tenant ID of other tenants.
	Tenant_Id string `q:"tenant_id"`

	//Specifies the route ID.
	RouteTableID string `q:"id"`

	//Specifies the subnet ID.
	SubnetID string `q:"subnet_id"`
}

type ListOptsBuilder interface {
	ToRouteTableListQuery() (string, error)
}

// ToRouteListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToRouteTableListQuery() (string, error) {
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
		query, err := opts.ToRouteTableListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return RouteTablePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToRouteTablesCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new routes. There are
// no required values.
type CreateOpts struct {
	Name        string `json:"name,omitempty" required:"true"`
	Destination string `json:"destination,omitempty" required:"true"`
	Tenant_Id   string `json:"tenant_id,omitempty"`
	VPC_ID      string `json:"vpc_id,omitempty" required:"true"`
}

// ToRouteCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToRouteTablesCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "route")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// logical routes. When it is created, the routes does not have an internal
// interface - it is not associated to any routes.
//
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToRouteTablesCreateMap()
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

