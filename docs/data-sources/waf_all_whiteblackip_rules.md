---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_all_whiteblackip_rules"
description: |-
  Use this data source to get list of the WAF blacklist and whitelist protective rules under all policies.
---

# huaweicloud_waf_all_whiteblackip_rules

Use this data source to get list of the WAF blacklist and whitelist protective rules under all policies.

## Example Usage

```hcl
data "huaweicloud_waf_all_whiteblackip_rules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policyids` - (Optional, String) Specifies the ID of the policy.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  If you want to query resources under all enterprise projects, set this parameter to **all_granted_eps**.
  Defaults to **0**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The list of the WAF blacklist and whitelist protective rules.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `id` - The ID of the rule.

* `name` - The name of the rule.

* `policyid` - The ID of the policy.

* `timestamp` - The creation time of the rule, in milliseconds.

* `description` - The description of the rule.

* `status` - The status of the rule.
  The valid values are as follows:
  + `0`: Disabled.
  + `1`: Enabled.

* `addr` - The IP address or IP segment.

* `white` - The protection action.
  The valid values are as follows:
  + `0`: Block.
  + `1`: Allow.
  + `2`: Only record.

* `ip_group` - The IP address group.

  The [ip_group](#items_ip_group_struct) structure is documented below.

* `time_mode` - The effective mode. The default value is **permanent** (takes effect immediately).

* `start` - The start time when the rule takes effect.
  This parameter is valid only when the effective mode is customized.

* `terminal` - The end time when the rule takes effect.
  This parameter is valid only when the effective mode is customized.

<a name="items_ip_group_struct"></a>
The `ip_group` block supports:

* `id` - The ID of the IP address group.

* `name` - The name of the IP address group.

* `size` - The number of IP addresses or IP segments contained in the IP address group.
