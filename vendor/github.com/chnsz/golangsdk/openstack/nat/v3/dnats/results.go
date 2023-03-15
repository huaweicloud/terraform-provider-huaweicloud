package dnats

// Rule is a struct that represents the private DNAT rule detail.
type Rule struct {
	// The ID of the DNAT rule.
	ID string `json:"id"`
	// The project ID.
	ProjectId string `json:"project_id"`
	// The description of the DNAT rule, which contain maximum of `255` characters, and angle brackets (< and >) are
	// not allowed.
	Description string `json:"description"`
	// The ID of the transit IP.
	TransitIpId string `json:"transit_ip_id"`
	// The ID of the gateway to which the private DNAT rule belongs.
	GatewayId string `json:"gateway_id"`
	// The network interface ID of the transit IP for private NAT.
	NetworkInterfaceId string `json:"network_interface_id"`
	// The backend type of the DNAT rule.
	// The values of the backend type are as follow:
	// + COMPUTE
	// + VIP
	// + ELB
	// + ELBv3
	// + CUSTOMIZE
	Type string `json:"type"`
	// The protocol type. The valid values (and the related protocol numbers) are 'UDP/udp (6)', 'TCP/tcp' (17) and
	// 'ANY/any (0)'.
	Protocol string `json:"protocol"`
	// The private IP address of the backend instance.
	PrivateIpAddress string `json:"private_ip_address"`
	// The port of the backend instance.
	InternalServicePort string `json:"internal_service_port"`
	// The port of the transit IP.
	TransitServicePort string `json:"transit_service_port"`
	// The ID of the enterprise project to which the private DNAT rule belongs.
	EnterpriseProjectId string `json:"enterprise_project_id"`
	// The creation time of the private DNAT rule.
	CreatedAt string `json:"created_at"`
	// The latest update time of the private DNAT rule.
	UpdatedAt string `json:"updated_at"`
}

type createResp struct {
	// The DNAT rule detail.
	Rule Rule `json:"dnat_rule"`
	// The request ID.
	RequestId string `json:"request_id"`
}

type queryResp struct {
	// The DNAT rule detail.
	Rule Rule `json:"dnat_rule"`
	// The request ID.
	RequestId string `json:"request_id"`
}

type updateResp struct {
	// The DNAT rule detail.
	Rule Rule `json:"dnat_rule"`
	// The request ID.
	RequestId string `json:"request_id"`
}
