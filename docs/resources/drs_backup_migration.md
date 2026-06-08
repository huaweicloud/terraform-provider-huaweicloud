---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_backup_migration"
description: |-
  Manages DRS backup migration resource within HuaweiCloud.
---

# huaweicloud_drs_backup_migration

Manages DRS backup migration resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "instance_id" {}
variable "bucket_name" {}

resource "huaweicloud_drs_backup_migration" "test" {
  base_info {
    name        = var.name
    engine_type = "sqlserver"
    description = "test backup migration"
  }

  target_db_info {
    target_instance_id = var.instance_id
  }

  backup_info {
    file_source = "OBS"
    bucket_name = var.bucket_name

    files {
      name     = "test.bak"
      obs_path = "testFolder/testPath"
    }
  }

  options {
    is_cover              = "false"
    is_default_restore    = "true"
    is_delete_backup_file = "false"
    is_last_backup        = true
    is_precheck           = true
    recovery_mode         = "full"
    db_names              = ["db-test-name"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `base_info` - (Required, List) Specifies the basic information of the backup migration job.
  The maximum number of elements is `1`.
  The [base_info](#block--base_info) structure is documented below.

* `target_db_info` - (Required, List, NonUpdatable) Specifies the target RDS for SQL Server instance information.
  The maximum number of elements is `1`.
  The [target_db_info](#block--target_db_info) structure is documented below.

* `backup_info` - (Required, List, NonUpdatable) Specifies the backup file information.
  The maximum number of elements is `1`.
  The [backup_info](#block--backup_info) structure is documented below.

* `options` - (Required, List, NonUpdatable) Specifies the backup migration configuration parameters.
  The maximum number of elements is `1`.
  The [options](#block--options) structure is documented below.

<a name="block--base_info"></a>
The `base_info` block supports:

* `name` - (Required, String) Specifies the job name. The name consists of `4` to `50` characters, starting with
  a letter. Only letters, digits, underscores (_) and hyphens (-) are allowed.

* `engine_type` - (Required, String, NonUpdatable) Specifies the database engine type.
  The options are as follows:
  + **sqlserver**: RDS for SQL Server engine.

* `description` - (Optional, String) Specifies the description of the job.

* `tags` - (Optional, List, NonUpdatable) Specifies the tag information.
  The [tags](#block--tags) structure is documented below.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.

<a name="block--tags"></a>
The `tags` block supports:

* `key` - (Optional, String, NonUpdatable) Specifies the tag key. The maximum length is 36 characters. Only letters, digits,
  underscores (_), hyphens (-) and Chinese characters are allowed.

* `value` - (Optional, String, NonUpdatable) Specifies the tag value. The value can contain letters, digits, spaces and
  _ . : / = + - @.

<a name="block--target_db_info"></a>
The `target_db_info` block supports:

* `target_instance_id` - (Required, String, NonUpdatable) Specifies the target RDS for SQL Server instance ID.

* `ms_file_stream_status` - (Optional, String, NonUpdatable) Specifies whether the target instance has FileStream
  mode enabled. This can be obtained through the RDS for SQL Server details API.

* `file_id` - (Optional, String, NonUpdatable) Specifies the file ID of the RDS for SQL Server backup file.
  It is mandatory when performing a full restoration from RDS. This can be obtained from the RDS backup management
  page.

<a name="block--backup_info"></a>
The `backup_info` block supports:

* `file_source` - (Required, String, NonUpdatable) Specifies the backup file source. The options are as follows:
  + **OBS**: The backup file is stored in OBS.
  + **RDS**: The backup file is from an RDS instance.

* `files` - (Required, List, NonUpdatable) Specifies the backup file information list.
  The [files](#block--files) structure is documented below.

* `bucket_name` - (Optional, String, NonUpdatable) Specifies the OBS bucket name. It is mandatory when `file_source`
  is **OBS**. The length ranges from `3` to `63` characters. Only lowercase letters, digits, hyphens (-) and dots (.)
  are allowed, and it must start and end with a letter or digit. IP addresses are not allowed.

<a name="block--files"></a>
The `files` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the backup file name.

* `obs_path` - (Optional, String, NonUpdatable) Specifies the file path in the OBS bucket. It is mandatory when
  `file_source` is **OBS**.

* `rds_version` - (Optional, String, NonUpdatable) Specifies the database version of the bak file. It is mandatory
  when `file_source` is **RDS**.

* `rds_source_instance_id` - (Optional, String, NonUpdatable) Specifies the instance to which the bak file belongs.
  It is mandatory when `file_source` is **RDS**.

<a name="block--options"></a>
The `options` block supports:

* `is_last_backup` - (Required, Bool, NonUpdatable) Specifies whether this is the last backup. A typical incremental
  recovery process involves multiple incremental backup recoveries. Each incremental backup recovery keeps the target
  database in a restoring state where it is not readable or writable, until the last incremental backup recovery is
  completed. The scenarios for setting this field to **true** are as follows:
  + One-time full migration, after which no incremental recovery will be performed.
  + During incremental recovery, when it is determined to be the last incremental backup in the final cutover phase.

* `is_precheck` - (Required, Bool, NonUpdatable) Specifies whether to perform a pre-check. The options are as follows:
  + **true**: Perform pre-check.
  + **false**: Do not perform pre-check.

* `recovery_mode` - (Required, String, NonUpdatable) Specifies the recovery mode. The options are as follows:
  + **full**: Full migration.
  + **incre**: Incremental migration.

* `is_cover` - (Optional, String, NonUpdatable) Specifies whether to overwrite the target database. Defaults to **false**.
  The options are as follows:
  + **true**: Overwrite the target database.
  + **false**: Do not overwrite the target database.

* `is_default_restore` - (Optional, String, NonUpdatable) Specifies whether to restore all databases by default.
  Defaults to **true**. The options are as follows:
  + **true**: Restore all databases in the backup file.
  + **false**: Specify the database names to be restored in the `db_names` field.

* `is_delete_backup_file` - (Optional, String, NonUpdatable) Specifies whether to delete the OBS backup file downloaded
  to the RDS for SQL Server disk when the job ends. Defaults to **true**. The options are as follows:
  + **true**: Delete the backup file.
  + **false**: Do not delete the backup file.

* `db_names` - (Optional, List, NonUpdatable) Specifies the database names.

* `reset_db_name_map` - (Optional, Map, NonUpdatable) Specifies the mapping of new database names. The key is the
  original database name and the value is the new database name. This feature ignores the original database name in
  the backup file and restores it to the specified new database name through DRS. The conditions for using this
  feature are as follows:
  + There is only one database in the backup file.
  + The backup file is a full backup type (recovery mode is **full**), and it is a one-time recovery
    (`is_last_backup` is **true**).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `status` - The status of the backup migration job. The options are as follows:
  + **success**: The job succeeded.
  + **failed**: The job failed.

* `create_time` - The creation time of the job.

* `finish_time` - The completion time of the job.

* `new_db_names` - The new database names after backup restoration mapping.

* `instance_name` - The RDS instance name.

* `error_log` - The failure reason during migration.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

## Import

The DRS backup migration can be imported by `id`. e.g.

```bash
$ terraform import huaweicloud_drs_backup_migration.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `backup_info.0.files`.
It is generally recommended running **terraform plan** after importing a job. You can then
decide if changes should be applied to the job, or the resource definition should be updated to align with the job. Also
you can ignore changes as below.

```hcl
resource "huaweicloud_drs_backup_migration" "test" {
    ...

  lifecycle {
    ignore_changes = [
      backup_info.0.files,
    ]
  }
}
```
