---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_instance_no_index_tables"
description: |-
  Use this data source to query the tables without index of an RDS instance.
---

# huaweicloud_rds_instance_no_index_tables

Use this data source to query the tables without index of an RDS instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_instance_no_index_tables" "test" {
  instance_id = var.instance_id
  table_type  = "no_primary_key"
  newest      = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `table_type` - (Required, String) Specifies the type of the table. Value options: **no_primary_key**.

* `newest` - (Required, Bool) Specifies whether the query should focus on retrieving the latest unindexed table.
  Value options:
  + **true**: The query retrieves the latest table without an index.
  + **false**: The latest unindexed table is not retrieved during the query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tables` - Indicates the list of tables without index of the instance.

  The [tables](#tables_struct) structure is documented below.

* `last_diagnose_timestamp` - Indicates the last diagnose timestamp.

<a name="tables_struct"></a>
The `tables` block supports:

* `db_name` - Indicates the database name.

* `schema_name` - Indicates the schema name.

* `table_name` - Indicates the table namee.
