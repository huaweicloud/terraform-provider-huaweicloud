---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_connections"
description: ""
---

# huaweicloud_vpn_connections

Use this data source to get a list of VPN connections.

## Example Usage

```hcl
variable "status" {}
variable "name" {}
variable "vpn_type" {}
variable "gateway_id" {}
variable "gateway_ip" {}

data "huaweicloud_vpn_connections" "services" {
  status     = var.status
  name       = var.name
  vpn_type   = var.vpn_type
  gateway_id = var.gateway_id
  gateway_ip = var.gateway_ip
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the VPN customer gateways.
  If omitted, the provider-level region will be used.

* `connection_id` - (Optional, String) Specifies the ID of the VPN connection.

* `name` - (Optional, String) Specifies the name of the VPN connection.

* `gateway_id` - (Optional, String) Specifies the gateway ID of the VPN connection.

* `gateway_ip` - (Optional, String) Specifies the gateway IP of the VPN connection.

* `status` - (Optional, String) Specifies the status of the VPN connection.

* `vpn_type` - (Optional, String) Specifies the VPN type of the VPN connection.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `connections` - All resource connection that match the filter parameters.
  The [connections](#connections) structure is documented below.

<a name="connections"></a>
The `connections` block supports:

* `id` - Indicates the ID of the connection.

* `name` - Indicates the name of the connection.

* `status` - Indicates the status of the connection.

* `vpn_type` - Indicates the VPN type of the connection.

* `gateway_id` - Indicates the gateway ID of the connection.

* `gateway_ip` - Indicates the gateway IP of the connection.

* `customer_gateway_id` - Indicates the customer gateway ID of the connection.

* `peer_subnets` - Indicates the peer subnets of the connection.

* `tunnel_local_address` - Indicates the tunnel local address of the connection.

* `tunnel_peer_address` - Indicates the tunnel peer address of the connection.

* `enable_nqa` - Indicates the enable nqa of the connection.

* `enterprise_project_id` - Indicates the enterprise project ID of the connection.

* `connection_monitor_id` - Indicates the connection monitor ID of the connection.

* `ha_role` - Indicates the ha role of the connection.

* `created_at` - The created time.

* `updated_at` - The last updated time.

* `policy_rules` - Indicates the policy rules information of the connection.

* `ipsecpolicy` - Indicates the ipsecpolicy information of the connection.

* `ikepolicy` - Indicates the ikepolicy information of the connection.

  The [policy_rules](#policy_Rules) structure is documented below.

<a name="policy_Rules"></a>
The `policy_rules` block supports:

* `rule_index` - Indicates the rule index of the policy rules.

* `source` - Indicates the source of the policy rules certificate.

* `destination` - Indicates the destination of the policy rules certificate.

  The [ipsecpolicy](#ipsecpolicy) structure is documented below.

<a name="ipsecpolicy"></a>
The `ipsecpolicy` block supports:

* `authentication_algorithm` - Indicates the authentication algorithm of the ipsecpolicy certificate.

* `encryption_algorithm` - Indicates the encryption algorithm of the ipsecpolicy certificate.

* `pfs` - Indicates the pfs of the ipsecpolicy certificate.

* `transform_protocol` - Indicates the transform protocol of the ipsecpolicy certificate.

* `lifetime_seconds` - Indicates the lifetime seconds of the ipsecpolicy certificate.

* `encapsulation_mode` - Indicates the encapsulation mode of the ipsecpolicy certificate.

  The [ikepolicy](#ikepolicy) structure is documented below.

<a name="ikepolicy"></a>
The `ikepolicy` block supports:

* `ike_version` - Indicates the ike version of the ikepolicy certificate.

* `phase1_negotiation_mode` - Indicates the phase1 negotiation mode of the ikepolicy certificate.

* `authentication_algorithm` - Indicates the authentication algorithm of the ikepolicy certificate.

* `encryption_algorithm` - Indicates the encryption algorithm of the ikepolicy certificate.

* `dh_group` - Indicates the dh group of the ikepolicy certificate.

* `authentication_method` - Indicates the souauthentication methodrce of the ikepolicy certificate.

* `lifetime_method` - Indicates the lifetime method of the ikepolicy certificate.

* `local_id_type` - Indicates the local ID type of the ikepolicy certificate.

* `local_id` - Indicates the local ID of the ikepolicy certificate.

* `peer_id_type` - Indicates the peer ID type of the ikepolicy certificate.

* `peer_id` - Indicates the peer ID of the ikepolicy certificate.

* `dpd` - Indicates the dpd information of the ikepolicy certificate.

 The [dpd](#dpd) structure is documented below.

<a name="dpd"></a>
The `ikedpdpolicy` block supports:

* `timeout` - Indicates the timeout of the dpd certificate.

* `interval` - Indicates the interval of the dpd certificate.

* `msg` - Indicates the msg of the dpd certificate.
