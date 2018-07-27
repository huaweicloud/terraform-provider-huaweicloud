---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_v1"
sidebar_current: "docs-huaweicloud-datasource-vpc-v1"
description: |-
  Get information on an huaweicloud VPC.
---

# huaweicloud_vpc_v1

huaweicloud_vpc_v1 provides details about a specific VPC.

This resource can prove useful when a module accepts a vpc id as an input variable and needs to, for example, determine the CIDR block of that VPC.

## Example Usage

The following example shows how one might accept a VPC id as a variable and use this data source to obtain the data necessary to create a subnet within it.

```hcl

variable "vpc_name" {}

data "huaweicloud_vpc_v1" "vpc" {
  name = "${var.vpc_name}"
}

```

## Argument Reference

The arguments of this data source act as filters for querying the available VPCs in the current region. The given filters must match exactly one VPC whose data will be exported as attributes.

* `region` - (Optional) The region in which to obtain the V1 VPC client. A VPC client is needed to retrieve VPCs. If omitted, the region argument of the provider is used.

* `id` - (Optional) The id of the specific VPC to retrieve.

* `status` - (Optional) The current status of the desired VPC. Can be either CREATING, OK, DOWN, PENDING_UPDATE, PENDING_DELETE, or ERROR.

* `name` - (Optional) A unique name for the VPC. The name must be unique for a tenant. The value is a string of no more than 64 characters and can contain digits, letters, underscores (_), and hyphens (-).

* `cidr` - (Optional) The cidr block of the desired VPC.



## Attributes Reference

The following attributes are exported:

* `id` - ID of the VPC.

* `name` -  See Argument Reference above.

* `status` - See Argument Reference above.

* `cidr` - See Argument Reference above.

* `routes` - The list of route information with destination and nexthop fields.

* `shared` - Specifies whether the cross-tenant sharing is supported.

* `region` - See Argument Reference above.

