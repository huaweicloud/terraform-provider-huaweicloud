---
layout: "huaweicloud"
page_title: "Huaweicloud: huaweicloud_cci_network_v1"
sidebar_current: "docs-huaweicloud-resource-cci-network-v1"
description: |-
  Provides Cloud Container Instance(CCI) resource.
---

# huaweicloud_cci_network_v1

Provides a CCI resource.
This is an alternative to `huaweicloud_cci_network_v1`


## Example Usage

 ```hcl
variable "sg_id" { }
variable "project_id" { }
variable "domain_id" { }
variable "vpc_id" { }
variable "net_id" { }
variable "subnet_id" { }
	
resource "huaweicloud_cci_network_v1" "net_1" {
  name           = "cci-net"
  namespace      = "test-ns"
  security_group = var.sg_id
  project_id     = var.project_id
  domain_id      = var.domain_id
  vpc_id         = var.vpc_id
  network_id     = var.net_id
  subnet_id      = var.subnet_id
  available_zone = "cn-north-1a"
  cidr           = "192.168.0.0/24"
}
```

## Argument Reference

The following arguments are supported:


* `name` - (Required) CCI Network name. Changing this parameter will create a new resource.

* `namespace` - (Required) CCI Network namespace. Changing this parameter will create a new resource.

* `security_group` - (Required) ID of the security group to which the subnet of the network belongs. Changing this parameter will create a new resource.

* `project_id` - (Required) Project ID of the tenant. Changing this parameter will create a new resource.

* `domain_id` - (Required) Domain ID of the tenant. Changing this parameter will create a new resource.

* `vpc_id` - (Required) ID of the VPC to which the network belongs. Changing this parameter will create a new resource.

* `network_id` - (Required) Network ID of the VPC subnet in which the network belongs. Changing this parameter will create a new resource.

* `subnet_id` - (Required) ID of the VPC subnet to which the network belongs. Changing this parameter will create a new resource.

* `available_zone` - (Required) AZ to which the VPC subnet of the network belongs. Changing this parameter will create a new resource.

* `cidr` - (Required) Network segment of the VPC subnet to which the network belongs. Changing this parameter will create a new resource.


## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference.

  * `id` -  Id of the instance resource.

