package gateways

import "github.com/chnsz/golangsdk/openstack/common/tags"

// Gateway is the structure represents the private NAT gateway details.
type Gateway struct {
	// The ID of the private NAT gateway.
	ID string `json:"id"`
	// The project ID to which the private NAT gateway belongs.
	ProjectId string `json:"project_id"`
	// The name of the private NAT gateway.
	// The valid length is limited from `1` to `64`, only letters, digits, hyphens (-) and underscores (_) are allowed.
	Name string `json:"name"`
	// The description of the private NAT gateway, which contain maximum of `255` characters, and
	// angle brackets (<) and (>) are not allowed.
	Description string `json:"description"`
	// The specification of the private NAT gateway.
	// The valid values are as follows:
	// + **Small**: Small type, which supports up to `20` rules, `200 Mbit/s` bandwidth, `20,000` PPS and `2,000` SNAT
	//   connections.
	// + **Medium**: Medium type, which supports up to `50` rules, `500 Mbit/s` bandwidth, `50,000` PPS and `5,000` SNAT
	//   connections.
	// + **Large**: Large type, which supports up to `200` rules, `2 Gbit/s` bandwidth, `200,000` PPS and `20,000` SNAT
	//   connections.
	// + **Extra-Large**: Extra-large type, which supports up to `500` rules, `5 Gbit/s` bandwidth, `500,000` PPS and
	//   `50,000` SNAT connections.
	Spec string `json:"spec"`
	// The current status of the private NAT gateway.
	Status string `json:"status"`
	// The creation time of the private NAT gateway.
	CreatedAt string `json:"created_at"`
	// The latest update time of the private NAT gateway.
	UpdatedAt string `json:"updated_at"`
	// The VPC configuration of the private NAT gateway.
	DownLinkVpcs []DownLinkVpcResp `json:"downlink_vpcs"`
	// The key/value pairs to associate with the NAT geteway.
	Tags []tags.ResourceTag `json:"tags"`
	// The enterprise project ID to which the private NAT gateway belongs.
	EnterpriseProjectId string `json:"enterprise_project_id"`
}

// DownLinkVpcResp is an object that represents the VPC configuration to which private NAT gateway belongs.
type DownLinkVpcResp struct {
	// The subnet ID to which the private NAT gateway belongs.
	SubnetId string `json:"virsubnet_id"`
	// The VPC ID to which the private NAT gateway belongs.
	VpcId string `json:"vpc_id"`
}

type createResp struct {
	// The gateway detail.
	Gateway Gateway `json:"gateway"`
}

type queryResp struct {
	// The gateway detail.
	Gateway Gateway `json:"gateway"`
}

type updateResp struct {
	// The gateway detail.
	Gateway Gateway `json:"gateway"`
}
