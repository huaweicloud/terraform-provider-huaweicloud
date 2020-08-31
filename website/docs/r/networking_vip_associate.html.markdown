---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_networking_vip_associate"
sidebar_current: "docs-huaweicloud-resource-networking-vip-associate"
description: |-
  Manages a Vip associate resource within HuaweiCloud.
---

# huaweicloud\_networking\_vip_associate

Manages a Vip associate resource within HuaweiCloud.
This is an alternative to `huaweicloud_networking_vip_associate_v2`

## Example Usage

```hcl
data "huaweicloud_vpc_subnet" "mynet" {
  name = "subnet-default"
}

resource "huaweicloud_networking_vip" "myvip" {
  network_id = data.huaweicloud_vpc_subnet.mynet.id
  subnet_id  = data.huaweicloud_vpc_subnet.mynet.subnet_id
}

resource "huaweicloud_networking_vip_associate" "vip_associated" {
  vip_id   = huaweicloud_networking_vip.myvip.id
  port_ids = [var.port_1, var.port_2]
}
```

## Argument Reference

The following arguments are supported:

* `vip_id` - (Required) The ID of vip to attach the port to.
    Changing this creates a new vip associate.

* `port_ids` - (Required) An array of one or more IDs of the ports to attach the vip to.
    Changing this creates a new vip associate.

## Attributes Reference

The following attributes are exported:

* `vip_id` - See Argument Reference above.
* `port_ids` - See Argument Reference above.
* `vip_subnet_id` - The ID of the subnet this vip connects to.
* `vip_ip_address` - The IP address in the subnet for this vip.
