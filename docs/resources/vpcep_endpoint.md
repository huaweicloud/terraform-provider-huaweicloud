---
subcategory: "VPC Endpoint (VPCEP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcep_endpoint"
description: |
  Provides a resource to manage a VPC endpoint resource.
---

# huaweicloud_vpcep_endpoint

Provides a resource to manage a VPC endpoint resource.

## Example Usage

### Access to the public interface service

```hcl
variable "public_interface_service_id" {}
variable "vpc_id" {}
variable "network_id" {}

resource "huaweicloud_vpcep_endpoint" "test" {
  service_id       = var.public_interface_service_id
  vpc_id           = var.vpc_id
  network_id       = var.network_id
  enable_dns       = true
  enable_whitelist = true
  whitelist        = ["192.168.0.0/24"]
}
```

### Access to the private interface service

```hcl
variable "service_vpc_id" {}
variable "vm_port" {}
variable "vpc_id" {}
variable "network_id" {}

resource "huaweicloud_vpcep_service" "test" {
  name        = "demo-service"
  server_type = "VM"
  vpc_id      = var.service_vpc_id
  port_id     = var.vm_port

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
}

resource "huaweicloud_vpcep_endpoint" "test" {
  service_id  = huaweicloud_vpcep_service.test.id
  vpc_id      = var.vpc_id
  network_id  = var.network_id
  enable_dns  = true
  description = "test description"
}
```

### Access to the gateway service without policy statement

```hcl
variable "gateway_service_id" {}
variable "vpc_id" {}

resource "huaweicloud_vpcep_endpoint" "test" {
  service_id  = var.gateway_service_id
  vpc_id      = var.vpc_id
  description = "test description"
}
```

### Access to the gateway service with policy statement

```hcl
variable "gateway_service_id" {}
variable "vpc_id" {}

resource "huaweicloud_vpcep_endpoint" "test" {
  service_id  = var.gateway_service_id
  vpc_id      = var.vpc_id
  description = "test description"

  policy_statement = <<EOF
  [
    {
      "Effect": "Allow",
      "Action": [
        "obs:bucket:ListBucket"
      ],
      "Resource": [
        "obs:*:*:*:*/*",
        "obs:*:*:*:*"
      ]
    },
    {
      "Effect": "Deny",
      "Action": [
        "obs:object:DeleteObject"
      ],
      "Resource": [
        "obs:*:*:*:*/*",
        "obs:*:*:*:*"
      ]
    }
  ]
EOF
}
```

### Access to the interface service with policy statement

