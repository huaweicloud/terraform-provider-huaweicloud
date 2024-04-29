---
subcategory: "VPC Endpoint (VPCEP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcep_endpoint"
description: ""
---

# huaweicloud_vpcep_endpoint

Provides a resource to manage a VPC endpoint resource.

## Example Usage

### Access to the public service

```hcl
variable "vpc_id" {}
variable "network_id" {}

data "huaweicloud_vpcep_public_services" "cloud_service" {
  service_name = "dis"
}

resource "huaweicloud_vpcep_endpoint" "myendpoint" {
  service_id       = data.huaweicloud_vpcep_public_services.cloud_service.services[0].id
  vpc_id           = var.vpc_id
  network_id       = var.network_id
  enable_dns       = true
  enable_whitelist = true
  whitelist        = ["192.168.0.0/24"]
}
```

### Access to the private service

```hcl
variable "service_vpc_id" {}
variable "vm_port" {}
variable "vpc_id" {}
variable "network_id" {}

resource "huaweicloud_vpcep_service" "demo" {
  name        = "demo-service"
  server_type = "VM"
  vpc_id      = var.service_vpc_id
  port_id     = var.vm_port

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
}

resource "huaweicloud_vpcep_endpoint" "demo" {
  service_id  = huaweicloud_vpcep_service.demo.id
  vpc_id      = var.vpc_id
  network_id  = var.network_id
  enable_dns  = true
  description = "test description"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the VPC endpoint. If omitted, the provider-level
  region will be used. Changing this creates a new VPC endpoint.

* `service_id` - (Required, String, ForceNew) Specifies the ID of the VPC endpoint service.
  The VPC endpoint service could be private or public. Changing this creates a new VPC endpoint.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the VPC where the VPC endpoint is to be created. Changing
  this creates a new VPC endpoint.

* `network_id` - (Optional, String, ForceNew) Specifies the network ID of the subnet in the VPC specified by `vpc_id`.
  Changing this creates a new VPC endpoint.

  -> This field is required when creating a VPC endpoint for connecting an interface VPC endpoint service.

* `ip_address` - (Optional, String, ForceNew) Specifies the IP address for accessing the associated VPC endpoint
  service. Only IPv4 addresses are supported. Changing this creates a new VPC endpoint.

* `enable_dns` - (Optional, Bool, ForceNew) Specifies whether to create a private domain name. The default value is
  true. Changing this creates a new VPC endpoint.

  -> This field is valid only when creating a VPC endpoint for connecting an interface VPC endpoint service.

* `description` - (Optional, String, ForceNew) Specifies the description of the VPC endpoint.

  Changing this creates a new VPC endpoint.

* `routetables` - (Optional, List, ForceNew) Specifies the IDs of the route tables associated with the VPC endpoint.
  Changing this creates a new VPC endpoint.

  -> This field is valid only when creating a VPC endpoint for connecting a gateway VPC endpoint service.
    The default route table will be used when this field is not specified.

* `enable_whitelist` - (Optional, Bool) Specifies whether to enable access control. The default value is
  false.

* `whitelist` - (Optional, List) Specifies the list of IP address or CIDR block which can be accessed to the
  VPC endpoint. This field is valid when `enable_whitelist` is set to **true**. The max length of whitelist is 20.

* `tags` - (Optional, Map) The key/value pairs to associate with the VPC endpoint.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of the VPC endpoint.

* `status` - The status of the VPC endpoint. The value can be **accepted**, **pendingAcceptance** or **rejected**.

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
