---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_eip_associate"
description: ""
---

# huaweicloud_vpc_eip_associate

Associates an EIP to a specified IP address or port.

## Example Usage

### Associate with a fixed IP

```hcl
variable "public_address" {}
variable "network_id" {}

resource "huaweicloud_vpc_eip_associate" "associated" {
  public_ip  = var.public_address
  network_id = var.network_id
  fixed_ip   = "192.168.0.100"
}
```

### Associate with a port

```hcl
variable "network_id" {}

data "huaweicloud_networking_port" "myport" {
  network_id = var.network_id
  fixed_ip   = "192.168.0.100"
}

resource "huaweicloud_vpc_eip" "myeip" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_vpc_eip_associate" "associated" {
  public_ip = huaweicloud_vpc_eip.myeip.address
  port_id   = data.huaweicloud_networking_port.myport.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to associate the EIP. If omitted, the provider-level
  region will be used. Changing this creates a new resource.

* `public_ip` - (Required, String, ForceNew) Specifies the EIP address to associate. Changing this creates a new resource.

* `fixed_ip` - (Optional, String, ForceNew) Specifies a private IP address to associate with the EIP.
  Changing this creates a new resource.

* `network_id` - (Optional, String, ForceNew) Specifies the ID of the network to which the **fixed_ip** belongs.
  It is mandatory when `fixed_ip` is set. Changing this creates a new resource.

* `port_id` - (Optional, String, ForceNew) Specifies an existing port ID to associate with the EIP.
  This parameter and `fixed_ip` are alternative. Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
* `mac_address` - The MAC address of the private IP.
* `status` - The status of EIP, should be **BOUND**.
* `public_ipv6` - The IPv6 address of the private IP.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minute.
* `delete` - Default is 5 minute.

## Import

EIP associations can be imported using the `id` of the EIP, e.g.

```bash
$ terraform import huaweicloud_vpc_eip_associate.eip 2c7f39f3-702b-48d1-940c-b50384177ee1
```
