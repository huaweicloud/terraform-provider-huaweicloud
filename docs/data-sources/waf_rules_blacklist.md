---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rules_blacklist"
description: |-
  Use this data source to get a list of blacklist and whitelist rules.
---

# huaweicloud_waf_rules_blacklist

Use this data source to get a list of blacklist and whitelist rules.

## Example Usage

```hcl
variable "policy_id" {}
variable "rule_id" {}

data "huaweicloud_waf_rules_blacklist" "test" {
  policy_id = var.policy_id
  rule_id   = var.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policy_id` - (Required, String) Specifies the ID of the policy to which the blacklist and whitelist rules belong.

* `rule_id` - (Optional, String) Specifies the ID of the blacklist or whitelist rule.

* `name` - (Optional, String) Specifies the name of the blacklist or whitelist rule.

* `status` - (Optional, String) Specifies the status of the blacklist or whitelist rule.
  The valid values are as follows:
  + **0**: The blacklist and whitelist rule is active.
  + **1**: The blacklist and whitelist rule is disabled.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the protection policies belong.
  If omitted, will query the blacklist and whitelist rules under the default enterprise project for enterprise users.

* `action` - (Optional, String) Specifies the protective action of the blacklist and whitelist rule.
  The valid values are as follows:
  + **0**: Intercept the request.
  + **1**: Release the request.
  + **2**: Record the request only.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of the blacklist and whitelist rules.

  The [rules](#rules_struct) structure is documented below.

<a name="rules_struct"></a>
The `rules` block supports:

* `id` - The ID of the blacklist or whitelist rule.

* `policy_id` - The ID of the policy to which the blacklist and whitelist rule belongs.

* `name` - The name of the blacklist or whitelist rule.

* `status` - The status of the blacklist or whitelist rule.

* `description` - The description of the blacklist or whitelist rule.

* `action` - The protective action of the blacklist and whitelist rule.

* `ip_address` - The IP address included in the blacklist and whitelist rule.

* `address_group` - The IP address group included in the blacklist and whitelist rule.

  The [address_group](#rules_address_group_struct) structure is documented below.

* `created_at` - The creation time of the blacklist and whitelist rule.

<a name="rules_address_group_struct"></a>
The `address_group` block supports:

* `id` - The ID of the IP address group.

* `name` - The name of the IP address group.

* `size` - The number of IP addresses or IP address ranges in the IP address group.
