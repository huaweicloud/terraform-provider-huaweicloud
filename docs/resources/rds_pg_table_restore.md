---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_table_restore"
description: |-
  Manages an RDS instance PostgreSQL table restore resource within HuaweiCloud.
---

# huaweicloud_rds_pg_table_restore

Manages an RDS instance PostgreSQL table restore resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_pg_table_restore" "test" {
  instance_id  = var.instance_id
  restore_time = 1754954459000

  databases {
    database = "test1"

    schemas {
      schema = "test1"

      tables {
        old_name = "table1"
        new_name = "table1_update"
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which the RDS instance exists. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of RDS PostgreSQL instance.

* `restore_time` - (Required, Int, NonUpdatable) Specifies the restoration time point. A timestamp in milliseconds is used.

* `databases` - (Required, List, NonUpdatable) Specifies the list of databases to restore.  
  The [databases](#databases_struct) structure is documented below.

<a name="databases_struct"></a>
The `databases` block supports:

* `database` - (Required, String, NonUpdatable) Specifies the name of the database that contains the tables to restore.

* `schemas` - (Required, List, NonUpdatable) Specifies a list of schemas within the database.
  The [schemas](#schemas_struct) structure is documented below.

<a name="schemas_struct"></a>
The `schemas` block supports:

* `schema` - (Required, String, NonUpdatable) Specifies the name of the schema containing the tables to be restored.

* `tables` - (Required, List, NonUpdatable) Specifies a list of tables to be restored.
  The [tables](#tables_struct) structure is documented below.

<a name="tables_struct"></a>
The `tables` block supports:

* `old_name` - (Required, String, NonUpdatable) Specifies the name of the table before restoration.

* `new_name` - (Required, String, NonUpdatable) Specifies the name of the table after restoration.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID. The value is the restore job ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
