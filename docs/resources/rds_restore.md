---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_restore"
description: |-
  Manages an RDS instance restore resource within HuaweiCloud.
---

# huaweicloud_rds_restore

Manages an RDS instance restore resource within HuaweiCloud.

## Example Usage

### restore by backup_id

```hcl
variable "target_instance_id" {}
variable "source_instance_id" {}
variable "backup_id" {}

resource "huaweicloud_rds_restore" "test" {
  target_instance_id = var.target_instance_id
  source_instance_id = var.source_instance_id
  type               = "backup"
  backup_id          = var.backup_id
}
```

### restore by timestamp

```hcl
variable "target_instance_id" {}
variable "source_instance_id" {}

resource "huaweicloud_rds_restore" "test" {
  target_instance_id = var.target_instance_id
  source_instance_id = var.source_instance_id
  type               = "timestamp"
  restore_time       = 1673852043000
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the rds instance resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `target_instance_id` - (Required, String, NonUpdatable) Specifies the target instance ID.

* `source_instance_id` - (Required, String, NonUpdatable) Specifies the source instance ID.

* `type` - (Optional, String, NonUpdatable) Specifies the restoration type. Value options:
  + **backup**: indicates using backup files for restoration.
  + **timestamp**: indicates the point-in-time restoration mode.

* `backup_id` - (Optional, String, NonUpdatable) Specifies the ID of the backup to be restored. This parameter must be
  specified when `type` is set to **backup** or left empty.

* `restore_time` - (Optional, Int, NonUpdatable) Specifies the time point of data restoration in the UNIX timestamp format.
  The unit is millisecond and the time zone is UTC. This parameter must be specified when `type` is set to **timestamp**.

* `database_name` - (Optional, Map, NonUpdatable) Specifies the databases that will be restored. This parameter applies only
  to the SQL Server DB engine. The key is the old database name, the value is the new database name. If this parameter is
  specified, you can restore all or specific databases and rename new databases. If this parameter is not specified, all
  databases are restored by default. You can enter multiple new database names and separate them with commas (,). The new
  database names can contain but cannot be the same as the original database names. Note the following when you are
  specifying new database names:
  + New database names must be different from the original database names. If they are left blank, the original database
    names will be used for restoration by default.
  + The case-sensitivity settings of the new databases are the same as those of the original databases. Make sure the new
    database names are unique.
  + The total number of new and existing databases on the existing or original DB instances where data is restored cannot
    exceed the database quota specified by **rds_databases_quota**.
  + New database names cannot contain the following fields (case-insensitive): **rdsadmin**, **master**, **msdb**,
    **tempdb**, **model** and **resource**.
  + New database names must consist of `1` to `64` characters, including only letters, digits, underscores (_), and
    hyphens (-). If you want to restore data to multiple new databases, separate them with commas (,).
  + New database names must be different from any database names on the original DB instance.
  + New database names must be different from any database names on the existing or original DB instances where data is
    restored.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID. The value is the restore job ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
