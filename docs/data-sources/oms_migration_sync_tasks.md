---
subcategory: "Object Storage Migration Service (OMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_oms_migration_sync_tasks"
description: |-
  Use this data source to get the list of synchronization tasks.
---

# huaweicloud_oms_migration_sync_tasks

Use this data source to get the list of synchronization tasks.

## Example Usage

```hcl
data "huaweicloud_oms_migration_sync_tasks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `status` - (Optional, String) Specifies the synchronization task status.  
  The valid values are as follows:
  + **SYNCHRONIZING**: Being synchronizing.
  + **STOPPED**: Already stopped.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The list of synchronization tasks.
  The [tasks](#tasks_struct) structure is documented below.

<a name="tasks_struct"></a>
The `tasks` block supports:

* `sync_task_id` - The synchronization task ID.

* `src_cloud_type` - The source cloud service provider.
  + **AWS**
  + **Azure**
  + **Aliyun**
  + **Tencent**
  + **HuaweiCloud**
  + **QingCloud**
  + **KingsoftCloud**
  + **Baidu**
  + **Qiniu**
  + **UCloud**

* `src_region` - The region where the source bucket is located.

* `src_bucket` - The name of the source bucket.

* `dst_bucket` - The name of the destination bucket.

* `dst_region` - The region where the destination bucket is located.

* `description` - The synchronization task description.

* `status` - The synchronization task status.

* `enable_kms` - Whether KMS encryption is enabled.

* `enable_metadata_migration` - Whether metadata migration is enabled.

* `enable_restore` - Whether automatic restoration of archived data is enabled.

* `app_id` - The app ID.

* `source_cdn` - Whether migration from CDN is enabled.
  The [source_cdn](#source_cdn_struct) structure is documented below.

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

* `consistency_check` - The method for checking whether objects are consistent after migration.

* `create_time` - The creation time of the synchronization task, in milliseconds.

* `last_start_time` - The most recent start time of the synchronization task.

<a name="source_cdn_struct"></a>
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
