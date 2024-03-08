---
subcategory: "Cloud Bastion Host (CBH)"
---

# huaweicloud_cbh_instance

Manages a CBH instance resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}
variable "password" {}
variable "availability_zone" {}

resource "huaweicloud_cbh_instance" "test" {
  flavor_id         = "cbh.basic.10"
  name              = var.name
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.security_group_id
  availability_zone = var.availability_zone
  password          = var.password
  charging_mode     = "prePaid"
  period_unit       = "month"
  period            = 1
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the CBH instance. The field can contain `1` to `64` characters.
  Only letters, digits, underscores (_), and hyphens (-) are allowed.

  Changing this parameter will create a new resource.

* `flavor_id` - (Required, String, ForceNew) Specifies the product ID of the CBH server.

  Changing this parameter will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of a VPC.

  Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of a subnet.

  Changing this parameter will create a new resource.

* `security_group_id` - (Required, String, ForceNew) Specifies the ID of the security group.

  Changing this parameter will create a new resource.

* `availability_zone` - (Required, String, ForceNew) Specifies the availability zone name.

  Changing this parameter will create a new resource.

* `password` - (Required, String) Specifies the password for logging in to the management console. The value of the field
  has the following restrictions:
  + The value of the field must contain `8` to `32` characters.
  + The value of the field must contain at least three of the following: letters, digits, and special characters
    (!@$%^-_=+[{}]:,./?~#*).
  + The value of the field cannot contain the username or the username spelled backwards.

* `charging_mode` - (Required, String, ForceNew) Specifies the charging mode of the CBH instance.
  The options are as follows:
  + **prePaid**: the yearly/monthly billing mode.

  Changing this parameter will create a new resource.

* `period_unit` - (Required, String, ForceNew) Specifies the charging period unit of the instance.
  Valid values are *month* and *year*.

  Changing this parameter will create a new resource.

* `period` - (Required, Int, ForceNew) Specifies the charging period of the CBH instance.
  If `period_unit` is set to **month**, the value ranges from 1 to 9.
  If `period_unit` is set to **year**, the value ranges from 1 to 3.

  Changing this parameter will create a new resource.

* `auto_renew` - (Optional, String) Specifies whether auto-renew is enabled.
  Valid values are **true** and **false**. Defaults to **false**.

* `subnet_address` - (Optional, String) Specifies the IP address of the subnet.
  If not specified, a new IP address will be assigned.

* `public_ip_id` - (Optional, String) Specifies the ID of the elastic IP.

* `ipv6_enable` - (Optional, Bool, ForceNew) Specifies whether the IPv6 network is enabled. Defaults to **false**.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `public_ip` - Indicates the elastic IP address.

* `private_ip` - Indicates the private IP address of the instance.

* `status` - Indicates the status of the instance.

* `version` - Indicates the current version of the instance image.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `delete` - Default is 30 minutes.

## Import

The CBH instance can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_cbh_instance.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `charging_mode`, `period`, `period_unit`,
`auto_renew`, `password`, `ipv6_enable`.
It is generally recommended running `terraform plan` after importing an instance.
You can then decide if changes should be applied to the instance, or the resource definition should be updated
to align with the instance. Also, you can ignore changes as below.

```
resource "huaweicloud_cbh_instance" "test" {
    ...

  lifecycle {
    ignore_changes = [
      charging_mode, period, period_unit, auto_renew, password, ipv6_enable,
    ]
  }
}
```
