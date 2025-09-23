package loadbalancers

import (
	"github.com/chnsz/golangsdk"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToLoadBalancerCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	// Human-readable name for the Loadbalancer. Does not have to be unique.
	Name string `json:"name,omitempty"`

	// Human-readable type for the Loadbalancer.
	LoadBalancerType string `json:"loadbalancer_type,omitempty"`

	// Human-readable description for the Loadbalancer.
	Description string `json:"description,omitempty"`

	// The IP address of the Loadbalancer.
	VipAddress string `json:"vip_address,omitempty"`

	// The network on which to allocate the Loadbalancer's address.
	VipSubnetID string `json:"vip_subnet_cidr_id,omitempty"`

	// The V6 network on which to allocate the Loadbalancer's address.
	IpV6VipSubnetID string `json:"ipv6_vip_virsubnet_id,omitempty"`

	// The UUID of a l4 flavor.
	L4Flavor string `json:"l4_flavor_id,omitempty"`

	// Guaranteed.
	Guaranteed *bool `json:"guaranteed,omitempty"`

	// The VPC ID.
	VpcID string `json:"vpc_id,omitempty"`

	// Availability Zone List.
	AvailabilityZoneList []string `json:"availability_zone_list" required:"true"`

	// The UUID of the enterprise project who owns the Loadbalancer.
	EnterpriseProjectID string `json:"enterprise_project_id,omitempty"`

	// The tags of the Loadbalancer.
	Tags []Tag `json:"tags,omitempty"`

	// The administrative state of the Loadbalancer. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`

	// The UUID of a l7 flavor.
	L7Flavor string `json:"l7_flavor_id,omitempty"`

	// IPv6 Bandwidth.
	IPV6Bandwidth *BandwidthRef `json:"ipv6_bandwidth,omitempty"`

	// Public IP IDs.
	PublicIPIds []string `json:"publicip_ids,omitempty"`

	// Public IP.
	PublicIP *PublicIP `json:"publicip,omitempty"`

	// ELB VirSubnet IDs.
	ElbSubnetIds []string `json:"elb_virsubnet_ids,omitempty"`

	// IP Target Enable.
	IPTargetEnable *bool `json:"ip_target_enable,omitempty"`

	// Deletion Protection Enable.
	DeletionProtectionEnable *bool `json:"deletion_protection_enable,omitempty"`

	// Prepaid configuration
	PrepaidOpts *PrepaidOpts `json:"prepaid_options,omitempty"`

	// Autoscaling configuration
	AutoScaling *AutoScaling `json:"autoscaling,omitempty"`

	// Protection status
	ProtectionStatus string `json:"protection_status,omitempty"`

	// Protection reason
	ProtectionReason string `json:"protection_reason,omitempty"`

	// Waf failure action
	WafFailureAction string `json:"waf_failure_action,omitempty"`

	// IpV6 Vip Address
	Ipv6VipAddress string `json:"ipv6_vip_address,omitempty"`
}

// BandwidthRef
type BandwidthRef struct {
	// Share Bandwidth ID
	ID string `json:"id" required:"true"`
}

// UBandwidthRef
type UBandwidthRef struct {
	// Share Bandwidth ID
	ID *string `json:"id"`
}

// PublicIP
type PublicIP struct {
	// IP Version.
	IPVersion int `json:"ip_version,omitempty"`

	// Network Type
	NetworkType string `json:"network_type" required:"true"`

	// Billing Info.
	BillingInfo string `json:"billing_info,omitempty"`

	// Description.
	Description string `json:"description,omitempty"`

	// Bandwidth
	Bandwidth Bandwidth `json:"bandwidth" required:"true"`
}

// Bandwidth
type Bandwidth struct {
	// ID
	Id string `json:"id,omitempty"`

	// Name
	Name string `json:"name,omitempty"`

	// Size
	Size int `json:"size,omitempty"`

	// Charge Mode
	ChargeMode string `json:"charge_mode,omitempty"`

	// Share Type
	ShareType string `json:"share_type,omitempty"`

	// Billing Info.
	BillingInfo string `json:"billing_info,omitempty"`
}

// Tag
type Tag struct {
	// Tag Key
	Key string `json:"key,omitempty"`
	// Tag Value
	Value string `json:"value,omitempty"`
}

// Prepaid configuration
type PrepaidOpts struct {
	PeriodType string `json:"period_type,omitempty"`
	PeriodNum  int    `json:"period_num,omitempty"`
	AutoRenew  bool   `json:"auto_renew,omitempty"`
	AutoPay    bool   `json:"auto_pay,omitempty"`
}

// AutoScaling configuration
type AutoScaling struct {
	Enable      bool   `json:"enable"`
	MinL7Flavor string `json:"min_l7_flavor_id,omitempty"`
}

// ToLoadBalancerCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToLoadBalancerCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "loadbalancer")
}

// Create is an operation which provisions a new loadbalancer based on the
// configuration defined in the CreateOpts struct. Once the request is
// validated and progress has started on the provisioning process, a
// CreateResult will be returned.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToLoadBalancerCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}

