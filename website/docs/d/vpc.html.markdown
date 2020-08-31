---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc"
sidebar_current: "docs-huaweicloud-datasource-vpc"
description: |-
  Get information on a VPC.
---

# huaweicloud\_vpc

Provides details about a specific VPC.
This is an alternative to `huaweicloud_vpc_v1`

## Example Usage

The following example shows how one might accept a VPC id as a variable and use this data source to obtain the data necessary to create a subnet within it.

```hcl

variable "vpc_name" {}

data "huaweicloud_vpc" "vpc" {
  name = "${var.vpc_name}"
}

```

## Argument Reference

The arguments of this data source act as filters for querying the available VPCs in the current region. The given filters must match exactly one VPC whose data will be exported as attributes.


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
