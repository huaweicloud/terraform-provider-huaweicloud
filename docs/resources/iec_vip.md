---
subcategory: "Intelligent EdgeCloud (IEC)"
---

# huaweicloud\_iec\_vip

Manages a VIP resource within HuaweiCloud IEC.

## Example Usage

```hcl
variable "iec_subnet_id" {}

resource "huaweicloud_iec_vip" "vip_test" {
  subnet_id = var.iec_subnet_id
}
```

## Argument Reference

The following arguments are supported:

* `subnet_id` - (Required, String, ForceNew) Specifies the subnet network id 
    of vip binding in which to allocate IP address for this vip.
    Changing this parameter creates a new vip resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the vip.

* `mac_address` - The MAC address of the vip.

* `fixed_ips` - An array of IP addresses binding to the vip.

## Timeouts

This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `delete` - Default is 10 minute.

## Import

IEC VIP can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_iec_vip.vip_test 61fd8d31-8f92-4526-a5f5-07ec303e69e7
```
