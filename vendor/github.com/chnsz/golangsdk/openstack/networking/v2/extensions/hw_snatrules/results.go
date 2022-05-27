package hw_snatrules

import (
	"github.com/chnsz/golangsdk"
)

// SnatRule is a struct that represents a snat rule
type SnatRule struct {
	ID                string `json:"id"`
	NatGatewayID      string `json:"nat_gateway_id"`
	NetworkID         string `json:"network_id"`
	TenantID          string `json:"tenant_id"`
	FloatingIPID      string `json:"floating_ip_id"`
	FloatingIPAddress string `json:"floating_ip_address"`
	Description       string `json:"description"`
	Status            string `json:"status"`
	AdminStateUp      bool   `json:"admin_state_up"`
	Cidr              string `json:"cidr"`
	SourceType        int    `json:"source_type"`
}

type commonResult struct {
	golangsdk.Result
}

func (r commonResult) Extract() (SnatRule, error) {
	var sr SnatRule
	err := r.Result.ExtractIntoStructPtr(&sr, "snat_rule")
	return sr, err
}

// CreateResult is a return struct of create method
type CreateResult struct {
	commonResult
}

// UpdateResult is a return struct of update method
type UpdateResult struct {
	commonResult
}

// GetResult is a return struct of get method
type GetResult struct {
	commonResult
}

// DeleteResult is a return struct of delete method
type DeleteResult struct {
	golangsdk.ErrResult
}
