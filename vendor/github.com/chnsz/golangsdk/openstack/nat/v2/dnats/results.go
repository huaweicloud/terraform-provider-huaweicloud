package dnats

// Rule is a struct that represents the DNAT rule detail.
type Rule struct {
	// The ID of the DNAT rule.
	ID string `json:"id"`
	// The project ID.
	TenantId string `json:"tenant_id"`
	// The description of the DNAT rule.
	Description string `json:"description"`
	// The port ID of network.
	PortId string `json:"port_id"`
	// The private IP address of a user.
	PrivateIp string `json:"private_ip"`
	// The ID of the gateway to which the DNAT rule belongs.
	GatewayId string `json:"nat_gateway_id"`
	// The IDs of floating IP connected by DNAT rule.
	FloatingIpId string `json:"floating_ip_id"`
	// The floating IP address connected by DNAT rule.
	FloatingIpAddress string `json:"floating_ip_address"`
	// The ID of the global EIP connected by the DNAT rule.
	GlobalEipId string `json:"global_eip_id"`
	// The global EIP address connected by the DNAT rule.
	GlobalEipAddress string `json:"global_eip_address"`
	// The current status of the DNAT rule.
	Status string `json:"status"`
	// The frozen status.
	AdminStateUp bool `json:"admin_state_up"`
	// The port used by Floating IP provide services for external systems.
	InternalServicePort int `json:"internal_service_port"`
	// The port used by ECSs or BMSs to provide services for external systems.
	ExternalServicePort int `json:"external_service_port"`
	// The port range used by Floating IP provide services for external systems.
	InternalServicePortRange string `json:"internal_service_port_range"`
	// The port range used by ECSs or BMSs to provide services for external systems.
	EXternalServicePortRange string `json:"external_service_port_range"`
	// The protocol type. The valid values are 'udp', 'tcp' and 'any'.
	Protocol string `json:"protocol"`
	// The creation time of the dnat rule.
	CreatedAt string `json:"created_at"`
}

type createResp struct {
	// The DNAT rule detail.
	Rule Rule `json:"dnat_rule"`
}

type queryResp struct {
	// The DNAT rule detail.
	Rule Rule `json:"dnat_rule"`
}

type updateResp struct {
	// The DNAT rule detail.
	Rule Rule `json:"dnat_rule"`
}
