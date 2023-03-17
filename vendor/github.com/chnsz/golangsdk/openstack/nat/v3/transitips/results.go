package transitips

import "github.com/chnsz/golangsdk/openstack/common/tags"

// TransitIp is the structure that represents the detail of the transit IP for private NAT.
type TransitIp struct {
	// The transit IP ID.
	ID string `json:"id"`
	// The project ID to which the transit IP (private NAT) belongs.
	ProjectId string `json:"project_id"`
	// The network interface ID of the transit IP for private NAT.
	NetworkInterfaceId string `json:"network_interface_id"`
	// The IP address of the transit IP.
	IpAddress string `json:"ip_address"`
	// The creation time.
	CreatedAt string `json:"created_at"`
	// The latest update time.
	UpdatedAt string `json:"updated_at"`
	// The ID of the subnet to which the transit IP belongs.
	SubnetId string `json:"virsubnet_id"`
	// The key/value pairs to associate with the transit IP.
	Tags []tags.ResourceTag `json:"tags"`
	// The ID of the private NAT gateway to which the transit IP belongs.
	GatewayId string `json:"gateway_id"`
	// The ID of the enterprise project to which the transit IP belongs.
	EnterpriseProjectId string `json:"enterprise_project_id"`
}

type createResp struct {
	// The transit IP detail.
	TransitIp TransitIp `json:"transit_ip"`
	// Request ID.
	RequestId string `json:"request_id"`
}

type queryResp struct {
	// The transit IP detail.
	TransitIp TransitIp `json:"transit_ip"`
	// Request ID.
	RequestId string `json:"request_id"`
}
