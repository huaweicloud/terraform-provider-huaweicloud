package loadbalancer_elbs

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
	//"fmt"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToLoadBalancerListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the Loadbalancer attributes you want to see returned. SortKey allows you to
// sort by a particular attribute. SortDir sets the direction, and is
// either `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	Name            string `q:"name"`
	Description     string `q:"description"`
	VpcID           string `q:"vpc_id"`
	Bandwidth       int    `q:"bandwidth"`
	Type            string `q:"type"`
	AdminStateUp    *bool  `q:"admin_state_up"`
	VipSubnetID     string `q:"vip_subnet_id"`
	AZ              string `q:"az"`
	ChargeMode      string `q:"charge_mode"`
	EipType         string `q:"eip_type"`
	SecurityGroupID string `q:"security_group_id"`
	VipAddress      string `q:"vip_address"`
	TenantID        string `q:"tenantId"`
	ID              string `q:"id"`
	Limit           int    `q:"limit"`
	Marker          string `q:"marker"`
	SortKey         string `q:"sort_key"`
	SortDir         string `q:"sort_dir"`
}

// ToLoadbalancerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToLoadBalancerListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// routers. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those routers that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(c *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	//fmt.Printf("url=%s.\n", url)
	if opts != nil {
		query, err := opts.ToLoadBalancerListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return LoadBalancerPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToLoadBalancerCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	// Required. Specifies the load balancer name.
	// The name is a string of 1 to 64 characters that consist of letters, digits, underscores (_), and hyphens (-).
	Name string `json:"name" required:"true"`
	// Optional. Provides supplementary information about the load balancer.
	// The value is a string of 0 to 128 characters and cannot contain angle brackets (<>).
	Description string `json:"description,omitempty"`
	// Required. Specifies the VPC ID.
	VpcID string `json:"vpc_id" required:"true"`
	// Optional. Specifies the bandwidth (Mbit/s). This parameter is mandatory when type is
	// set to External, and it is invalid when type is set to Internal.
	// The value ranges from 1 to 300.
	Bandwidth int `json:"bandwidth,omitempty"`
	// Required. Specifies the load balancer type.
	// The value can be Internal or External.
	Type string `json:"type" required:"true"`
	// Required.  Specifies the status of the load balancer.
	// Optional values:
	// 0 or false: indicates that the load balancer is stopped. Only tenants are allowed to enter these two values.
	// 1 or true: indicates that the load balancer is running properly.
	// 2 or false: indicates that the load balancer is frozen. Only tenants are allowed to enter these two values.
	AdminStateUp *bool `json:"admin_state_up" required:"true"`
	// Optional.  Specifies the subnet ID of backend ECSs. This parameter is mandatory when type is set to Internal.
	VipSubnetID string `json:"vip_subnet_id,omitempty"`
	// Optional.  Specifies the ID of the availability zone (AZ). This parameter is mandatory when type
	// is set to Internal, and it is invalid when type is set to External.
	AZ string `json:"az,omitempty"`
	// Optional.  Specifies the security group ID.
	// The value is a string of 1 to 200 characters that consists of uppercase and lowercase letters, digits, and hyphens (-).
	// This parameter is mandatory when type is set to Internal.
	SecurityGroupID string `json:"security_group_id,omitempty"`
	// Optional.  Specifies the IP address used by ELB for providing services. When type is set to External,
	// the parameter value is the elastic IP address. When type is set to Internal, the parameter value is
	// the private network IP address.
	// You can select an existing elastic IP address and create a public network load balancer.
	// When this parameter is configured, parameters bandwidth, charge_mode, and eip_type are invalid.
	VipAddress string `json:"vip_address,omitempty"`
	// Specifies the tenant ID. This parameter is mandatory when type is set to Internal
	TenantID string `json:"tenantid,omitempty"`
}

// ToLoadBalancerCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToLoadBalancerCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create is an operation which provisions a new loadbalancer based on the
// configuration defined in the CreateOpts struct. Once the request is
// validated and progress has started on the provisioning process, a
// CreateResult will be returned.
//
// Users with an admin role can create loadbalancers on behalf of other tenants by
// specifying a TenantID attribute different than their own.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToLoadBalancerCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	//fmt.Printf("Create (%+v): rootURL: %s, b=%+v.\n", c, rootURL(c), b)
	_, r.Err = c.Post(rootURL(c), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get retrieves a particular Loadbalancer based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// UpdateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Update operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type UpdateOptsBuilder interface {
	ToLoadBalancerUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is the common options struct used in this package's Update
// operation.
type UpdateOpts struct {
	// Required. Specifies the load balancer name.
	// The name is a string of 1 to 64 characters that consist of letters, digits, underscores (_), and hyphens (-).
	Name string `json:"name,omitempty"`
	// Optional. Provides supplementary information about the load balancer.
	// The value is a string of 0 to 128 characters and cannot contain angle brackets (<>).
	Description string `json:"description,omitempty"`
	// Optional. Specifies the bandwidth (Mbit/s). This parameter is mandatory when type is
	// set to External, and it is invalid when type is set to Internal.
	// The value ranges from 1 to 300.
	Bandwidth int `json:"bandwidth,omitempty"`
	// Required.  Specifies the status of the load balancer.
	// Optional values:
	// 0 or false: indicates that the load balancer is stopped. Only tenants are allowed to enter these two values.
	// 1 or true: indicates that the load balancer is running properly.
	// 2 or false: indicates that the load balancer is frozen. Only tenants are allowed to enter these two values.
	AdminStateUp *bool `json:"admin_state_up,omitempty"`
}

// ToLoadBalancerUpdateMap casts a UpdateOpts struct to a map.
func (opts UpdateOpts) ToLoadBalancerUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update is an operation which modifies the attributes of the specified LoadBalancer.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOpts) (r UpdateResult) {
	b, err := opts.ToLoadBalancerUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}

// Delete will permanently delete a particular LoadBalancer based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete2(resourceURL(c, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
