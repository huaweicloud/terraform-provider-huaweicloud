---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_networking_vip"
sidebar_current: "docs-huaweicloud-resource-networking-vip"
description: |-
  Manages a Vip resource within HuaweiCloud.
---

# huaweicloud\_networking\_vip

Manages a Vip resource within HuaweiCloud.
This is an alternative to `huaweicloud_networking_vip`

## Example Usage

```hcl
data "huaweicloud_vpc_subnet" "mynet" {
  name = "subnet-default"
}

resource "huaweicloud_networking_vip" "myvip" {
  network_id = data.huaweicloud_vpc_subnet.mynet.id
  subnet_id  = data.huaweicloud_vpc_subnet.mynet.subnet_id
}
```

## Argument Reference

The following arguments are supported:

* `network_id` - (Required) The Network ID of the VPC subnet to attach the vip to.
    Changing this creates a new vip.

* `subnet_id` - (Required) Subnet in which to allocate IP address for this vip.
    Changing this creates a new vip.

* `ip_address` - (Optional) IP address desired in the subnet for this vip.
    If you don't specify `ip_address`, an available IP address from
    the specified subnet will be allocated to this vip.

* `name` - (Optional) A unique name for the vip.

## Attributes Reference

The following attributes are exported:

* `network_id` - See Argument Reference above.
* `subnet_id` - See Argument Reference above.
* `ip_address` - See Argument Reference above.
* `name` - See Argument Reference above.
* `status` - The status of vip.
* `id` - The ID of the vip.
* `tenant_id` - The tenant ID of the vip.
* `device_owner` - The device owner of the vip.
