---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_restore"
description: |-
  Using this resource to restore a CBR backup within HuaweiCloud.
---

# huaweicloud_cbr_restore

Using this resource to restore a CBR backup within HuaweiCloud.

-> This resource is only a one-time action resource to restore a CBR backup. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.

-> For backup and recovery constraints and limitations, please refer to the following documents:
<br/>1. [Restoring from a Cloud Server Backup](https://support.huaweicloud.com/intl/en-us/usermanual-cbr/cbr_03_0032.html)
<br/>2. [Restoring from a Cloud Disk Backup](https://support.huaweicloud.com/intl/en-us/usermanual-cbr/cbr_03_0033.html)
<br/>3. [Restoring from a Desktop Backup](https://support.huaweicloud.com/intl/en-us/usermanual-cbr/cbr_03_0110.html)

## Example Usage

## Restoring an ECS backup

```hcl
variable "ecs_backup_id" {}
variable "ecs_server_id" {}
variable "evs_backup_id" {}
variable "evs_volume_id" {}

resource "huaweicloud_cbr_restore" "test" {
  backup_id = var.ecs_backup_id
  server_id = var.ecs_server_id
  power_on  = true

  mappings {
    backup_id = var.evs_backup_id
    volume_id = var.evs_volume_id
  }
}
```

## Restoring an EVS backup

```hcl
variable "evs_backup_id" {}
variable "evs_volume_id" {}

resource "huaweicloud_cbr_restore" "test" {
  backup_id = var.evs_backup_id
  volume_id = var.evs_volume_id
}
```

## Restoring a Workspace backup

```hcl
variable "workspace_backup_id" {}
variable "workspace_resource_id" {}

resource "huaweicloud_cbr_restore" "test" {
  backup_id   = var.workspace_backup_id
  resource_id = var.workspace_resource_id
  power_on    = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `backup_id` - (Required, String, NonUpdatable) Specifies the backup ID.

* `mappings` - (Optional, List, NonUpdatable) Specifies the restored mapping relationship. This parameter is mandatory for
  VM restoration and optional for disk restoration.
  You can obtain the disk information through data source `huaweicloud_cbr_backup`.

  The [mappings](#restore_mappings_struct) structure is documented below.

* `power_on` - (Optional, Bool, NonUpdatable) Whether the server is powered on after restoration. Defaults to **false**.

* `server_id` - (Optional, String, NonUpdatable) Specifies the ID of the target VM to be restored. This parameter is
  mandatory for VM restoration.

* `volume_id` - (Optional, String, NonUpdatable) Specifies the ID of the target disk to be restored. This parameter is
  mandatory for disk restoration.

* `resource_id` - (Optional, String, NonUpdatable) Specifies the ID of the resource to be restored.

* `details` - (Optional, List, NonUpdatable) Specifies the restoration details.

  The [details](#restore_details_struct) structure is documented below.

<a name="restore_details_struct"></a>
The `details` block supports:

* `destination_path` - (Required, String, NonUpdatable) Specifies the destination path.

<a name="restore_mappings_struct"></a>
The `mappings` block supports:

* `backup_id` - (Required, String, NonUpdatable) Specifies the disk backup ID.

* `volume_id` - (Required, String, NonUpdatable) Specifies the ID of the disk to which data is restored.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also `backup_id`).

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
