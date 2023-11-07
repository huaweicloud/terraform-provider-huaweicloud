package deprecated

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/fwaas_v2/firewall_groups"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/fwaas_v2/policies"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/fwaas_v2/routerinsertion"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/fwaas_v2/rules"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/layer3/routers"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/vpnaas/endpointgroups"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/vpnaas/ikepolicies"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/vpnaas/services"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/vpnaas/siteconnections"
	"github.com/chnsz/golangsdk/openstack/networking/v2/networks"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"
	"github.com/chnsz/golangsdk/openstack/networking/v2/subnets"
)

// NetworkCreateOpts represents the attributes used when creating a new network.
type NetworkCreateOpts struct {
	networks.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// ToNetworkCreateMap casts a CreateOpts struct to a map.
// It overrides networks.ToNetworkCreateMap to add the ValueSpecs field.
func (opts NetworkCreateOpts) ToNetworkCreateMap() (map[string]interface{}, error) {
	return BuildRequest(opts, "network")
}

// SubnetCreateOpts represents the attributes used when creating a new subnet.
type SubnetCreateOpts struct {
	subnets.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// ToSubnetCreateMap casts a CreateOpts struct to a map.
// It overrides subnets.ToSubnetCreateMap to add the ValueSpecs field.
func (opts SubnetCreateOpts) ToSubnetCreateMap() (map[string]interface{}, error) {
	b, err := BuildRequest(opts, "subnet")
	if err != nil {
		return nil, err
	}

	if m := b["subnet"].(map[string]interface{}); m["gateway_ip"] == "" {
		m["gateway_ip"] = nil
	}

	return b, nil
}

// FloatingIPCreateOpts represents the attributes used when creating a new floating ip.
type FloatingIPCreateOpts struct {
	floatingips.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// ToFloatingIPCreateMap casts a CreateOpts struct to a map.
// It overrides floatingips.ToFloatingIPCreateMap to add the ValueSpecs field.
func (opts FloatingIPCreateOpts) ToFloatingIPCreateMap() (map[string]interface{}, error) {
	return BuildRequest(opts, "floatingip")
}

// RouterCreateOpts represents the attributes used when creating a new router.
type RouterCreateOpts struct {
	routers.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// ToRouterCreateMap casts a CreateOpts struct to a map.
// It overrides routers.ToRouterCreateMap to add the ValueSpecs field.
func (opts RouterCreateOpts) ToRouterCreateMap() (map[string]interface{}, error) {
	return BuildRequest(opts, "router")
}

// PortCreateOpts represents the attributes used when creating a new port.
type PortCreateOpts struct {
	ports.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// ToPortCreateMap casts a CreateOpts struct to a map.
// It overrides ports.ToPortCreateMap to add the ValueSpecs field.
func (opts PortCreateOpts) ToPortCreateMap() (map[string]interface{}, error) {
	return BuildRequest(opts, "port")
}

// VpnIPSecPolicyCreateOpts represents the attributes used when creating a new IPSec policy.
type VpnIPSecPolicyCreateOpts struct {
	ipsecpolicies.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// VpnServiceCreateOpts represents the attributes used when creating a new VPN service.
type VpnServiceCreateOpts struct {
	services.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// VpnEndpointGroupCreateOpts represents the attributes used when creating a new endpoint group.
type VpnEndpointGroupCreateOpts struct {
	endpointgroups.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// VpnIKEPolicyCreateOpts represents the attributes used when creating a new IKE policy.
type VpnIKEPolicyCreateOpts struct {
	ikepolicies.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// VpnIKEPolicyLifetimeCreateOpts represents the attributes used when creating a new lifetime for an IKE policy.
type VpnIKEPolicyLifetimeCreateOpts struct {
	ikepolicies.LifetimeCreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// VpnSiteConnectionCreateOpts represents the attributes used when creating a new IPSec site connection.
type VpnSiteConnectionCreateOpts struct {
	siteconnections.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// BuildRequest takes an opts struct and builds a request body for
// golangsdk to execute
func BuildRequest(opts interface{}, parent string) (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	b = AddValueSpecs(b)

	return map[string]interface{}{parent: b}, nil
}

// AddValueSpecs expands the 'value_specs' object and removes 'value_specs'
// from the reqeust body.
func AddValueSpecs(body map[string]interface{}) map[string]interface{} {
	if body["value_specs"] != nil {
		for k, v := range body["value_specs"].(map[string]interface{}) {
			body[k] = v
		}
		delete(body, "value_specs")
	}

	return body
}

// MapValueSpecs converts ResourceData into a map
func MapValueSpecs(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("value_specs").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

// MapResourceProp converts ResourceData property into a map
func MapResourceProp(d *schema.ResourceData, prop string) map[string]interface{} {
	m := make(map[string]interface{})
	for key, val := range d.Get(prop).(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

// FirewallGroup is an HuaweiCloud firewall group.
type FirewallGroup struct {
	firewall_groups.FirewallGroup
	routerinsertion.FirewallGroupExt
}

// FirewallGroupCreateOpts represents the attributes used when creating a new firewall.
type FirewallGroupCreateOpts struct {
	firewall_groups.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// ToFirewallCreateMap casts a FirewallGroupCreateOpts struct to a map.
// It overrides firewalls.ToFirewallCreateMap to add the ValueSpecs field.
func (opts FirewallGroupCreateOpts) ToFirewallCreateMap() (map[string]interface{}, error) {
	return BuildRequest(opts, "firewall_group")
}

// FirewallGroupUpdateOpts represents the attributes used when updating a firewall
type FirewallGroupUpdateOpts struct {
	firewall_groups.UpdateOptsBuilder
}

// ToFirewallUpdateMap casts a FirewallGroupUpdateOpts struct to a map.
func (opts FirewallGroupUpdateOpts) ToFirewallUpdateMap() (map[string]interface{}, error) {
	return BuildRequest(opts, "firewall")
}

// PolicyCreateOpts represents the attributes used when creating a new firewall policy.
type PolicyCreateOpts struct {
	policies.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// ToFirewallPolicyCreateMap casts a PolicyCreateOpts struct to a map.
// It overrides policies.ToFirewallPolicyCreateMap to add the ValueSpecs field.
func (opts PolicyCreateOpts) ToFirewallPolicyCreateMap() (map[string]interface{}, error) {
	return BuildRequest(opts, "firewall_policy")
}

// RuleCreateOpts represents the attributes used when creating a new firewall rule.
type RuleCreateOpts struct {
	rules.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// ToRuleCreateMap casts a CreateOpts struct to a map.
// It overrides rules.ToRuleCreateMap to add the ValueSpecs field.
func (opts RuleCreateOpts) ToRuleCreateMap() (map[string]interface{}, error) {
	b, err := BuildRequest(opts, "firewall_rule")
	if err != nil {
		return nil, err
	}

	if m := b["firewall_rule"].(map[string]interface{}); m["protocol"] == "any" {
		m["protocol"] = nil
	}

	return b, nil
}
