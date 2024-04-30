---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_gateways"
description: ""
---

# huaweicloud_vpn_gateways

Use this data source to get the list of VPN gateways.

## Example Usage

```hcl
variable "gateway_name" {}

data "huaweicloud_vpn_gateways" "test" {
  name = var.gateway_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the gateway.

* `gateway_id` - (Optional, String) Specifies the ID of the gateway.

* `network_type` - (Optional, String) Specifies the network type of the gateway.
  The value can be: **public** and **private**.

* `attachment_type` - (Optional, String) Specifies the attachment type of the gateway.
  The value can be: **vpc** and **er**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the gateway.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `gateways` - The list of gateways.
  The [gateways](#Gateway_Gateways) structure is documented below.

<a name="Gateway_Gateways"></a>
The `gateways` block supports:

* `id` - The ID of the gateway

* `name` - The name of the gateway.

* `network_type` - The network type of the gateway.

* `status` - The status of the gateway.

* `attachment_type` - The attachment type.

* `vpc_id` - The ID of the VPC to which the VPN gateway is connected.

* `er_id` - The ID of the ER to which the VPN gateway is connected.

* `er_attachment_id` - The ER attachment ID.

* `local_subnets` - The local subnets.

* `connect_subnet` - The VPC network segment used by the VPN gateway.

* `bgp_asn` - The ASN number of BGP

* `flavor` - The flavor of the VPN gateway.

* `availability_zones` - The availability zone IDs.

* `connection_number` - The max number of connections.

* `used_connection_number` - The number of used connections.

* `used_connection_group` - The number of used connection groups.

* `enterprise_project_id` - The enterprise project ID

* `eips` - The EIPs used by the fateway.
  The [eips](#Gateway_eips) structure is documented below.

* `created_at` - The create time.

* `updated_at` - The update time.

* `access_vpc_id` - The ID of the access VPC.

* `access_subnet_id` - The ID of the access subnet.

* `access_private_ips` - The list of private access IPs.

* `ha_mode` - The HA mode.
  The value can be: **active-active** and **active-standby**.

<a name="Gateway_eips"></a>
The `eips` block supports:

* `bandwidth_billing_info` - The bandwidth billing info.

* `bandwidth_id` - The bandwidth ID.

* `bandwidth_name` - The bandwidth name.

* `bandwidth_size` - Bandwidth size in Mbit/s.

* `billing_info` - The billing info.

* `charge_mode` - The charge mode of the bandwidth.

* `id` - The public IP ID.

* `ip_address` - The public IP address.

* `ip_version` - The public IP version.

* `type` - The EIP type.
