---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_v1"
sidebar_current: "docs-huaweicloud-resource-vpc-v1"
description: |-
  Manages a V1 VPC resource within HuaweiCloud.
---

# huaweicloud_vpc_v1

Manages a VPC resource within HuaweiCloud.

## Example Usage

```hcl

variable "vpc_name" {
  default = "huaweicloud_vpc"
}

variable "vpc_cidr" {
  default = "192.168.0.0/16"
}

resource "huaweicloud_vpc_v1" "vpc_v1" {
  name = "${var.vpc_name}"
  cidr = "${var.vpc_cidr}"
}

```

## Argument Reference

The following arguments are supported:

* `cidr` - (Required) The range of available subnets in the VPC. The value ranges from 10.0.0.0/8 to 10.255.255.0/24, 172.16.0.0/12 to 172.31.255.0/24, or 192.168.0.0/16 to 192.168.255.0/24.

* `region` - (Optional) The region in which to obtain the V1 VPC client. A VPC client is needed to create a VPC. If omitted, the region argument of the provider is used. Changing this creates a new VPC.

* `name` - (Required) The name of the VPC. The name must be unique for a tenant. The value is a string of no more than 64 characters and can contain digits, letters, underscores (_), and hyphens (-). Changing this updates the name of the existing VPC.



## Attributes Reference

The following attributes are exported:

* `id` -  ID of the VPC.

* `name` -  See Argument Reference above.

* `cidr` - See Argument Reference above.

* `status` - The current status of the desired VPC. Can be either CREATING, OK, DOWN, PENDING_UPDATE, PENDING_DELETE, or ERROR.

* `shared` - Specifies whether the cross-tenant sharing is supported.

* `region` - See Argument Reference above.

## Import

VPCs can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_vpc_v1.vpc_v1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
