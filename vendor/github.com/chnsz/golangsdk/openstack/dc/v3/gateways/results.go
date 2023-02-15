package gateways

import (
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

// createResp is the structure that represents the API response of the 'Create' method, which contains virtual gateway
// details.
type createResp struct {
	// The response detail of the virtual gateway.
	VirtualGateway VirtualGateway `json:"virtual_gateway"`
}

// VirtualGateway is the structure that represents the details of the virtual gateway.
type VirtualGateway struct {
	// The ID of the virtual gateway.
	ID string `json:"id"`
	// The ID of the VPC connected to the virtual gateway.
	VpcId string `json:"vpc_id"`
	// The project ID to which the virtual gateway belongs.
	TenantId string `json:"tenant_id"`
	// Specifies the name of the virtual gateway.
	// The valid length is limited from 0 to 64, only chinese and english letters, digits, hyphens (-), underscores (_)
	// and dots (.) are allowed.
	// The name must start with a chinese or english letter, and the Chinese characters must be in **UTF-8** or
	// **Unicode** format.
	Name string `json:"name"`
	// Specifies the description of the virtual gateway.
	// The description contain a maximum of 64 characters and the angle brackets (< and >) are not allowed.
	// Chinese characters must be in **UTF-8** or **Unicode** format.
	Description string `json:"description"`
	// The type of virtual gateway.
	Type string `json:"type"`
	// The list of IPv4 subnets from the virtual gateway to access cloud services, which is usually the CIDR block of
	// the VPC.
	LocalEpGroup []string `json:"local_ep_group"`
	// The list of IPv6 subnets from the virtual gateway to access cloud services, which is usually the CIDR block of
	// the VPC.
	LocalEpGroupIpv6 []string `json:"local_ep_group_ipv6"`
	// The current status of the virtual gateway.
	Status string `json:"status"`
	// The local BGP ASN of the virtual gateway.
	BgpAsn int `json:"bgp_asn"`
	// The enterprise project ID to which the virtual gateway belongs.
	EnterpriseProjectId string `json:"enterprise_project_id"`
	// The key/value pairs to associate with the virtual gateway.
	Tags []tags.ResourceTag `json:"tags"`
}

// getResp is the structure that represents the API response of the 'Get' method, which contains virtual gateway
// details.
type getResp struct {
	// The response detail of the virtual gateway.
	VirtualGateway VirtualGateway `json:"virtual_gateway"`
}
