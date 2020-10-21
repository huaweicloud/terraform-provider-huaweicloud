package natgateways

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// CreateOptsBuilder is an interface must satisfy to be used as Create
// options.
type CreateOptsBuilder interface {
	ToNatGatewayCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new nat gateway
// resource.
type CreateOpts struct {
	Name                string `json:"name" required:"true"`
	Description         string `json:"description,omitempty"`
	Spec                string `json:"spec" required:"true"`
	RouterID            string `json:"router_id" required:"true"`
	InternalNetworkID   string `json:"internal_network_id" required:"true"`
	TenantID            string `json:"tenant_id,omitempty"`
	EnterpriseProjectID string `json:"enterprise_project_id,omitempty"`
}

type ListOpts struct {
	Limit               int    `q:"limit"`
	ID                  string `q:"id"`
	Name                string `q:"name"`
	TenantId            string `q:"tenant_id"`
	Description         string `q:"description"`
	Spec                string `q:"spec"`
	RouterID            string `q:"router_id"`
	InternalNetworkID   string `q:"internal_network_id"`
	Status              string `q:"status"`
	AdminStateUp        *bool  `q:"admin_state_up"`
	CreatedAt           string `q:"created_at"`
	EnterpriseProjectID string `q:"enterprise_project_id"`
}

// ToNatGatewayCreateMap allows CreateOpts to satisfy the CreateOptsBuilder
// interface
func (opts CreateOpts) ToNatGatewayCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "nat_gateway")
}

// Create is a method by which can create a new nat gateway
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToNatGatewayCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{201},
	})
	return
}

// Get is a method by which can get the detailed information of the specified
// nat gateway.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// Delete is a method by which can be able to delete a nat gateway
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}

// UpdateOptsBuilder is the interface type must satisfy to be used as Update
// options.
type UpdateOptsBuilder interface {
	ToNatGatewayUpdateMap() (map[string]interface{}, error)
}

//UpdateOpts is a struct which represents the request body of update method
type UpdateOpts struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Spec        string `json:"spec,omitempty"`
}

func (opts UpdateOpts) ToNatGatewayUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "nat_gateway")
}

func (opts ListOpts) ToNatGatewayListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

type ListOptsBuilder interface {
	ToNatGatewayListQuery() (string, error)
}

//Update allows nat gateway resources to be updated.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToNatGatewayUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

func List(c *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToNatGatewayListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return NatGatewayPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
