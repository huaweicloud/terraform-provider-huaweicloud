---
subcategory: "Deprecated"
---

# huaweicloud\_vbs\_backup

!> **Warning:** It has been deprecated.

The VBS Backup data source provides details about a specific VBS Backup.
This is an alternative to `huaweicloud_vbs_backup`

## Example Usage

```hcl
variable "backup_id" {}

data "huaweicloud_vbs_backup" "mybackup" {
  id = var.backup_id
}
```

## Argument Reference
The following arguments are supported:

* `region` - (Optional) The region in which to obtain the VBS Backup. If omitted, the provider-level region will work as default.

* `id` - (Optional) The id of the vbs backup.

* `name` - (Optional) The name of the vbs backup.

* `volume_id` - (Optional) The source volume ID of the backup.

* `snapshot_id` - (Optional) ID of the snapshot associated with the backup.

* `status` - (Optional) The status of the VBS backup.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `description` - The description of the vbs backup.

* `availability_zone` - The AZ where the backup resides.

* `size` - The size of the vbs backup.

* `container` - The container of the backup.

* `service_metadata` - The metadata of the vbs backup.
