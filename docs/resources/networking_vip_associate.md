---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_networking_vip_associate"
description: ""
---

# huaweicloud_networking_vip_associate

Using this resource, one or more NICs (to which the ECS instance belongs) can be bound to the VIP.

-> A VIP can only have one resource.

## Example Usage

```hcl
variable "vip_id" {}
variable "nic_port_ids" {
  type = list(string)
}

resource "huaweicloud_networking_vip_associate" "vip_associated" {
  vip_id   = var.vip_id
  port_ids = var.nic_port_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the vip associate resource. If omitted, the
  provider-level region will be used.

* `vip_id` - (Required, String, ForceNew) The ID of vip to attach the ports to.

* `port_ids` - (Required, List) An array of one or more IDs of the ports to attach the vip to.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.
* `vip_subnet_id` - The ID of the subnet this vip connects to.
* `vip_ip_address` - The IP address in the subnet for this vip.
* `ip_addresses` - The IP addresses of ports to attach the vip to.

## Import

Vip associate can be imported using the `vip_id` and port IDs separated by slashes (no limit on the number of
port IDs), e.g.

```bash
$ terraform import huaweicloud_networking_vip_associate.vip_associated vip_id/port1_id/port2_id
```
