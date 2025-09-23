---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_account_privileges"
description: |-
  Manages an RDS PostgreSQL account privileges resource within HuaweiCloud.
---

# huaweicloud_rds_pg_account_privileges

Manages an RDS PostgreSQL account privileges resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "user_name" {}

resource "huaweicloud_rds_pg_account_privileges" "test" {
  instance_id            = var.instance_id
  user_name              = var.user_name
  role_privileges        = ["CREATEROLE","LOGIN","REPLICATION"]
  system_role_privileges = ["pg_signal_backend","root"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS PostgreSQL instance.

* `user_name` - (Required, String, NonUpdatable) Specifies the username of the account.

* `role_privileges` - (Optional, List) Specifies the list of role privileges. Value options: **CREATEDB**,
  **CREATEROLE**, **LOGIN**, **REPLICATION**.

* `system_role_privileges` - (Optional, List) Specifies the list of system role privileges. Value options:
  **pg_monitor**, **pg_signal_backend**, **root**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of account which is formatted `<instance_id>/<user_name>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The RDS PostgreSQL privileges roles can be imported using the `instance_id` and `user_name` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_rds_pg_account_privileges.test <instance_id>/<user_name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `system_role_privileges`. It is generally
recommended running `terraform plan` after importing the RDS PostgreSQL account privileges. You can then decide if
changes should be applied to the RDS PostgreSQL account privileges, or the RDS PostgreSQL account privileges definition
should be updated to align with the account. Also you can ignore changes as below.

```hcl
resource "huaweicloud_rds_pg_account_privileges" "test" {
    ...

  lifecycle {
    ignore_changes = [
      system_role_privileges,
    ]
  }
}
```
