package huaweicloud

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dns/v2/recordsets"
	"github.com/chnsz/golangsdk/openstack/dns/v2/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

// RecordSetCreateOpts represents the attributes used when creating a new DNS record set.
type RecordSetCreateOpts struct {
	recordsets.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// ToRecordSetCreateMap casts a CreateOpts struct to a map.
// It overrides recordsets.ToRecordSetCreateMap to add the ValueSpecs field.
func (opts RecordSetCreateOpts) ToRecordSetCreateMap() (map[string]interface{}, error) {
	b, err := BuildRequest(opts, "")
	if err != nil {
		return nil, err
	}

	if m, ok := b[""].(map[string]interface{}); ok {
		return m, nil
	}

	return nil, fmtp.Errorf("Expected map but got %T", b[""])
}

// ZoneCreateOpts represents the attributes used when creating a new DNS zone.
type ZoneCreateOpts struct {
	zones.CreateOpts
	ValueSpecs map[string]interface{} `json:"value_specs,omitempty"`
}

// ToZoneCreateMap casts a CreateOpts struct to a map.
// It overrides zones.ToZoneCreateMap to add the ValueSpecs field.
func (opts ZoneCreateOpts) ToZoneCreateMap() (map[string]interface{}, error) {
	b, err := BuildRequest(opts, "")
	if err != nil {
		return nil, err
	}

	if m, ok := b[""].(map[string]interface{}); ok {
		if opts.TTL > 0 {
			m["ttl"] = opts.TTL
		}

		return m, nil
	}

	return nil, fmtp.Errorf("Expected map but got %T", b[""])
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
