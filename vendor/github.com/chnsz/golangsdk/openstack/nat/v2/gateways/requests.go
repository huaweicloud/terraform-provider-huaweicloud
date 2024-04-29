package gateways

import "github.com/chnsz/golangsdk"

// CreateOpts is the structure used to create a new gateway for NAT service.
type CreateOpts struct {
	// The gateway name.
	Name string `json:"name" required:"true"`
	// The ID of the VPC to which the gateway belongs.
	VpcId string `json:"router_id" required:"true"`
	// The network ID that VPC have.
	InternalNetworkId string `json:"internal_network_id" required:"true"`
	// The gateway specification.
	Spec string `json:"spec" required:"true"`
	// The gateway description.
	Description string `json:"description,omitempty"`
	// The private IP address of the public NAT gateway.
	// The IP address is assigned by the VPC subnet.
	NgportIpAddress string `json:"ngport_ip_address,omitempty"`
	// The enterprise project ID to which the gateway belongs.
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create a gateway using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*Gateway, error) {
	b, err := golangsdk.BuildRequestBody(opts, "nat_gateway")
	if err != nil {
		return nil, err
	}

	var r createResp
	_, err = c.Post(rootURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Gateway, err
}

// Get is a method used to obtain the gateway detail by its ID.
func Get(c *golangsdk.ServiceClient, gatewayId string) (*Gateway, error) {
	var r queryResp
	_, err := c.Get(resourceURL(c, gatewayId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Gateway, err
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// The project ID.
	TenantId string `q:"tenant_id"`
	// The gateway ID.
	ID string `q:"id"`
	// The enterprise project ID to which the gateway belongs.
	EnterpriseProjectId string `q:"enterprise_project_id"`
	// The gateway description.
	Description string `q:"description"`
	// The creation time.
	CreatedAt string `q:"created_at"`
	// The gateway name.
	Name string `q:"name"`
	// The status of the gateway name.
	Status string `q:"status"`
	// The gateway specification.
	Spec string `q:"spec"`
	// The frozen status.
	AdminStateUp string `q:"admin_state_up"`
	// The network ID that VPC have.
	InternalNetworkId string `q:"internal_network_id"`
	// The ID of the VPC to which the gateway belongs.
	VpcId string `q:"router_id"`
	// The ID of the VPC to which the gateway belongs.
	Limit int `q:"limit"`
}

// List is a method to query all gateways using given parameters.
func List(client *golangsdk.ServiceClient, opts ListOpts) ([]Gateway, error) {
	url := rootURL(client)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r listResp
	_, err = client.Get(url, &r, &golangsdk.RequestOpts{
		MoreHeaders: client.MoreHeaders,
	})
	if err != nil {
		return nil, err
	}
	return r.Gateways, nil
}

// UpdateOpts is the structure used to modify an existing gateway.
type UpdateOpts struct {
	// The gateway name.
	Name string `json:"name,omitempty"`
	// The gateway description.
	Description *string `json:"description,omitempty"`
	// The gateway specification.
	Spec string `json:"spec,omitempty"`
}

// Update is a method used to modify an existing gateway using given parameters.
func Update(c *golangsdk.ServiceClient, gatewayId string, opts UpdateOpts) (*Gateway, error) {
	b, err := golangsdk.BuildRequestBody(opts, "nat_gateway")
	if err != nil {
		return nil, err
	}

	var r updateResp
	_, err = c.Put(resourceURL(c, gatewayId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Gateway, err
}

// Delete is a method to remove the specified gateway using its ID.
func Delete(c *golangsdk.ServiceClient, gatewayId string) error {
	_, err := c.Delete(resourceURL(c, gatewayId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
