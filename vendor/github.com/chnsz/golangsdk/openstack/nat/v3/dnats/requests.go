package dnats

import "github.com/chnsz/golangsdk"

// CreateOpts is the structure used to create a new private DNAT rule.
type CreateOpts struct {
	// The ID of the gateway to which the private DNAT rule belongs.
	GatewayId string `json:"gateway_id" required:"true"`
	// The ID of the transit IP for private NAT.
	TransitIpId string `json:"transit_ip_id" required:"true"`
	// The description of the DNAT rule, which contain maximum of `255` characters, and angle brackets (< and >) are
	// not allowed.
	Description string `json:"description,omitempty"`
	// The network interface ID of the transit IP for private NAT.
	NetworkInterfaceId string `json:"network_interface_id,omitempty"`
	// The protocol type of the private DNAT rule.
	// The valid values (and the related protocol numbers) are 'UDP/udp (6)', 'TCP/tcp' (17) and 'ANY/any (0)'.
	Protocol string `json:"protocol,omitempty"`
	// The private IP address of the backend instance.
	PrivateIpAddress string `json:"private_ip_address,omitempty"`
	// The port of the backend instance.
	InternalServicePort string `json:"internal_service_port,omitempty"`
	// The port of the transit IP.
	TransitServicePort string `json:"transit_service_port,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create a private DNAT rule using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*Rule, error) {
	b, err := golangsdk.BuildRequestBody(opts, "dnat_rule")
	if err != nil {
		return nil, err
	}

	var r createResp
	_, err = c.Post(rootURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Rule, err
}

// Get is a method used to obtain the private DNAT rule detail by its ID.
func Get(c *golangsdk.ServiceClient, ruleId string) (*Rule, error) {
	var r queryResp
	_, err := c.Get(resourceURL(c, ruleId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Rule, err
}

// UpdateOpts is the structure used to modify an existing private DNAT rule.
type UpdateOpts struct {
	// The ID of the transit IP for private NAT.
	TransitIpId string `json:"transit_ip_id,omitempty"`
	// The description of the DNAT rule, which contain maximum of `255` characters, and angle brackets (< and >) are
	// not allowed.
	Description *string `json:"description,omitempty"`
	// The network interface ID of the transit IP for private NAT.
	NetworkInterfaceId string `json:"network_interface_id,omitempty"`
	// The protocol type. The valid values (and the related protocol numbers) are 'UDP/udp (6)', 'TCP/tcp' (17) and
	// 'ANY/any (0)'.
	Protocol string `json:"protocol,omitempty"`
	// The private IP address of the backend instance.
	PrivateIpAddress string `json:"private_ip_address,omitempty"`
	// The port of the backend instance.
	InternalServicePort string `json:"internal_service_port,omitempty"`
	// The port of the transit IP.
	TransitServicePort string `json:"transit_service_port,omitempty"`
}

// Update is a method used to modify an existing DNAT rule using given parameters.
func Update(c *golangsdk.ServiceClient, ruleId string, opts UpdateOpts) (*Rule, error) {
	b, err := golangsdk.BuildRequestBody(opts, "dnat_rule")
	if err != nil {
		return nil, err
	}

	var r updateResp
	_, err = c.Put(resourceURL(c, ruleId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Rule, err
}

// Delete is a method to remove the specified DNAT rule using its ID.
func Delete(c *golangsdk.ServiceClient, ruleId string) error {
	_, err := c.Delete(resourceURL(c, ruleId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
