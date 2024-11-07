---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_restored_tables"
description: |-
  Use this data source to get the tables that can be restored to a specified point in time.
---

# huaweicloud_rds_restored_tables

Use this data source to get the tables that can be restored to a specified point in time.

## Example Usage

```hcl
variable "instance_id" {}
variable "restore_time" {}

data "huaweicloud_rds_restored_tables" "test" {
  engine       = "postgresql"
  instance_ids = [var.instance_id]
  restore_time = var.restore_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `engine` - (Required, String) Specifies the database engine. The supported engines are as follows, not case-sensitive:
  **postgresql**.

* `instance_ids` - (Required, List) Specifies the RDS instance IDs.

* `restore_time` - (Required, String) Specifies the restoration time point. A timestamp in milliseconds is used.

* `instance_name_like` - (Optional, String) Specifies the instance name, which can be used for fuzzy query.

* `database_name_like` - (Optional, String) Specifies the database name, which can be used for fuzzy query.

* `table_name_like` - (Optional, String) Specifies the table name, which can be used for fuzzy query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `table_limit` - Indicates the maximum number of tables that can be restored.

* `instances` - Indicates the instance information.

  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `id` - Indicates the instance ID.

* `name` - Indicates the instance name.

* `total_tables` - Indicates the number of tables that can be restored.

* `databases` - Indicates the database information.

  The [databases](#instances_databases_struct) structure is documented below.

<a name="instances_databases_struct"></a>
The `databases` block supports:

* `name` - Indicates the database name.

* `total_tables` - Indicates the number of tables that can be restored.

* `schemas` - Indicates the schema information.
  The [schemas](#instances_databases_schemas_struct) structure is documented below.

<a name="instances_databases_schemas_struct"></a>
The `schemas` block supports:

* `name` - Indicates the schema name.

* `total_tables` - Indicates the number of tables that can be restored.

* `tables` - Indicates the table information.
  The [tables](#instances_databases_schemas_tables_struct) structure is documented below.

<a name="instances_databases_schemas_tables_struct"></a>
The `tables` block supports:

* `name` - Indicates the table name.
