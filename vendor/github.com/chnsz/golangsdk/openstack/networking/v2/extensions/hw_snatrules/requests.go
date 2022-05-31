package hw_snatrules

import (
	"github.com/chnsz/golangsdk"
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
	FloatingIPID string `json:"floating_ip_id" required:"true"`
	Description  string `json:"description,omitempty"`
	NetworkID    string `json:"network_id,omitempty"`
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
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder is an interface must satisfy to be used as Update
// options.
type UpdateOptsBuilder interface {
	ToSnatRuleUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains all the values needed to update a snat rule
// resource.
type UpdateOpts struct {
	NatGatewayID      string  `json:"nat_gateway_id" required:"true"`
	FloatingIPAddress string  `json:"public_ip_address,omitempty"`
	Description       *string `json:"description,omitempty"`
}

// ToSnatRuleUpdateMap allows UpdateOpts to satisfy the UpdateOptsBuilder
// interface
func (opts UpdateOpts) ToSnatRuleUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "snat_rule")
}

// Update is a method by which can update a snat rule
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToSnatRuleUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, nil)
	return
}

// Get is a method by which can get the detailed information of the specified
// snat rule.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// Delete is a method by which can be able to delete a snat rule
func Delete(c *golangsdk.ServiceClient, id, natGatewayID string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURLDelete(c, id, natGatewayID), nil)
	return
}
