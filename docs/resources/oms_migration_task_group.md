---
subcategory: "Object Storage Migration Service (OMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_oms_migration_task_group"
description: ""
---

# huaweicloud_oms_migration_task_group

Manages an OMS migration task group resource within HuaweiCloud.

## Example Usage

### OMS Migration Task Group with PREFIX

```hcl
variable "source_region" {}
variable "source_bucket" {}
variable "source_access_key" {}
variable "source_secret_key" {}
variable "dest_region" {}
variable "dest_bucket" {}
variable "dest_access_key" {}
variable "dest_secret_key" {}

resource "huaweicloud_oms_migration_task_group" "test" {
  source_object {
    data_source = "Aliyun"
    region      = var.source_region
    bucket      = var.source_bucket
    access_key  = var.source_access_key
    secret_key  = var.source_secret_key
    object      = [""]
  }

  destination_object {
    region     = var.dest_region
    bucket     = var.dest_bucket
    access_key = var.dest_access_key
    secret_key = var.dest_secret_key
  }

  type                           = "PREFIX"
  enable_kms                     = true
  description                    = "test task group"
  migrate_since                  = "2023-01-02 15:04:05"
  object_overwrite_mode          = "CRC64_COMPARISON_OVERWRITE"
  consistency_check              = "crc64"
  enable_requester_pays          = true
  enable_failed_object_recording = true

  bandwidth_policy {
    max_bandwidth = 1
    start         = "15:00"
    end           = "16:00"
  }

  bandwidth_policy {
    max_bandwidth = 2
    start         = "16:00"
    end           = "17:00"
  }
}
```

### OMS Migration Task Group with LIST

```hcl
variable "source_region" {}
variable "source_bucket" {}
variable "source_access_key" {}
variable "source_secret_key" {}
variable "dest_region" {}
variable "dest_bucket" {}
variable "dest_access_key" {}
variable "dest_secret_key" {}
variable "list_file_bucket" {}
variable "list_file_key" {}

resource "huaweicloud_oms_migration_task_group" "test" {
  source_object {
    data_source      = "HuaweiCloud"
    region           = var.source_region
    bucket           = var.source_bucket
    access_key       = var.source_access_key
    secret_key       = var.source_secret_key
    list_file_bucket = var.list_file_bucket
    list_file_key    = var.list_file_key
  }

  destination_object {
    region     = var.dest_region
    bucket     = var.dest_bucket
    access_key = var.dest_access_key
    secret_key = var.dest_secret_key
  }

  type        = "LIST"
  description = "test task group"
}
```

### OMS Migration Task Group with URL_LIST

