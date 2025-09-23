---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_restored_tables"
description: |-
  Use this data source to get the available tables for table-level point-in-time recovery.
---

# huaweicloud_gaussdb_mysql_restored_tables

Use this data source to get the available tables for table-level point-in-time recovery.

## Example Usage

```hcl
variable "instance_id" {}
variable "restore_time" {}

data "huaweicloud_gaussdb_mysql_restored_tables" "test" {
  instance_id     = var.instance_id
  restore_time    = var.restore_time
  last_table_info = "true"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB MySQL instance,

* `restore_time` - (Required, String) Specifies the backup time, in timestamp format.

* `last_table_info` - (Required, String) Specifies  whether data is restored to the most recent table.
  + **true**: most recent table.
  + **false**: time-specific table

* `database_name` - (Optional, String) Specifies the database name, which is used for fuzzy match.

* `table_name` - (Optional, String) Specifies the table name, which is used for fuzzy match.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `databases` - Indicates the database information.

  The [databases](#databases_struct) structure is documented below.

<a name="databases_struct"></a>
The `databases` block supports:

* `name` - Indicates the database name.

* `total_tables` - Indicates the total number of tables.

* `tables` - Indicates the table information.

  The [tables](#databases_tables_struct) structure is documented below.

<a name="databases_tables_struct"></a>
The `tables` block supports:

* `name` - Indicates the table name.
