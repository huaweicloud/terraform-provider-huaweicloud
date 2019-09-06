package bandwidths

import (
	"github.com/huaweicloud/golangsdk"
)

type PrePaid struct {
	OrderID string `json:"order_id"`
}

type PostPaid struct {
	Name                string         `json:"name"`
	Size                int            `json:"size"`
	ID                  string         `json:"id"`
	ShareType           string         `json:"share_type"`
	ChargeMode          string         `json:"charge_mode"`
	BandwidthType       string         `json:"bandwidth_type"`
	TenantID            string         `json:"tenant_id"`
	PublicipInfo        []PublicipInfo `json:"publicip_info"`
	EnterpriseProjectID string         `json:"enterprise_project_id"`
	//BillingInfo         string         `json:"billing_info"`
}

type PublicipInfo struct {
	PublicipID        string `json:"publicip_id"`
	PublicIPAddress   string `json:"publicip_address"`
	Publicipv6Address string `json:"publicipv_6_address"`
	IpVersion         int    `json:"ip_version"`
	PublicipType      string `json:"publicip_type"`
}

type UpdateResult struct {
	golangsdk.Result
}

func (r UpdateResult) ExtractOrderID() (PrePaid, error) {
	var s PrePaid
	err := r.ExtractInto(&s)
	return s, err
}

func (r UpdateResult) Extract() (PostPaid, error) {
	var s struct {
		Bandwidth PostPaid `json:"bandwidth"`
	}
	err := r.ExtractInto(&s)
	return s.Bandwidth, err
}

type CreateResult struct {
	golangsdk.Result
}

func (r CreateResult) Extract() (*BandWidth, error) {
	var entity BandWidth
	err := r.ExtractIntoStructPtr(&entity, "bandwidth")
	return &entity, err
}

type BatchCreateResult struct {
	golangsdk.Result
}

func (r BatchCreateResult) Extract() (*[]BandWidth, error) {
	var entity []BandWidth
	err := r.ExtractIntoSlicePtr(&entity, "bandwidths")
	return &entity, err
}

type BandWidth struct {
	// Specifies the bandwidth name. The value is a string of 1 to 64
	// characters that can contain letters, digits, underscores (_), and hyphens (-).
	Name string `json:"name"`

	// Specifies the bandwidth size. The value ranges from 1 Mbit/s to
	// 300 Mbit/s.
	Size int `json:"size"`

	// Specifies the bandwidth ID, which uniquely identifies the
	// bandwidth.
	ID string `json:"id"`

	// Specifies whether the bandwidth is shared or exclusive. The
	// value can be PER or WHOLE.
	ShareType string `json:"share_type"`

	// Specifies the elastic IP address of the bandwidth.  The
	// bandwidth, whose type is set to WHOLE, supports up to 20 elastic IP addresses. The
	// bandwidth, whose type is set to PER, supports only one elastic IP address.
	PublicipInfo []PublicIpinfo `json:"publicip_info"`

	// Specifies the tenant ID of the user.
	TenantId string `json:"tenant_id"`

	// Specifies the bandwidth type.
	BandwidthType string `json:"bandwidth_type"`

	// Specifies the charging mode (by traffic or by bandwidth).
	ChargeMode string `json:"charge_mode"`

	// Specifies the billing information.
	BillingInfo string `json:"billing_info"`

	// Enterprise project id
	EnterpriseProjectID string `json:"enterprise_project_id"`

	// Status
	Status string `json:"status"`
}

type PublicIpinfo struct {
	// Specifies the tenant ID of the user.
	PublicipId string `json:"publicip_id"`

	// Specifies the elastic IP address.
	PublicipAddress string `json:"publicip_address"`

	// Specifies the elastic IP v6 address.
	Publicipv6Address string `json:"publicipv6_address"`

	// Specifies the elastic IP version.
	IPVersion int `json:"ip_version"`

	// Specifies the elastic IP address type. The value can be
	// 5_telcom, 5_union, or 5_bgp.
	PublicipType string `json:"publicip_type"`
}

type DeleteResult struct {
	golangsdk.ErrResult
}
