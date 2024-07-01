---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_sub_network_interfaces"
description: |-
  Use this data source to get a list of VPC supplementary network interfaces.
---

# huaweicloud_vpc_sub_network_interfaces

Use this data source to get a list of VPC supplementary network interfaces.

## Example Usage

```hcl
variable "subnet_id" {}
variable "parent_id" {}

data "huaweicloud_vpc_sub_network_interfaces" "test" {
  subnet_id = var.subnet_id
  parent_id = var.parent_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `interface_id` - (Optional, String) Specifies the ID of the supplementary network interface.

* `vpc_id` - (Optional, String) Specifies the ID of the VPC to which the supplementary network interface belongs.

* `subnet_id` - (Optional, String) Specifies the ID of the subnet to which the supplementary network interface belongs.

* `parent_id` - (Optional, String) Specifies the ID of the elastic network interface
  to which the supplementary network interface belongs.

* `ip_address` - (Optional, String) Specifies the private IPv4 address of the supplementary network interface.

* `mac_address` - (Optional, String) Specifies the MAC address of the supplementary network interface.

* `description` - (Optional, List) Specifies the description of the supplementary network interface.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `sub_network_interfaces` - The list of supplementary network interfaces.

  The [sub_network_interfaces](#sub_network_interfaces_struct) structure is documented below.

<a name="sub_network_interfaces_struct"></a>
The `sub_network_interfaces` block supports:

* `id` - The ID of supplementary network interface.

* `vpc_id` - The ID of the VPC to which the supplementary network interface belongs.

* `subnet_id` - The ID of the subnet to which the supplementary network interface belongs.

* `parent_id` - The ID of the elastic network interface to which the supplementary network interface belongs.

* `parent_device_id` - The ID of the parent device.

* `security_groups` - The list of the security groups IDs to which the supplementary network interface belongs.

* `vlan_id` - The vlan ID of the supplementary network interface.

* `ip_address` - The private IPv4 address of the supplementary network interface.

* `mac_address` - The MAC address of the supplementary network interface.

* `ipv6_ip_address` - The IPv6 address of the supplementary network interface.

* `description` - The description of the supplementary network interface.

* `security_enabled` - Whether the IPv6 address is it enabled of the supplementary network interface.

* `project_id` - The ID of the project to which the supplementary network interface belongs.

* `tags` - The tags of a supplementary network interface.

* `created_at` - The time when the supplementary network interface is created.
