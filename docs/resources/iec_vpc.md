---
subcategory: "Intelligent EdgeCloud (IEC)"
---

# huaweicloud\_iec\_vpc

Manages a VPC resource within HuaweiCloud IEC.

## Example Usage

```hcl
variable "iec_vpc_name" {
  default = "iec-vpc-test"
}

variable "iec_vpc_cidr" {
  default = "192.168.0.0/16"
}

resource "huaweicloud_iec_vpc" "vpc_test" {
  name = var.iec_vpc_name
  cidr = var.iec_vpc_cidr
  mode = "SYSTEM"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of the IEC VPC. The name must be unique.
    The value is a string of no more than 64 characters and can contain digits, 
    letters, underscores (_), and hyphens (-).

* `cidr` - (Required, String) The range of available subnets in the VPC. 
    The ranges of *SYSTEM* mode is from 10.0.0.0/8 to 10.255.0.0/16, 
    172.16.0.0/12 to 172.31.0.0/16, or 192.168.0.0/16.
    The ranges of *CUSTOMER* mode is from 10.0.0.0/8 to 10.255.255.0/24, 
    172.16.0.0/12 to 172.31.255.0/24, or 192.168.0.0/16 to 192.168.255.0/24.

* `mode` - (Optional, String) Specifies the mode of the iec vpc. "SYSTEM" and 
    "CUSTOMER" are supported.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` -  ID of the VPC.

* `subnet_num` - Specifies the number of the subnets. 

## Timeouts

This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `delete` - Default is 3 minute.

## Import

VPCs can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_iec_vpc.vpc_test 5741168f-437a-11eb-b721-fa163e8ac569
```
