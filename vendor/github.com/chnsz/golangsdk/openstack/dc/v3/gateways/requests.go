package gateways

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

// CreateOpts is the structure used to create a virtual gateway.
type CreateOpts struct {
	// The ID of the VPC connected to the virtual gateway.
	VpcId string `json:"vpc_id" required:"true"`
	// The list of IPv4 subnets from the virtual gateway to access cloud services, which is usually the CIDR block of
	// the VPC.
	LocalEpGroup []string `json:"local_ep_group" required:"true"`
	// The list of IPv6 subnets from the virtual gateway to access cloud services, which is usually the CIDR block of
	// the VPC.
	LocalEpGroupIpv6 []string `json:"local_ep_group_ipv6,omitempty"`
	// Specifies the name of the virtual gateway.
	// The valid length is limited from 0 to 64, only chinese and english letters, digits, hyphens (-), underscores (_)
	// and dots (.) are allowed.
	// The name must start with a chinese or english letter, and the Chinese characters must be in **UTF-8** or
	// **Unicode** format.
	Name string `json:"name,omitempty"`
	// Specifies the description of the virtual gateway.
	// The description contain a maximum of 64 characters and the angle brackets (< and >) are not allowed.
	// Chinese characters must be in **UTF-8** or **Unicode** format.
	Description string `json:"description,omitempty"`
	// The local BGP ASN of the virtual gateway.
	BgpAsn int `json:"bgp_asn,omitempty"`
	// The enterprise project ID to which the virtual gateway belongs.
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
	// The key/value pairs to associate with the virtual gateway.
	Tags []tags.ResourceTag `json:"tags,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create a virtual gateway using given parameters.
func Create(client *golangsdk.ServiceClient, opts CreateOpts) (*VirtualGateway, error) {
	b, err := golangsdk.BuildRequestBody(opts, "virtual_gateway")
	if err != nil {
		return nil, err
	}

	var r createResp
	_, err = client.Post(rootURL(client), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.VirtualGateway, err
}

// UpdateOpts is the structure used to update the configuration of the virtual gateway.
type UpdateOpts struct {
	// The list of IPv4 subnets from the virtual gateway to access cloud services, which is usually the CIDR block of
	// the VPC.
	LocalEpGroup []string `json:"local_ep_group" required:"true"`
	// The list of IPv6 subnets from the virtual gateway to access cloud services, which is usually the CIDR block of
	// the VPC.
	LocalEpGroupIpv6 []string `json:"local_ep_group_ipv6,omitempty"`
	// Specifies the name of the virtual gateway.
	// The valid length is limited from 0 to 64, only chinese and english letters, digits, hyphens (-), underscores (_)
	// and dots (.) are allowed.
	// The name must start with a chinese or english letter, and the Chinese characters must be in **UTF-8** or
	// **Unicode** format.
	Name string `json:"name,omitempty"`
	// Specifies the description of the virtual gateway.
	// The description contain a maximum of 64 characters and the angle brackets (< and >) are not allowed.
	// Chinese characters must be in **UTF-8** or **Unicode** format.
	Description *string `json:"description,omitempty"`
}

// Update is a method used to update the specified virtual gateway using given parameters.
func Update(client *golangsdk.ServiceClient, gatewayId string, opts UpdateOpts) (*VirtualGateway, error) {
	b, err := golangsdk.BuildRequestBody(opts, "virtual_gateway")
	if err != nil {
		return nil, err
	}

	var r createResp
	_, err = client.Put(resourceURL(client, gatewayId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.VirtualGateway, err
}

// Get is a method used to obtain the details of the virtual gateway using its ID.
func Get(client *golangsdk.ServiceClient, gatewayId string) (*VirtualGateway, error) {
	var r getResp
	_, err := client.Get(resourceURL(client, gatewayId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.VirtualGateway, err
}

// Delete is a method used to remove an existing virtual gateway using its ID.
func Delete(client *golangsdk.ServiceClient, gatewayId string) error {
	_, err := client.Delete(resourceURL(client, gatewayId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
