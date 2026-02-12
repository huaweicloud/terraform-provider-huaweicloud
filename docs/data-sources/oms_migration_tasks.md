---
subcategory: "Object Storage Migration Service (OMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_oms_migration_tasks"
description: |-
  Use this data source to get the list of migration tasks.
---

# huaweicloud_oms_migration_tasks

Use this data source to get the list of migration task groups.

## Example Usage

```hcl
data "huaweicloud_oms_migration_tasks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `group_id` - (Optional, String) Specifies the migration task group ID.

* `status` - (Optional, Int) Specifies the migration task status.
  The valid values are as follows:
  + **1**: Waiting for scheduling.
  + **2**: Migrating.
  + **3**: Paused.
  + **4**: Failed.
  + **5**: Succeeded.
  + **7**: Pausing.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The list of migration task groups.
  The [tasks](#oms_migration_tasks_struct) structure is documented below.

<a name="oms_migration_tasks_struct"></a>
The `tasks` block supports:

* `bandwidth_policy` - The traffic limiting rules.
  The [bandwidth_policy](#oms_migration_tasks_bandwidth_policy_struct) structure is documented below.

* `complete_size` - The size of the objects that have been processed in the task, in bytes.

* `description` - The task description.

* `dst_node` - The destination information.
  The [dst_node](#oms_migration_tasks_dst_node_struct) structure is documented below.

* `enable_failed_object_recording` - Whether the function of recording failed objects is enabled.

* `enable_kms` - Whether KMS is used to encrypt the data to be stored in the destination OBS bucket.

* `enable_metadata_migration` - Whether metadata migration is enabled.

* `enable_restore` - Whether automatic restoration of archived data is enabled.

* `error_reason` - The task failure cause.
  The value is an empty string if the task is not in the Migration failed state.
  The [error_reason](#oms_migration_tasks_error_reason_struct) structure is documented below.

* `fail_num` - The number of objects that fail to be migrated.

* `failed_object_record` - The record of objects that failed to be migrated.
  The [failed_object_record](#oms_migration_tasks_failed_object_record_struct) structure is documented below.

* `group_id` - The ID of the migration task group to which the task belongs.

* `id` - The task ID.

* `is_query_over` - Whether the statistics of source objects in the migration task have been scanned.

* `left_time` - The remaining time of the task, in milliseconds.

* `migrate_since` - The specified timestamp will be migrated, in milliseconds.

* `migrate_speed` - The migration speed, in byte/s.

* `name` - The task name.

* `progress` - The task progress.

* `real_size` - The total size of the migrated objects, in bytes.

* `skipped_num` - The number of objects skipped during migration.

* `src_node` - The source information.
  The [src_node](#oms_migration_tasks_src_node_struct) structure is documented below.

* `start_time` - The start time of the migration task, in milliseconds.

* `status` - The migration task status.

* `successful_num` - The number of successfully migrated objects.

* `task_type` - The task type.
  + **list**
  + **object**
  + **prefix**
  + **url_list**

* `group_type` - The migration group type.
  + **NORMAL_TASK**: General migration tasks.
  + **SYNC_TASK**: Migration tasks generated for a synchronization task.
  + **GROUP_TASK**: Migration tasks in a task group.

* `total_num` - The total number of objects that need to be migrated in the task.

* `total_size` - The size of objects that need to be migrated in the task, in bytes.

* `total_time` - The total time used, in milliseconds.

* `smn_info` - The SMN notification results.
  The [smn_info](#oms_migration_tasks_smn_info_struct) structure is documented below.

* `source_cdn` - Whether migration from CDN is enabled.
  The [source_cdn](#oms_migration_tasks_source_cdn_struct) structure is documented below.

* `success_record_error_reason` - The error code returned for the failure in recording the list of objects  
  that are successfully migrated.

* `skip_record_error_reason` - The error code returned for the failure in recording the list of objects that  
  are ignored.

* `object_overwrite_mode` - How a source object handles its paired destination object,
  either overwriting the object or skipping the migration.  
  + **NO_OVERWRITE**
  + **SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE**
  + **CRC64_COMPARISON_OVERWRITE**
  + **FULL_OVERWRITE**

* `dst_storage_policy` - The destination storage class.
  + **STANDARD**
  + **IA**
  + **ARCHIVE**
  + **DEEP_ARCHIVE**
  + **SRC_STORAGE_MAPPING**

* `consistency_check` - The method for checking whether objects are consistent before and after migration.
  + **size_last_modified**
  + **crc64**
  + **no_check**

* `enable_requester_pays` - Whether Requester Pays is enabled.

* `task_priority` - The migration task priority.
  + **HIGH**: High priority.
  + **MEDIUM**: Medium priority.
  + **LOW**: Low priority.

<a name="oms_migration_tasks_bandwidth_policy_struct"></a>
The `bandwidth_policy` block supports:

* `max_bandwidth` - The maximum traffic bandwidth allowed in the specified period. The unit is byte/s.

* `start` - The time when traffic limiting is started. The format is **hh:mm**.

* `end` - The time when traffic limiting is ended. The format is **hh:mm**.

<a name="oms_migration_tasks_dst_node_struct"></a>
The `dst_node` block supports:

* `bucket` - The name of the destination bucket.

* `region` - The region where the destination bucket is located.

* `save_prefix` - The path prefix used to organize object locations in the destination bucket.

<a name="oms_migration_tasks_error_reason_struct"></a>
The `error_reason` block supports:

* `error_code` - The error code returned when a migration fails.

* `error_msg` - The migration failure cause.

<a name="oms_migration_tasks_src_node_struct"></a>
The `src_node` block supports:

* `bucket` - The name of the source bucket.

* `cloud_type` - The source cloud service provider.
  + **AWS**
  + **AZURE**
  + **ALIYUN**
  + **TENCENT**
  + **HUAWEICLOUD**
  + **QINGCLOUD**
  + **KINGSOFTCLOUD**
  + **BAIDU**
  + **QINIU**
  + **URLSOURCE**
  + **UCLOUD**
  + **GOOGLE**

* `region` - The region where the source bucket is located.

* `app_id` - The app ID.

* `object_key` - The name prefixes of objects to be migrated.

* `list_file` - The configurations of the list file.
  The [list_file](#oms_migration_tasks_list_file_struct) structure is documented below.

<a name="oms_migration_tasks_list_file_struct"></a>
The `list_file` block supports:

* `list_file_key` - The object names in the object or URL list file.

* `obs_bucket` - The name of the OBS bucket for storing object list files.

* `list_file_num` - The number of stored object list files.

<a name="oms_migration_tasks_smn_info_struct"></a>
The `smn_info` block supports:

* `notify_result` - Whether SMN messages are sent successfully after a migration task is complete.

* `notify_error_message` - The error codes presenting why SMN messages failed to be sent.

* `topic_name` - The SMN topic name.

<a name="oms_migration_tasks_source_cdn_struct"></a>
The `source_cdn` block supports:

* `domain` - The domain name used to obtain objects to be migrated.

* `protocol` - The protocol type.
  + **http**
  + **https**

* `authentication_type` - The authentication type.
  + **NONE**
  + **QINIU_PRIVATE_AUTHENTICATION**
  + **ALIYUN_OSS_A**
  + **ALIYUN_OSS_B**
  + **ALIYUN_OSS_C**
  + **KSYUN_PRIVATE_AUTHENTICATION**
  + **AZURE_SAS_TOKEN**
  + **TENCENT_COS_A**
  + **TENCENT_COS_B**
  + **TENCENT_COS_C**
  + **TENCENT_COS_D**

<a name="oms_migration_tasks_failed_object_record_struct"></a>
The `failed_object_record` block supports:

* `result` - Whether retransmission of failed objects is supported.

* `list_file_key` - The path for storing the list of failed objects.

* `error_code` - The error code returned when the list of failed objects fails to be uploaded.
