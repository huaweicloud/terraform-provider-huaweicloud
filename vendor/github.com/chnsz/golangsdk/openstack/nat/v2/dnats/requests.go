package dnats

import "github.com/chnsz/golangsdk"

// CreateOpts is the structure used to create a new DNAT rule.
type CreateOpts struct {
	// The ID of the gateway to which the DNAT rule belongs.
	GatewayId string `json:"nat_gateway_id" required:"true"`
	// The IDs of floating IP connected by DNAT rule.
	FloatingIpId string `json:"floating_ip_id,omitempty"`
	// The ID of the global EIP connected by the DNAT rule.
	GlobalEipId string `json:"global_eip_id,omitempty"`
	// The protocol type. The valid values are 'udp', 'tcp' and 'any'.
	Protocol string `json:"protocol" required:"true"`
	// The port used by Floating IP provide services for external systems.
	InternalServicePort *int `json:"internal_service_port" required:"true"`
	// The port used by ECSs or BMSs to provide services for external systems.
	ExternalServicePort *int `json:"external_service_port" required:"true"`
	// The port range used by Floating IP provide services for external systems.
	InternalServicePortRange string `json:"internal_service_port_range,omitempty"`
	// The port range used by ECSs or BMSs to provide services for external systems.
	EXternalServicePortRange string `json:"external_service_port_range,omitempty"`
	// The description of the DNAT rule.
	Description string `json:"description,omitempty"`
	// The port ID of network.
	PortId string `json:"port_id,omitempty"`
	// The private IP address of a user.
	PrivateIp string `json:"private_ip,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create a DNAT rule using given parameters.
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

// Get is a method used to obtain the DNAT rule detail by its ID.
func Get(c *golangsdk.ServiceClient, gatewayId string) (*Rule, error) {
	var r queryResp
	_, err := c.Get(resourceURL(c, gatewayId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Rule, err
}

// UpdateOpts is the structure used to modify an existing DNAT rule.
type UpdateOpts struct {
	// The ID of the gateway to which the DNAT rule belongs.
	GatewayId string `json:"nat_gateway_id" required:"true"`
	// The description of the DNAT rule.
	Description *string `json:"description,omitempty"`
	// The port ID of network.
	PortId string `json:"port_id,omitempty"`
	// The private IP address of a user.
	PrivateIp string `json:"private_ip,omitempty"`
	// The protocol type. The valid values are 'udp', 'tcp' and 'any'.
	Protocol string `json:"protocol,omitempty"`
	// The IDs of floating IP connected by DNAT rule.
	FloatingIpId string `json:"floating_ip_id,omitempty"`
	// The ID of the global EIP connected by the DNAT rule.
	GlobalEipId string `json:"global_eip_id,omitempty"`
	// The port used by Floating IP provide services for external systems.
	InternalServicePort *int `json:"internal_service_port,omitempty"`
	// The port used by ECSs or BMSs to provide services for external systems.
	ExternalServicePort *int `json:"external_service_port,omitempty"`
	// The port range used by Floating IP provide services for external systems.
	InternalServicePortRange string `json:"internal_service_port_range,omitempty"`
	// The port range used by ECSs or BMSs to provide services for external systems.
	ExternalServicePortRange string `json:"external_service_port_range,omitempty"`
}

// Update is a method used to modify an existing DNAT rule using given parameters.
func Update(c *golangsdk.ServiceClient, gatewayId string, opts UpdateOpts) (*Rule, error) {
	b, err := golangsdk.BuildRequestBody(opts, "dnat_rule")
	if err != nil {
		return nil, err
	}

	var r updateResp
	_, err = c.Put(resourceURL(c, gatewayId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Rule, err
}

// Delete is a method to remove the specified DNAT rule using its ID.
func Delete(c *golangsdk.ServiceClient, gatewayId, ruleId string) error {
	_, err := c.Delete(deleteURL(c, gatewayId, ruleId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
