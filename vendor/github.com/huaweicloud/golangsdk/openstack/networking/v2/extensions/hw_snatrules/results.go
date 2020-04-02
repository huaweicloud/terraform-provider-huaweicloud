package hw_snatrules

import (
	"github.com/huaweicloud/golangsdk"
)

// SnatRule is a struct that represents a snat rule
type SnatRule struct {
	ID                string `json:"id"`
	NatGatewayID      string `json:"nat_gateway_id"`
	NetworkID         string `json:"network_id"`
	TenantID          string `json:"tenant_id"`
	FloatingIPID      string `json:"floating_ip_id"`
	FloatingIPAddress string `json:"floating_ip_address"`
	Status            string `json:"status"`
	AdminStateUp      bool   `json:"admin_state_up"`
	Cidr              string `json:"cidr"`
	SourceType        int    `json:"source_type"`
}

// GetResult is a return struct of get method
type GetResult struct {
	golangsdk.Result
}

func (r GetResult) Extract() (SnatRule, error) {
	var SR SnatRule
	err := r.Result.ExtractIntoStructPtr(&SR, "snat_rule")
	return SR, err
}

// CreateResult is a return struct of create method
type CreateResult struct {
	golangsdk.Result
}

func (r CreateResult) Extract() (SnatRule, error) {
	var SR SnatRule
	err := r.Result.ExtractIntoStructPtr(&SR, "snat_rule")
	return SR, err
}

// DeleteResult is a return struct of delete method
type DeleteResult struct {
	golangsdk.ErrResult
}
