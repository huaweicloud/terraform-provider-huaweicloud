---
subcategory: "Cloud Application Engine (CAE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cae_vpc_egress"
description: |-
  Manage a CAE environment to access VPC resource within HuaweiCloud.
---

# huaweicloud_cae_vpc_egress

Manage a CAE environment to access VPC resource within HuaweiCloud.

## Example Usage

```hcl
variable "environment_id" {}
variable "route_table_id" {}
variable "cidr" {}

resource "huaweicloud_cae_vpc_egress" "test" {
  environment_id = var.environment_id
  route_table_id = var.route_table_id
  cidr           = var.cidr
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `environment_id` - (Required, String, ForceNew) Specifies the ID of the CAE environment.
  Changing this creates a new resource.

* `route_table_id` - (Required, String, ForceNew) Specifies the ID of the route table corresponding to the subnet to which
  the CAE environment belongs.  
  Changing this creates a new resource.

* `cidr` - (Required, String, ForceNew) Specifies the destination CIDR of the routing table corresponding to the subnet
  to which the CAE environment belongs.  
  Changing this creates a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the ID of the enterprise project to which the
  VPC egress belongs.  
  Changing this creates a new resource.

-> If the `environment_id` belongs to the non-default enterprise project, this parameter is required and is only valid
  for enterprise users.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, in UUID format.

## Import

The resource can be imported using `environment_id`, `route_table_id`, and `cidr`, separated by commas (,), e.g.

```bash
$ terraform import huaweicloud_cae_vpc_egress.test <environment_id>,<route_table_id>,<cidr>
```

For the VPC egress with the non-default enterprise project ID, its enterprise project ID need to be specified
additionanlly when importing. All fields are separated by commas (,), e.g.

```bash
$ terraform import huaweicloud_cae_vpc_egress.test <environment_id>,<route_table_id>,<cidr>,<enterprise_project_id>
```
