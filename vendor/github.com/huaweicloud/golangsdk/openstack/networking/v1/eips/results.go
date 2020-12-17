package eips

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

//ApplyResult is a struct which represents the result of apply public ip
type ApplyResult struct {
	golangsdk.Result
}

func (r ApplyResult) Extract() (PublicIp, error) {
	var ip struct {
		Ip PublicIp `json:"publicip"`
	}
	err := r.Result.ExtractInto(&ip)
	return ip.Ip, err
}

//PublicIp is a struct that represents a public ip
type PublicIp struct {
	ID                  string `json:"id"`
	Status              string `json:"status"`
	Type                string `json:"type"`
	PublicAddress       string `json:"public_ip_address"`
	PrivateAddress      string `json:"private_ip_address"`
	PortID              string `json:"port_id"`
	TenantID            string `json:"tenant_id"`
	CreateTime          string `json:"create_time"`
	BandwidthID         string `json:"bandwidth_id"`
	BandwidthSize       int    `json:"bandwidth_size"`
	BandwidthShareType  string `json:"bandwidth_share_type"`
	EnterpriseProjectID string `json:"enterprise_project_id"`
}

//GetResult is a return struct of get method
type GetResult struct {
	golangsdk.Result
}

func (r GetResult) Extract() (PublicIp, error) {
	var Ip struct {
		Ip PublicIp `json:"publicip"`
	}
	err := r.Result.ExtractInto(&Ip)
	return Ip.Ip, err
}

//DeleteResult is a struct of delete result
type DeleteResult struct {
	golangsdk.ErrResult
}

//UpdateResult is a struct which contains the result of update method
type UpdateResult struct {
	golangsdk.Result
}

func (r UpdateResult) Extract() (PublicIp, error) {
	var ip PublicIp
	err := r.Result.ExtractIntoStructPtr(&ip, "publicip")
	return ip, err
}

type PublicIPPage struct {
	pagination.LinkedPageBase
}

func (r PublicIPPage) NextPageURL() (string, error) {
	publicIps, err := ExtractPublicIPs(r)
	if err != nil {
		return "", err
	}
	return r.WrapNextPageURL(publicIps[len(publicIps)-1].ID)
}

func ExtractPublicIPs(r pagination.Page) ([]PublicIp, error) {
	var s struct {
		PublicIPs []PublicIp `json:"publicips"`
	}
	err := r.(PublicIPPage).ExtractInto(&s)
	return s.PublicIPs, err
}

// IsEmpty checks whether a NetworkPage struct is empty.
func (r PublicIPPage) IsEmpty() (bool, error) {
	s, err := ExtractPublicIPs(r)
	return len(s) == 0, err
}
