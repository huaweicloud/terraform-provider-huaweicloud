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
