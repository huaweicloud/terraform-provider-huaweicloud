---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_sqlserver_database"
description: ""
---

# huaweicloud_rds_sqlserver_database

Manages RDS SQLServer database resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_sqlserver_database" "test" {
  instance_id = var.instance_id
  name        = "test"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS SQLServer instance.

* `name` - (Required, String, NonUpdatable) Specifies the database name. The database name can contain 1 to 64 characters,
  and can include letters, digits, hyphens (-), underscores (_), and periods (.). It cannot start or end with an RDS for
  SQL Server system database name. RDS for SQL Server system databases include **master**, **msdb**, **model**,
  **tempdb**, **resource**, and **rdsadmin**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of database which is formatted `<instance_id>/<name>`.

* `character_set` - Indicates the character set used by the database.

* `state` - Indicates the database status. Its value can be any of the following:
  + **Creating**: The database is being created.
  + **Running**: The database is running.
  + **Deleting**: The database is being deleted.
  + **Not Exists**: The database does not exist.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The RDS sqlserver database can be imported using the `instance_id` and `name` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_rds_sqlserver_database.test <instance_id>/<name>
```
