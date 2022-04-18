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
	// Name
	Name string `json:"name" required:"true"`

	// Size
	Size int `json:"size" required:"true"`

	// Charge Mode
	ChargeMode string `json:"charge_mode" required:"true"`

	// Share Type
	ShareType string `json:"share_type" required:"true"`

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

// GetStatuses will return the status of a particular LoadBalancer.
func GetStatuses(c *golangsdk.ServiceClient, id string) (r GetStatusesResult) {
	_, r.Err = c.Get(statusRootURL(c, id), &r.Body, nil)
	return
}
