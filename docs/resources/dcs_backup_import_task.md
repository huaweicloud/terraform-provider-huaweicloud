---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_backup_import_task"
description: |-
  Manages a DCS backup import task resource within HuaweiCloud.
---

# huaweicloud_dcs_backup_import_task

Manages a DCS backup import task resource within HuaweiCloud.

## Example Usage

### create backup import task by OBS bucket

```hcl
variable "target_dcs_instance_id" {}

resource "huaweicloud_dcs_backup_import_task" "test" {
  task_name        = "test_task_name"
  migration_type   = "backupfile_import"
  migration_method = "full_amount_migration"
  description      = "terraform test"

  backup_files {
    file_source = "self_build_obs"
    bucket_name = "test-dcs"

    files {
      file_name = "appendonly.aof"
    }
    files {
      file_name = "test_redis_backup.rdb"
    }
  }

  target_instance{
    id       = var.target_dcs_instance_id
    password = "test_1234"
  }
}
```

### create backup import task by backup ID

```hcl
variable "backup_id" {}
variable "target_dcs_instance_id" {}

resource "huaweicloud_dcs_backup_import_task" "test" {
  task_name        = "test_task_name"
  migration_type   = "backupfile_import"
  migration_method = "full_amount_migration"
  description      = "terraform test"

  backup_files {
    file_source = "backup_record"
    backup_id   = var.backup_id
  }

  target_instance{
    id       = var.target_dcs_instance_id
    password = "test_1234"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `task_name` - (Required, String, NonUpdatable) Specifies the backup import task name.

* `migration_type` - (Required, String, NonUpdatable) Specifies the migration mode. Value options:
  + **backupfile_import**: importing backup files

* `migration_method` - (Required, String, NonUpdatable) Specifies the type of the migration. Value options:
  + **full_amount_migration**: full migration
  + **incremental_migration**: incremental migration

* `backup_files` - (Required, List, NonUpdatable) Specifies the backup files to be imported when the migration mode is
  backup file import.
  The [backup_files](#backup_files_struct) structure is documented below.

* `target_instance` - (Required, List, NonUpdatable) Specifies the target Redis information.
  The [target_instance](#target_instance_struct) structure is documented below.

* `description` - (Optional, String, NonUpdatable) Specifies the backup import task description.

<a name="backup_files_struct"></a>
The `backup_files` block supports:

* `file_source` - (Required, String, NonUpdatable) Specifies the data source, which can be an OBS bucket or a backup
  record. Value options:
  + **self_build_obs**: OBS bucket
  + **backup_record**: backup record

* `bucket_name` - (Optional, String, NonUpdatable) Specifies the OBS bucket name. It is mandatory when `file_source`
  is **self_build_obs**.

* `files` - (Optional, List, NonUpdatable) Specifies the list of backup files to be imported. It is mandatory when
  `file_source` is **self_build_obs**.
  The [files](#files_struct) structure is documented below.

* `backup_id` - (Optional, String, NonUpdatable) Specifies the backup record ID. It is mandatory when `file_source` is
  **backup_record**.

<a name="files_struct"></a>
The `files` block supports:

* `file_name` - (Required, String, NonUpdatable) Specifies the name of a backup file.

* `size` - (Optional, String, NonUpdatable) Specifies the file size in bytes.

* `update_at` - (Optional, String, NonUpdatable) Specifies the time when the file was last modified. The format is
  **YYYY-MM-DDTHH:MM:SS**.

<a name="target_instance_struct"></a>
The `target_instance` block supports:

* `id` - (Required, String, NonUpdatable) Specifies the Redis instance ID.

* `password` - (Optional, String, NonUpdatable) Specifies the Redis password. If a password of the DCS instance is set,
  it is mandatory.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the migration task status.

* `created_at` - Indicates the time when the migration task is created.

* `updated_at` - Indicates the time when the migration task is complete.

* `released_at` - Indicates the time when the migration ECS is released.

* `target_instance` - Indicates the target Redis information.
  The [target_instance](#target_instance_struct) structure is documented below.

<a name="target_instance_struct"></a>
The `target_instance` block supports:

* `name` - Indicates the Redis name.

## Timeouts

This resource provides the following timeout configuration option:

* `create` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The DCS backup import task can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_dcs_backup_import_task.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `target_instance.0.password`. It is generally
recommended running `terraform plan` after importing the resource. You can then decide if changes should be applied to
the resource, or the resource definition should be updated to align with the task. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_dcs_backup_import_task" "test" {
    ...

  lifecycle {
    ignore_changes = [
      target_instance.0.password,
    ]
  }
}
```
