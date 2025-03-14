---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_gateway"
description: ""
---

# huaweicloud_vpn_gateway

Manages a VPN gateway resource within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "name" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "eip_id1" {}
variable "eip_id2" {}

data "huaweicloud_vpn_gateway_availability_zones" "test" {
  flavor          = "professional1"
  attachment_type = "vpc"
}

resource "huaweicloud_vpn_gateway" "test" {
  name               = var.name
  vpc_id             = var.vpc_id
  local_subnets      = ["192.168.0.0/24", "192.168.1.0/24"]
  connect_subnet     = var.subnet_id
  availability_zones = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  eip1 {
    id = var.eip_id1
  }

  eip2 {
    id = var.eip_id2
  }
}
```

### Creating a VPN gateway with creating new EIPs

```hcl
variable "name" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "bandwidth_name1" {}
variable "bandwidth_name2" {}

data "huaweicloud_vpn_gateway_availability_zones" "test" {
  flavor          = "professional1"
  attachment_type = "vpc"
}

resource "huaweicloud_vpn_gateway" "test" {
  name               = var.name
  vpc_id             = var.vpc_id
  local_subnets      = ["192.168.0.0/24", "192.168.1.0/24"]
  connect_subnet     = var.subnet_id
  availability_zones = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  eip1 {
    bandwidth_name = var.bandwidth_name1
    type           = "5_bgp"
    bandwidth_size = 5
    charge_mode    = "traffic"
  }

  eip2 {
    bandwidth_name = var.bandwidth_name2
    type           = "5_bgp"
    bandwidth_size = 5
    charge_mode    = "traffic"
  }
}
```

### Creating a private VPN gateway with Enterprise Router

```hcl
variable "name" {}
variable "er_id" {}
variable "access_vpc_id" {}
variable "access_subnet_id" {}
variable "access_private_ip_1" {}
variable "access_private_ip_2" {}

data "huaweicloud_vpn_gateway_availability_zones" "test" {
  flavor          = "professional1"
  attachment_type = "er"
}

resource "huaweicloud_vpn_gateway" "test" {
  name               = var.name
  network_type       = "private"
  attachment_type    = "er"
  er_id              = var.er_id
  availability_zones = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  access_vpc_id      = var.access_vpc_id
  access_subnet_id   = var.access_subnet_id
  
  access_private_ip_1 = var.access_private_ip_1
  access_private_ip_2 = var.access_private_ip_2
}
```

### Creating a GM VPN gateway with certificate

```hcl
variable "vpc_id" {}
variable "cidr" {}
variable "subnet_id" {}

data "huaweicloud_vpn_gateway_availability_zones" "test" {
  attachment_type = "er"
  flavor          = "GM"
}

