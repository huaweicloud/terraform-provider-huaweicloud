---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vbs_backup_v2"
sidebar_current: "docs-huaweicloud-datasource-vbs-backup-v2"
description: |-
  Provides details about a specific VBS Backup.
---

# Data Source: huaweicloud_vbs_backup_v2

The VBS Backup data source provides details about a specific VBS Backup.

## Example Usage

```hcl
variable "backup_id" {}

data "huaweicloud_vbs_backup_v2" "mybackup" {
  id = "${var.backup_id}"
}
```

## Argument Reference
The following arguments are supported:

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
