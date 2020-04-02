package hw_snatrules

import (
	"github.com/huaweicloud/golangsdk"
)

// CreateOptsBuilder is an interface must satisfy to be used as Create
// options.
type CreateOptsBuilder interface {
	ToSnatRuleCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new snat rule
// resource.
type CreateOpts struct {
	NatGatewayID string `json:"nat_gateway_id" required:"true"`
	NetworkID    string `json:"network_id,omitempty"`
	FloatingIPID string `json:"floating_ip_id" required:"true"`
	Cidr         string `json:"cidr,omitempty"`
	SourceType   int    `json:"source_type,omitempty"`
}

// ToSnatRuleCreateMap allows CreateOpts to satisfy the CreateOptsBuilder
// interface
func (opts CreateOpts) ToSnatRuleCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "snat_rule")
}

// Create is a method by which can create a new snat rule
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToSnatRuleCreateMap()
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
// snat rule.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// Delete is a method by which can be able to delete a snat rule
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}
