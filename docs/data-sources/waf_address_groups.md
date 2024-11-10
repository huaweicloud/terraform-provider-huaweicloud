---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_address_groups"
description: |-
  Use this data source to get a list of WAF address groups.
---

# huaweicloud_waf_address_groups

Use this data source to get a list of WAF address groups.

## Example Usage

```hcl
variable enterprise_project_id {}

data "huaweicloud_waf_address_groups" "test" {
  name                  = "test-name"
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the address group. Fuzzy search is supported.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  For enterprise users, if omitted, default enterprise project will be used.

* `ip_address` - (Optional, String) Specifies the IP address or IP address ranges. If this parameter is specified,
  the address group that contains the specified IP address or IP address ranges are queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - The list of WAF address groups.

The `groups` block supports:

* `id` - The WAF address group ID.

* `name` - The name of the address group.

* `ip_addresses` - The IP addresses or IP address ranges.

* `description` - The description of the address group.

* `share_count` - The total number of the users share the address group.

* `accept_count` - The number of the users accept the sharing.

* `process_status` - The status of processing.

* `rules` - The list of rules that use the IP address group.

The `rules` block supports:

* `rule_id` - The ID of rule.

* `rule_name` - The name of rule.

* `policy_id` - The ID of policy.

* `policy_name` - The name of policy.
