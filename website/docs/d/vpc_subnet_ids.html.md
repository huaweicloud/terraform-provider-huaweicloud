---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_subnet_ids"
sidebar_current: "docs-huaweicloud-datasource-subnet-ids"
description: |-
  Provides a list of subnet Ids for a VPC
---

# huaweicloud\_vpc\_subnet\_ids

Provides a list of subnet ids for a vpc_id
This is an alternative to `huaweicloud_vpc_subnet_ids_v1`

This resource can be useful for getting back a list of subnet ids for a vpc.

## Example Usage

The following example shows outputing all cidr blocks for every subnet id in a vpc.

```hcl
data "huaweicloud_vpc_subnet_ids" "subnet_ids" {
  vpc_id = var.vpc_id
}

data "huaweicloud_vpc_subnet" "subnet" {
  count = length(data.huaweicloud_vpc_subnet_ids.subnet_ids.ids)
  id    = tolist(data.huaweicloud_vpc_subnet_ids.subnet_ids.ids)[count.index]
}

output "subnet_cidr_blocks" {
  value = [for s in data.huaweicloud_vpc_subnet.subnet: "${s.name}: ${s.id}: ${s.cidr}"]
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` (Required) - Specifies the VPC ID used as the query filter.

## Attributes Reference

The following attributes are exported:

* `ids` - A set of all the subnet ids found. This data source will fail if none are found.
