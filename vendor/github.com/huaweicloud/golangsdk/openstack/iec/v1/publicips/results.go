package publicips

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/common"
	"github.com/huaweicloud/golangsdk/pagination"
)

type CommonPublicIP struct {
	// Specifies the ID of the elastic IP address, which uniquely
	// identifies the elastic IP address.
	ID string `json:"id"`

	// Specifies the status of the elastic IP address.
	Status string `json:"status"`

	// Specifies the obtained elastic IP address.
	PublicIpAddress string `json:"public_ip_address"`

	// Value range: 4, 6, respectively, to create ipv4 and ipv6, when not created ipv4 by default
	IPVersion int `json:"ip_version"`

	// Specifies the time for applying for the elastic IP address.
	CreateTime string `json:"create_time"`

	// Specifies the Siteid.
	SiteID string `json:"site_id"`

	// SiteInfo
	SiteInfo string `json:"site_info"`

	//Operator information
	Operator common.Operator `json:"operator"`
}

type PublicIPCreateResp struct {
	CommonPublicIP

	// Specifies the bandwidth size.
	BandwidthSize int `json:"bandwidth_size"`

	Region string `json:"region,omitempty"`
}

type PublicIPUpdateResp struct {
	CommonPublicIP

	// Specifies the private IP address bound to the elastic IP address.
	PrivateIpAddress string `json:"private_ip_address"`

	// Specifies the port ID.
	PortID string `json:"port_id"`

	// Specifies the bandwidth size.
	BandwidthSize int `json:"bandwidth_size"`

	// Specifies the bandwidth ID of the elastic IP address.
	BandwidthID string `json:"bandwidth_id"`

	// Specifies whether the bandwidth is shared or exclusive.
	BandwidthShareType string `json:"bandwidth_share_type"`

	// Specifies the bandwidth name.
	BandwidthName string `json:"bandwidth_name"`
}

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

func (r CreateResult) Extract() (*PublicIPCreateResp, error) {
	var entity PublicIPCreateResp
	err := r.ExtractIntoStructPtr(&entity, "publicip")
	return &entity, err
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type GetResult struct {
	commonResult
}

func (r GetResult) Extract() (*common.PublicIP, error) {
	var entity common.PublicIP
	err := r.ExtractIntoStructPtr(&entity, "publicip")
	return &entity, err
}

type PublicIPPage struct {
	pagination.LinkedPageBase
}

func ExtractPublicIPs(r pagination.Page) ([]common.PublicIP, error) {
	var s struct {
		PublicIPs []common.PublicIP `json:"publicips"`
	}
	err := r.(PublicIPPage).ExtractInto(&s)
	return s.PublicIPs, err
}

// IsEmpty checks whether a NetworkPage struct is empty.
func (r PublicIPPage) IsEmpty() (bool, error) {
	s, err := ExtractPublicIPs(r)
	return len(s) == 0, err
}

type UpdateResult struct {
	commonResult
}

func (r UpdateResult) Extract() (*PublicIPUpdateResp, error) {
	var entity PublicIPUpdateResp
	err := r.ExtractIntoStructPtr(&entity, "publicip")
	return &entity, err
}
