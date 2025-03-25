package loadbalancers

import (
	"github.com/chnsz/golangsdk"
)

// LoadBalancer is the primary load balancing configuration object that
// specifies the virtual IP address on which client traffic is received, as well
// as other details such as the load balancing method to be use, protocol, etc.
type LoadBalancer struct {
	// The unique ID for the LoadBalancer.
	ID string `json:"id"`

	// Human-readable description for the Loadbalancer.
	Description string `json:"description"`

	// The provisioning status of the LoadBalancer.
	// This value is ACTIVE, PENDING_CREATE or ERROR.
	ProvisioningStatus string `json:"provisioning_status"`

	// The administrative state of the Loadbalancer.
	// A valid value is true (UP) or false (DOWN).
	AdminStateUp bool `json:"admin_state_up"`

	// The name of the provider.
	Provider string `json:"provider"`

	// Pools are the pools related to this Loadbalancer.
	Pools []PoolRef `json:"pools"`

	// Listeners are the listeners related to this Loadbalancer.
	Listeners []ListenerRef `json:"listeners"`

	// The operating status of the LoadBalancer. This value is ONLINE or OFFLINE.
	OperatingStatus string `json:"operating_status"`

	// The IP address of the Loadbalancer.
	VipAddress string `json:"vip_address"`

	// The UUID of the subnet on which to allocate the virtual IP for the
	// Loadbalancer address.
	VipSubnetCidrID string `json:"vip_subnet_cidr_id"`

	// Human-readable name for the LoadBalancer. Does not have to be unique.
	Name string `json:"name"`

	LoadBalancerType string `json:"loadbalancer_type"`

	// Owner of the LoadBalancer.
	ProjectID string `json:"project_id"`

	// The UUID of the port associated with the IP address.
	VipPortID string `json:"vip_port_id"`

	// The UUID of a flavor if set.
	Tags []Tag `json:"tags"`

	// Guaranteed.
	Guaranteed bool `json:"guaranteed"`

	// The VPC ID.
	VpcID string `json:"vpc_id"`

	// EIP Info.
	Eips []EipInfo `json:"eips"`

	// Ipv6 Vip Address.
	Ipv6VipAddress string `json:"ipv6_vip_address"`

	// Ipv6 Vip Virsubnet ID.
	Ipv6VipVirsubnetID string `json:"ipv6_vip_virsubnet_id"`

	// Ipv6 Vip Port ID.
	Ipv6VipPortID string `json:"ipv6_vip_port_id"`

	// Availability Zone List.
	AvailabilityZoneList []string `json:"availability_zone_list"`

	// The UUID of the enterprise project who owns the Loadbalancer.
	EnterpriseProjectID string `json:"enterprise_project_id"`

	// Billing Info.
	BillingInfo string `json:"billing_info"`

	// L4 Flavor ID.
	L4FlavorID string `json:"l4_flavor_id"`

	// L4 Scale Flavor ID.
	L4ScaleFlavorID string `json:"l4_scale_flavor_id"`

	// L7 Flavor ID.
	L7FlavorID string `json:"l7_flavor_id"`

	// L7 Scale Flavor ID.
	L7ScaleFlavorID string `json:"l7_scale_flavor_id"`

	// Gateway flavor ID.
	GwFlavorId string `json:"gw_flavor_id"`

	// Public IP Info.
	PublicIps []PublicIpInfo `json:"publicips"`

	// Elb Virsubnet IDs.
	ElbVirsubnetIDs []string `json:"elb_virsubnet_ids"`

	// Elb Virsubnet Type.
	ElbVirsubnetType string `json:"elb_virsubnet_type"`

	// Ip Target Enable.
	IpTargetEnable bool `json:"ip_target_enable"`

	// Frozen Scene.
	FrozenScene string `json:"frozen_scene"`

	// Ipv6 Bandwidth.
	IPV6Bandwidth BandwidthRef `json:"ipv6_bandwidth"`

	// Update protection status
	ProtectionStatus string `json:"protection_status"`

	// Update protection reason
	ProtectionReason string `json:"protection_reason"`

	// Deletion Protection Enable.
	DeletionProtectionEnable bool `json:"deletion_protection_enable"`

	// Autoscaling configuration
	AutoScaling AutoScaling `json:"autoscaling"`

	// Waf failure action
	WafFailureAction string `json:"waf_failure_action"`

	// Charge Mode
	ChargeMode string `json:"charge_mode"`

	// Public Border Group
	PublicBorderGroup string `json:"public_border_group"`

	// Creation time
	CreatedAt string `json:"created_at"`

	// Update time
	UpdatedAt string `json:"updated_at"`
}

// EipInfo
type EipInfo struct {
	// Eip ID
	EipID string `json:"eip_id"`
	// Eip Address
	EipAddress string `json:"eip_address"`
	// Eip Address
	IpVersion int `json:"ip_version"`
}

// PoolRef
type PoolRef struct {
	ID string `json:"id"`
}

// ListenerRef
type ListenerRef struct {
	ID string `json:"id"`
}

// PublicIPInfo
type PublicIpInfo struct {
	// Public IP ID
	PublicIpID string `json:"publicip_id"`
	// Public IP Address
	PublicIpAddress string `json:"publicip_address"`
	// IP Version
	IpVersion int `json:"ip_version"`
}

// StatusTree represents the status of a loadbalancer.
type StatusTree struct {
	Loadbalancer *LoadBalancer `json:"loadbalancer"`
}

// Prepaid response
type PrepaidResponse struct {
	LoadBalancerID string `json:"loadbalancer_id"`
	OrderID        string `json:"order_id"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a loadbalancer.
func (r commonResult) Extract() (*LoadBalancer, error) {
	var s struct {
		LoadBalancer *LoadBalancer `json:"loadbalancer"`
	}
	err := r.ExtractInto(&s)
	return s.LoadBalancer, err
}

// Extract is a function that accepts a result and extracts a loadbalancer.
func (r commonResult) ExtractPrepaid() (*PrepaidResponse, error) {
	var s PrepaidResponse
	err := r.ExtractInto(&s)
	return &s, err
}

// GetStatusesResult represents the result of a GetStatuses operation.
// Call its Extract method to interpret it as a StatusTree.
type GetStatusesResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts the status of
// a Loadbalancer.
func (r GetStatusesResult) Extract() (*StatusTree, error) {
	var s struct {
		Statuses *StatusTree `json:"statuses"`
	}
	err := r.ExtractInto(&s)
	return s.Statuses, err
}

// ChangeResult represents the result of a ChangeChargingMode operation.
// Call its Extract method to get the order ID.
type ChangeResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts the order ID.
func (r ChangeResult) Extract() (string, error) {
	var s struct {
		LoadBalancerIdList []string `json:"loadbalancer_id_list"`
		EipIdList          []string `json:"eip_id_list"`
		OrderId            string   `json:"order_id"`
	}
	err := r.ExtractInto(&s)
	return s.OrderId, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a LoadBalancer.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a LoadBalancer.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a LoadBalancer.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
