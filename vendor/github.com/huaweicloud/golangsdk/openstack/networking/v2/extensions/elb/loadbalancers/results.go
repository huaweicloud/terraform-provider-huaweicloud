package loadbalancers

import (
	"github.com/huaweicloud/golangsdk"
)

type LoadBalancer struct {
	VipAddress      string `json:"vip_address"`
	UpdateTime      string `json:"update_time"`
	CreateTime      string `json:"create_time"`
	ID              string `json:"id"`
	Status          string `json:"status"`
	BandWidth       int    `json:"bandwidth"`
	VpcID           string `json:"vpc_id"`
	AdminStateUp    int    `json:"admin_state_up"`
	VipSubnetID     string `json:"vip_subnet_id"`
	Type            string `json:"type"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	SecurityGroupID string `json:"security_group_id"`
}

type GetResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a router.
func (r GetResult) Extract() (*LoadBalancer, error) {
	s := &LoadBalancer{}
	return s, r.ExtractInto(s)
}
