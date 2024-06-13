package interfaces

import (
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

// createResp is the structure that represents the API response of the 'Create' method, which contains request ID and
// the virtual interface details.
type createResp struct {
	RequestId string `json:"request_id"`
	// The response detail of the virtual gateway.
	VirtualInterface VirtualInterface `json:"virtual_interface"`
}

// VirtualInterface is the structure that represents the details of the virtual interface.
type VirtualInterface struct {
	// The ID of the virtual interface.
	ID string `json:"id"`
	// Specifies the name of the virtual interface.
	// The valid length is limited from 0 to 64, only chinese and english letters, digits, hyphens (-), underscores (_)
	// and dots (.) are allowed.
	// The name must start with a chinese or english letter, and the Chinese characters must be in **UTF-8** or
	// **Unicode** format.
	Name string `json:"name"`
	// Manage status.
	// The valid values are 'true' and 'false'.
	AdminStateUp bool `json:"admin_state_up"`
	// The ingress bandwidth size of the virtual interface.
	Bandwidth int `json:"bandwidth"`
	// The creation time of the virtual interface.
	CreatedAt string `json:"create_time"`
	// The latest update time of the virtual interface.
	UpdatedAt string `json:"update_time"`
	// Specifies the description of the virtual interface.
	// The description contain a maximum of 128 characters and the angle brackets (< and >) are not allowed.
	// Chinese characters must be in **UTF-8** or **Unicode** format.
	Description string `json:"description"`
	// The ID of the direct connection associated with the virtual interface.
	DirectConnectId string `json:"direct_connect_id"`
	// The service type of the virtual interface.
	ServiceType string `json:"service_type"`
	// The current status of the virtual interface.
	// The valid values are as follows:
	// + ACTIVE
	// + DOWN
	// + BUILD
	// + ERROR
	// + PENDING_CREATE
	// + PENDING_UPDATE
	// + PENDING_DELETE
	// + DELETED
	// + AUTHORIZATION
	// + REJECTED
	Status string `json:"status"`
	// The ID of the target tenant ID, which is used for cross tenant virtual interface creation.
	TenantId string `json:"tenant_id"`
	// The type of the virtual interface.
	Type string `json:"type"`
	// The ID of the virtual gateway to which the virtual interface is connected.
	VgwId string `json:"vgw_id"`
	// The VLAN for constom side.
	Vlan int `json:"vlan"`
	// The route specification of the remote VIF network.
	RouteLimit int `json:"route_limit"`
	// Whether to enable the Bidirectional Forwarding Detection (BFD) function.
	EnableBfd bool `json:"enable_bfd"`
	// Whether to enable the Network Quality Analysis (NQA) function.
	EnableNqa bool `json:"enable_nqa"`
	// The ID of the Intelligent EdgeSite (IES) associated with the virtual interface.
	IesId string `json:"ies_id"`
	// The ID of the link aggregation group (LAG) associated with the virtual interface.
	LagId string `json:"lag_id"`
	// The ID of the local gateway (LGW) associated with the virtual interface.
	LgwId string `json:"lgw_id"`
	// The local BGP ASN in client side.
	BgpAsn int `json:"bgp_asn"`
	// The (MD5) password for the local BGP.
	BgpMd5 string `json:"bgp_md5"`
	// The attributed Device ID.
	DeviceId string `json:"device_id"`
	// The IPv4 address of the virtual interface in cloud side.
	LocalGatewayV4Ip string `json:"local_gateway_v4_ip"`
	// The IPv4 address of the virtual interface in client side.
	RemoteGatewayV4Ip string `json:"remote_gateway_v4_ip"`
	// The address family type.
	AddressFamily string `json:"address_family"`
	// The IPv6 address of the virtual interface in cloud side.
	LocalGatewayV6Ip string `json:"local_gateway_v6_ip"`
	// The IPv6 address of the virtual interface in client side.
	RemoteGatewayV6Ip string `json:"remote_gateway_v6_ip"`
	// The CIDR list of remote subnets.
	RemoteEpGroup []string `json:"remote_ep_group"`
	// The CIDR list of subnets in service side.
	ServiceEpGroup []string `json:"service_ep_group"`
	// The route mode of the virtual interface.
	RouteMode string `json:"route_mode"`
	// Whether limit rate.
	RateLimit bool `json:"rate_limit"`
	// The VLAN for constom side.
	VifPeers []VifPeer `json:"vif_peers"`
	// The Peer details of the VIF.
	ExtendAttribute VifExtendAttribute `json:"extend_attribute"`
	// The enterprise project ID to which the virtual interface belongs.
	EnterpriseProjectId string `json:"enterprise_project_id"`
	// The key/value pairs to associate with the virtual interface.
	Tags []tags.ResourceTag `json:"tags"`
}

// VifPeer is the structure that represents the related information of each peer.
type VifPeer struct {
	// Resource ID.
	ID string `json:"id"`
	// The ID of the target tenant ID, which is used for cross tenant virtual interface creation.
	TenantId string `json:"tenant_id"`
	// Specifies the name of the VIF peer.
	// The valid length is limited from 0 to 64, only chinese and english letters, digits, hyphens (-), underscores (_)
	// and dots (.) are allowed.
	// The name must start with a chinese or english letter, and the Chinese characters must be in **UTF-8** or
	// **Unicode** format.
	Name string `json:"name"`
	// Specifies the description of the virtual interface.
	// The description contain a maximum of 128 characters and the angle brackets (< and >) are not allowed.
	// Chinese characters must be in **UTF-8** or **Unicode** format.
	Description string `json:"description"`
	// The address family type.
	AddressFamily string `json:"address_family"`
	// Local gateway IP.
	LocalGatewayIp string `json:"local_gateway_ip"`
	// Remote gateway IP.
	RemoteGatewayIp string `json:"remote_gateway_ip"`
	// The routing mode, which can be static or bgp.
	RouteMode string `json:"route_mode"`
	// BGP ASN.
	BgpAsn int `json:"bgp_asn"`
	// BGP MD5 password.
	BgpMd5 string `json:"bgp_md5"`
	// The CIDR list of remote subnets.
	RemoteEpGroup []string `json:"remote_ep_group"`
	// The CIDR list of subnets in service side.
	ServiceEpGroup []string `json:"service_ep_group"`
	// Attributed Device ID.
	DeviceId string `json:"device_id"`
	// Whether to enable BFD.
	EnableBfd bool `json:"enable_bfd"`
	// Whether to enable NQA.
	EnableNqa bool `json:"enable_nqa"`
	// Attributed Device ID.
	BgpRouteLimit int `json:"bgp_route_limit"`
	// Attributed Device ID.
	BgpStatus string `json:"bgp_status"`
	// The status of the virtual interface peer.
	Status string `json:"status"`
	// The virtual interface ID corresponding to the VIF peer.
	VifId string `json:"vif_id"`
	// The number of received BGP routes if BGP routing is used.
	ReceiveRouteNum int `json:"receive_route_num"`
}

// VifExtendAttribute is the structure that represents the reliability detection information of BFD/NQA.
type VifExtendAttribute struct {
	// The availability detection types for virtual interface.
	// + nqa
	// + bfd
	HaType string `json:"ha_type"`
	// The specific configuration mode detected for virtual interface.
	// + auto_single
	// + auto_multi
	// + static_single
	// + static_multi
	// + enhance_nqa
	HaMode string `json:"ha_mode"`
	// The detection retries.
	DetectMultiplier string `json:"detect_multiplier"`
	// The receive time interval for detection.
	MinRxInterval string `json:"min_rx_interval"`
	// The transmit time interval for detection.
	MinTxInterval string `json:"min_tx_interval"`
	// The identifier of the detected remote side, used for static BFD.
	RemoteDisclaim string `json:"remote_disclaim"`
	// The identifier of the detected local side, used for static BFD.
	LocalDisclaim string `json:"local_disclaim"`
}

// getResp is the structure that represents the API response of the 'Get' method, which contains request ID and the
// virtual interface details.
type getResp struct {
	RequestId string `json:"request_id"`
	// The response detail of the virtual interface.
	VirtualInterface VirtualInterface `json:"virtual_interface"`
}

// createResp is the structure that represents the API response of the 'Update' method, which contains request ID and
// the virtual interface details.
type updateResp struct {
	RequestId string `json:"request_id"`
	// The response detail of the virtual gateway.
	VirtualInterface VirtualInterface `json:"virtual_interface"`
}
