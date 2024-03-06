package snats

import (
	"github.com/chnsz/golangsdk"
)

// CreateOpts is the structure that used to create a new SNAT rule.
type CreateOpts struct {
	// The ID of the gateway to which the SNAT rule belongs.
	GatewayId string `json:"nat_gateway_id" required:"true"`
	//  IDs of floating IPs connected by SNAT rules (separated by commas).
	FloatingIpId string `json:"floating_ip_id,omitempty"`
	// The IDs (separated by commas) of global EIPs connected by SNAT rule.
	GlobalEipId string `json:"global_eip_id,omitempty"`
	// The description of the SNAT rule.
	Description string `json:"description,omitempty"`
	// The network IDs of subnet connected by SNAT rule (VPC side).
	NetworkId string `json:"network_id,omitempty"`
	// The CIDR block connected by SNAT rule (DC side).
	Cidr string `json:"cidr,omitempty"`
	// The resource type of the SNAT rule.
	// The valid values are as follows:
	// + 0: VPC side.
	// + 1: DC side.
	SourceType int `json:"source_type,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create an SNAT rule using given parameters.
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

// Get is a method used to obtain the detail of the SNAT rule by its ID.
func Get(c *golangsdk.ServiceClient, ruleId string) (*Rule, error) {
	var r queryResp
	_, err := c.Get(resourceURL(c, ruleId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Rule, err
}

// CreateOpts is the structure that used to update the configuration of the SNAT rule.
type UpdateOpts struct {
	// The ID of the gateway to which the SNAT rule belongs.
	GatewayId string `json:"nat_gateway_id" required:"true"`
	// The floating IP addresses (separated by commas) connected by SNAT rule.
	FloatingIpAddress string `json:"public_ip_address,omitempty"`
	// The IDs (separated by commas) of global EIPs connected by SNAT rule.
	GlobalEipId string `json:"global_eip_id,omitempty"`
	// The description of the SNAT rule.
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

// Delete is a method used to delete an existing SNAT rule.
func Delete(c *golangsdk.ServiceClient, gatewayId, ruleId string) error {
	_, err := c.Delete(deleteURL(c, gatewayId, ruleId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
