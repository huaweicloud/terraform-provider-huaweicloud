---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_servers"
description: |-
  Use this data source to get the list of VPN servers.
---

# huaweicloud_vpn_servers

Use this data source to get the list of VPN servers.

## Example Usage

```hcl
data "huaweicloud_vpn_servers" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `vpn_servers` - The VPN server list.

  The [vpn_servers](#vpn_servers_struct) structure is documented below.

<a name="vpn_servers_struct"></a>
The `vpn_servers` block supports:

* `p2c_vgw_id` - The ID of a P2C VPN gateway.

* `id` - The server ID.

* `client_auth_type` - The client authentication mode.

* `tunnel_protocol` - A tunnel protocol.

* `status` - The server status.

* `local_subnets` - The local CIDR block list.

* `server_certificate` - The server certificate information.

  The [server_certificate](#vpn_servers_server_certificate_struct) structure is documented below.

* `client_ca_certificates` - The client CA certificate information.

  The [client_ca_certificates](#vpn_servers_client_ca_certificates_struct) structure is documented below.

* `ssl_options` - The SSL options information.

  The [ssl_options](#vpn_servers_ssl_options_struct) structure is documented below.

* `client_cidr` - The client CIDR block.

* `created_at` - The creation time.

* `updated_at` - The update time.

<a name="vpn_servers_server_certificate_struct"></a>
The `server_certificate` block supports:

* `id` - The server certificate ID.

* `name` - The server certificate name.

* `serial_number` - The serial number of the server certificate.

* `expiration_time` - The expiration time of the server certificate.

* `signature_algorithm` - The signature algorithm of the server certificate.

* `issuer` - The issuer of the server certificate.

* `subject` - The subject of the server certificate.

<a name="vpn_servers_client_ca_certificates_struct"></a>
The `client_ca_certificates` block supports:

* `id` - The client CA certificate ID.

* `name` - The name of the client CA certificate.

* `issuer` - The issuer of the client CA certificate.

* `signature_algorithm` - The signature algorithm of the client CA certificate.

* `subject` - The subject of the client CA certificate.

* `serial_number` - The serial number of the client CA certificate.

* `expiration_time` - The expiration time of the client CA certificate.

* `created_at` - The creation time of the client CA certificate.

* `updated_at` - The update time of the client CA certificate.

<a name="vpn_servers_ssl_options_struct"></a>
The `ssl_options` block supports:

* `protocol` - The protocol.

* `port` - The port.

* `is_compressed` - Whether compression is enabled.

* `encryption_algorithm` - The encryption algorithm.

* `authentication_algorithm` - The authentication algorithm.