// Get retrieves a particular Loadbalancer based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToLoadBalancerUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is the common options struct used in this package's Update
// operation.
type UpdateOpts struct {
	// Human-readable name for the Loadbalancer. Does not have to be unique.
	Name string `json:"name,omitempty"`

	// Human-readable description for the Loadbalancer.
	Description *string `json:"description,omitempty"`

	// The administrative state of the Loadbalancer. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`

	// The IP address of the Loadbalancer.
	VipAddress string `json:"vip_address,omitempty"`

	// The network on which to allocate the Loadbalancer's address.
	VipSubnetID *string `json:"vip_subnet_cidr_id"`

	// The V6 network on which to allocate the Loadbalancer's address.
	IpV6VipSubnetID *string `json:"ipv6_vip_virsubnet_id"`

	// The UUID of a l4 flavor.
	L4Flavor string `json:"l4_flavor_id,omitempty"`

	// The UUID of a l7 flavor.
	L7Flavor string `json:"l7_flavor_id,omitempty"`

	// Human-readable type for the Loadbalancer.
	LoadBalancerType string `json:"loadbalancer_type,omitempty"`

	// IPv6 Bandwidth.
	IPV6Bandwidth *UBandwidthRef `json:"ipv6_bandwidth,omitempty"`

	// ELB VirSubnet IDs.
	ElbSubnetIds []string `json:"elb_virsubnet_ids,omitempty"`

	// IP Target Enable.
	IPTargetEnable *bool `json:"ip_target_enable,omitempty"`

	// Deletion Protection Enable.
	DeletionProtectionEnable *bool `json:"deletion_protection_enable,omitempty"`

	// Prepaid configuration
	PrepaidOpts *PrepaidOpts `json:"prepaid_options,omitempty"`

	// Autoscaling configuration
	AutoScaling *AutoScaling `json:"autoscaling,omitempty"`

	// Update protection status
	ProtectionStatus string `json:"protection_status,omitempty"`

	// Update protection reason
	ProtectionReason *string `json:"protection_reason,omitempty"`

	// Waf failure action
	WafFailureAction string `json:"waf_failure_action,omitempty"`

	// IpV6 Vip Address
	Ipv6VipAddress string `json:"ipv6_vip_address,omitempty"`
}

// ToLoadBalancerUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToLoadBalancerUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "loadbalancer")
}

// Update is an operation which modifies the attributes of the specified
// LoadBalancer.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOpts) (r UpdateResult) {
	b, err := opts.ToLoadBalancerUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	// if loadbalancer_type is gateway, it indicates ipv4_subnet_id has not been changed, the value should be nil
	// so remove loadbalancer_type and ipv4_subnet_id from the request body
	if v, ok := b["loadbalancer"].(map[string]interface{})["loadbalancer_type"]; ok && v.(string) == "gateway" {
		delete(b["loadbalancer"].(map[string]interface{}), "vip_subnet_cidr_id")
		delete(b["loadbalancer"].(map[string]interface{}), "loadbalancer_type")
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}

// Delete will permanently delete a particular LoadBalancer based on its
// unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}

// ForceDelete will delete the LoadBalancer and the sub resource(LoadBalancer, listeners, unbind associated pools)
func ForceDelete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceForceDeleteURL(c, id), nil)
	return
}

// GetStatuses will return the status of a particular LoadBalancer.
func GetStatuses(c *golangsdk.ServiceClient, id string) (r GetStatusesResult) {
	_, r.Err = c.Get(statusRootURL(c, id), &r.Body, nil)
	return
}

type UpdateAvailabilityZone interface {
	ToAvailabilityZoneUpdateMap() (map[string]interface{}, error)
}

// Availability Zone List.
type AvailabilityZoneOpts struct {
	AvailabilityZoneList []string `json:"availability_zone_list" required:"true"`
}

// ToAvailabilityZoneUpdateMap builds a request body from AvailabilityZoneOpts.
func (opts AvailabilityZoneOpts) ToAvailabilityZoneUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// AddAvailabilityZone will add availability zone list
func AddAvailabilityZone(c *golangsdk.ServiceClient, id string, opts AvailabilityZoneOpts) (r UpdateResult) {
	b, err := opts.ToAvailabilityZoneUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(updateAvailabilityZoneURL(c, id, "batch-add"), b, &r.Body, &golangsdk.RequestOpts{})
	return
}

// RemoveAvailabilityZone will remove availability zone list
func RemoveAvailabilityZone(c *golangsdk.ServiceClient, id string, opts AvailabilityZoneOpts) (r UpdateResult) {
	b, err := opts.ToAvailabilityZoneUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(updateAvailabilityZoneURL(c, id, "batch-remove"), b, &r.Body, &golangsdk.RequestOpts{})
	return
}

// Charging info.
type ChangeChargingModeOpts struct {
	LoadBalancerIds []string       `json:"loadbalancer_ids" required:"true"`
	ChargingMode    string         `json:"charge_mode" required:"true"`
	PrepaidOptions  PrepaidOptions `json:"prepaid_options,omitempty"`
}

type PrepaidOptions struct {
	IncludePublicIp *bool  `json:"include_publicip,omitempty"`
	PeriodType      string `json:"period_type" required:"true"`
	PeriodNum       int    `json:"period_num,omitempty"`
	AutoRenew       string `json:"auto_renew,omitempty"`
	AutoPay         bool   `json:"auto_pay,omitempty"`
}

// ChangeChargingMode will change the charging mode of the loadbalancer
func ChangeChargingMode(c *golangsdk.ServiceClient, opts ChangeChargingModeOpts) (r ChangeResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(changeChargingModeURL(c), b, &r.Body, &golangsdk.RequestOpts{})
	return
}
