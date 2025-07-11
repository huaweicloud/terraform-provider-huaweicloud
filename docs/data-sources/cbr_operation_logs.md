---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_operation_logs"
description: |-
  Use this data source to get CBR operation logs within HuaweiCloud.
---

# huaweicloud_cbr_operation_logs

Use this data source to get CBR operation logs within HuaweiCloud.

## Example Usage

```hcl
variable enterprise_project_id {}
variable operation_type {}
variable provider_id {}
variable resource_id {}
variable resource_name {}

data "huaweicloud_cbr_operation_logs" "test" {
  enterprise_project_id = var.enterprise_project_id
  operation_type        = var.operation_type
  provider_id           = var.provider_id
  resource_id           = var.resource_id
  resource_name         = var.resource_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `end_time` - (Optional, String) Specifies the end time of a task. The time is in the
  **%YYYY-%mm-%ddT%HH:%MM:%SSZ** format.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID. **all_granted_eps**
  indicates querying the IDs of all enterprise projects on which the user has permissions.

* `operation_type` - (Optional, String) Specifies the task type. Possible values are **backup**, **copy**,
  **replication**, **restore**, **delete**, **sync**, **vault_delete**, or **remove_resource**.

* `provider_id` - (Optional, String) Specifies the backup provider ID, which specifies whether the backup object is a
  server or disk.

* `resource_id` - (Optional, String) Specifies the backup resource ID which the target resource is associated.

* `resource_name` - (Optional, String) Specifies the backup resource name which the target resource is associated.

* `start_time` - (Optional, String) Specifies the start time of a task. The time is in the
  **%YYYY-%mm-%ddT%HH:%MM:%SSZ** format.

* `status` - (Optional, String) Specifies the task status. Possible values are **success**, **skipped**, **failed**,
  **running**, **timeout** or **waiting**.

* `vault_id` - (Optional, String) Specifies the ID of the vault with which the target resource is associated.

* `vault_name` - (Optional, String) Specifies the name of the vault with which the target resource is associated.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `operation_logs` - The task list.
  The [operation_logs](#operation_logs_struct) structure is documented below.

<a name="operation_logs_struct"></a>
The `operation_logs` block supports:

* `checkpoint_id` - The backup record ID.

* `created_at` - The creat time of the task.

* `ended_at` - The end time of the task.

* `error_info` - The task error message.
  The [error_info](#error_info_struct) structure is documented below.

* `extra_info` - The task extension information.
  The [extra_info](#extra_info_struct) structure is documented below.

* `id` - The task ID.

* `operation_type` - The task type. Possible values are **backup**, **copy**, **replication**, **restore**, **delete**,
  **sync**, **vault_delete**, or **remove_resource**.

* `policy_id` - The policy ID.

* `provider_id` - The backup provider ID, which specifies whether the backup object is a server or disk.

* `started_at` - The start time of the task.

* `status` - The task status. Possible values are **success**, **skipped**, **failed**, **running**, **timeout**
  or **waiting**.

* `updated_at` - The modification time of the task.

* `vault_id` - The ID of the vault with which the target resource is associated.

* `vault_name` - The name of the vault with which the target resource is associated.

<a name="error_info_struct"></a>
The `error_info` block supports:

* `code` - The error code. For details, see "CBR Error Codes" in Cloud Backup and Recovery User Guide.
  [reference](https://support.huaweicloud.com/intl/en-us/api-cbr/ErrorCode.html)

* `message` - The error message.

<a name="extra_info_struct"></a>
The `extra_info` block supports:

* `backup` - The extended parameters of backup.
  The [backup](#backup_struct) structure is documented below.

* `common` - The common parameters.
  The [common](#common_struct) structure is documented below.

* `delete` - The extended parameters of deletion.
  The [delete](#delete_struct) structure is documented below.

* `sync` - The extended parameters of synchronization.
  The [sync](#sync_struct) structure is documented below.

* `remove_resources` - The extended parameters of removing resources from a vault.
  The [remove_resources](#remove_resources_struct) structure is documented below.

* `replication` - The extended parameters of replication.
  The [replication](#replication_struct) structure is documented below.

* `resource` - The resource information.
  The [resource](#resource_struct) structure is documented below.

* `restore` - The extended parameters of restoration.
  The [restore](#restore_struct) structure is documented below.

* `vault_delete` - The extended parameters of deleting a vault.
  The [vault_delete](#vault_delete_struct) structure is documented below.

<a name="backup_struct"></a>
The `backup` block supports:

* `app_consistency_error_code` - The error code returned if application-consistent backup fails.
  For details, see "CBR Error Codes" in Cloud Backup and Recovery User Guide.
  [reference](https://support.huaweicloud.com/intl/en-us/api-cbr/ErrorCode.html)

* `app_consistency_error_message` - The error message returned if application-consistent backup fails.

* `app_consistency_status` - The application-consistent backup status.

* `backup_id` - The backup ID.

* `backup_name` - The backup name.

* `incremental` - Whether incremental backup is used.

<a name="common_struct"></a>
The `common` block supports:

* `progress` - The progress of the query task. The value ranges from `0` to `100`.

* `request_id` - The request ID.

* `task_id` - The backup task ID.

<a name="delete_struct"></a>
The `delete` block supports:

* `backup_id` - The backup ID.

* `backup_name` - The backup name.

<a name="sync_struct"></a>
The `sync` block supports:

* `sync_backup_num` - The number of synchronized backups.

* `delete_backup_num` - The number of deleted backups.

* `err_sync_backup_num` - The number of backups that failed to be synchronized.

<a name="remove_resources_struct"></a>
The `remove_resources` block supports:

* `fail_count` - The number of resources that fail to be deleted.

* `total_count` - The number of deleted backups.

* `resources` - The resource information.
  The [resource](#resource_struct) structure is documented below.

<a name="replication_struct"></a>
The `replication` block supports:

* `destination_backup_id` - The destination backup ID.

* `destination_checkpoint_id` - The destination restore point ID.

* `destination_project_id` - The destination project ID.

* `destination_region` - The destination region.

* `source_backup_id` - The source backup ID.

* `source_checkpoint_id` - The source restore point ID.

* `source_project_id` - The source project ID.

* `source_region` - The source region.

* `source_backup_name` - The source backup name.

* `destination_backup_name` - The destination backup name.

<a name="resource_struct"></a>
The `resource` block supports:

* `extra_info` - The extra information of the resource.
  The [resource_extra_info](#resource_extra_info_struct) structure is documented below.

* `id` - The ID of the resource to be backed up.

* `name` - The name of the resource to be backed up. The value consists of `0` to `255` characters.

* `type` - The resource type. Possible values are **OS::Nova::Server**, **OS::Cinder::Volume**,
  **OS::Ironic::BareMetalServer**, **OS::Native::Server**, **OS::Sfs::Turbo** or **OS::Workspace::DesktopV2**.

<a name="resource_extra_info_struct"></a>
The `resource_extra_info` block supports:

* `exclude_volumes` - The IDs of the disks that will not be backed up. This parameter is used when servers are added to
  a vault, which include all server disks. But some disks do not need to be backed up. Or in case that a server was
  previously added and some disks on this server do not need to be backed up.

<a name="restore_struct"></a>
The `restore` block supports:

* `backup_id` - The backup ID.

* `backup_name` - The backup name.

* `target_resource_id` - The ID of the resource to be restored.

* `target_resource_name` - The name of the resource to be restored.

<a name="vault_delete_struct"></a>
The `vault_delete` block supports:

* `fail_delete_count` - The number of resources that fail to be deleted in this task.

* `total_delete_count` - The number of backups deleted in this task.
