---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_mysql_account"
description: ""
---

# huaweicloud_rds_mysql_account

Manages RDS Mysql account resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "account_password" {}

resource "huaweicloud_rds_mysql_account" "test" {
  instance_id = var.instance_id
  name        = "test"
  password    = var.account_password
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the rds account resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the rds instance id.

* `name` - (Required, String, NonUpdatable) Specifies the username of the db account. Only lowercase letters, digits,
  hyphens (-), and underscores (_) are allowed.
  + If the database version is MySQL 5.6, the username consists of 1 to 16 characters.
  + If the database version is MySQL 5.7 or 8.0, the username consists of 1 to 32 characters.

* `password` - (Required, String) Specifies the password of the db account. The parameter must be 8 to 32 characters
  long and contain only letters(case-sensitive), digits, and special characters(~!@#$%^*-_=+?,()&). The value must be
  different from `name` or `name` spelled backwards.

* `hosts` - (Optional, List, NonUpdatable) Specifies the IP addresses that are allowed to access your DB instance.
  + If the IP address is set to %, all IP addresses are allowed to access your instance.
  + If the IP address is set to 10.10.10.%, all IP addresses in the subnet 10.10.10.X are allowed to access
    your instance.
  + Multiple IP addresses can be added.

* `description` - (Optional, String) Specifies remarks of the database account. The parameter must be 1 to 512
  characters long and is supported only for MySQL 8.0.25 and later versions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of account which is formatted `<instance_id>/<name>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

RDS account can be imported using the `instance_id` and `name` separated by a slash, e.g.:

```bash
$ terraform import huaweicloud_rds_mysql_account.account_1 <instance_id>/<name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `password`. It is generally recommended
running `terraform plan` after importing the RDS Mysql account. You can then decide if changes should be applied to
the RDS Mysql account, or the resource definition should be updated to align with the RDS Mysql account. Also you
can ignore changes as below.

```hcl
resource "huaweicloud_rds_mysql_account" "account_1" {
    ...

  lifecycle {
    ignore_changes = [
      password
    ]
  }
}
```
