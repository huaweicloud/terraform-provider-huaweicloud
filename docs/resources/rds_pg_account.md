---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_account"
description: ""
---

# huaweicloud_rds_pg_account

Manages RDS PostgreSQL account resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "account_password" {}

resource "huaweicloud_rds_pg_account" "test" {
  instance_id = var.instance_id
  name        = "test_account_name"
  password    = var.account_password
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS PostgreSQL instance.

* `name` - (Required, String, NonUpdatable) Specifies the username of the DB account. The username contains 1 to 63
  characters, including letters, digits, and underscores (_). It cannot start with pg or a digit and must be different
  from system usernames. System users include **rdsAdmin**, **rdsMetric**, **rdsBackup**, **rdsRepl**, **rdsProxy**,
  and **rdsDdm**.

* `password` - (Required, String) Specifies the password of the DB account. The value must be 8 to 32 characters long
  and contain at least three types of the following characters: uppercase letters, lowercase letters, digits, and special
  characters (~!@#%^*-_=+?,). The value cannot contain the username or the username spelled backwards.

* `description` - (Optional, String) Specifies the remarks of the DB account. The parameter must be 1 to 512 characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of account which is formatted `<instance_id>/<name>`.

* `attributes` - Indicates the permission attributes of a user.
  The [attributes](#PgAccount_Attributes) structure is documented below.

<a name="PgAccount_Attributes"></a>
The `attributes` block supports:

* `rol_super` - Indicates whether a user has the super-user permission.

* `rol_inherit` - Indicates whether a user automatically inherits the permissions of the role to which the user belongs.

* `rol_create_role` - Indicates whether a user can create other sub-users.

* `rol_create_db` - Indicates whether a user can create a database.

* `rol_can_login` - Indicates whether a user can log in to the database.

* `rol_conn_limit` - Indicates the maximum number of concurrent connections to a DB instance.

* `rol_replication` - Indicates whether the user is a replication role.

* `rol_bypass_rls` - Indicates whether a user bypasses each row-level security policy.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The RDS PostgreSQL account can be imported using the `instance_id` and `name` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_rds_pg_account.test <instance_id>/<name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `password`. It is generally recommended
running `terraform plan` after importing the RDS PostgreSQL account. You can then decide if changes should be applied to
the RDS PostgreSQL account, or the resource definition should be updated to align with the RDS PostgreSQL account. Also
you can ignore changes as below.

```hcl
resource "huaweicloud_rds_pg_account" "account_1" {
    ...

  lifecycle {
    ignore_changes = [
      password
    ]
  }
}
```
