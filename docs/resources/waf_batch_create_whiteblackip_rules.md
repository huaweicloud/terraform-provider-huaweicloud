---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_batch_create_whiteblackip_rules"
description: |-
  Manages a resource to batch create whitelist/blacklist IP rules within HuaweiCloud WAF.
---

# huaweicloud_waf_batch_create_whiteblackip_rules

Manages a resource to batch create whitelist/blacklist IP rules within HuaweiCloud WAF.

-> All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be used.

-> This resource is a one-time action resource for batch creating whitelist/blacklist IP rules. Deleting this resource
   will not remove the created rules, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "enterprise_project_id" {}
variable "addr" {}
variable "policy_ids" {
  type = list(string)
}

resource "huaweicloud_waf_batch_create_whiteblackip_rules" "test" {
  name                  = "test_rule"
  white                 = 1
  policy_ids            = var.policy_ids
  addr                  = var.addr
  description           = "test description"
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the rule name.
  The name can contain only letters, digits, hyphens (-), underscores (_), and periods (.),
  and cannot exceed `64` characters.

* `white` - (Required, Int, NonUpdatable) Specifies the protection action.
  The value can be:
  + `0`: block
  + `1`: allow
  + `2`: log only

* `policy_ids` - (Required, List, NonUpdatable) Specifies the list of policy IDs to which the rule will be applied.

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

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  Value `0` indicates the default enterprise project.
  Value **all_granted_eps** indicates all enterprise projects to which the user has been granted access.
  Defaults to `0`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
