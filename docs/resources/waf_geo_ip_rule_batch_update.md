---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_geo_ip_rule_batch_update"
description: |-
  Manages a resource to batch update the geo IP rules within HuaweiCloud.
---

# huaweicloud_waf_geo_ip_rule_batch_update

Manages a resource to batch update the geo IP rules within HuaweiCloud.

-> All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be used.

-> This resource is only a one-time action resource using to batch update geo IP rules. Deleting this resource will
not change the current geo IP rule configurations, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "policy_id" {}
variable "rule_ids"  {
  type = list(string)
}

resource "huaweicloud_waf_geo_ip_rule_batch_update" "test" {
  geoip = "US|CA|JP"

  policy_rule_ids {
    policy_id = var.policy_id
    rule_ids  = var.rule_ids
  }

  status = 1
  name   = "updated_geo_rule"
  white  = 2
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `geoip` - (Required, String, NonUpdatable) Specifies the geo location codes blocked by the rule, separated by '|'.

* `policy_rule_ids` - (Optional, List, NonUpdatable) Specifies the policy and rule details.
  The [policy_rule_ids](#geo_ip_policy_rules) structure is documented below.

* `status` - (Optional, Int, NonUpdatable) Specifies the status of the geo IP rule.
  The valid values are as follows:
  + `0`: Disabled
  + `1`: Enabled

* `name` - (Optional, String, NonUpdatable) Specifies the name of the geo IP rule.

* `white` - (Optional, Int, NonUpdatable) Specifies the protection action.
  The valid values are as follows:
  + `1`: Allow
  + `2`: Block

<a name="geo_ip_policy_rules"></a>
The `policy_rule_ids` block supports:

* `policy_id` - (Optional, String, NonUpdatable) Specifies the policy ID to which the geo IP rule belongs.

* `rule_ids` - (Optional, List, NonUpdatable) Specifies the ID list of the geo IP rule.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
