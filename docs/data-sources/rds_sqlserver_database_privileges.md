---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_sqlserver_database_privileges"
description: ""
---

# huaweicloud_rds_sqlserver_database_privileges

Use this data source to get the list of RDS SQLServer database privileges.

## Example Usage

```hcl
var "instance_id" {}
var "db_name" {}

data "huaweicloud_rds_sqlserver_database_privileges" "test" {
  instance_id = var.instance_id
  db_name     = var.db_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS SQLServer instance.

* `db_name` - (Required, String) Specifies the database name.

* `user_name` - (Optional, String) Specifies the username of the database account.

* `readonly` - (Optional, Bool) Specifies whether the database permission is **read-only**. Values option:
  + **true**: indicates the read-only permission.
  + **false**: indicates the read and write permission.

  Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `users` - The list of the accounts that accociated with the database.
  The [users](#RDS_sqlserver_database_privileges) structure is documented below.

<a name="RDS_sqlserver_database_privileges"></a>
The `users` block supports:

* `name` - The username of the database account.

* `readonly` - The read-only database permission.
