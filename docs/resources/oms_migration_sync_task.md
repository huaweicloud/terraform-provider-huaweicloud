---
subcategory: "Object Storage Migration Service (OMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_oms_migration_sync_task"
description: ""
---

# huaweicloud_oms_migration_sync_task

Manages an OMS migration synchronization task resource within HuaweiCloud.

## Example Usage

```hcl
variable "source_region" {}
variable "source_bucket" {}
variable "source_access_key" {}
variable "source_secret_key" {}
variable "dest_region" {}
variable "dest_bucket" {}
variable "dest_access_key" {}
variable "dest_secret_key" {}

resource "huaweicloud_oms_migration_sync_task" "test" {
  region = var.dest_region

  src_cloud_type = "HuaweiCloud"
  src_region     = var.source_region
  src_bucket     = var.source_bucket
  src_ak         = var.source_access_key
  src_sk         = var.source_secret_key
  dst_bucket     = var.dest_bucket
  dst_ak         = var.dest_access_key
  dst_sk         = var.dest_secret_key
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. Which is also the region
  for the destination bucket. If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `src_region` - (Required, String, ForceNew) Specifies the region where the source bucket is located.
  Changing this creates a new resource.

* `src_bucket` - (Required, String, ForceNew) Specifies the name of the source bucket.
  Changing this creates a new resource.

* `src_ak` - (Required, String, ForceNew) Specifies the access key for accessing the source bucket.
  Changing this creates a new resource.

* `src_sk` - (Required, String, ForceNew) Specifies the secret key for accessing the source bucket.
  Changing this creates a new resource.

* `dst_bucket` - (Required, String, ForceNew) Specifies the name of the destination bucket.
  Changing this creates a new resource.

* `dst_ak` - (Required, String, ForceNew) Specifies the access key for accessing the destination bucket.
  Changing this creates a new resource.

* `dst_sk` - (Required, String, ForceNew) Specifies the secret key for accessing the destination bucket.
  Changing this creates a new resource.

* `src_cloud_type` - (Optional, String, ForceNew) Specifies the source cloud service provider. Value options:
  **AWS**, **Azure**, **Aliyun**, **Tencent**, **HuaweiCloud**, **QingCloud**, **KingsoftCloud**, **Baidu**,
  **Qiniu**, **URLSource** and **UCloud**. Default value: **Aliyun**. Changing this creates a new resource.

* `app_id` - (Optional, String, ForceNew) Specifies the APP ID. This parameter is mandatory when `src_cloud_type` is
  **Tencent**. Changing this creates a new resource.

* `consistency_check` - (Optional, String, ForceNew) Specifies the consistency check method, which is used to check
  whether objects are consistent before and after migration. All check methods take effect for only objects that are
  in the same encryption status in the source and destination buckets. The check method and results will be recorded
  in the object list. Value options:

  + **size_last_modified**: the system checks object consistency with object size and last modification time.
    If a source object is as large as but was last modified earlier than its paired destination object, the system
    considers the source object does not need to be migrated or has been already migrated successfully.

  + **crc64**: this option is only available for migration on Huawei Cloud or from Alibaba Cloud or Tencent Cloud. If
    a source object and its paired destination object have CRC64 checksums, the checksums are checked. Otherwise, their
    sizes and last modification times are checked.

  + **no_check**: this option is only available for migration of HTTP/HTTPS data. This option takes effect for source
    objects whose sizes cannot be obtained using the content-length field in the standard HTTP protocol. These source
    objects will overwrite their paired destination objects directly.
    If the size of a source object can be obtained, its size and last modification time will be checked.

  The default value is **size_last_modified**. Changing this creates a new resource.

* `description` - (Optional, String, ForceNew) Specifies the description of the synchronization task.
  Changing this creates a new resource.

* `enable_kms` - (Optional, Bool, ForceNew) Specifies whether to enable the KMS encryption function.
  Default value: **false**. Changing this creates a new resource.

* `enable_metadata_migration` - (Optional, Bool, ForceNew) Specifies whether metadata migration is enabled.
  Default value: **false**. Even if disabled, the ContentType metadata will still be migrated
  to ensure a successful migration. Changing this creates a new resource.

* `enable_restore` - (Optional, Bool, ForceNew) Specifies whether to automatically restore the archive data. If enabled,
  archive data is automatically restored and migrated. Default value: **false**. Changing this creates a new resource.

* `source_cdn` - (Optional, List, ForceNew) Specifies the CDN information. If this parameter is contained,
  using CDN to download source data is supported, the source objects to be migrated are obtained from the CDN domain
  name during migration. Changing this creates a new resource.
  The [source_cdn](#block--source_cdn) structure is documented below.

* `action` - (Optional, String) Specifies the action for migration synchronization task. Value options:

  + **start**: Start a migration synchronization task.
  + **stop**:  Pause a migration synchronization task.

<a name="block--source_cdn"></a>
The `source_cdn` block supports:

* `domain` - (Required, String, ForceNew) Specifies the domain name from which to obtain objects to be migrated.
  Changing this creates a new resource.

* `protocol` - (Required, String, ForceNew) Specifies the protocol type. Value options: **http** and **https**.
  Changing this creates a new resource.

* `authentication_type` - (Optional, String, ForceNew) Specifies the authentication type. Value options:

  + **NONE**
  + **QINIU_PRIVATE_AUTHENTICATION**
  + **ALIYUN_OSS_A**
  + **ALIYUN_OSS_B**
  + **ALIYUN_OSS_C**
  + **KSYUN_PRIVATE_AUTHENTICATION**
  + **TENCENT_COS_A**
  + **TENCENT_COS_B**
  + **TENCENT_COS_C**
  + **TENCENT_COS_D**

  Default value: **NONE**. Changing this creates a new resource.

* `authentication_key` - (Optional, String, ForceNew) Specifies the CDN authentication key.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` -  The resource ID.

