package snats

import "github.com/chnsz/golangsdk"

// CreateOpts is the structure used to create a new private SNAT rule.
type CreateOpts struct {
	// The ID of the gateway to which the private SNAT rule belongs.
	GatewayId string `json:"gateway_id" required:"true"`
	// The ID list of the transit IPs for private NAT.
	TransitIpIds []string `json:"transit_ip_ids" required:"true"`
	// The CIDR block of the match rule.
	// Exactly one of cidr and virsubnet_id must be set.
	Cidr string `json:"cidr,omitempty"`
	// The subnet ID of the match rule.
	// Exactly one of cidr and virsubnet_id must be set.
	SubnetId string `json:"virsubnet_id,omitempty"`
	// The description of the private SNAT rule, which contain maximum of `255` characters, and angle brackets (< and >)
	// are not allowed.
	Description string `json:"description,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create a private SNAT rule using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*Rule, error) {
	b, err := golangsdk.BuildRequestBody(opts, "snat_rule")
	if err != nil {
		return nil, err
	}

	var r createResp
	_, err = c.Post(rootURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Rule, err
}

// Get is a method used to obtain the private SNAT rule detail by its ID.
func Get(c *golangsdk.ServiceClient, ruleId string) (*Rule, error) {
	var r queryResp
	_, err := c.Get(resourceURL(c, ruleId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Rule, err
}

// UpdateOpts is the structure used to modify an existing private SNAT rule.
type UpdateOpts struct {
	// The ID list of the transit IPs for private NAT.
	TransitIpIds []string `json:"transit_ip_ids,omitempty"`
	// The description of the private SNAT rule, which contain maximum of `255` characters, and angle brackets (< and >)
	// are not allowed.
	Description *string `json:"description,omitempty"`
}

// Update is a method used to modify an existing SNAT rule using given parameters.
func Update(c *golangsdk.ServiceClient, ruleId string, opts UpdateOpts) (*Rule, error) {
	b, err := golangsdk.BuildRequestBody(opts, "snat_rule")
	if err != nil {
		return nil, err
	}

	var r updateResp
	_, err = c.Put(resourceURL(c, ruleId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Rule, err
}

// Delete is a method to remove the specified SNAT rule using its ID.
func Delete(c *golangsdk.ServiceClient, ruleId string) error {
	_, err := c.Delete(resourceURL(c, ruleId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
