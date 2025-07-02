---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_database_privilege"
description: |-
  Manages an RDS PostgreSQL database privilege resource within HuaweiCloud.
---

# huaweicloud_rds_pg_database_privilege

Manages an RDS PostgreSQL database privilege resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "db_name" {}
variable "user_name_1" {}
variable "user_name_2" {}
variable "schema_name_1" {}
variable "schema_name_2" {}

resource "huaweicloud_rds_pg_database_privilege" "test" {
  instance_id = var.instance_id
  db_name     = var.db_name

  users {
    name        = var.user_name_1
    readonly    = true
    schema_name = var.schema_name_1
  }

  users {
    name        = var.user_name_2
    readonly    = false
    schema_name = var.schema_name_2
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the RDS database privilege resource. If omitted,
  the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the RDS instance ID.

* `db_name` - (Required, String, NonUpdatable) Specifies the database name.

* `users` - (Required, List) Specifies the account that associated with the database.
  The [users](#users_struct) structure is documented below.

<a name="users_struct"></a>
The `users` block supports:

* `name` - (Required, String) Specifies the username of the database account.

* `readonly` - (Required, Bool) Specifies the read-only permission. The value can be:
  + **true**: indicates the read-only permission.
  + **false**: indicates the read and write permission.

* `schema_name` - (Required, String) Specifies the name of the schema.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of database privilege which is formatted `<instance_id>/<db_name>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.
