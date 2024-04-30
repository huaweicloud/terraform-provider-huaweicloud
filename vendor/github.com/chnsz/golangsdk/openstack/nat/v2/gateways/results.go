package gateways

type createResp struct {
	// The gateway detail.
	Gateway Gateway `json:"nat_gateway"`
}

// Gateway is the structure that represents the detail of the NAT gateway.
type Gateway struct {
	// The gateway ID.
	ID string `json:"id"`
	// The project ID to which the gateway belongs.
	TenantId string `json:"tenant_id"`
	// The gateway name.
	Name string `json:"name"`
	// The gateway description.
	Description string `json:"description"`
	// The gateway specification.
	Spec string `json:"spec"`
	// The gateway status.
	// The valid values are as follows:
	// + ACTIVE
	// + PENDING_CREATE
	// + PENDING_UPDATE
	// + PENDING_DELETE
	// + INACTIVE
	Status string `json:"status"`
	// The frozen status.
	AdminStateUp bool `json:"admin_state_up"`
	// The creation time.
	CreatedAt string `json:"created_at"`
	// The ID of the VPC to which the gateway belongs.
	RouterId string `json:"router_id"`
	// The network ID that VPC have.
	InternalNetworkId string `json:"internal_network_id"`
	// The private IP address of the public NAT gateway.
	// The IP address is assigned by the VPC subnet.
	NgportIpAddress string `json:"ngport_ip_address"`
	// The enterprise project ID to which the gateway belongs.
	EnterpriseProjectId string `json:"enterprise_project_id"`
}

type queryResp struct {
	// The gateway detail.
	Gateway Gateway `json:"nat_gateway"`
}

type updateResp struct {
	// The gateway detail.
	Gateway Gateway `json:"nat_gateway"`
}

type listResp struct {
	// The list of the gateway details.
	Gateways []Gateway `json:"nat_gateways"`
}
