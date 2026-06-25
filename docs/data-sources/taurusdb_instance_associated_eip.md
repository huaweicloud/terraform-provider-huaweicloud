---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_instance_associated_eip"
description: |-
  Use this data source to query the EIP information associated with a TaurusDB instance within HuaweiCloud.
---

# huaweicloud_taurusdb_instance_associated_eip

Use this data source to query the EIP information associated with a TaurusDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_taurusdb_instance_associated_eip" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the TaurusDB instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `can_enable_public_access` - Whether public access can be enabled.

* `eip_id` - The EIP ID.

* `type` - The network type of the EIP.
  Valid values are as follows:
  + **5_bgp**: dynamic BGP.
  + **5_sbgp**: static BGP.
  + **5_youxuanbgp**: premium BGP.

* `port_id` - The port ID.

* `public_ip_address` - The EIP address.

* `private_ip_address` - The private IP address bound to the EIP.

* `status` - The EIP status.
  Valid values are as follows:
  + **FREEZED**: The EIP is frozen.
  + **BIND_ERROR**: The EIP failed to be bound.
  + **BINDING**: The EIP is being bound.
  + **PENDING_DELETE**: The EIP is being released.
  + **PENDING_CREATE**: The EIP is being created.
  + **NOTIFYING**: The EIP is being created.
  + **NOTIFY_DELETE**: The EIP is being released.
  + **PENDING_UPDATE**: The EIP is being updated.
  + **DOWN**: The EIP has not been bound.
  + **ACTIVE**: The EIP has been bound.
  + **ELB**: The EIP has been bound to a load balancer.
  + **VPN**: The EIP has been bound to a VPN.
  + **ERROR**: The EIP is failed.

* `create_time` - The time when the EIP was assigned.

* `bandwidth_id` - The bandwidth ID.

* `bandwidth_name` - The bandwidth name.

* `bandwidth_size` - The bandwidth size in Mbit/s.

* `bandwidth_share_type` - The bandwidth type. Valid values are **PER** (dedicated bandwidth)
  and **WHOLE** (shared bandwidth).

* `profile` - The additional parameters, including the order ID and product ID.
  If profile is not empty, the EIP is billed on a yearly/monthly basis.

  The [profile](#profile_struct) structure is documented below.

<a name="profile_struct"></a>
The `profile` block supports:

* `order_id` - The order ID.

* `product_id` - The product ID.
