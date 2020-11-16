---
subcategory: "Virtual Private Cloud (VPC)"
---

# huaweicloud\_vpc

Manages a VPC resource within HuaweiCloud.
This is an alternative to `huaweicloud_vpc_v1`

## Example Usage

```hcl
variable "vpc_name" {
  default = "huaweicloud_vpc"
}

variable "vpc_cidr" {
  default = "192.168.0.0/16"
}

resource "huaweicloud_vpc" "vpc" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc" "vpc_with_tags" {
  name = var.vpc_name
  cidr = var.vpc_cidr

  tags = {
    foo = "bar"
    key = "value"
  }
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the vpc resource. If omitted, the provider-level region will work as default. Changing this creates a new resource. Changing this creates a new vpc resource.

* `cidr` - (Required) The range of available subnets in the VPC. The value ranges from 10.0.0.0/8 to 10.255.255.0/24, 172.16.0.0/12 to 172.31.255.0/24, or 192.168.0.0/16 to 192.168.255.0/24.

* `name` - (Required) The name of the VPC. The name must be unique for a tenant. The value is a string of no more than 64 characters and can contain digits, letters, underscores (_), and hyphens (-). Changing this updates the name of the existing VPC.

* `tags` - (Optional) The key/value pairs to associate with the vpc.

* `enterprise_project_id` - (Optional) The enterprise project id of the vpc. Changing this creates a new vpc.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the VPC.

* `name` -  See Argument Reference above.

* `cidr` - See Argument Reference above.

* `status` - The current status of the desired VPC. Can be either CREATING, OK, DOWN, PENDING_UPDATE, PENDING_DELETE, or ERROR.

* `shared` - Specifies whether the cross-tenant sharing is supported.

* `routes` - Specifies the route information. Structure is documented below.

* `region` - See Argument Reference above.

* `tags` - See Argument Reference above.

The `routes` block contains:

* `destination` - Specifies the destination network segment of a route.

* `nexthop` - Specifies the next hop of a route.

## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `delete` - Default is 3 minute.

## Import

VPCs can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_vpc.vpc_v1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
