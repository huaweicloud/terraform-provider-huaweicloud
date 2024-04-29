---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_sub_network_interface"
description: ""
---

# huaweicloud_vpc_sub_network_interface

Manages a supplementary network interface resource within HuaweiCloud.

## Example Usage

```hcl
variable "subnet_id" {}
variable "parent_id" {}
variable "vlan_id" {}

resource "huaweicloud_vpc_sub_network_interface" "test" {
  subnet_id   = var.subnet_id
  parent_id   = var.parent_id
  vlan_id     = var.vlan_id
  description = "create a supplementary network interface"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of the subnet to which the supplementary network
  interface belongs.
  Changing this creates a new resource.

* `parent_id` - (Required, String, ForceNew) Specifies the ID of the elastic network interface to which the
  supplementary network interface belongs.  
  Changing this creates a new resource.

* `security_group_ids` - (Optional, List) Specifies the list of the security groups IDs to which the supplementary
  network interface belongs.

* `description` - (Optional, String) Specifies the description of the supplementary network interface.

* `vlan_id` - (Optional, String) Specifies the vlan ID of the supplementary network interface.
  The valid value is range from `1` t0 `4094`.

* `ip_address` - (Optional, String) Specifies the private IPv4 address of the supplementary network interface.

* `ipv6_enable` - (Optional, Bool) Specifies the IPv6 address is it enabled of the supplementary network interface.

* `ipv6_ip_address` - (Optional, String) Specifies the IPv6 address of the supplementary network interface.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `mac_address` - The MAC address of the supplementary network interface.

* `parent_device_id` - The ID of the ECS to which the supplementary network interface belongs.

* `vpc_id` - The ID of the VPC to which the supplementary network interface belongs.

* `status` - The status of the supplementary network interface.

* `created_at` - The create time of the supplementary network interface.

## Import

The supplementary network interface can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_vpc_sub_network_interface.test <id>
```
