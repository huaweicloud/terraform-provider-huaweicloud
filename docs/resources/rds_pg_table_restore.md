---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_table_restore"
description: |-
  Manages table-level point-in-time-recovery for PostgreSQL RDS instances.
---

# huaweicloud_rds_pg_table_restore

Manages table-level point-in-time-recovery for PostgreSQL RDS instances.

## Example Usage

```hcl
variable "instance_id" {}
variable "database" {}
variable "schema" {}
variable "old_name" {}
variable "new_name" {}

data "huaweicloud_rds_restore_time_ranges" "test" {
  instance_id = var.instance_id
}

resource "huaweicloud_rds_pg_table_restore" "test" {
  instances {
    instance_id  = var.instance_id
    restore_time = data.huaweicloud_rds_restore_time_ranges.test.restore_time[0].start_time  

    databases {
      database = var.database

      schemas {
        schema = var.schema

        tables {
          old_name = var.old_name
          new_name = var.new_name
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the rds instance resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `instances` - (Required, List, ForceNew) A list of RDS PostgreSQL instances where the table restore operation will   be performed.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of RDS MySQL instance.

  Changing this creates a new resource.

* `restore_time` - (Required, Int, ForceNew) Specifies the restoration time point. A timestamp in milliseconds is used.

  Changing this creates a new resource.

* `is_fast_restore` - (Optional, Bool, ForceNew) Specifies whether to use fast restoration.

  Changing this creates a new resource.

* `databases` - (Required, List, ForceNew) Specifies the tables that will be restored.
  The [databases](#databases_struct) structure is documented below.

  Changing this creates a new resource.

<a name="databases_struct"></a>
The `databases` block supports:

* `database` - (Required, String, ForceNew) Specifies the database name.

  Changing this creates a new resource.

* `tables` - (Required, List, ForceNew) Specifies the tables.
  The [tables](#tables_struct) structure is documented below.

  Changing this creates a new resource.

<a name="tables_struct"></a>
The `tables` block supports:

* `old_name` - (Required, String, ForceNew) Specifies the name of the table before restoration.

  Changing this creates a new resource.

* `new_name` - (Required, String, ForceNew) Specifies the name of the table after restoration.

  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID. The value is the restore job ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
