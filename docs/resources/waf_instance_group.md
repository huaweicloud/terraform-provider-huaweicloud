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
  If omitted, the provider-level region will be used. Changing this setting will create a new certificate.

* `name` - (Required, String) Specifies the instance group name.
  The maximum length is 64 characters. Only letters, digits and underscores (_) are allowed.

* `vpc_id` - (Required, String, ForceNew) Specifies the id of the VPC that the WAF dedicated instances belongs to.

* `description` - (Optional, String) Specifies the description of the instance group.

-> **NOTE**: The following arguments can not be changed from their default values until the ELB mode WAF instances
are created. The arguments are:`body_limit`, `header_limit`, `connection_timeout`, `write_timeout`, `read_timeout` and
`load_balances`.

* `body_limit` - (Optional, Int) Specifies the body limit of the forwarding policy.
  The value ranges from 2,000 to 8,000. Defaults is 4,000.

* `header_limit` - (Optional, Int) Specifies the header limit of the forwarding policy.
  The value ranges from 4,000 to 20,000. Defaults is 8,000.

* `connection_timeout` - (Optional, Int) Specifies the time for connection timeout in the forwarding policy.
  The value ranges from 1 to 20. Defaults is 10.

* `write_timeout` - (Optional, Int) Specifies the time for writing timeout in the forwarding policy.
  The value ranges from 1 to 20. Defaults is 10.

* `read_timeout` - (Optional, Int) Specifies the time for reading timeout in the forwarding policy.
  The value ranges from 1 to 20. Defaults is 10.

* `load_balances` - (Optional, List) Specifies the IDs of the ELB instances that has been bound to the instance group.
  It is an array of ELB instance IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The rule ID in UUID format.

## Import

The instance group can be imported using the ID, e.g.:

```sh
terraform import huaweicloud_waf_instance_group.group_1 0be1e69d-1987-4d9c-9dc5-fc7eed592398
```