```hcl
variable "gateway_service_id" {}
variable "vpc_id" {}
variable "network_id" {}

resource "huaweicloud_vpcep_endpoint" "test" {
  service_id  = var.gateway_service_id
  vpc_id      = var.vpc_id
  network_id  = var.network_id
  description = "test description"

  policy_document = <<EOF
  {
    "Version": "5.0",
    "Statement": [
      {
        "Action": [
          "*"
        ],
        "Resource": [
          "*"
        ],
        "Effect": "Allow",
        "Principal": "*"
      }
    ]
  }
  EOF
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the VPC endpoint. If omitted, the provider-level
  region will be used. Changing this creates a new VPC endpoint.

* `service_id` - (Required, String, ForceNew) Specifies the ID of the VPC endpoint service.
  The VPC endpoint service could be private interface service, public interface service or gateway service.
  + For private interface service, the value of `service_id` can be obtained through resource `huaweicloud_vpcep_service`
    or datasource `huaweicloud_vpcep_services`.
  + For public interface service, the value of `service_id` can be obtained through datasource
    `huaweicloud_vpcep_public_services`.
  + For gateway service, due to API reasons, the current provider's capabilities do not support the creation of gateway
    VPC endpoint services. Please try to obtain `service_id` through datasource `huaweicloud_vpcep_public_services` or
    look for VPCEP operation and maintenance help to find the gateway service ID.

  Changing this creates a new VPC endpoint.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the VPC where the VPC endpoint is to be created. Changing
  this creates a new VPC endpoint.

* `network_id` - (Optional, String, ForceNew) Specifies the network ID of the subnet in the VPC specified by `vpc_id`.
  This field is required when creating a VPC endpoint for connecting an interface VPC endpoint service.
  The use of this field has the following restrictions:
  + The subnet CIDR block of the VPC cannot overlap with **198.19.128.0/17**.
  + The destination address of the custom route in the VPC route table cannot overlap with **198.19.128.0/17**.

  Changing this creates a new VPC endpoint.

* `ip_address` - (Optional, String, ForceNew) Specifies the IP address for accessing the associated VPC endpoint
  service. Only IPv4 addresses are supported. This field is required when creating a VPC endpoint for connecting an
  interface VPC endpoint service.

  Changing this creates a new VPC endpoint.

* `enable_dns` - (Optional, Bool, ForceNew) Specifies whether to create a private domain name. The default value is
  **true**. This field is valid only when creating a VPC endpoint for connecting an interface VPC endpoint service.

  Changing this creates a new VPC endpoint.

* `description` - (Optional, String, ForceNew) Specifies the description of the VPC endpoint. The value can contain
  characters such as letters and digits, but cannot contain less than signs (<) and great than signs (>).

  Changing this creates a new VPC endpoint.

* `routetables` - (Optional, List) Specifies the IDs of the route tables associated with the VPC endpoint.
  This field is valid only when creating a VPC endpoint for connecting a gateway VPC endpoint service.
  The default route table will be used when this field is not specified.

* `enable_whitelist` - (Optional, Bool) Specifies whether to enable access control. The default value is **false**.

* `whitelist` - (Optional, List) Specifies the list of IP address or CIDR block which can be accessed to the
  VPC endpoint. This field is valid when `enable_whitelist` is set to **true**. The max length of whitelist is 20.
  This field is valid only when creating a VPC endpoint for connecting an interface VPC endpoint service.

* `policy_statement` - (Optional, String) Specifies the policy of the gateway VPC endpoint. The value is a string in
  JSON array format. This parameter is only available when `enable_policy` of the VPC endpoint services for
  Object Storage Service (OBS) and Scalable File Service (SFS) is set to **true**.

  -> Please refer to [official document](https://support.huaweicloud.com/intl/en-us/usermanual-iam/iam_01_0017.html) for
  the data structure rules of the policy. Just pay attention to the fields `Effect`, `Action` and `Resource`.

* `policy_document` - (Optional, String) Specifies the policy of the interface VPC endpoint. The value is a string in
  JSON array format. This parameter is only available when `enable_policy` set to **true**. This parameter is not
  available for Object Storage Service (OBS) and Scalable File Service (SFS).

  -> Please refer to [official document](https://support.huaweicloud.com/intl/en-us/usermanual-iam/iam_01_0017.html) for
  the data structure rules of the policy. Just pay attention to the fields `Effect`, `Action` and `Resource`.

* `ip_version` - (Optional, String, ForceNew) Specifies the IP version of the VPC endpoint.
  Changing this will create a new resource.
  The valid values are as follows:
  + **ipv4**: The VPC endpoint IP address can only be an IPv4 address.
  + **dualstack**: The VPC endpoint IP address can be an IPv4 address or IPv6 address.

* `ipv6_address` - (Optional, String, ForceNew) Specifies the IPv6 address for accessing the connected VPC
  endpoint service.
  Changing this will create a new resource.
  
  -> The IPv6 address must be a subnet of the IPv6 network segment of the subnet associated with the VPC endpoint.
  <br/>If you not specifies, the IPv6 address generated by the system will be used.

-> 1.Only professional VPC endpoint supports `ip_version` and `ipv6_address` parameters.
  <br/>2.Currently, professional VPC endpoints are available in the **cn-east-4**, **me-east-1**, **cn-east-5**,
  and **af-north-1** regions.

* `tags` - (Optional, Map) The key/value pairs to associate with the VPC endpoint.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of the VPC endpoint.

* `status` - The status of the VPC endpoint. The value can be **pendingAcceptance**, **creating**, **accepted**,
  **rejected**, **failed**, **deleting**.

* `service_name` - The name of the VPC endpoint service.

* `service_type` - The type of the VPC endpoint service.

* `packet_id` - The packet ID of the VPC endpoint.

* `private_domain_name` - The domain name for accessing the associated VPC endpoint service. This parameter is only
  available when enable_dns is set to true.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

VPC endpoint can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_vpcep_endpoint.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `enable_dns`.

It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_vpcep_endpoint" "test" {
  ...

  lifecycle {
    ignore_changes = [
      enable_dns,
    ]
  }
}
```
