---
subcategory: "Intelligent EdgeCloud (IEC)"
---

# huaweicloud\_iec\_security\_group

Manages a security group resource within HuaweiCloud IEC.
This is an alternative to `huaweicloud_iec_security_group_v1`

## Example Usage

```hcl
resource "huaweicloud_iec_security_group" "secgroup_test" {
  name        = "my_secgroup_test"
  description = "My security group"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) A unique name for the security group.

* `description` - (Optional, String) Description for the security group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a iec security group resource ID in UUID format.

* `security_group_rules` - An Array of one or more security group rules. The security_group_rules object structure is documented below.

The `security_group_rules` block supports:

* `id` - The id of the iec security group rules.
* `security_group_id` - The id of the iec security group rules.
* `description` - The description for the iec security group rules.
* `direction` - The direction of the iec security group rules.
* `ethertype` - The layer 3 protocol type.
* `port_range_max` - The higher part of the allowed port range.
* `port_range_min` - The lower part of the allowed port range.
* `protocol` - The layer 4 protocol type.
* `remote_ip_prefix` - The remote CIDR of the iec security group rules.
* `remote_group_id` - The remote group id of the iec security group rules.

## Default Security Group Rules

In most cases, HuaweiCloud will create some ingress and egress rules for each new iec security group. These iec security group rules will not be managed by Terraform.


## Timeouts

This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `delete` - Default is 10 minute.

## Import

Security Groups can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_iec_security_group.secgroup_test 2a02d1d3-437c-11eb-b721-fa163e8ac569
```
