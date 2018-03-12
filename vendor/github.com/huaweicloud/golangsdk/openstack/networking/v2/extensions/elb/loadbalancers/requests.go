package loadbalancers

import (
	"log"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/elb"
	"github.com/huaweicloud/golangsdk/openstack/utils"
)

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
	Name            string `json:"name" required:"true"`
	Description     string `json:"description,omitempty"`
	VpcID           string `json:"vpc_id" required:"true"`
	BandWidth       int    `json:"bandwidth,omitempty"`
	Type            string `json:"type" required:"true"`
	AdminStateUp    int    `json:"admin_state_up" required:"true"`
	VipSubnetID     string `json:"vip_subnet_id,omitempty"`
	AZ              string `json:"az,omitempty"`
	ChargeMode      string `json:"charge_mode,omitempty"`
	EipType         string `json:"eip_type,omitempty"`
	SecurityGroupID string `json:"security_group_id,omitempty"`
	VipAddress      string `json:"vip_address,omitempty"`
	TenantID        string `json:"tenantId,omitempty"`
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
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r elb.JobResult) {
	b, err := opts.ToLoadBalancerCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	log.Printf("[DEBUG] create ELB-LoadBalancer url:%q, body=%#v", rootURL(c), b)
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
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
	// Optional. Human-readable name for the Loadbalancer. Does not have to be unique.
	Name string `json:"name,omitempty"`
	// Optional. Human-readable description for the Loadbalancer.
	Description string `json:"description"`

	BandWidth int `json:"bandwidth,omitempty"`
	// Optional. The administrative state of the Loadbalancer. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp int `json:"admin_state_up"`
}

// ToLoadBalancerUpdateMap casts a UpdateOpts struct to a map.
func (opts UpdateOpts) ToLoadBalancerUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update is an operation which modifies the attributes of the specified LoadBalancer.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOpts, not_pass_param []string) (r elb.JobResult) {
	b, err := opts.ToLoadBalancerUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	utils.DeleteNotPassParams(&b, not_pass_param)
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will permanently delete a particular LoadBalancer based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r elb.JobResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Delete2(resourceURL(c, id), &r.Body, reqOpt)
	return
}
