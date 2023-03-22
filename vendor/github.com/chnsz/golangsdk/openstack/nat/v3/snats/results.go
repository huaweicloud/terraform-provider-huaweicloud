package snats

// Rule is a struct that represents the private SNAT rule detail.
type Rule struct {
	// The ID of the SNAT rule.
	ID string `json:"id"`
	// The project ID.
	ProjectId string `json:"project_id"`
	// The ID of the gateway to which the private SNAT rule belongs.
	GatewayId string `json:"gateway_id"`
	// The CIDR block of the match rule.
	// Exactly one of cidr and virsubnet_id must be set.
	Cidr string `json:"cidr"`
	// The subnet ID of the match rule.
	// Exactly one of cidr and virsubnet_id must be set.
	SubnetId string `json:"virsubnet_id"`
	// The description of the SNAT rule, which contain maximum of `255` characters, and angle brackets (< and >) are
	// not allowed.
	Description string `json:"description"`
	// The ID of the transit IP.
	TransitIpAssociations []AssociatedTransitIp `json:"transit_ip_associations"`
	// The ID of the enterprise project to which the private SNAT rule belongs.
	EnterpriseProjectId string `json:"enterprise_project_id"`
	// The creation time of the private SNAT rule.
	CreatedAt string `json:"created_at"`
	// The latest update time of the private SNAT rule.
	UpdatedAt string `json:"updated_at"`
}

type AssociatedTransitIp struct {
	// The ID of the transit IP.
	ID string `json:"transit_ip_id"`
	// The address of the transit IP.
	Address string `json:"transit_ip_address"`
}

type createResp struct {
	// The SNAT rule detail.
	Rule Rule `json:"snat_rule"`
	// The request ID.
	RequestId string `json:"request_id"`
}

type queryResp struct {
	// The SNAT rule detail.
	Rule Rule `json:"snat_rule"`
	// The request ID.
	RequestId string `json:"request_id"`
}

type updateResp struct {
	// The SNAT rule detail.
	Rule Rule `json:"snat_rule"`
	// The request ID.
	RequestId string `json:"request_id"`
}
