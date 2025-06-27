---
subcategory: "Object Storage Migration Service (OMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_oms_migration_task"
description: ""
---

# huaweicloud_oms_migration_task

Manages an OMS migration task resource within HuaweiCloud.

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
variable "topic_urn" {}

resource "huaweicloud_oms_migration_task" "test" {
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

  type        = "object"
  description = "test task"

  bandwidth_policy {
    max_bandwidth = 2
    start         = "15:00"
    end           = "16:00"
  }

  smn_config {
    topic_urn          = var.topic_urn
    trigger_conditions = ["FAILURE", "SUCCESS"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `source_object` - (Required, List, ForceNew) Specifies the source information. The [object](#source_object_object)
  structure is documented below. Changing this creates a new resource.

* `destination_object` - (Required, List, ForceNew) Specifies the destination information. The [object](#destination_object_object)
  structure is documented below. Changing this creates a new resource.

* `type` - (Required, String, ForceNew) Specifies the task type. The value can be:
  + **list**: indicates migrating objects using an object list.
  + **url_list**: indicates migrating objects using a URL object list.
  + **object**: indicates migrating selected files or folders.
  + **prefix**: indicates migrating objects with specified prefixes.
  
  Changing this creates a new resource.

* `start_task` - (Optional, Bool) Specifies whether to start the task. Default value: **true**.

* `enable_kms` - (Optional, Bool, ForceNew) Specifies whether to enable the KMS encryption function.
  Default value: **false**. Changing this creates a new resource.

* `description` - (Optional, String, ForceNew) Specifies the description of the task.
  Changing this creates a new resource.

* `migrate_since` - (Optional, String, ForceNew) Specifies a time in format **yyyy-MM-dd HH:mm:ss**,
  e.g. **2006-01-02 15:04:05**. The system migrates only the objects that are modified after the specified time.
  No time is specified by default. Changing this creates a new resource.

* `enable_restore` - (Optional, Bool, ForceNew) Specifies whether to automatically restore the archive data. If enabled,
  archive data is automatically restored and migrated. Default value: **false**. Changing this creates a new resource.

* `enable_failed_object_recording` - (Optional, Bool, ForceNew) Specifies whether to record failed objects. If this
  function is enabled, information about objects that fail to be migrated will be stored in the destination bucket.
  Default value: **true**. Changing this creates a new resource.

* `bandwidth_policy` - (Optional, List) Specifies the traffic limit rules. Each element in the array
  corresponds to the maximum bandwidth of a time segment. A maximum of 5 time segments are allowed, and the time
  segments must not overlap. The [object](#bandwidth_policy_object) structure is  documented below.

* `source_cdn` - (Optional, List, ForceNew) Specifies the CDN information. If this parameter is contained,
  using CDN to download source data is supported, the source objects to be migrated are obtained from the CDN domain
  name during migration. The [object](#source_cdn_object) structure is documented below.
  Changing this creates a new resource.

* `smn_config` - (Optional, List, ForceNew) Specifies the SMN message sending configuration.
  The [object](#smn_config_object) structure is  documented below. Changing this creates a new resource.

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

* `enable_metadata_migration` - (Optional, Bool, ForceNew) Specifies whether metadata migration is enabled. Even if this
  function is disabled, the ContentType metadata will still be migrated to ensure a successful migration.
  Default value: **false**. Changing this creates a new resource.

* `task_priority` - (Optional, String, ForceNew) Specifies the task priority.
  The value can be **HIGH**, **MEDIUM**, or **LOW**. Changing this creates a new resource.

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

* `region` - (Optional, String, ForceNew) Specifies the region where the source bucket is located. `region` is mandatory
  when `type` is not **url_list**. Changing this creates a new resource.

* `bucket` - (Optional, String, ForceNew) Specifies the name of the source bucket. `bucket` is mandatory when `type`
  is not **url_list**. Changing this creates a new resource.

* `access_key` - (Optional, String, ForceNew) Specifies the access key for accessing the source bucket. This parameter
  is mandatory when `type` is not **url_list**. Changing this creates a new resource.

* `secret_key` - (Optional, String, ForceNew) Specifies the secret key for accessing the destination bucket. This
  parameter is mandatory when `type` is not **url_list**. Changing this creates a new resource.

* `security_token` - (Optional, String, ForceNew) Specifies the temporary token for accessing the source bucket.
  Changing this creates a new resource.

* `object` - (Optional, List, ForceNew) Specifies the list of object keys.
  + If `type` is set to **object**, this parameter specifies the names of the objects to be migrated. The strings
  ending with a slash (/) indicate the folders to be migrated, and the strings not ending with a slash (/) indicate the
  files to be migrated.
  + If `type` is set to **prefix**, this parameter indicates the name prefixes of the objects to be migrated.
  Set this parameter to [""] to migrate the entire bucket
  
  Changing this creates a new resource.

* `data_source` - (Optional, String, ForceNew) Specifies the source cloud service provider. If `type` is
  **url_list**,set this parameter to **URLSource**. The value can be **AWS**, **Azure**, **Aliyun**, **Tencent**,
  **HuaweiCloud**, **QingCloud**, **KingsoftCloud**, **Baidu**, **Qiniu**, **URLSource** and **UCloud**.
  The default value is **Aliyun**. Changing this creates a new resource.

* `app_id` - (Optional, String, ForceNew) Specifies the APP ID. This parameter is mandatory when `data_source` is  
  **Tencent**. Changing this creates a new resource.

* `list_file_bucket` - (Optional, String, ForceNew) Specifies the name of the OBS bucket for storing the object list files.
  `list_file_bucket` is mandatory when `type` is set to **list** or **url_list**. Changing this creates a new resource.
  
  -> Ensure that the OBS bucket is in the same region as the destination bucket, or the task will fail to be created.

* `list_file_key` - (Optional, String, ForceNew) Specifies the object name of the list file or URL list file.
  `list_file_key` is mandatory when `type` is set to **list** or **url_list**. Changing this creates a new resource.

* `list_file_num` - (Optional, String, ForceNew) Specifies the number of stored object list files.
  Changing this creates a new resource.

* `json_auth_file` - (Optional, String, ForceNew) Specifies the file used for Google Cloud Storage authentication.
  Changing this creates a new resource.

<a name="destination_object_object"></a>
The `destination_object` block supports:

* `region` - (Required, String, ForceNew) Specifies the region where the destination bucket is located.
  Changing this creates a new resource.

* `bucket` - (Required, String, ForceNew) Specifies the name of the destination bucket.
  Changing this creates a new resource.

* `access_key` - (Required, String, ForceNew) Specifies the access key for accessing the destination bucket.
  Changing this creates a new resource.

* `secret_key` - (Required, String, ForceNew) Specifies the secret key for accessing the destination bucket.
  Changing this creates a new resource.

* `security_token` - (Optional, String, ForceNew) Specifies the temporary token for accessing the destination bucket.
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

<a name="smn_config_object"></a>
The `smn_config` block supports:

* `topic_urn` - (Required, String, ForceNew) Specifies the SMN message topic URN bound to a migration task.
  Changing this creates a new resource.

* `trigger_conditions` - (Required, List, ForceNew) Specifies the trigger conditions of sending messages using SMN.
  The value can be:
  + **FAILURE**: indicates that an SMN message will be sent after the migration task fails.
  + **SUCCESS**: indicates that an SMN message will be sent after the migration task succeeds.

  Changing this creates a new resource.

* `language` - (Optional, String, ForceNew) Specifies the SMN message language. The value can be **zh-cn** or
  **en-us**. Default value: **en-us**. Changing this creates a new resource.

* `message_template_name` - (Optional, String, ForceNew) Specifies the message template name.
  If this parameter is specified, SMN messages are sent using the specified template.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the migration task.

* `name` - The name of the migration task.

* `status` - The status the migration task. The value can be:
  + **1**: Waiting to migrate.
  + **2**: Migrating.
  + **3**: Migration paused.
  + **4**: Migration failed.
  + **5**: Migration succeeded.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
