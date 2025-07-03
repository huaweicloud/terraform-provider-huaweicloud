---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_sqlserver_account"
description: ""
---

# huaweicloud_rds_sqlserver_account

Manages RDS SQLServer account resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_sqlserver_account" "test" {
  instance_id = var.instance_id
  name        = "test_account_name"
  password    = "Test@12345678"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS SQLServer instance.

* `name` - (Required, String, NonUpdatable) Specifies the username of the DB account. The username consists of 1 to 128
  characters and must be different from system usernames. System users include **rdsadmin**, **rdsuser**, **rdsbackup**,
  and **rdsmirror**.

* `password` - (Required, String) Specifies the password of the DB account. It consists of 8 to 128 characters and
  contains at least three types of the following characters: uppercase letters, lowercase letters, digits, and special
  characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of account which is formatted `<instance_id>/<name>`.

* `state` - Indicates the DB user status. Its value can be any of the following:
  + **unavailable**: The database user is unavailable.
  + **available**: The database user is available.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The RDS sqlserver account can be imported using the `instance_id` and `name` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_rds_sqlserver_account.test <instance_id>/<name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `password`. It is generally recommended
running `terraform plan` after importing a RDS SQLServer account. You can then decide if changes should be applied to
the RDS SQLServer account, or the resource definition should be updated to align with the RDS SQLServer account. Also
you can ignore changes as below.

```hcl
resource "huaweicloud_rds_sqlserver_account" "test" {
    ...

  lifecycle {
    ignore_changes = [
      password,
    ]
  }
}
```