resource "huaweicloud_vpn_gateway" "test" {
  name               = "test"
  vpc_id             = var.vpc_id
  flavor             = "GM"
  network_type       = "private"
  local_subnets      = [var.cidr]
  connect_subnet     = var.subnet_id
  availability_zones = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  certificate {
    name              = "test"
    content           = "-----BEGIN CERTIFICATE-----\nTHIS IS YOUR CERT CONTENT\n-----END CERTIFICATE-----"
    private_key       = "-----BEGIN EC PRIVATE KEY-----\nTHIS IS YOUR PRIVATE KEY CONTENT\n-----END EC PRIVATE KEY-----"
    certificate_chain = "-----BEGIN CERTIFICATE-----\nTHIS IS YOUR CERTIFICATE CHAIN CONTENT\n-----END CERTIFICATE-----"
    enc_certificate   = "-----BEGIN CERTIFICATE-----\nTHIS IS YOUR ENC CERTIFICATE CONTENT\n-----END CERTIFICATE-----"
    enc_private_key   = "-----BEGIN EC PRIVATE KEY-----\nTHIS IS YOUR ENC PRIVATE KEY CONTENT\n-----END EC PRIVATE KEY-----"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) The name of the VPN gateway.  
  The valid length is limited from `1` to `64`, only letters, digits, hyphens (-) and underscores (_) are allowed.

* `availability_zones` - (Required, List, ForceNew) The list of availability zone IDs.

  Changing this parameter will create a new resource.

* `flavor` - (Optional, String) The flavor of the VPN gateway.
  + The value at creation can be **Basic**, **Professional1**, **Professional2** or **GM**. Defaults to **Professional1**.
  + The value during update can be **Basic**, **Professional1** or **Professional2**.

* `attachment_type` - (Optional, String, ForceNew) The attachment type. The value can be **vpc** or **er**.
  Defaults to **vpc**.

  Changing this parameter will create a new resource.

* `network_type` - (Optional, String, ForceNew) The network type. The value can be **public** or **private**.
  Defaults to **public**.

  Changing this parameter will create a new resource.

* `vpc_id` - (Optional, String, ForceNew) The ID of the VPC to which the VPN gateway is connected.
  This parameter is mandatory when `attachment_type` is **vpc**.

  Changing this parameter will create a new resource.

* `local_subnets` - (Optional, List) The list of local subnets.
  This parameter is mandatory when `attachment_type` is **vpc**.

* `connect_subnet` - (Optional, String, ForceNew) The Network ID of the VPC subnet used by the VPN gateway.
  This parameter is mandatory when `attachment_type` is **vpc**.

  Changing this parameter will create a new resource.

* `er_id` - (Optional, String, ForceNew) The enterprise router ID to attach with to VPN gateway.
  This parameter is mandatory when `attachment_type` is **er**.

  Changing this parameter will create a new resource.

* `ha_mode` - (Optional, String, ForceNew) The HA mode of VPN gateway. Valid values are **active-active** and
  **active-standby**. The default value is **active-active**.

  Changing this parameter will create a new resource.

* `eip1` - (Optional, List, ForceNew) The master 1 IP in active-active VPN gateway or the master IP
  in active-standby VPN gateway. This parameter is mandatory when `network_type` is **public** or left empty.
  The [object](#Gateway_CreateRequestEip) structure is documented below.

  Changing this parameter will create a new resource.

* `eip2` - (Optional, List, ForceNew) The master 2 IP in active-active VPN gateway or the slave IP
  in active-standby VPN gateway. This parameter is mandatory when `network_type` is **public** or left empty.
  The [object](#Gateway_CreateRequestEip) structure is documented below.

  Changing this parameter will create a new resource.

* `access_vpc_id` - (Optional, String, ForceNew) The access VPC ID.
  The default value is the value of `vpc_id`.

  Changing this parameter will create a new resource.

* `access_subnet_id` - (Optional, String, ForceNew) The access subnet ID.
  The default value is the value of `connect_subnet`.

  Changing this parameter will create a new resource.

* `access_private_ip_1` - (Optional, String, ForceNew) The private IP 1 in private network type VPN gateway.
  It is the master IP 1 in **active-active** HA mode, and the master IP in **active-standby** HA mode.
  Must declare the **access_private_ip_2** at the same time, and can not use the same IP value.

  Changing this parameter will create a new resource.

* `access_private_ip_2` - (Optional, String, ForceNew) The private IP 2 in private network type VPN gateway.
  It is the master IP 2 in **active-active** HA mode, and the slave IP in **active-standby** HA mode.
  Must declare the **access_private_ip_1** at the same time, and can not use the same IP value.

  Changing this parameter will create a new resource.

* `asn` - (Optional, Int, ForceNew) The ASN number of BGP. The value ranges from `1` to `4,294,967,295`.
  Defaults to `64,512`.

  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String) The enterprise project ID.

* `certificate` - (Optional, List) The GM certificate of the **GM** flavor gateway.
  The [object](#Gateway_certificate) structure is documented below.

* `tags` - (Optional, Map) Specifies the tags of the VPN gateway.

* `delete_eip_on_termination` - (Optional, Bool) Whether to delete the EIP when the VPN gateway is deleted.
  Defaults to **true**.

<a name="Gateway_CreateRequestEip"></a>
The `eip1` or `eip2` block supports:

* `id` - (Optional, String, ForceNew) The public IP ID.

  Changing this parameter will create a new resource.

* `type` - (Optional, String, ForceNew) The EIP type. The value can be **5_bgp** and **5_sbgp**.

  Changing this parameter will create a new resource.

* `bandwidth_name` - (Optional, String, ForceNew) The bandwidth name.  
  The valid length is limited from `1` to `64`, only letters, digits, hyphens (-) and underscores (_) are allowed.

  Changing this parameter will create a new resource.

* `bandwidth_size` - (Optional, Int, ForceNew) Bandwidth size in Mbit/s. When the `flavor` is **Basic**, the value
  cannot be greater than `100`. When the `flavor` is **Professional1**, the value cannot be greater than `300`.
  When the `flavor` is **Professional2**, the value cannot be greater than `1,000`.

  Changing this parameter will create a new resource.

* `charge_mode` - (Optional, String, ForceNew) The charge mode of the bandwidth. The value can be **bandwidth** and **traffic**.

  Changing this parameter will create a new resource.

  ~> You can use `id` to specify an existing EIP or use `type`, `bandwidth_name`, `bandwidth_size` and `charge_mode` to
    create a new EIP.

<a name="Gateway_certificate"></a>
The `certificate` block supports:

* `name` - (Required, String) The name of the gateway certificate.

* `content` - (Required, String) The content of the gateway certificate.

* `private_key` - (Required, String) The private of the gateway certificate.

* `certificate_chain` - (Required, String) The certificate chain of the gateway certificate.

* `enc_certificate` - (Required, String) The enc certificate of the gateway certificate.

* `enc_private_key` - (Required, String) The enc private key of the gateway certificate.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the VPN gateway

* `status` - The status of VPN gateway.

* `created_at` - The create time.

* `updated_at` - The update time.

* `used_connection_group` - The number of used connection groups.

* `used_connection_number` - The number of used connections.

* `er_attachment_id` - The ER attachment ID.

* `eip1` - The master 1 IP in active-active VPN gateway or the master IP in active-standby VPN gateway.
  The [object](#Gateway_GetResponseEip) structure is documented below.

* `eip2` - The master 2 IP in active-active VPN gateway or the slave IP in active-standby VPN gateway.
  The [object](#Gateway_GetResponseEip) structure is documented below.

* `certificate` - The GM certificate of the **GM** flavor gateway.
  The [object](#Gateway_certificate_attr) structure is documented below.

<a name="Gateway_GetResponseEip"></a>
The `eip1` or `eip2` block supports:

* `bandwidth_id` - The bandwidth ID.

* `ip_address` - The public IP address.

* `ip_version` - The public IP version.

<a name="Gateway_certificate_attr"></a>
The `certificate` block supports:

* `certificate_id` - The certificate ID.

* `status` - The status of the certificate.

* `issuer` - The issuer of the certificate.

* `signature_algorithm` - The signature algorithm of the certificate.

* `certificate_serial_number` - The serial number of the certificate.

* `certificate_subject` - The subject of the certificate.

* `certificate_expire_time` - The expire time of the certificate.

* `certificate_chain_serial_number` - The serial number of the certificate chain.

* `certificate_chain_subject` - The subject of the certificate chain.

* `certificate_chain_expire_time` - The expire time of the certificate.

* `enc_certificate_subject` - The subject of the enc certificate.

* `enc_certificate_expire_time` - The expire time of the enc certificate.

* `enc_certificate_serial_number` - The serial number of the enc certificate.

* `created_at` - The create time of the gateway certificate.

* `updated_at` - The update time of the gateway certificate.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The gateway can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_vpn_gateway.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attribute is `delete_eip_on_termination`. It is generally
recommended running `terraform plan` after importing the resource. You can then decide if changes should be applied
to the gateway, or the resource definition should be updated to align with the gateway.
Also you can ignore changes as below.

```hcl
resource "huaweicloud_vpn_gateway" "test" {
    ...

  lifecycle {
    ignore_changes = [
      delete_eip_on_termination
    ]
  }
}
```
