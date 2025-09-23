---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_backups"
description: |-
  Use this data source to query the CBR backups within HuaweiCloud.
---

# huaweicloud_cbr_backups

Use this data source to query the CBR backups within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cbr_backups" "test" {}
```

## Argument reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the backup detail.
  If omitted, the provider-level region will be used.

* `checkpoint_id` - (Optional, String) Specifies the restore point ID.

* `dec` - (Optional, Bool) Specifies the dedicated cloud tag, which only takes effect in dedicated cloud scenarios.

* `end_time` - (Optional, String) Specifies the time when the backup ends. In `%YYYY-%mm-%ddT%HH:%MM:%SSZ` format.
  For example, **2018-02-01T12:00:00Z**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  **all_granted_eps** indicates querying the IDs of all enterprise projects on which the user has permissions.

* `image_type` - (Optional, String) Specifies the backup type. The value can be **backup** or **replication**.

* `incremental` - (Optional, Bool) Specifies whether incremental backup is used.

* `member_status` - (Optional, String) Specifies the backup sharing status.

* `name` - (Optional, String) Specifies the backup name.

* `own_type` - (Optional, String) Specifies the owning type of backup. Valid values are **private**, **shared**,
  and **all_granted**. Private backups are queried by default.

* `parent_id` - (Optional, String) Specifies the parent backup ID.

* `resource_az` - (Optional, String) Specifies the resource availability zones.

* `resource_id` - (Optional, String) Specifies the resource ID.

* `resource_name` - (Optional, String) Specifies the resource name.

* `resource_type` - (Optional, String) Specifies the resource type. Valid values are **OS::Nova::Server**,
  **OS::Cinder::Volume**, **OS::Ironic::BareMetalServer**, **OS::Native::Server**, **OS::Sfs::Turbo**,
  **OS::Workspace::DesktopV2**.

* `show_replication` - (Optional, Bool) Specifies whether to show replication records.

* `sort` - (Optional, String) Specifies the sort key. A group of properties separated by commas (,) and sorting directions.
  The value is in the format of `<key1>[:<direction>],<key2>[:<direction>]`, where the value of direction is asc
  (ascending order) or desc (descending order). If a direction is not specified, the default sorting direction is desc.
  The value of sort can contain a maximum of `255` characters. The key can be as follows: **created_at**, **updated_at**,
  **name**, **status**, **protected_at**, and **id**.

* `start_time` - (Optional, String) Specifies the time when the backup starts. In `%YYYY-%mm-%ddT%HH:%MM:%SSZ` format.
  For example, **2018-02-01T12:00:00Z**.

* `status` - (Optional, String) Specifies the status. When the API is called, multiple statuses can be transferred for
  filtering, separated by commas (,). For example, **available,error**.

* `used_percent` - (Optional, String) Specifies the using percent of the occupied vault capacity. The value ranges from
  `1` to `100`. For example, if used_percent is set to `80`, all backups who occupied `80%` or more of the vault capacity
  are displayed.

* `vault_id` - (Optional, String) Specifies the vault ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `backups` - The backup list.

  The [backups](#backups_struct) structure is documented below.

<a name="backups_struct"></a>
The `backups` block supports:

* `checkpoint_id` - The restore point ID.

* `created_at` - The creation time.

* `description` - The backup description.

* `expired_at` - The expiration time.

* `extend_info` - The extended information.

  The [extend_info](#extend_info_struct) structure is documented below.

* `id` - The backup ID.

* `image_type` - The backup type.

* `name` - The backup name.

* `parent_id` - The parent backup ID.

* `project_id` - The project ID.

* `protected_at` - The backup time.

* `resource_az` - The resource availability zone.

* `resource_id` - The resource ID.

* `resource_name` - The resource name.

* `resource_size` - The resource size, in GB.

* `resource_type` - The resource type.

* `status` - The backup status.

* `updated_at` - The update time.

* `vault_id` - The vault ID.

* `replication_records` - The replication record.

  The [replication_records](#replication_records_struct) structure is documented below.

* `enterprise_project_id` - The enterprise project ID.

* `provider_id` - The backup provider ID.

* `children` - The children backup list. This field is JSON format string.
  Its JSON structure is the same as the field `backups`.

* `incremental` - Whether incremental backup is used.

* `version` - The backup snapshot type.

<a name="extend_info_struct"></a>
The `extend_info` block supports:

* `auto_trigger` - Whether the backup is automatically generated.

* `bootable` - Whether the backup is a system disk backup.

* `snapshot_id` - Snapshot ID of the disk backup.

* `support_lld` - Whether to allow lazy loading for fast restoration.

* `supported_restore_mode` - The restoration mode.

* `os_images_data` - The ID list of images created using backups.

  The [os_images_data](#os_images_data_struct) structure is documented below.

* `contain_system_disk` - Whether the VM backup data contains system disk data.

* `encrypted` - Whether the backup is encrypted.

* `system_disk` - Whether the disk is a system disk.

* `is_multi_az` - Whether multi-AZ backup redundancy is used.

<a name="os_images_data_struct"></a>
The `os_images_data` block supports:

* `image_id` - The image ID.

<a name="replication_records_struct"></a>
The `replication_records` block supports:

* `created_at` - The start time of the replication.

* `destination_backup_id` - The ID of the destination backup used for replication.

* `destination_checkpoint_id` - The record ID of the destination backup used for replication.

* `destination_project_id` - The ID of the replication destination project.

* `destination_region` - The replication destination region.

* `destination_vault_id` - The destination vault ID.

* `extra_info` - The additional information of the replication.

  The [extra_info](#extra_info_struct) structure is documented below.

* `id` - The replication record ID.

* `source_backup_id` - The ID of the source backup used for replication.

* `source_checkpoint_id` - The ID of the source backup record used for replication.

* `source_project_id` - The ID of the replication source project.

* `source_region` - The replication source region.

* `status` - The replication status.

* `vault_id` - The ID of the vault where the backup resides.

<a name="extra_info_struct"></a>
The `extra_info` block supports:

* `progress` - The replication progress.

* `fail_code` - The error code. This field is empty if the operation is successful.

* `fail_reason` - The error cause.

* `auto_trigger` - Whether replication is automatically scheduled.

* `destinatio_vault_id` - The destination vault ID.
