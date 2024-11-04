---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_p2c_gateway_connections"
description: |-
  Use this data source to get the list of P2C VPN gateway connections.
---

# huaweicloud_vpn_p2c_gateway_connections

Use this data source to get the list of P2C VPN gateway connections.

## Example Usage

```hcl
variable "gateway_id" {}

data "huaweicloud_vpn_p2c_gateway_connections" "test" {
  p2c_gateway_id = var.gateway_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `p2c_gateway_id` - (Required, String) Specifies the ID of a P2C VPN gateway instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `connections` - The connection list.

  The [connections](#connections_struct) structure is documented below.

<a name="connections_struct"></a>
The `connections` block supports:

* `connection_id` - The connection ID.

* `client_ip` - The IP address of a client.

* `client_user_name` - The username of a client.

* `inbound_packets` - The number of inbound packets.

* `inbound_bytes` - The number of inbound bytes.

* `outbound_packets` - The number of outbound packets.

* `outbound_bytes` - The number of outbound bytes.

* `client_virtual_ip` - The virtual IP address of a client.

* `connection_established_time` - The time when a connection is established.

* `timestamp` - The timestamp.
