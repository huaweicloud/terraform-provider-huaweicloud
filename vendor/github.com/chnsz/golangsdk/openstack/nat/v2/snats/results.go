package snats

// Rule is a struct that represents the SNAT rule detail.
type Rule struct {
	// The ID of SNAT rule.
	ID string `json:"id"`
	// The ID of the gateway to which the SNAT rule belongs.
	GatewayId string `json:"nat_gateway_id"`
	// The network IDs of subnet connected by SNAT rule (VPC side).
	NetworkId string `json:"network_id"`
	// The project ID.
	TenantId string `json:"tenant_id"`
	// The IDs (separated by commas) of the floating IP that the SNAT rule have.
	FloatingIpId string `json:"floating_ip_id"`
	// The floating IP addresses (separated by commas) connected by SNAT rule.
	FloatingIpAddress string `json:"floating_ip_address"`
	// The IDs (separated by commas) of the global EIPs connected by SNAT rule.
	GlobalEipId string `json:"global_eip_id"`
	// The global EIP addresses (separated by commas) connected by SNAT rule.
	GlobalEipAddress string `json:"global_eip_address"`
	// The description of the SNAT rule.
	Description string `json:"description"`
	// The status of the SNAT rule.
	Status string `json:"status"`
	// The frozen status.
	AdminStateUp bool `json:"admin_state_up"`
	// The CIDR block connected by SNAT rule (DC side).
	Cidr string `json:"cidr"`
	// The resource type of the SNAT rule.
	// The valid values are as follows:
	// + 0: VPC side.
	// + 1: DC side.
	SourceType int `json:"source_type"`
}

type createResp struct {
	// The SNAT rule detail.
	Rule Rule `json:"snat_rule"`
}

type queryResp struct {
	// The SNAT rule detail.
	Rule Rule `json:"snat_rule"`
}

type updateResp struct {
	// The SNAT rule detail.
	Rule Rule `json:"snat_rule"`
}
