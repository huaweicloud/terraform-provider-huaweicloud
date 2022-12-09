package interfaces

import (
	"github.com/chnsz/golangsdk"
)

// CreateOpts is the structure used to create a virtual interface.
type CreateOpts struct {
	// The type of the virtual interface.
	Type string `json:"type" required:"true"`
	// The VLAN for constom side.
	Vlan int `json:"vlan" required:"true"`
	// The ingress bandwidth size of the virtual interface.
	Bandwidth int `json:"bandwidth" required:"true"`
	// The ID of the virtual gateway to which the virtual interface is connected.
	VgwId string `json:"vgw_id" required:"true"`
	// The route mode of the virtual interface.
	RouteMode string `json:"route_mode" required:"true"`
	// The CIDR list of remote subnets.
	RemoteEpGroup []string `json:"remote_ep_group" required:"true"`
	// The CIDR list of subnets in service side.
	ServiceEpGroup []string `json:"service_ep_group,omitempty"`
	// Specifies the name of the virtual interface.
	// The valid length is limited from 1 to 64, only chinese and english letters, digits, hyphens (-), underscores (_)
	// and dots (.) are allowed.
	// The name must start with a chinese or english letter, and the Chinese characters must be in **UTF-8** or
	// **Unicode** format.
	Name string `json:"name,omitempty"`
	// Specifies the description of the virtual interface.
	// The description contain a maximum of 128 characters and the angle brackets (< and >) are not allowed.
	// Chinese characters must be in **UTF-8** or **Unicode** format.
	Description string `json:"description,omitempty"`
	// The ID of the direct connection associated with the virtual interface.
	DirectConnectId string `json:"direct_connect_id,omitempty"`
	// The service type of the virtual interface.
	ServiceType string `json:"service_type,omitempty"`
	// The IPv4 address of the virtual interface in cloud side.
	LocalGatewayV4Ip string `json:"local_gateway_v4_ip,omitempty"`
	// The IPv4 address of the virtual interface in client side.
	RemoteGatewayV4Ip string `json:"remote_gateway_v4_ip,omitempty"`
	// The address family type.
	AddressFamily string `json:"address_family,omitempty"`
	// The IPv6 address of the virtual interface in cloud side.
	LocalGatewayV6Ip string `json:"local_gateway_v6_ip,omitempty"`
	// The IPv6 address of the virtual interface in client side.
	RemoteGatewayV6Ip string `json:"remote_gateway_v6_ip,omitempty"`
	// The local BGP ASN in client side.
	BgpAsn int `json:"bgp_asn,omitempty"`
	// The (MD5) password for the local BGP.
	BgpMd5 string `json:"bgp_md5,omitempty"`
	// Whether to enable the Bidirectional Forwarding Detection (BFD) function.
	EnableBfd bool `json:"enable_bfd,omitempty"`
	// Whether to enable the Network Quality Analysis (NQA) function.
	EnableNqa bool `json:"enable_nqa,omitempty"`
	// The ID of the link aggregation group (LAG) associated with the virtual interface.
	LagId string `json:"lag_id,omitempty"`
	// The ID of the target tenant ID, which is used for cross tenant virtual interface creation.
	ResourceTenantId string `json:"resource_tenant_id,omitempty"`
	// The enterprise project ID to which the virtual interface belongs.
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create a virtual interface using given parameters.
func Create(client *golangsdk.ServiceClient, opts CreateOpts) (*VirtualInterface, error) {
	b, err := golangsdk.BuildRequestBody(opts, "virtual_interface")
	if err != nil {
		return nil, err
	}

	var r createResp
	_, err = client.Post(rootURL(client), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.VirtualInterface, err
}

// Get is a method used to obtain the details of the virtual interface using its ID.
func Get(client *golangsdk.ServiceClient, interfaceId string) (*VirtualInterface, error) {
	var r getResp
	_, err := client.Get(resourceURL(client, interfaceId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.VirtualInterface, err
}

// UpdateOpts is the structure used to update an existing virtual interface.
type UpdateOpts struct {
	// Specifies the name of the virtual interface.
	// The valid length is limited from 0 to 64, only chinese and english letters, digits, hyphens (-), underscores (_)
	// and dots (.) are allowed.
	// The name must start with a chinese or english letter, and the Chinese characters must be in **UTF-8** or
	// **Unicode** format.
	Name string `json:"name,omitempty"`
	// Specifies the description of the virtual interface.
	// The description contain a maximum of 128 characters and the angle brackets (< and >) are not allowed.
	// Chinese characters must be in **UTF-8** or **Unicode** format.
	Description *string `json:"description,omitempty"`
	// The ingress bandwidth size of the virtual interface.
	Bandwidth int `json:"bandwidth,omitempty"`
	// The CIDR list of remote subnets.
	RemoteEpGroup []string `json:"remote_ep_group,omitempty"`
	// The CIDR list of subnets in service side.
	ServiceEpGroup []string `json:"service_ep_group,omitempty"`
	// Whether to enable the Bidirectional Forwarding Detection (BFD) function.
	EnableBfd *bool `json:"enable_bfd,omitempty"`
	// Whether to enable the Network Quality Analysis (NQA) function.
	EnableNqa *bool `json:"enable_nqa,omitempty"`
	// The status of the virtual interface to be changed.
	Status string `json:"status,omitempty"`
}

// Update is a method used to update the specified virtual interface using given parameters.
func Update(client *golangsdk.ServiceClient, interfaceId string, opts UpdateOpts) (*VirtualInterface, error) {
	b, err := golangsdk.BuildRequestBody(opts, "virtual_interface")
	if err != nil {
		return nil, err
	}

	var r updateResp
	_, err = client.Put(resourceURL(client, interfaceId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.VirtualInterface, err
}

// Delete is a method used to remove an existing virtual interface using its ID.
func Delete(client *golangsdk.ServiceClient, interfaceId string) error {
	_, err := client.Delete(resourceURL(client, interfaceId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
