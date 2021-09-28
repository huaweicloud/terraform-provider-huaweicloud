package deprecated

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/vpnaas/endpointgroups"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/vpnaas/ikepolicies"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/vpnaas/services"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/vpnaas/siteconnections"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
