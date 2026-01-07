---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_batch_update_whiteblackip_rules"
description: |-
  Manages a resource to batch update white-black IP rules within HuaweiCloud WAF.
---

# huaweicloud_waf_batch_update_whiteblackip_rules

Manages a resource to batch update white-black IP rules within HuaweiCloud WAF.

-> All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be used.

-> This resource is a one-time action resource for batch updating white-black IP rules. Deleting this resource
   will not remove the created rules, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
resource "huaweicloud_waf_batch_update_whiteblackip_rules" "test" {
  name        = "test"
  white       = 1
  addr        = "127.0.0.1"
  description = "test description"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the white-black IP rule.

* `white` - (Required, Int, NonUpdatable) Specifies the protection action.
  The value can be:
  + `0`: block
  + `1`: allow
  + `2`: log only

* `policy_rule_ids` - (Required, List, NonUpdatable) Specifies an array of policies and rule IDs to associate protection
  policies with their corresponding rule sets.
  The [policy_rule_ids](#Rule_policy_rule_ids) structure is documented below.

* `addr` - (Optional, String, NonUpdatable) Specifies the IP address or CIDR block.
  This parameter and `ip_group_id` are alternative.

* `ip_group_id` - (Optional, String, NonUpdatable) Specifies the IP address group ID.
  This parameter and `addr` are alternative.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the rule.

* `time_mode` - (Optional, String, NonUpdatable) Specifies the time mode.
  The value can be:
  + **permanent**: The rule takes effect immediately (default).
  + **customize**: The rule takes effect during the specified time period.

* `start` - (Optional, Int, NonUpdatable) Specifies the start timestamp.
  This field is valid only when `time_mode` is set to **customize**.

* `terminal` - (Optional, Int, NonUpdatable) Specifies the end timestamp.
  This field is valid only when `time_mode` is set to **customize**.

<a name="Rule_policy_rule_ids"></a>
The `policy_rule_ids` block supports:

* `policy_id` - (Required, String, NonUpdatable) Specifies the policy ID.

* `rule_ids` - (Required, List, NonUpdatable) Specifies the rule IDs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
