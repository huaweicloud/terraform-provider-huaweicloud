---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_internet_gateway"
description: ""
---

# huaweicloud_vpc_internet_gateway

Manages a VPC internet gateway resource within HuaweiCloud.

## Example Usage

```hcl
variable "vpc_id" {}
variable "igw_name" {}

resource "huaweicloud_vpc_internet_gateway" "test" {
  vpc_id    = var.vpc_id
  name      = var.igw_name
  add_route = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the IGW.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID which the IGW associated with. A VPC can only associate
  with one IGW. Changing this creates a new resource.

* `subnet_id` - (Optional, String, ForceNew) Specifies the subnet ID which the IGW associated with.
  Changing this creates a new resource.

* `name` - (Optional, String) Specifies the IGW name.

* `add_route` - (Optional, Bool, ForceNew) Specifies whether to add a default route pointing to the IGW in the default
  route table of the VPC with the destination address 0.0.0.0/0. Changing this creates a new resource.

* `enable_ipv6` - (Optional, Bool) Specifies whether to enable IPv6. It's not allow change true to false. Make sure the
  subnet is enable IPv6 before setting to true.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The create time of the IGW.

* `updated_at` - The update time of the IGW.

## Import

The VPC internet gateway can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_vpc_internet_gateway.test <id>
```
