---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_networking_vip_v2"
sidebar_current: "docs-huaweicloud-resource-networking-vip-v2"
description: |-
  Manages a V2 vip resource within HuaweiCloud.
---

# huaweicloud\_networking\_vip_v2

Manages a V2 vip resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_router_interface_v2" "router_interface_1" {
  router_id = "${huaweicloud_networking_router_v2.router_1.id}"
  subnet_id = "${huaweicloud_networking_subnet_v2.subnet_1.id}"
}

resource "huaweicloud_networking_router_v2" "router_1" {
  name = "router_1"
  external_network_id = "0a2228f2-7f8a-45f1-8e09-9039e1d09975"
}

resource "huaweicloud_networking_vip_v2" "vip_1" {
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  subnet_id = "${huaweicloud_networking_subnet_v2.subnet_1.id}"
}
```

## Argument Reference

The following arguments are supported:

* `network_id` - (Required) The ID of the network to attach the vip to.
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
