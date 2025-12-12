---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_virtual_interface"
description: |-
  Manages a virtual interface resource within HuaweiCloud.
---

# huaweicloud_dc_virtual_interface

Manages a virtual interface resource within HuaweiCloud.

## Example Usage

### Create a DC virtual interface with VGW service type

```hcl
variable "direct_connect_id" {}
variable "vgw_id" {}
variable "interface_name" {}

resource "huaweicloud_dc_virtual_interface" "test" {
  direct_connect_id = var.direct_connect_id
  vgw_id            = var.vgw_id
  name              = var.interface_name
  type              = "private"
  route_mode        = "static"
  vlan              = 522
  bandwidth         = 5

  remote_ep_group = [
    "1.1.1.0/30",
  ]

  address_family       = "ipv4"
  local_gateway_v4_ip  = "1.1.1.1/30"
  remote_gateway_v4_ip = "1.1.1.2/30"
  enable_bfd           = true
  
  
  extend_attribute   {
    ha_type           = "bfd"
    ha_mode           = "auto_single"
    min_rx_interval   = 1000
    min_tx_interval   = 1000
    detect_multiplier = 3
  }
}
```

### Create a DC virtual interface with GDGW service type

```hcl
variable "direct_connect_id" {}
variable "interface_name" {}
variable "gateway_id" {}

resource "huaweicloud_dc_virtual_interface" "test" {
  direct_connect_id = var.direct_connect_id
  name              = var.interface_name
  type              = "private"
  route_mode        = "static"
  vlan              = 76
  bandwidth         = 5
  service_type      = "GDGW"
  gateway_id        = var.gateway_id

  remote_ep_group = [
    "1.1.1.0/30",
  ]

  address_family       = "ipv4"
  local_gateway_v4_ip  = "1.1.1.1/30"
  remote_gateway_v4_ip = "1.1.1.2/30"
  enable_bfd           = true
  
  extend_attribute   {
    ha_type           = "bfd"
    ha_mode           = "auto_single"
    min_rx_interval   = 1000
    min_tx_interval   = 1000
    detect_multiplier = 3
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the virtual interface is located.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `direct_connect_id` - (Required, String, ForceNew) Specifies the ID of the direct connection associated with the
  virtual interface.
  Changing this will create a new resource.

* `name` - (Required, String) Specifies the name of the virtual interface.
  The valid length is limited from `1` to `64`, only chinese and english letters, digits, hyphens (-), underscores (_)
  and dots (.) are allowed.
  The Chinese characters must be in **UTF-8** or **Unicode** format.

* `type` - (Required, String, ForceNew) Specifies the type of the virtual interface.
  The valid value is **private**.
  Changing this will create a new resource.

* `route_mode` - (Required, String, ForceNew) Specifies the route mode of the virtual interface.
  The valid values are **static** and **bgp**.
  Changing this will create a new resource.

* `vlan` - (Required, Int, ForceNew) Specifies the customer VLAN to be connected.
  The valid value is range from `0` to `3,999`.
  Changing this will create a new resource.

* `bandwidth` - (Required, Int) Specifies the bandwidth of the virtual interface. The unit is Mbit/s.
  The size range depends on the direct connection.

* `remote_ep_group` - (Required, List) Specifies the CIDR list of remote subnets.
  A CIDR that contains CIDRs of local subnet (corresponding to the parameter `local_gateway_v4_ip` or
  `local_gateway_v6_ip`) and remote subnet (corresponding to the parameter `remote_gateway_v4_ip` or
  `remote_gateway_v6_ip`) must exist in the list.

* `priority` - (Optional, String) The priority of a virtual interface. The value can be **normal** or **low**.
  If the priorities are the same, the virtual interfaces work in load balancing mode.
  If the priorities are different, the virtual interfaces work in active/standby pairs.
  Outbound traffic is preferentially forwarded to the normal virtual interface with a higher priority.
  This option is only supported by virtual interfaces that use BGP routing.

* `service_ep_group` - (Optional, List) Specifies the subnets that access Internet services through a connection.
  This field is required in public network connections.

* `description` - (Optional, String) Specifies the description of the virtual interface.
  The description contain a maximum of `128` characters and the angle brackets (< and >) are not allowed.
  Chinese characters must be in **UTF-8** or **Unicode** format.

* `service_type` - (Optional, String, ForceNew) Specifies the service type of the virtual interface.
  The valid values are **VGW**, **GDGW** and **LGW**. The default value is **VGW**.
  Changing this will create a new resource.

* `vgw_id` - (Optional, String, ForceNew) Specifies the ID of the virtual gateway to which the virtual interface is
  connected.

* `gateway_id` - (Optional, String, ForceNew) Specifies the ID of the gateway associated with the virtual
  interface (the ID of the global DC gateway).
  This field is required when `service_type` is set to **GDGW**.

  Changing this will create a new resource.

* `local_gateway_v4_ip` - (Optional, String, ForceNew) Specifies the IPv4 address of the virtual interface in cloud
  side. This parameter is mandatory if `address_family` is set to **ipv4**.
  Changing this will create a new resource.

* `remote_gateway_v4_ip` - (Optional, String, ForceNew) Specifies the IPv4 address of the virtual interface in client
  side. This parameter is mandatory if `address_family` is set to **ipv4**.
  Changing this will create a new resource.

* `address_family` - (Optional, String, ForceNew) Specifies the service type of the virtual interface.
  The valid values are **ipv4** and **ipv6**. Defaults to **ipv4**.
  Changing this will create a new resource.

* `local_gateway_v6_ip` - (Optional, String, ForceNew) Specifies the IPv6 address of the virtual interface in cloud
  side. This parameter is mandatory if `address_family` is set to **ipv6**.
  Changing this will create a new resource.

* `remote_gateway_v6_ip` - (Optional, String, ForceNew) Specifies the IPv6 address of the virtual interface in client
  side. This parameter is mandatory if `address_family` is set to **ipv6**.
  Changing this will create a new resource.

-> The CIDRs of `local_gateway_v4_ip` and `remote_gateway_v4_ip` (or `local_gateway_v6_ip` and `remote_gateway_v6_ip`)
  must be in the same subnet.

* `asn` - (Optional, Int, ForceNew) Specifies the local BGP ASN of the virtual interface.
  The valid value is range from `1` to `4,294,967,295`, except `64,512`.
  Changing this will create a new resource.

* `bgp_md5` - (Optional, String, ForceNew) Specifies the (MD5) password for the local BGP.
  Changing this will create a new resource.

* `enable_bfd` - (Optional, Bool) Specifies whether to enable the Bidirectional Forwarding Detection (BFD) function.
  Defaults to `false`.

* `enable_nqa` - (Optional, Bool) Specifies whether to enable the Network Quality Analysis (NQA) function.
  Defaults to `false`.

-> The values of parameter `enable_bfd` and `enable_nqa` cannot be `true` at the same time.

* `lag_id` - (Optional, String, ForceNew) Specifies the ID of the link aggregation group (LAG) associated with the
  virtual interface.
  Changing this will create a new resource.

* `resource_tenant_id` - (Optional, String, ForceNew) Specifies the project ID of another tenant in the same region
  which is used to create virtual interface across tenant. After the across tenant virtual interface is successfully
  created, the target tenant needs to accept the virtual interface request for the virtual interface to take effect.
  Changing this will create a new resource.

  -> 1. When `resource_tenant_id` is specified, `vgw_id` must be the target tenant virtual gateway id.
  <br/>2. When `resource_tenant_id` is specified, the tags can only be configured after the target tenant accepts the
  virtual interface request and the virtual interface takes effect.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the virtual
  interface belongs. This field is valid only when `resource_tenant_id` is not specified.
  Changing this will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the virtual interface.

* `extend_attribute` - (Optional, List) The extended parameter information.
  The [extend_attribute](#DCVirtualInterface_extend_attribute) structure is documented below.

<a name="DCVirtualInterface_extend_attribute"></a>
The `extend_attribute` block supports:

* `ha_type` - (Optional, String) The availability detection type of the virtual interface. The value can be **bfd** or **nqa**

* `ha_mode` - (Optional, String) The availability detection mode.

  Valid values are:
  + **auto_single**: auto single bfd detection.
  + **auto_multi**: auto multi bfd detection.
  + **static_single**: static single bfd detection.
  + **static_multi**: static multi bfd detection.
  + **enhance_nqa**: enhance nqa detection.

* `detect_multiplier` - (Optional, Int) The number of detection retries.
  The valid value is range from `3` to `20`. The default value is `3`.

* `min_rx_interval` - (Optional, Int) The interval for receiving detection packets.
  The valid value is range from `200` to `5000`. The default value is `1000`.

* `min_tx_interval` - (Optional, Int) The interval for sending detection packets.
  The valid value is range from `200` to `5000`. The default value is `1000`.

* `remote_disclaim` - (Optional, Int) The remote identifier of the static BFD session.
  The valid value is range from `1` to `16384`.

* `local_disclaim` - (Optional, Int) The local identifier of the static BFD session.
  The valid value is range from `1` to `16384`.

* `ipv6_remote_disclaim` - (Optional, Int) The remote identifier of the static IPv6 BFD session.
  The valid value is range from `1` to `16384`.

* `ipv6_local_disclaim` - (Optional, Int) The local identifier of the static IPv6 BFD session.
  The valid value is range from `1` to `16384`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the virtual interface.

* `device_id` - The attributed device ID.

* `status` - The current status of the virtual interface.

* `bgp_route_limit` - The BGP route configuration.

* `ies_id` - The edge site ID.

* `lgw_id` - The ID of the local gateway, which is used in IES scenarios.

* `rate_limit` - Whether rate limiting is enabled on a virtual interface.

* `reason` - The error information if the status of a line is Error.

* `route_limit` - The remote subnet route configurations of the virtual interface.

* `created_at` - The creation time of the virtual interface.

* `updated_at` - The latest update time of the virtual interface.

* `vif_peers` - The peer information of the virtual interface.
  The [vif_peers](#DCVirtualInterface_vif_peers) structure is documented below.

* `extend_attribute` - The extended parameter information.
  The [extend_attribute](#DCVirtualInterface_extend_attribute) structure is documented below.

<a name="DCVirtualInterface_vif_peers"></a>
The `vif_peers` block supports:

* `id` - The VIF peer resource ID.

* `name` - The name of the virtual interface peer.

* `description` - The description of the virtual interface peer.

* `address_family` - The address family type of the virtual interface, which can be **IPv4** or **IPv6**.

* `local_gateway_ip` - The address of the virtual interface peer used on the cloud.

* `remote_gateway_ip` - The address of the virtual interface peer used in the on-premises data center.

* `route_mode` - The routing mode, which can be **static** or **bgp**.

* `bgp_asn` - The ASN of the BGP peer.

* `bgp_md5` - The MD5 password of the BGP peer.

* `device_id` - The ID of the device that the virtual interface peer belongs to.

* `enable_bfd` - Whether to enable BFD.

* `enable_nqa` - Whether to enable NQA.

* `bgp_route_limit` - The BGP route configuration.

* `bgp_status` - The BGP protocol status of the virtual interface peer. If the virtual interface peer uses **static**
  routing, the status is null.

* `status` - The status of the virtual interface peer.

* `vif_id` - The ID of the virtual interface corresponding to the virtual interface peer.

* `receive_route_num` - The number of received BGP routes if **bgp** routing is used. If **static** routing is used,
  this parameter is meaningless and the value is **-1**.

* `remote_ep_group` - The remote subnet list, which records the CIDR blocks used in the on-premises data center.

* `service_ep_group` - The list of public network addresses that can be accessed by the on-premises data center.

<a name="DCVirtualInterface_extend_attribute"></a>
The `extend_attribute` block supports:

* `ha_type` - The availability detection type of the virtual interface.

* `ha_mode` - The availability detection mode.

* `detect_multiplier` - The number of detection retries.

* `min_rx_interval` - The interval for receiving detection packets.

* `min_tx_interval` - The interval for sending detection packets.

* `remote_disclaim` - The remote identifier of the static BFD session.

* `local_disclaim` - The local identifier of the static BFD session.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

Virtual interfaces can be imported using their `id`, e.g.

```shell
$ terraform import huaweicloud_dc_virtual_interface.test 5bb22e82-5b07-4845-bd1b-b064eca92e0a
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include:
`resource_tenant_id`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```bash
resource "huaweicloud_dc_virtual_interface" "test" {
    ...

  lifecycle {
    ignore_changes = [
      resource_tenant_id,
    ]
  }
}
```