* `created_at` - Indicates the time when the migration synchronization task was created.

* `dst_storage_policy` - Indicates the destination storage class. The value can be:

  + **STANDARD**
  + **IA**
  + **ARCHIVE**
  + **DEEP_ARCHIVE**
  + **SRC_STORAGE_MAPPING**

* `last_start_at` - Indicates the last time when the migration synchronization task started.

* `monthly_acceptance_request` - Indicates the number of objects requested to be synchronized in the current month.

* `monthly_failure_object` - Indicates the number of objects that failed to be synchronized in the current month.

* `monthly_size` - Indicates the total size of synchronized objects in the current month, in bytes.

* `monthly_skip_object` - Indicates the number of objects that were ignored in the current month.

* `monthly_success_object` - Indicates the number of objects that were successfully synchronized in the current month.

* `object_overwrite_mode` - Indicates the type of the source object to overwrite its paired destination object.
  The value can be:

  + **NO_OVERWRITE**: Indicates the system never allows override. The system always skips source objects and keeps
    their paired destination objects.

  + **SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE**: Indicates the system allows override based on the results of size or
    modification time checks. If a source object is not as large as or was last modified more recently than its paired
    destination object, the source object will overwrite the destination object. Otherwise, the source object will be
    skipped.

  + **CRC64_COMPARISON_OVERWRITE**: Indicates the system allows override if the source and destination objects have
    different CRC64 checksums. This option is only available for migration on Huawei Cloud or from Alibaba Cloud or
    Tencent Cloud. If a source object has a CRC64 checksum different from the paired destination object, the source
    object will overwrite the destination object. Otherwise, the source object will be skipped.
    If any of them doesn't have a CRC64 checksum, their sizes and last modification times are checked.

  + **FULL_OVERWRITE**: Indicates the system always allows override. The system always allows source objects to
    overwrite their paired destination objects.

* `status` - Indicates the status of the migration synchronization task. The value can be:

  + **SYNCHRONIZING**:synchronizing.
  + **STOPPED**:stopped.
