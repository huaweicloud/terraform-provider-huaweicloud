---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_backup_sync"
description: |-
  Using this resource to synchronize VMware backups within HuaweiCloud.
---

# huaweicloud_cbr_backup_sync

Using this resource to synchronize VMware backups within HuaweiCloud.

-> This resource is only a one-time action resource to synchronize a VMware backup. Deleting this resource will
not clear the corresponding backup synchronization record, but will only remove the resource information from the
tfstate file.

## Example Usage

```hcl
variable "backup_id" {}
variable "backup_name" {}
variable "bucket_name" {}
variable "image_path" {}
variable "resource_id" {}
variable "resource_name" {}
variable "resource_type" {}
variable "created_at" {}

resource "huaweicloud_cbr_backup_sync" "test" {
  backup_id     = var.backup_id
  backup_name   = var.backup_name
  bucket_name   = var.bucket_name
  image_path    = var.image_path
  resource_id   = var.resource_id
  resource_name = var.resource_name
  resource_type = var.resource_type
  created_at    = var.created_at
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `backup_id` - (Required, String, NonUpdatable) Specifies the ID of the backup to be synchronized.

* `backup_name` - (Required, String, NonUpdatable) Specifies the name of the backup.

* `bucket_name` - (Required, String, NonUpdatable) Specifies the name of the bucket where the backup is stored.

* `image_path` - (Required, String, NonUpdatable) Specifies the path of the backup image in the bucket.

* `resource_id` - (Required, String, NonUpdatable) Specifies the ID of the resource to be backed up.

* `resource_name` - (Required, String, NonUpdatable) Specifies the name of the resource to be backed up.

* `resource_type` - (Required, String, NonUpdatable) Specifies the type of the resource to be backed up.

* `created_at` - (Required, Int, NonUpdatable) Specifies the timestamp when the backup was created,
  for example `1,548,898,428`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the backup synchronization resource.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
