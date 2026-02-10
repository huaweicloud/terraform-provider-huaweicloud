---
subcategory: "Object Storage Migration Service (OMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_oms_migration_task_groups"
description: |-
  Use this data source to get the list of migration task groups.
---

# huaweicloud_oms_migration_task_groups

Use this data source to get the list of migration task groups.

## Example Usage

```hcl
data "huaweicloud_oms_migration_task_groups" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `status` - (Optional, String) Specifies the migration task group status.
  The valid values are as follows:
  + **0**: Waiting.
  + **1**: Executing/creating.
  + **2**: Running monitor task.
  + **3**: Paused.
  + **4**: Creation failed.
  + **5**: Migration failed.
  + **6**: Migration completed.
  + **7**: Pausing.
  + **8**: Waiting for deletion.
  + **9**: Deleting.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `taskgroups` - The list of migration task groups.
  The [taskgroups](#oms_migration_task_groups_struct) structure is documented below.

<a name="oms_migration_task_groups_struct"></a>
The `taskgroups` block supports:

* `group_id` - The task group ID.

* `status` - The migration task group status.

* `error_reason` - The task group failure cause.
  The value is an empty string if the task is not in the Migration failed state.
  The [error_reason](#oms_migration_task_groups_error_reason_struct) structure is documented below.

* `src_node` - The source information.
  The [src_node](#oms_migration_task_groups_src_node_struct) structure is documented below.

* `description` - The task group description.

* `dst_node` - The destination information.
  The [dst_node](#oms_migration_task_groups_dst_node_struct) structure is documented below.

* `enable_metadata_migration` - Whether metadata migration is enabled.

* `enable_failed_object_recording` - Whether the function of recording failed objects is enabled.

* `enable_restore` - Whether automatic restoration of archived data is enabled.

* `enable_kms` - Whether KMS is used to encrypt the data to be stored in the destination OBS bucket.

* `task_type` - The task type.
  + **LIST**
  + **URL_LIST**
  + **PREFIX**

* `bandwidth_policy` - The traffic limiting rules.
  The [bandwidth_policy](#oms_migration_task_groups_bandwidth_policy_struct) structure is documented below.

* `smn_config` - The configuration of SMN message sending.
  The [smn_config](#oms_migration_task_groups_smn_config_struct) structure is documented below.

* `source_cdn` - Whether migration from CDN is enabled.
  The [source_cdn](#oms_migration_task_groups_source_cdn_struct) structure is documented below.

* `migrate_since` - The specified timestamp will be migrated, in milliseconds.

* `migrate_speed` - The migration speed, in byte/s.

* `total_time` - The method for checking whether objects are consistent after migration.

* `start_time` - The start time of the migration task group, in milliseconds.

* `total_task_num` - The total number of migration tasks in the task group.

* `create_task_num` - The number of created migration tasks in the task group.

* `failed_task_num` - The number of failed migration tasks in the task group.

* `complete_task_num` - The number of completed migration tasks in the task group.

* `paused_task_num` - The number of paused migration tasks in the task group.

* `executing_task_num` - The number of migration tasks being executed in the task group.

* `waiting_task_num` - The number of waiting migration tasks in the task group.

* `total_num` - The total number of objects to be migrated in the migration task group.

* `create_complete_num` - The total number of objects included in the created migration tasks.

* `success_num` - The number of migrated objects.

* `fail_num` - The number of failed objects.

* `skip_num` - The number of skipped objects.

* `total_size` - The size of the objects that have been migrated, in bytes.

* `create_complete_size` - The total size of objects migrated in the created migration tasks, in bytes.

* `complete_size` - The total size of migrated objects, in bytes.

* `failed_object_record` - The record of failed objects.
  The [failed_object_record](#oms_migration_task_groups_failed_object_record_struct) structure is documented below.

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

<a name="oms_migration_task_groups_error_reason_struct"></a>
The `error_reason` block supports:

* `error_code` - The error code returned when a migration fails.

* `error_msg` - The migration failure cause.

<a name="oms_migration_task_groups_src_node_struct"></a>
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
  The [list_file](#oms_migration_task_groups_list_file_struct) structure is documented below.

<a name="oms_migration_task_groups_list_file_struct"></a>
The `list_file` block supports:

* `list_file_key` - The object names in the object or URL list file.

* `obs_bucket` - The name of the OBS bucket for storing object list files.

* `list_file_num` - The number of stored object list files.

<a name="oms_migration_task_groups_dst_node_struct"></a>
The `dst_node` block supports:

* `bucket` - The name of the destination bucket.

* `region` - The region where the destination bucket is located.

* `save_prefix` - The path prefix used to organize object locations in the destination bucket.

<a name="oms_migration_task_groups_bandwidth_policy_struct"></a>
The `bandwidth_policy` block supports:

* `max_bandwidth` - The maximum traffic bandwidth allowed in the specified period. The unit is byte/s.

* `start` - The time when traffic limiting is started. The format is **hh:mm**.

* `end` - The time when traffic limiting is ended. The format is **hh:mm**.

<a name="oms_migration_task_groups_smn_config_struct"></a>
The `smn_config` block supports:

* `notify_result` - Whether SMN messages are sent successfully after a migration task is complete.

* `notify_error_message` - The error codes presenting why SMN messages failed to be sent.

* `topic_name` - The SMN topic name.

<a name="oms_migration_task_groups_source_cdn_struct"></a>
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

<a name="oms_migration_task_groups_failed_object_record_struct"></a>
The `failed_object_record` block supports:

* `result` - Whether retransmission of failed objects is supported.

* `list_file_key` - The path for storing the list of failed objects.

* `error_code` - The error code returned when the list of failed objects fails to be uploaded.
