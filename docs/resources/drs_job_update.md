---
subcategory: "DRS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_job_update"
description: |-
  Use this resource to update DRS job details within HuaweiCloud.
---

# huaweicloud_drs_job_update

Use this resource to update DRS job details within HuaweiCloud.

-> This resource is a one-time action resource for updating DRS job details. Deleting this resource will
   not rollback the update operation, but will only remove the resource information from the tfstate file.

## Example Usage

### Update sync type policy

```hcl
variable "job_id" {}

resource "huaweicloud_drs_job_update" "test" {
  job_id = var.job_id
  type   = "policy"

  policy_config {
    sync_type_policy = "tableData"
  }
}
```

### Update job name and description

```hcl
variable "job_id" {}

resource "huaweicloud_drs_job_update" "test" {
  job_id = var.job_id
  type   = "name"

  base_info {
    name        = "drs-test-rename"
    description = "updated by terraform"
  }
}
```

### Update speed limit

```hcl
variable "job_id" {}

resource "huaweicloud_drs_job_update" "test" {
  job_id = var.job_id
  type   = "speed_limit"

  speed_limit {
    speed      = "10"
    start_time = "08:00"
    end_time   = "22:00"
    is_effect  = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the DRS job is located.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String, NonUpdatable) Specifies the DRS job ID to be updated.

* `type` - (Required, String, NonUpdatable) Specifies the update type. The valid values are as follows:
  + **endpoint**: Update source/target database connection information.
  + **description**: Update job description.
  + **name**: Update job name.
  + **notify**: Update alarm notification settings.
  + **expired_days**: Update job expiration days.
  + **policy**: Update synchronization policy. Only specified fields are updated, unspecified fields are preserved.
  + **speed_limit**: Update speed limit configuration.
  + **tuning_params**: Update tuning parameters.
  + **db_object**: Update database object selection (async, returns 202).
  + **user_migration**: Update user migration settings.
  + **precheck**: Trigger pre-check.
  + **db_param**: Update database parameters.
  + **db_password**: Update database password (triggers kernel restart).
  + **re_create**: Re-create the job after 3 days of configuration.
  + **all**: Update all parameters.
  + **resource_group**: Update resource group.
  + **application_id**: Update application ID.
  + **dr_info**: Update disaster recovery information.
  + **replay_config**: Update replay simulation parameters.

* `policy_config` - (Optional, List, NonUpdatable) Specifies the synchronization policy configuration.
  Only valid when `type` is `policy`. Unspecified fields are preserved (not overwritten).
  The [policy_config](#drs_job_update_policy_config) structure is documented below.

* `base_info` - (Optional, List, NonUpdatable) Specifies the base information of the job.
  Only valid when `type` is `name` or `description`.
  The [base_info](#drs_job_update_base_info) structure is documented below.

* `speed_limit` - (Optional, List, NonUpdatable) Specifies the speed limit configuration.
  Only valid when `type` is `speed_limit`.
  The [speed_limit](#drs_job_update_speed_limit) structure is documented below.

<a name="drs_job_update_policy_config"></a>
The `policy_config` block supports:

* `sync_type_policy` - (Optional, String) Specifies the synchronization data type. The valid values are as follows:
  + **tableStructure**: Synchronize table structure.
  + **constraintData**: Synchronize index/constraints.
  + **tableData**: Synchronize table data.
  
  Multiple values can be combined with commas.

* `conflict_policy` - (Optional, String) Specifies the conflict policy for incremental synchronization.
  Valid values: **overwrite**, **ignore**, **stop**.

* `filter_ddl_policy` - (Optional, String) Specifies the DDL filter policy. Valid value: **drop_database**.

* `ddl_trans` - (Optional, Bool) Specifies whether to synchronize DDL during incremental phase.

* `index_trans` - (Optional, Bool) Specifies whether to synchronize non-unique indexes.

* `struct_trans` - (Optional, Bool) Specifies whether to synchronize structure.

* `data_sync_topology_type` - (Optional, String) Specifies the data synchronization topology type.
  Valid values: **one2one**, **one2many**, **many2one**.

* `support_ddl_info` - (Optional, String) Specifies the supported DDL types.

* `full_mode` - (Optional, String) Specifies the full synchronization mode.

* `export_snapshot` - (Optional, Bool) Specifies whether to export snapshot.

* `replace_definer` - (Optional, Bool) Specifies whether to replace definer.

* `increment_read_mode` - (Optional, String) Specifies the incremental read mode.

* `slot_name` - (Optional, String) Specifies the logical replication slot name.
  Only modifiable in CONFIGURATION status.

* `transformation_name_case_policy` - (Optional, String) Specifies the object name case transformation policy.

* `is_create_table_with_index` - (Optional, Bool) Specifies whether to create index simultaneously when creating table.

* `kafka_data_format` - (Optional, String) Specifies the Kafka data format.

* `kafka_compression_type` - (Optional, String) Specifies the Kafka compression type.

* `partition_policy` - (Optional, String) Specifies the Kafka partition policy.

* `topic_policy` - (Optional, String) Specifies the Kafka topic policy.

* `topic_name_format` - (Optional, String) Specifies the Kafka topic name format.

* `topic` - (Optional, String) Specifies the Kafka topic.

* `partitions_num` - (Optional, Int) Specifies the number of Kafka partitions.

* `replication_factor` - (Optional, Int) Specifies the Kafka replication factor.

* `filter_user_names` - (Optional, String) Specifies the filtered user names.

<a name="drs_job_update_base_info"></a>
The `base_info` block supports:

* `name` - (Optional, String) Specifies the job name (4-50 characters).

* `description` - (Optional, String) Specifies the job description (max 256 characters).

<a name="drs_job_update_speed_limit"></a>
The `speed_limit` block supports:

* `speed` - (Required, String) Specifies the speed limit value in MB/s (1-9999).

* `start_time` - (Required, String) Specifies the start time in HH:MM format (UTC).

* `end_time` - (Required, String) Specifies the end time in HH:MM format (UTC, must end with :59).

* `is_effect` - (Optional, Bool) Specifies whether the speed limit is effective.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as the job ID.
