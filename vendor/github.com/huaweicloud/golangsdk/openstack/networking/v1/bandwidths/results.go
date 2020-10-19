package bandwidths

import (
	"github.com/huaweicloud/golangsdk"
)

//BandWidth is a struct that represents a bandwidth
type BandWidth struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Size      int    `json:"size"`
	ShareType string `json:"share_type"`
	//PublicIPInfo  string `json:"publicip_info"`
	TenantID      string `json:"tenant_id"`
	BandwidthType string `json:"bandwidth_type"`
	ChargeMode    string `json:"charge_mode"`

	PublicipInfo []PublicIpinfo `json:"publicip_info"`

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

//GetResult is a return struct of get method
type GetResult struct {
	golangsdk.Result
}

func (r GetResult) Extract() (BandWidth, error) {
	var BW struct {
		BW BandWidth `json:"bandwidth"`
	}
	err := r.Result.ExtractInto(&BW)
	return BW.BW, err
}

//UpdateResult is a struct which contains the result of update method
type UpdateResult struct {
	golangsdk.Result
}

func (r UpdateResult) Extract() (BandWidth, error) {
	var bw BandWidth
	err := r.Result.ExtractIntoStructPtr(&bw, "bandwidth")
	return bw, err
}

//ListResult is a struct which contains the result of list method
type ListResult struct {
	golangsdk.Result
}

func (r ListResult) Extract() ([]BandWidth, error) {
	var s []BandWidth
	err := r.ExtractIntoSlicePtr(&s, "bandwidths")
	return s, err
}
