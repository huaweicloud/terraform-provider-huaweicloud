---
subcategory: "Intelligent EdgeCloud (IEC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iec_vip"
description: ""
---

# huaweicloud_iec_vip

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

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of the network to which the vip belongs.
  Changing this parameter creates a new vip resource.

* `ip_address` - (Optional, String, ForceNew) Specifies the IP address desired in the subnet for this vip.
  If you don't specify it, an available IP address from the specified subnet will be allocated to this vip.
  Changing this parameter creates a new vip resource.

* `port_ids` - (Optional, List) Specifies an array of IDs of the ports to attach the vip to.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the vip.

* `mac_address` - The MAC address of the vip.

* `allowed_addresses` - An array of IP addresses of the ports to attach the vip to.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

IEC VIP can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_iec_vip.vip_test 61fd8d31-8f92-4526-a5f5-07ec303e69e7
```
