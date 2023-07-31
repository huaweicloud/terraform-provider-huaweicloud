---
subcategory: "Web Application Firewall (WAF)"
---

# huaweicloud_waf_instance_groups

Use this data source to get a list of WAF instance groups.

## Example Usage

```hcl
data "huaweicloud_waf_instance_groups" "groups_1" {
  name = "example_name"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to obtain the WAF instance groups.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) The name of WAF instance group used for matching.
  The value is not case-sensitive and supports fuzzy matching.

* `vpc_id` - (Optional, String) The id of the VPC that the WAF dedicated instances belongs to.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `groups` - A list of WAF instance groups.

The `groups` block supports:

* `region` - The region in which to create the instance group.

* `name` - The instance group name.

* `vpc_id` - The id of the VPC that the WAF dedicated instances belongs to.

* `description` - Description of the instance group.

* `body_limit` - The body limit of the forwarding policy.

* `header_limit` - The header limit of the forwarding policy.

* `connection_timeout` - The time for connection timeout in the forwarding policy.

* `write_timeout` - The time for writing timeout in the forwarding policy.

* `read_timeout` - The time for reading timeout in the forwarding policy.

* `load_balancers` - The IDs of the ELB instances that has been bound to the instance group.
