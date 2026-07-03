---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_instance_db_table_info"
description: |
  Use this data source to query database and table information of a TaurusDB instance within Huaweicloud.
---

# huaweicloud_taurusdb_instance_db_table_info

Use this data source to query database and table information of a TaurusDB instance within Huaweicloud.

## Example Usage

### Query database names of an instance

```hcl
variable "instance_id" {}

data "huaweicloud_taurusdb_instance_db_table_info" "test" {
  instance_id = var.instance_id
}
```

### Query table names of a specific database

```hcl
variable "instance_id" {}
variable "database_name" {}

data "huaweicloud_taurusdb_instance_db_table_info" "test" {
  instance_id   = var.instance_id
  database_name = var.database_name
}
```

### Query table metadata of a specific table

```hcl
variable "instance_id" {}
variable "database_name" {}
variable "table_name" {}

data "huaweicloud_taurusdb_instance_db_table_info" "test" {
  instance_id   = var.instance_id
  database_name = var.database_name
  table_name    = var.table_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the resource. If omitted, the provider-level region
  will be used.

* `instance_id` - (Required, String) Specifies the ID of the TaurusDB instance.

* `database_name` - (Optional, String) Specifies the database name. If specified, `table_names` will be returned.
  If not specified, `database_names` will be returned.

* `table_name` - (Optional, String) Specifies the table name. If specified along with `database_name`,
  `table_meta_infos` will be returned. This parameter requires `database_name` to be set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `database_names` - Indicates the list of database names. This attribute is populated only when `database_name`
  is not specified.

* `table_names` - Indicates the list of table names. This attribute is populated only when `database_name`
  is specified and `table_name` is not specified.

* `table_meta_infos` - Indicates the list of table metadata information. This attribute is populated only when
  both `database_name` and `table_name` are specified.
  The [table_meta_infos](#table_meta_infos_struct) structure is documented below.

<a name="table_meta_infos_struct"></a>
The `table_meta_infos` block contains:

* `column_name` - Indicates the column name.

* `column_type` - Indicates the data type of a column.

* `column_key` - Indicates whether the column is an index column.

* `column_default` - Indicates the default value of a column.

* `is_nullable` - Indicates whether the column data can be NULL.

* `extra` - Indicates extra information, for example, whether the column is an auto-increment column.
