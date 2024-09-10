---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_connection"
description: ""
---

# huaweicloud_vpn_connection

Manages a VPN connection resource within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "name" {}
variable "peer_subnet" {}
variable "gateway_id" {}
variable "gateway_ip" {}
variable "customer_gateway_id" {}

resource "huaweicloud_vpn_connection" "test" {
  name                = var.name
  gateway_id          = var.gateway_id
  gateway_ip          = var.gateway_ip
  customer_gateway_id = var.customer_gateway_id
  peer_subnets        = [var.peer_subnet]
  vpn_type            = "static"
  psk                 = "Test@123"
}
```

### VPN connection with policy

```hcl
variable "name" {}
variable "peer_subnet" {}
variable "gateway_id" {}
variable "gateway_ip" {}
variable "customer_gateway_id" {}

resource "huaweicloud_vpn_connection" "test" {
  name                = var.name
  gateway_id          = var.gateway_id
  gateway_ip          = var.gateway_ip
  customer_gateway_id = var.customer_gateway_id
  peer_subnets        = [var.peer_subnet]
  vpn_type            = "static"
  psk                 = "Test@123"

  ikepolicy {
    authentication_algorithm = "sha2-256"
    authentication_method    = "pre-share"
    encryption_algorithm     = "aes-128"
    ike_version              = "v2"
    lifetime_seconds         = 86400
    pfs                      = "group14"
  }

  ipsecpolicy {
    authentication_algorithm = "sha2-256"
    encapsulation_mode       = "tunnel"
    encryption_algorithm     = "aes-128"
    lifetime_seconds         = 3600
    pfs                      = "group14"
    transform_protocol       = "esp"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) The name of the VPN connection.

* `gateway_id` - (Required, String, ForceNew) The VPN gateway ID.

  Changing this parameter will create a new resource.

* `gateway_ip` - (Required, String, ForceNew) The VPN gateway IP ID.

  Changing this parameter will create a new resource.

* `vpn_type` - (Required, String, ForceNew) The connection type. The value can be **policy**, **static** or **bgp**.

  Changing this parameter will create a new resource.

* `customer_gateway_id` - (Required, String) The customer gateway ID.

* `psk` - (Required, String) The pre-shared key.

* `peer_subnets` - (Optional, List) The CIDR list of customer subnets. This parameter must be empty
  when the `attachment_type` of the VPN gateway is set to **er** and `vpn_type` is set to **policy** or **bgp**.
  This parameter is mandatory in other scenarios.

* `tunnel_local_address` - (Optional, String) The local tunnel address.

* `tunnel_peer_address` - (Optional, String) The peer tunnel address.

* `enable_nqa` - (Optional, Bool) Whether to enable NQA check. Defaults to **false**.

* `ikepolicy` - (Optional, List) The IKE policy configurations.
The [ikepolicy](#Connection_CreateRequestIkePolicy) structure is documented below.

* `ipsecpolicy` - (Optional, List) The IPsec policy configurations.
The [ipsecpolicy](#Connection_CreateRequestIpsecPolicy) structure is documented below.

* `policy_rules` - (Optional, List) The policy rules. Only works when vpn_type is set to **policy**
The [policy_rules](#Connection_PolicyRule) structure is documented below.

* `tags` - (Optional, Map) Specifies the tags of the VPN connection.

* `ha_role` - (Optional, String, ForceNew) Specifies the mode of the VPN connection.
  The valid values are **master** and **slave**, defaults to **master**.
  This parameter is optional when you create a connection for a VPN gateway in **active-active** mode.
  When you create a connection for a VPN gateway in **active-standby** mode, **master** indicates
  the active connection, and **slave** indicates the standby connection.
  In **active-active** mode, this field must be set to **master** for the connection established
  using the active EIP or active private IP address of the VPN gateway, and must be set to **slave**
  for the connection established using active EIP 2 or active private IP address 2 of the VPN gateway.

  Changing this parameter will create a new resource.

<a name="Connection_CreateRequestIkePolicy"></a>
The `ikepolicy` block supports:

* `authentication_algorithm` - (Optional, String) The authentication algorithm. The value can be **sha1**, **md5**,
  **sha2-256**, **sha2-384**, **sha2-512**. Defaults to **sha2-256**. **sha1** and **md5** are less secure,
  please use them with caution.

* `encryption_algorithm` - (Optional, String) The encryption algorithm. The value can be **3des**, **aes-128**, **aes-192**,
  **aes-256**, **aes-128-gcm-16**, **aes-256-gcm-16**, **aes-128-gcm-128**, **aes-256-gcm-128**. Defaults to **aes-128**.
  **3des** is less secure, please use it with caution.

* `ike_version` - (Optional, String) The IKE negotiation version. The value can be **v1** and **v2**. Defaults to **v2**.

* `lifetime_seconds` - (Optional, Int) The life cycle of SA in seconds. The value ranges from `60` to `604,800`.
  Defaults to `86,400`. When the life cycle expires, IKE SA will be automatically updated.

* `local_id_type` - (Optional, String) The local ID type. The value can be **ip** or **fqdn**. Defaults to **ip**.

* `local_id` - (Optional, String) The local ID.

* `peer_id_type` - (Optional, String) The peer ID type. The value can be **ip**, **fqdn** or **any**. Defaults to **ip**.

* `peer_id` - (Optional, String) The peer ID.

* `phase1_negotiation_mode` - (Optional, String) The negotiation mode, only works when the ike_version is v1.
  The value can be **main** or **aggressive**. Defaults to **main**.

* `authentication_method` - (Optional, String, ForceNew) The authentication method during IKE negotiation.
  The value can be **pre-share** and **digital-envelope-v2**. Defaults to **pre-share**.

* `dh_group` - (Optional, String) Specifies the DH group used for key exchange in phase 1.
  The value can be **group1**, **group2**, **group5**, **group14**, **group15**, **group16**, **group19**, **group20**,
  or **group21**. Exercise caution when using **group1**, **group2**, **group5**,
  or **group14** as they have low security. Defaults to **group15**.

* `dpd` - (Optional, List) Specifies the dead peer detection (DPD) object.
  The [dpd](#Connection_DPD) structure is documented below.

<a name="Connection_DPD"></a>
The `dpd` block supports:

* `timeout` - (Optional, Int) Specifies the interval for retransmitting DPD packets.
  The value ranges from `2` to `60`, in seconds. Defaults to `15`.

* `interval` - (Optional, Int) Specifies the DPD idle timeout period.
  The value ranges from `10` to `3,600`, in seconds. Defaults to `30`.

* `msg` - (Optional, String) Specifies the format of DPD packets. The value can be:
  + **seq-hash-notify**: indicates that the payload of DPD packets is in the sequence of hash-notify;
  + **seq-notify-hash**: indicates that the payload of DPD packets is in the sequence of notify-hash;

  Defaults to **seq-hash-notify**.

<a name="Connection_CreateRequestIpsecPolicy"></a>
The `ipsecpolicy` block supports:

* `authentication_algorithm` - (Optional, String) The authentication algorithm. The value can be **sha1**, **md5**,
  **sha2-256**, **sha2-384**, **sha2-512**. Defaults to **sha2-256**. **sha1** and **md5** are less secure,
  please use them with caution.

* `encryption_algorithm` - (Optional, String) The encryption algorithm. The value can be **3des**, **aes-128**, **aes-192**,
  **aes-256**, **aes-128-gcm-16**, **aes-256-gcm-16**, **aes-128-gcm-128**, **aes-256-gcm-128**. Defaults to **aes-128**.
  **3des** is less secure, please use it with caution.

* `pfs` - (Optional, String) The DH key group used by PFS. The value can be **group1**, **group2**, **group5**, **group14**
  **group16**, **group19**, **group20**, **group21**. Defaults to **group14**.

* `lifetime_seconds` - (Optional, Int) The lifecycle time of Ipsec tunnel in seconds.
  The value ranges from `60` to `604,800`. Defaults to `3600`.

* `transform_protocol` - (Optional, String) The transform protocol. Only **esp** supported for now.
  Defaults to **esp**.

* `encapsulation_mode` - (Optional, String) The encapsulation mode, only **tunnel** supported for now.
  Defaults to **tunnel**.

<a name="Connection_PolicyRule"></a>
The `policy_rules` block supports:

* `rule_index` - (Optional, Int) The rule index.

* `destination` - (Optional, List) The list of destination CIDRs.

* `source` - (Optional, String) The source CIDR.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the VPN connection.

* `enterprise_project_id` - The enterprise project ID.

* `created_at` - The create time.

* `updated_at` - The update time.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The connection can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_vpn_connection.test <id>
```
