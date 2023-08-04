package bandwidths

import (
	"github.com/chnsz/golangsdk"
)

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

type UpdateResult struct {
	golangsdk.Result
}

type BandWidthWithOrder struct {
	BandWidth BandWidth `json:"bandwidth"`
	OrderID   string    `json:"order_id"`
}

func (r UpdateResult) Extract() (BandWidthWithOrder, error) {
	var s BandWidthWithOrder
	err := r.ExtractInto(&s)
	return s, err
}

type ChangeResult struct {
	golangsdk.Result
}

func (r ChangeResult) Extract() (string, error) {
	var s struct {
		BandwidthIDs []string `json:"bandwidth_ids"`
		OrderID      string   `json:"order_id"`
		RequestID    string   `json:"request_id"`
	}
	err := r.ExtractInto(&s)
	return s.OrderID, err
}

type DeleteResult struct {
	golangsdk.ErrResult
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
