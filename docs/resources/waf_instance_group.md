---
subcategory: "Web Application Firewall (WAF)"
---

# huaweicloud_waf_instance_group

Manages WAF instance groups within HuaweiCloud. The groups are used to bind the ELB instance to the ELB mode WAF.

## Example Usage

```hcl
variable "vpc_id" {}

resource "huaweicloud_waf_instance_group" "group_1" {
  name   = "example_name"
  vpc_id = var.vpc_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the instance group.
  If omitted, the provider-level region will be used. Changing this setting will create a new resource.

* `name` - (Required, String) Specifies the instance group name.
  The maximum length is 64 characters. Only letters, digits and underscores (_) are allowed.

* `vpc_id` - (Required, String, ForceNew) Specifies the id of the VPC that the WAF dedicated instances belongs to.

* `description` - (Optional, String) Specifies the description of the instance group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The rule ID in UUID format.

* `body_limit` - The body limit of the forwarding policy.

* `header_limit` - The header limit of the forwarding policy.

* `connection_timeout` - The time for connection timeout in the forwarding policy.

* `write_timeout` - The time for writing timeout in the forwarding policy.

* `read_timeout` - The time for reading timeout in the forwarding policy.

* `load_balancers` - The IDs of the ELB instances that has been bound to the instance group.

## Import

The instance group can be imported using the ID, e.g.:

```sh
terraform import huaweicloud_waf_instance_group.group_1 0be1e69d-1987-4d9c-9dc5-fc7eed592398
```
