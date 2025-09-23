---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_network_interface"
description: ""
---

# huaweicloud_vpc_network_interface

Manages a VPC network interface resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "subnet_id" {}

resource "huaweicloud_vpc_network_interface" "test" {
  name      = var.name
  subnet_id = var.subnet_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of the subnet to which the network interface belongs.
  Changing this creates a new resource.

* `fixed_ip_v4` - (Optional, String, ForceNew) Specifies the network interface IPv4 address.
  Changing this creates a new resource.

* `name` - (Optional, String) Specifies the network interface name.

* `security_group_ids` - (Optional, List) Specifies an array of one or more security group IDs.

* `allowed_addresses` - (Optional, List) Specifies an array of IP addresses that can be active on the
  network interface. If the IP address is "1.1.1.1/0", it means that the source/destination address
  check switch is turned off.

* `dhcp_lease_time` - (Optional, String) Specifies the DHCP lease time. The value format of value is "Xh",
  the value of "X" is "-1" or from "1" to "30000". If the value is "-1", the DHCP lease time is infinite.

* `tags` - (Optional, Map) Specifies the network interface tags in the format of key-value pairs.
  This parameter can only be used in **cn-south-2** for now.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `availability_zone` - Indicates the availability zone to which the network interface belongs.

* `device_id` - Indicates the ID of the device to which the network interface belongs.

* `device_owner` - Indicates the belonged device, which can be the DHCP server, router, load balancer, or Nova.

* `dns_name` - Indicates the default private network DNS name of the primary NIC.

* `enable_efi` - Indicates whether to enable EFI.

* `instance_id` - Indicates the ID of the instance to which the network interface belongs.

* `instance_type` - Indicates the type of the instance to which the network interface belongs.

* `ipv6_bandwidth_id` - Indicates the Shared bandwidth ID bound to IPv6 network interface.

* `mac_address` - Indicates the network interface MAC address.

* `port_security_enabled` - Indicates whether the security option is enabled for the network interface.
  If the option is not enabled, the security group and DHCP snooping do not take effect.

* `status` - Indicates the network interface status.

## Import

The network interface can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_vpc_network_interface.test <id>
```
