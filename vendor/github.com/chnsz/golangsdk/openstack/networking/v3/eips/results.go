package eips

import (
	"github.com/chnsz/golangsdk/openstack/networking/v1/bandwidths"
	"github.com/chnsz/golangsdk/pagination"
)

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

// listResp is the structure that represents the public ip list and page detail.
type listResp struct {
	// The list of the public ips.
	PublicIPs []PublicIp `json:"publicips"`
	// The page information.
	PageInfo pageInfo `json:"page_info"`
}

// pageInfo is the structure that represents the page information.
type pageInfo struct {
	// The next marker information.
	NextMarker string `json:"next_marker"`
}

type PublicIPPage struct {
	pagination.MarkerPageBase
}

func ExtractPublicIPs(r pagination.Page) ([]PublicIp, error) {
	var s listResp
	err := r.(PublicIPPage).ExtractInto(&s)
	return s.PublicIPs, err
}

// IsEmpty checks whether a NetworkPage struct is empty.
func (r PublicIPPage) IsEmpty() (bool, error) {
	s, err := ExtractPublicIPs(r)
	return len(s) == 0, err
}

// LastMarker method returns the last public ip ID in a public ip page.
func (p PublicIPPage) LastMarker() (string, error) {
	tasks, err := ExtractPublicIPs(p)
	if err != nil {
		return "", err
	}
	if len(tasks) == 0 {
		return "", nil
	}
	return tasks[len(tasks)-1].ID, nil
}
