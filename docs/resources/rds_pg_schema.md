---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_schema"
description: |-
  Manages an RDS PostgreSQL schema resource within HuaweiCloud.
---

# huaweicloud_rds_pg_schema

Manages an RDS PostgreSQL schema resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_pg_schema" "test" {
  instance_id = var.instance_id
  db_name     = "test_db"
  schema_name = "test_schema"
  owner       = "test_account"
}
```

~> Deleting RDS PostgreSQL schema resource is not supported, it will only be removed from the state.

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS PostgreSQL instance.

* `db_name` - (Required, String, NonUpdatable) Specifies the database name. The value contains 1 to 63 characters,
  including letters, digits, and underscores (_). It cannot start with pg or a digit, and must be different from RDS for
  PostgreSQL template library names. RDS for PostgreSQL template libraries include **postgres**, **template0**, and
  **template1**.

* `schema_name` - (Required, String, NonUpdatable) Specifies the schema name. The value contains 1 to 63 characters,
  including letters, digits, and underscores (_). It cannot start with pg or a digit, and must be different from RDS for
  PostgreSQL template database names and existing schema names.

* `owner` - (Required, String, NonUpdatable) Specifies the database owner. The value contains 1 to 63 characters. It
  cannot start with pg or a digit, and must be different from system usernames. System users include **rdsAdmin**,
  **rdsMetric**, **rdsBackup**, **rdsRepl**, **rdsProxy** and **rdsDdm**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is formatted `<instance_id>/<db_name>/<schema_name>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

## Import

The RDS postgresql schema can be imported using the `instance_id`,  `db_name` and `schema_name` separated by slashes, e.g.

```bash
$ terraform import huaweicloud_rds_pg_schema.test <instance_id>/<db_name>/<schema_name>
```
