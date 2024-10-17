---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_table_restore"
description: |-
  Use this resource restore tables to GaussDB MySQL instance within HuaweiCloud.
---

# huaweicloud_gaussdb_mysql_table_restore

Use this resource restore tables to GaussDB MySQL instance within HuaweiCloud.

-> **NOTE:** Deleting restoration record is not supported. If you destroy a resource of restoration record,
the restoration record is only removed from the state, but it remains in the cloud. And the instance doesn't return to
the state before restoration.

## Example Usage

```hcl
variable "instance_id" {}
variable "restore_time" {}
variable "backup_id" {}

resource "huaweicloud_gaussdb_mysql_table_restore" "test" {
  instance_id     = var.instance_id
  restore_time    = var.restore_time
  last_table_info = "true"
  
  restore_tables {
    database = "test_db"
    tables {
      old_name = "table_old_1"
      new_name = "table_new_1"
    }
    tables {
      old_name = "table_old_2"
      new_name = "table_new_2"
    }
  }
  restore_tables {
    database = "test_db1"
    tables {
      old_name = "table_old_1"
      new_name = "table_new_1"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the GaussDB mysql table restore resource. If
  omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the GaussDB mysql instance ID.

  Changing this creates a new resource.

* `restore_time` - (Required, String, ForceNew) Specifies the backup time, in timestamp format.

  Changing this creates a new resource.

* `restore_tables` - (Required, List, ForceNew) Specifies the database information.
  The [restore_tables](#restore_tables_struct) structure is documented below.

  Changing this creates a new resource.

* `last_table_info` - (Optional, String, ForceNew) Specifies whether the data is restored to the most recent table.
  Value options:
  + **true**: most recent table.
  + **false (default value)**: time-specific table.

  Changing this creates a new resource.

<a name="restore_tables_struct"></a>
The `restore_tables` block supports:

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
