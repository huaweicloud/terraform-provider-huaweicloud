package natgateways

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// NatGateway is a struct that represents a nat gateway
type NatGateway struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	RouterID          string `json:"router_id"`
	InternalNetworkID string `json:"internal_network_id"`
	TenantID          string `json:"tenant_id"`
	Spec              string `json:"spec"`
	Status            string `json:"status"`
	AdminStateUp      bool   `json:"admin_state_up"`
}

// GetResult is a return struct of get method
type GetResult struct {
	golangsdk.Result
}

func (r GetResult) Extract() (NatGateway, error) {
	var NatGW NatGateway
	err := r.Result.ExtractIntoStructPtr(&NatGW, "nat_gateway")
	return NatGW, err
}

// CreateResult is a return struct of create method
type CreateResult struct {
	golangsdk.Result
}

func (r CreateResult) Extract() (NatGateway, error) {
	var NatGW NatGateway
	err := r.Result.ExtractIntoStructPtr(&NatGW, "nat_gateway")
	return NatGW, err
}

// UpdateResult is a return struct of update method
type UpdateResult struct {
	golangsdk.Result
}

func (r UpdateResult) Extract() (NatGateway, error) {
	var NatGW NatGateway
	err := r.Result.ExtractIntoStructPtr(&NatGW, "nat_gateway")
	return NatGW, err
}

// DeleteResult is a return struct of delete method
type DeleteResult struct {
	golangsdk.ErrResult
}

type NatGatewayPage struct {
	pagination.LinkedPageBase
}

func (r NatGatewayPage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"nat_gateways_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}

func (r NatGatewayPage) IsEmpty() (bool, error) {
	is, err := ExtractNatGateways(r)
	return len(is) == 0, err
}

func ExtractNatGateways(r pagination.Page) ([]NatGateway, error) {
	var s struct {
		NatGateways []NatGateway `json:"nat_gateways"`
	}
	err := (r.(NatGatewayPage)).ExtractInto(&s)
	return s.NatGateways, err
}
