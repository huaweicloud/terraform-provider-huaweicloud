---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_database_restore"
description: |-
  Manages an RDS instance PostgreSQL database restore resource within HuaweiCloud.
---

# huaweicloud_rds_pg_database_restore

Manages an RDS instance PostgreSQL database restore resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_pg_database_restore" "test" {
  instance_id  = var.instance_id
  restore_time = 1754954459000

  databases {
    old_name = "database_name"
    new_name = "database_name_update"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of RDS PostgreSQL instance.

* `restore_time` - (Required, Int, NonUpdatable) Specifies the restoration time point.
  A timestamp in milliseconds is used.

* `databases` - (Required, List, NonUpdatable) Specifies databases to restore.  
  The [databases](#databases_struct) structure is documented below.

* `is_fast_restore` - (Optional, Bool, NonUpdatable) Specifies whether to use fast restoration.

<a name="databases_struct"></a>
The `databases` block supports:

* `old_name` - (Required, String, NonUpdatable) Specifies the name of the database before restoration.

* `new_name` - (Required, String, NonUpdatable) Specifies the name of the database after restoration.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID. The value is the restore job ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
