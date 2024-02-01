package eips

import "github.com/chnsz/golangsdk/openstack/networking/v1/bandwidths"

type PublicIp struct {
	ID                  string   `json:"id"`
	Status              string   `json:"status"`
	Type                string   `json:"publicip_pool_name"`
	PublicAddress       string   `json:"public_ip_address"`
	PublicIpv6Address   string   `json:"public_ipv6_address"`
	EnterpriseProjectID string   `json:"enterprise_project_id"`
	IpVersion           int      `json:"ip_version"`
	Alias               string   `json:"alias"`
	BillingInfo         string   `json:"billing_info"`
	CreatedAt           string   `json:"created_at"`
	UpdatedAt           string   `json:"updated_at"`
	Tags                []string `json:"tags"`

	AssociateInstanceType string `json:"associate_instance_type"`
	AssociateInstanceID   string `json:"associate_instance_id"`

	Bandwidth bandwidths.BandWidth `json:"bandwidth"`

	// Vnic return the port info if EIP associate to an instance with port
	Vnic Vnic `json:"vnic"`
}

type Vnic struct {
	PrivateAddress string `json:"private_ip_address"`
	PortID         string `json:"port_id"`
	IntsanceID     string `json:"instance_id"`
	IntsanceType   string `json:"instance_type"`
}

func (r GetResult) Extract() (PublicIp, error) {
	var getResp struct {
		IP PublicIp `json:"publicip"`
	}
	err := r.Result.ExtractInto(&getResp)
	return getResp.IP, err
}