```hcl
variable "dest_region" {}
variable "dest_bucket" {}
variable "dest_access_key" {}
variable "dest_secret_key" {}
variable "list_file_bucket" {}
variable "list_file_key" {}

resource "huaweicloud_oms_migration_task_group" "test" {
  source_object {
    data_source      = "URLSource"
    list_file_bucket = var.list_file_bucket
    list_file_key    = var.list_file_key
  }

  destination_object {
    region     = var.dest_region
    bucket     = var.dest_bucket
    access_key = var.dest_access_key
    secret_key = var.dest_secret_key
  }

  type        = "URL_LIST"
  description = "test task group"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `source_object` - (Required, List, ForceNew) Specifies the source information. The [object](#source_object_object)
  structure is documented below. Changing this creates a new resource.

* `destination_object` - (Required, List, ForceNew) Specifies the destination information.
  The [object](#destination_object_object) structure is documented below. Changing this creates a new resource.

* `type` - (Required, String, ForceNew) Specifies the task group type. The value can be:
  + **LIST**: indicates that the system will migrate the objects specified in the object list.
  + **URL_LIST**: indicates that the system will migrate the objects specified in the URL list.
  + **PREFIX**: indicates that the system will migrate the objects with a specific prefix.

  The default value is **PREFIX**. Changing this creates a new resource.

* `enable_kms` - (Optional, Bool, ForceNew) Specifies whether to enable the KMS encryption function.
  Default value: **false**. Changing this creates a new resource.

* `description` - (Optional, String, ForceNew) Specifies the description of the task group. The message cannot
  exceed 255 characters. The following special characters are not allowed: ^<>&"'.
  Changing this creates a new resource.

* `migrate_since` - (Optional, String, ForceNew) Specifies a time in format **yyyy-MM-dd HH:mm:ss**,
  e.g. **2006-01-02 15:04:05**. The system migrates only the objects that are modified after the specified time.
  No time is specified by default. Changing this creates a new resource.

* `object_overwrite_mode` - (Optional, String, ForceNew) Specifies whether to skip a source object or allow the source
  object to overwrite its paired destination object. Value options are as follows:

  + **NO_OVERWRITE**: indicates the system never allows override. The system always skips source objects and keeps
  their paired destination objects.

  + **SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE**: indicates the system allows override based on the results of size or
  modification time checks. If a source object is not as large as or was last modified more recently than its paired
  destination object, the source object will overwrite the destination object. Otherwise, the source object will be
  skipped.

  + **CRC64_COMPARISON_OVERWRITE**: indicates the system allows override if the source and destination objects have
  different CRC64 checksums. This option is only available for migration on Huawei Cloud or from Alibaba Cloud or
  Tencent Cloud. If a source object has a CRC64 checksum different from the paired destination object, the source
  object will overwrite the destination object. Otherwise, the source object will be skipped.
  If any of them doesn't have a CRC64 checksum, their sizes and last modification times are checked.

  + **FULL_OVERWRITE**: indicates the system always allows override. The system always allows source objects to
  overwrite their paired destination objects.

  The default value is **SIZE_LAST_MODIFIED_COMPARISON_OVERWRITE**. Changing this creates a new resource.

* `consistency_check` - (Optional, String, ForceNew) Specifies the consistency check method, which is used to check
  whether objects are consistent before and after migration. All check methods take effect for only objects that are
  in the same encryption status in the source and destination buckets. The check method and results will be recorded
  in the object list. Value options are as follows:

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

* `enable_requester_pays` - (Optional, Bool, ForceNew) Specifies whether to let the requester make payment.
  After enabled, the requester pays the request and data transmission fees.
  Default value: **false**. Changing this creates a new resource.

* `action` - (Optional, String) Specifies the action for migration task group. Value options are as follows:
  + **retry**: Restart a failed migration task group. This will migrate all objects in the failed tasks again.
  + **start**: Resume a paused migration task group.
  + **stop**:  Pause a migration task group when it is being monitored.

  -> **NOTE:** The usage of `action` has some limitations. Only creation failed task groups can be retried
  and only paused task groups can be started and only monitoring task groups can be stopped.

* `enable_failed_object_recording` - (Optional, Bool, ForceNew) Specifies whether to record failed objects.
  Default value: **true**. Changing this creates a new resource.

* `bandwidth_policy` - (Optional, List) Specifies the traffic limit rules. Each element in the array
  corresponds to the maximum bandwidth of a time segment. A maximum of 5 time segments are allowed, and the time
  segments must not overlap. The [object](#bandwidth_policy_object) structure is documented below.

* `source_cdn` - (Optional, List, ForceNew) Specifies the CDN information. If this parameter is contained,
  using CDN to download source data is supported, the source objects to be migrated are obtained from the CDN domain
  name during migration. The [object](#source_cdn_object) structure is documented below.
  Changing this creates a new resource.

* `enable_metadata_migration` - (Optional, Bool, ForceNew) Specifies whether metadata migration is enabled. Even if this
  function is disabled, the ContentType metadata will still be migrated to ensure a successful migration.
  Default value: **false**. Changing this creates a new resource.

* `dst_storage_policy` - (Optional, String, ForceNew) Specifies the destination storage class.
  This parameter is required only when the destination is Huawei Cloud OBS. The default value is STANDARD.
  + **STANDARD**: OBS Standard storage.
  + **IA**: OBS Infrequent Access storage.
  + **ARCHIVE**: OBS Archive storage
  + **DEEP_ARCHIVE**: OBS Deep Archive storage
  + **SRC_STORAGE_MAPPING**: converts the source storage class into an OBS storage class based on the predefined rules.
  Changing this creates a new resource.

<a name="source_object_object"></a>
The `source_object` block supports:

* `region` - (Optional, String, ForceNew) Specifies the region where the source bucket is located. `region` is
  mandatory when `type` is not **URL_LIST**. Changing this creates a new resource.

* `bucket` - (Optional, String, ForceNew) Specifies the name of the source bucket. `bucket` is mandatory when `type`
  is not **URL_LIST**. Changing this creates a new resource.

* `access_key` - (Optional, String, ForceNew) Specifies the access key for accessing the source bucket. This parameter
  is mandatory when `type` is not **URL_LIST**. Changing this creates a new resource.

* `secret_key` - (Optional, String, ForceNew) Specifies the secret key for accessing the source bucket. This
  parameter is mandatory when `type` is not **URL_LIST**. Changing this creates a new resource.

* `object` - (Optional, List, ForceNew) Specifies the name prefixes of objects to be migrated if `type` is set to
  **PREFIX**. If you want to migrate the entire bucket, set this parameter to [""].
  Changing this creates a new resource.

* `data_source` - (Optional, String, ForceNew) Specifies the source cloud service provider. If `type` is
  **URL_LIST**,set this parameter to **URLSource**. The value can be **AWS**, **Azure**, **Aliyun**, **Tencent**,
  **HuaweiCloud**, **QingCloud**, **KingsoftCloud**, **Baidu**, **Qiniu**, **URLSource** or **UCloud**.
  The default value is **Aliyun**. Changing this creates a new resource.

* `app_id` - (Optional, String, ForceNew) Specifies the APP ID. This parameter is mandatory when `data_source` is
  **Tencent**. Changing this creates a new resource.

* `list_file_bucket` - (Optional, String, ForceNew) Specifies the name of the OBS bucket for storing the object
  list files. `list_file_bucket` is mandatory when `type` is set to **LIST** or **URL_LIST**.
  Changing this creates a new resource.

  -> Ensure that the OBS bucket is in the same region as the destination bucket, or the task group will fail to be
  created.

* `list_file_key` - (Optional, String, ForceNew) Specifies the OBS bucket folder name of the list file or URL list file.
  `list_file_key` is mandatory when `type` is set to **LIST** or **URL_LIST**.

  + If `type` is **LIST**: You need to write the names of source objects to be migrated into an object list file
  and store the file in an OBS bucket on HUAWEI CLOUD. OMS migrates all objects specified in the object list file.

  + If `type` is **URL_LIST**: You need to write the URLs of the files to be migrated and their destination objects
  names into one or more .txt URL list files and store the files in an OBS bucket on HUAWEI CLOUD. You can store up
  to 2,000 list files in a fixed folder in the OBS bucket. Each list file cannot exceed 1 GB. OMS migrates all
  objects specified in the URL list files.

  Changing this creates a new resource.

  -> More details for the format requirements of list file. Please see
  the [User Guide](https://support.huaweicloud.com/intl/en-us/usermanual-oms/oms_01_0017.html).

<a name="destination_object_object"></a>
The `destination_object` block supports:

* `region` - (Required, String, ForceNew) Specifies the region where the destination bucket is located.
  The value must be the same as that of the service endpoint. Changing this creates a new resource.

* `bucket` - (Required, String, ForceNew) Specifies the name of the destination bucket.
  Changing this creates a new resource.

* `access_key` - (Optional, String, ForceNew) Specifies the access key for accessing the destination bucket.
  Changing this creates a new resource.

* `secret_key` - (Optional, String, ForceNew) Specifies the secret key for accessing the destination bucket.
  Changing this creates a new resource.

* `data_source` - (Optional, String, ForceNew) Specifies the destination data source. The default value is **HEC**.
  Changing this creates a new resource.

* `save_prefix` - (Optional, String, ForceNew) Specifies the path prefix in the destination bucket. The prefix is added
  before the object key to form a new key. Changing this creates a new resource.

<a name="bandwidth_policy_object"></a>
The `bandwidth_policy` block supports:

* `max_bandwidth` - (Required, Int) Specifies the maximum traffic bandwidth allowed in the specified time
  segment. The value ranges from `1` to `200`. The unit is MB/s.

* `start` - (Required, String) Specifies the start time of the traffic limit rule. The format is **hh:mm**,
  e.g. **12:03**.

* `end` - (Required, String) Specifies the end time of the traffic limit rule. The format is **hh:mm**,
  e.g. **12:03**.

<a name="source_cdn_object"></a>
The `source_cdn` block supports:

* `domain` - (Required, String, ForceNew) Specifies the domain name from which to obtain objects to be migrated.
  Changing this creates a new resource.

* `protocol` - (Required, String, ForceNew) Specifies the protocol type. Valid values are **HTTP** and **HTTPS**.
  Changing this creates a new resource.

* `authentication_type` - (Optional, String, ForceNew) Specifies the authentication type. Valid values are **NONE**,
  **QINIU_PRIVATE_AUTHENTICATION**, **ALIYUN_OSS_A**, **ALIYUN_OSS_B**, **ALIYUN_OSS_C**,
  **KSYUN_PRIVATE_AUTHENTICATION**, **TENCENT_COS_A**, **TENCENT_COS_B**, **TENCENT_COS_C**,
  **TENCENT_COS_D**. Default value: **None**. Changing this creates a new resource.

* `authentication_key` - (Optional, String, ForceNew) Specifies the CDN authentication key.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the migration task group.

* `status` - The migration task group status. The value can be:
  + **0**: Waiting.
  + **1**: Executing/Creating.
  + **2**: Monitoring.
  + **3**: Paused.
  + **4**: Creation failed.
  + **5**: Migration failed.
  + **6**: Migration completed.
  + **7**: Pausing.
  + **8**: Waiting to be deleted.
  + **9**: Deleted.

* `total_time` - The total amount of time used by the migration task group, in ms.
* `total_num` - The total number of objects to be migrated in the migration task group.
* `success_num` - The number of migrated objects.
* `fail_num` - The number of failed objects.
* `total_size` - The total size of migrated objects, in bytes.
* `complete_size` - The size (in bytes) of the objects that have been migrated.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.
