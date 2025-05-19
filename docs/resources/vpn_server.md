---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_server"
description: |-
  Manages a VPN server resource within HuaweiCloud.
---

# huaweicloud_vpn_server

Manages a VPN server resource within HuaweiCloud.

## Example Usage

### create a server, and client auth type is LOCAL_PASSWORD

```hcl
variable "p2c_vgw_id" {}
variable "local_subnets" {}
variable "client_cidr" {}
variable "server_certificate_id" {}

resource "huaweicloud_vpn_server" "test" {
  p2c_vgw_id       = var.p2c_vgw_id
  local_subnets    = var.local_subnets
  client_cidr      = var.client_cidr
  client_auth_type = "LOCAL_PASSWORD"

  server_certificate {
    id = var.server_certificate_id
  }

  ssl_options {
    protocol             = "TCP"
    port                 = 443
    encryption_algorithm = "AES-128-GCM"
    is_compressed        = false
  }
}
```

### create a server, and client auth type is CERT

```hcl
variable "p2c_vgw_id" {}
variable "local_subnets" {}
variable "client_cidr" {}
variable "server_certificate_id" {}

resource "huaweicloud_vpn_server" "test" {
  p2c_vgw_id       = var.p2c_vgw_id
  local_subnets    = var.local_subnets
  client_cidr      = var.client_cidr
  client_auth_type = "CERT"

  server_certificate {
    id = var.server_certificate_id
  }

  ssl_options {
    protocol             = "TCP"
    port                 = 443
    encryption_algorithm = "AES-128-GCM"
    is_compressed        = false
  }

  client_ca_certificates {
    name    = "test-cert"
    content = trimsuffix(<<EOF
-----BEGIN CERTIFICATE-----
YOUR CERT CONTENT
-----END CERTIFICATE-----
EOF
, "\n")
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `p2c_vgw_id` - (Required, String) Specifies the ID of a P2C VPN gateway instance.

* `local_subnets` - (Required, List) Specifies the list of local CIDR blocks. A local CIDR block is a destination CIDR
  block on the cloud to be accessed by client CIDR blocks through a VPN. The value is in the format of dotted decimal
  notation/mask, for example, **10.10.1.0/24**. By default, a maximum of **20** local CIDR blocks are supported. The
  local CIDR block cannot be **0.0.0.0/8**, **127.0.0.0/8**, **224.0.0.0/4**, or **240.0.0.0/4**.

* `client_cidr` - (Required, String) Specifies the client CIDR block. A virtual IP address on this CIDR block will be
  assigned to a client for establishing a connection. The value is in the format of dotted decimal notation/mask,
  for example, **192.168.1.0/24**. The client CIDR block cannot conflict with the routes in the default route table of
  the VPC to which the gateway belongs or any local CIDR block of the server. The number of available IP addresses in
  the client CIDR block must be greater than **four** times the maximum number of gateway connections. The client CIDR
  block cannot be **0.0.0.0/8**, **127.0.0.0/8**, **224.0.0.0/4**, **240.0.0.0/4**, or **169.254.0.0/16**.

* `tunnel_protocol` - (Optional, String) Specifies the tunnel protocol. Value can be **SSL**. Defaults to **SSL**.

* `client_auth_type` - (Optional, String) Specifies the client authentication mode.
  Value can be as follows:
  + **CERT**: certificate authentication
  + **LOCAL_PASSWORD**: password authentication (local)
  
  The default value is **LOCAL_PASSWORD**.

* `client_ca_certificates` - (Optional, List) Specifies the list of client CA certificates, which are used to
  authenticate client certificates. This parameter is mandatory when **SSL** is used as the `tunnel_protocol` and
  the `client_auth_type` is **CERT**. A maximum of **10** client CA certificates can be uploaded.

  The [client_ca_certificates](#block--client_ca_certificates) structure is documented below.

* `server_certificate` - (Optional, List) Specifies the server certificate info. This parameter is mandatory when
  **SSL** is used as the `tunnel_protocol`. It is recommended to use a certificate with a strong cryptographic
  algorithm, such as **RSA-3072** or **RSA-4096**.

  The [server_certificate](#block--server_certificate) structure is documented below.

* `ssl_options` - (Optional, List) Specifies the SSL options. This parameter is mandatory when **SSL** is used as the
  `tunnel_protocol`.

  The [ssl_options](#block--ssl_options) structure is documented below.

* `os_type` - (Optional, String) Specifies the OS type.
  Value can be **Windows**, **Linux**, **MacOS**, **Android** or **iOS**. The default value is **Windows**.

<a name="block--client_ca_certificates"></a>
The `client_ca_certificates` block supports:

* `content` - (Required, String) Specifies the certificate content.

* `name` - (Optional, String) Specifies the certificate name. If this parameter is left blank, the system automatically
  generates a certificate name. The value is a string of **1** to **64** characters, which can contain digits, letters,
  underscores (_), and hyphens (-).

<a name="block--server_certificate"></a>
The `server_certificate` block supports:

* `id` - (Required, String) Specifies the certificate ID, which is the ID of a certificated uploaded in the Cloud
  Certificate Manager (CCM).

<a name="block--ssl_options"></a>
The `ssl_options` block supports:

* `encryption_algorithm` - (Optional, String) Specifies the encryption algorithm.
  Value can be **AES-128-GCM**, **AES-256-GCM**. The default value is **AES-128-GCM**.

* `is_compressed` - (Optional, Bool) Specifies whether to compress data. The default value is false.

* `port` - (Optional, Int) Specifies the port number. Value can be **443**, **1194**. The default value is **443**.

* `protocol` - (Optional, String) Specifies the protocol. Value can be **TCP**. Defaults to **TCP**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `client_ca_certificates_uploaded` - Indicates the list of client CA certificates.

  The [client_ca_certificates_uploaded](#attrblock--client_ca_certificates_uploaded) structure is documented below.

* `server_certificate` - Indicates the server certificate info.

  The [server_certificate](#attrblock--server_certificate) structure is documented below.

* `ssl_options` - Indicates the SSL options.

  The [ssl_options](#attrblock--ssl_options) structure is documented below.

* `client_config` - The client config.

* `status` - The server status.

* `created_at` - The creation time.

* `updated_at` - The update time.

<a name="attrblock--client_ca_certificates_uploaded"></a>
The `client_ca_certificates_uploaded` block supports:

* `created_at` - The creation time of the client CA certificate.

* `expiration_time` - The expiration time of the client CA certificate.

* `id` - The client CA certificate ID.

* `issuer` - The issuer of the client CA certificate.

* `name` - Indicates the certificate name

* `serial_number` - The serial number of the client CA certificate.

* `signature_algorithm` - The signature algorithm of the client CA certificate.

* `subject` - The subject of the client CA certificate.

* `updated_at` - The update time of the client CA certificate.

<a name="attrblock--server_certificate"></a>
The `server_certificate` block supports:

* `expiration_time` - The expiration time of the server certificate.

* `issuer` - The issuer of the server certificate.

* `name` - The server certificate name.

* `serial_number` - The serial number of the server certificate.

* `signature_algorithm` - The signature algorithm of the server certificate.

* `subject` - The subject of the server certificate.

<a name="attrblock--ssl_options"></a>
The `ssl_options` block supports:

* `authentication_algorithm` - The authentication algorithm.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 20 minutes.

## Import

The server can be imported using `p2c_vgw_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_vpn_server.test <p2c_vgw_id>/<id>
```

Please add the followings if some attributes are missing when importing the resource.

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `client_ca_certificates` and `os_type`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the server, or the resource definition should be updated to
align with the server. Also you can ignore changes as below.

```hcl
resource "huaweicloud_vpn_server" "test" {
    ...

  lifecycle {
    ignore_changes = [
      client_ca_certificates, os_type,
    ]
  }
}
```
