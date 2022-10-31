package instances

import (
	"github.com/chnsz/golangsdk"
)

// CreateResult represents the result of a Create operation.
// Call its Extract method to get the response body.
type CreateResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts the response body.
func (r CreateResult) Extract() (OrderInfo, error) {
	var resp OrderInfo
	err := r.ExtractInto(&resp)
	return resp, err
}

// OrderInfo is the response body of a Create operation.
type OrderInfo struct {
	ID      string `json:"resourceId"`
	OrderID string `json:"orderId"`
}

// Instance represents a DataArts Studio instance.
type Instance struct {
	ID              string `json:"resource_id"`
	Name            string `json:"resource_name"`
	Type            string `json:"resource_type"`
	SpecCode        string `json:"resource_spec_code"`
	ProductID       string `json:"product_id"`
	VpcID           string `json:"vpc_id"`
	SubnetID        string `json:"net_id"`
	SecurityGroupID string `json:"security_group_id"`
	WorkspaceMode   string `json:"work_space_mode"`
	Status          int    `json:"status"`

	OrderID    string `json:"order_id"`
	OrderType  string `json:"order_type"`
	ChargeType string `json:"charge_type"`

	EffectiveTime int64  `json:"effective_time"`
	ExpireTime    int64  `json:"expire_time"`
	ExpireDays    string `json:"expire_days"`
	CreateUser    string `json:"create_user"`
	CreateTime    int64  `json:"create_time"`

	IsTrialOrder int `json:"is_trial_order"`
	IsAutoRenew  int `json:"is_auto_renew"`

	Region              string `json:"region_id"`
	AvailabilityZone    string `json:"availability_zone"`
	ProjectID           string `json:"project_id"`
	DomainID            string `json:"domain_id"`
	EnterpriseProjectID string `json:"eps_id"`
}
