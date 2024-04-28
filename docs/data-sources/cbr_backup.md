---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_backup"
description: ""
---

# huaweicloud_cbr_backup

Use this data source to query the backup detail using its ID within Huaweicloud.

## Example Usage

### Using backup ID to query the backup detail

```hcl
variable "backup_id" {}

data "huaweicloud_cbr_backup" "test" {
  id = "backup_id"
}
```

## Argument reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the backup detail.
  If omitted, the provider-level region will be used.

* `id` - (Required, String) Specifies the backup ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `parent_id` - The parent backup ID.

* `type` - The backup type.

* `name` - The backup name.

* `description` - The backup description.

* `checkpoint_id` - The restore point ID.

* `resource_id` - The backup resource ID.

* `resource_type` - The backup resource type.

* `resource_size` - The backup resource size, in GB.

* `resource_name` - The backup resource name.

* `resource_az` - The availability zone where the backup resource is located.

* `enterprise_project_id` - The enterprise project to which the backup resource belongs.

* `vault_id` - The vault to which the backup resource belongs.

* `replication_records` - The replication records.
  The [object](#cbr_backup_replication_records) structure is documented below.

* `children` - The backup list of the sub-backup resources.
  The [object](#cbr_backup_children) structure is documented below.

* `extend_info` - The extended information.
  The [object](#cbr_backup_extend_info) structure is documented below.

* `status` - The backup status.

* `created_at` - The creation time of the backup.

* `updated_at` - The latest update time of the backup.

* `expired_at` - The expiration time of the backup.

<a name="cbr_backup_replication_records"></a>
The `replication_records` block supports:

* `id` - The replication record ID.

* `destination_backup_id` - The ID of the destination backup used for replication.

* `destination_checkpoint_id` - The record ID of the destination backup used for replication.

* `destination_project_id` - The ID of the replication destination project.

* `destination_region` - The replication destination region.

* `destination_vault_id` - The destination vault ID.

* `source_backup_id` - The ID of the source backup used for replication.

* `source_checkpoint_id` - The ID of the source backup record used for replication.

* `source_project_id` - The ID of the replication source project.

* `source_region` - The replication source region.

* `status` - The replication status.

* `vault_id` - The ID of the vault where the backup resides.

* `extra_info` - The additional information of the replication.
  The [object](#cbr_backup_replication_record_extra_info) structure is documented below.

* `created_at` - The creation time of the replication.

<a name="cbr_backup_replication_record_extra_info"></a>
The `extra_info` block supports:

* `progress` - The replication progress.

* `fail_code` - The error code.

* `fail_reason` - The error cause.

* `auto_trigger` - Whether replication is automatically scheduled.

* `destination_vault_id` - The destination vault ID.

<a name="cbr_backup_children"></a>
The `children` block supports:

* `id` - The sub-backup ID.

* `name` - The sub-backup name.

* `description` - The sub-backup description.

* `type` - The sub-backup type.

* `checkpoint_id` - The restore point ID of the sub-backup resource.

* `resource_id` - The sub-backup resource ID.

* `resource_type` - The sub-backup resource type.

* `resource_size` - The sub-backup resource size, in GB.

* `resource_name` - The sub-backup resource name.

* `resource_az` - The availability zone where the backup sub-backup resource is located.

* `enterprise_project_id` - The enterprise project to which the backup sub-backup resource belongs.

* `vault_id` - The vault to which the backup resource belongs.

* `replication_records` - The replication records.

* `extend_info` - The extended information.

* `status` - The sub-backup status.

* `created_at` - The creation time of the sub-backup.

* `updated_at` - The latest update time of the sub-backup.

* `expired_at` - The expiration time of the sub-backup.

<a name="cbr_backup_extend_info"></a>
The `extend_info` block supports:

* `auto_trigger` - Whether the backup is automatically generated.

* `bootable` - Whether the backup is a system disk backup.

* `incremental` - Whether the backup is an incremental backup.

* `snapshot_id` - Snapshot ID of the disk backup.

* `support_lld` - Whether to allow lazy loading for fast restoration.

* `supported_restore_mode` - The restoration mode.

* `os_registry_images` - The ID list of images created using backups.

* `contain_system_disk` - Whether the VM backup data contains system disk data.

* `encrypted` - Whether the backup is encrypted.

* `is_system_disk` - Whether the disk is a system disk.
