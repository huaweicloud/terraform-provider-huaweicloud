---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_mysql_database_table_restore"
description: |-
  Manages an RDS instance MySQL database table restore resource within HuaweiCloud.
---

# huaweicloud_rds_mysql_database_table_restore

Manages an RDS instance MySQL database table restore resource within HuaweiCloud.

## Example Usage

### MySQL databases restore

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_mysql_database_table_restore" "test" {
  restore_time = 1673852043000
  instance_id  = var.instance_id

  databases {
    old_name = "test111"
    new_name = "test111_update"
  }
}
```

### MySQL tables restore

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_mysql_database_table_restore" "test" {
  restore_time = 1673852043000
  instance_id  = var.instance_id

  restore_tables {
    database = "test111"
    tables {
      old_name = "table111"
      new_name = "table111_update"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the rds instance resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of RDS MySQL instance.

  Changing this creates a new resource.

* `restore_time` - (Required, Int, ForceNew) Specifies the restoration time point. A timestamp in milliseconds is used.

  Changing this creates a new resource.

* `is_fast_restore` - (Optional, Bool, ForceNew) Specifies whether to use fast restoration.

  Changing this creates a new resource.

* `databases` - (Optional, List, ForceNew) Specifies the databases that will be restored.
  The [databases](#databases_struct) structure is documented below.

  Changing this creates a new resource.

* `restore_tables` - (Optional, List, ForceNew) Specifies the tables that will be restored.
  The [restore_tables](#restore_tables_struct) structure is documented below.

  Changing this creates a new resource.

-> Exactly one of `databases` and `restore_tables` must be set.

<a name="databases_struct"></a>
The `databases` block supports:

* `old_name` - (Required, String, ForceNew) Specifies the name of the database before restoration.

  Changing this creates a new resource.

* `new_name` - (Required, String, ForceNew) Specifies the name of the database after restoration.

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
