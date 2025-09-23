---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_restored_databases"
description: |-
  Use this data source to get the databases that can be restored to a specified point in time.
---

# huaweicloud_rds_restored_databases

Use this data source to get the databases that can be restored to a specified point in time.

## Example Usage

```hcl
variable "instance_id" {}
variable "restore_time" {}

data "huaweicloud_rds_restored_databases" "test" {
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
  **postgresql**, **mysql**.

* `instance_ids` - (Required, List) Specifies the RDS instance IDs.

* `restore_time` - (Required, String) Specifies the restoration time point. A timestamp in milliseconds is used.

* `instance_name_like` - (Optional, String) Specifies the instance name, which can be used for fuzzy query.

* `database_name_like` - (Optional, String) Specifies the database name, which can be used for fuzzy query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `database_limit` - Indicates the maximum number of databases that can be restored for a single instance. If the number
  of databases queried exceeds this limit, only the databases within this limit are returned to the response.

* `table_limit` - Indicates the maximum number of tables in all databases that can be restored for a single instance. If
  the number of tables queried exceeds this limit, only the databases whose total number of tables is within this limit
  are returned to the response.

* `instances` - Indicates the instance information.

  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `id` - Indicates the instance ID.

* `name` - Indicates the instance name.

* `total_tables` - Indicates the total number of tables in all restorable databases of the instance.

* `databases` - Indicates the database information.

  The [databases](#instances_databases_struct) structure is documented below.

<a name="instances_databases_struct"></a>
The `databases` block supports:

* `name` - Indicates the database name. Databases whose names contain Chinese characters will be filtered out and cannot
  be restored.

* `total_tables` - Indicates the total number of tables in the database.
