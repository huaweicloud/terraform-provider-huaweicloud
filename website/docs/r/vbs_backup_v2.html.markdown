---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vbs_backup_v2"
sidebar_current: "docs-huaweicloud-datasource-vbs-backup-v2"
description: |-
  Provides an VBS Backup resource.
---

# huaweicloud_vbs_backup_v2

Provides an VBS Backup resource.
 
# Example Usage

 ```hcl
variable "backup_name" {}

variable "volume_id" {}
 
resource "huaweicloud_vbs_backup_v2" "mybackup" {
  volume_id = "${var.volume_id}"
  name = "${var.backup_name}"
}
 ```

# Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the vbs backup. Changing the parameter creates a new backup.

* `volume_id` - (Required) The id of the disk to be backed up. Changing the parameter creates a new backup.

* `snapshot_id` - (Optional) The snapshot id of the disk to be backed up. Changing the parameter creates a new backup.

* `description` - (Optional) The description of the vbs backup. Changing the parameter creates a new backup.

**tags** **- (Optional)** List of tags to be configured for the backup resources. Changing the parameter creates a new backup.

* `key` - (Required) Specifies the tag key. Changing the parameter creates a new backup

* `value` - (Required) Specifies the tag value. Changing the parameter creates a new backup

# Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the vbs backup.

* `container` - The container of the backup.

* `status` - The status of the VBS backup.

* `availability_zone` - The AZ where the backup resides.

* `fail_reason` - Cause of the backup failure.

* `size` - The size of the vbs backup.

* `object_count` - Number of objects on Object Storage Service (OBS) for the disk data.

* `tenant_id` - The ID of the tenant to which the backup belongs.

* `service_metadata` - The metadata of the vbs backup.

 
# Import

VBS Backup can be imported using the `backup id`, e.g.

```
 $ terraform import huaweicloud_vbs_backup_v2.mybackup 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```