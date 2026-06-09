---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_backups_batch_delete"
description: |-
  Use this resource to batch delete TaurusDB manual backups within HuaweiCloud.
---

# huaweicloud_taurusdb_backups_batch_delete

Use this resource to batch delete TaurusDB manual backups within HuaweiCloud.

-> This resource is a one-time action resource for batch deleting TaurusDB manual backups. Deleting this resource will
   not restore the deleted backups, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "backup_ids" {
  type = list(string)
}

resource "huaweicloud_taurusdb_backups_batch_delete" "test" {
  backup_ids = var.backup_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to batch delete the TaurusDB backups.
  If omitted, the provider-level region will be used.

* `backup_ids` - (Required, List, NoneUpdatable) Specifies the IDs of the manual backups to be deleted.
  A maximum of 50 IDs are allowed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `success_count` - Indicates the number of backups that have been deleted.

* `failed_count` - Indicates the number of backups that failed to be deleted.

* `failed_results` - Indicates the backup deletion errors.
  The [failed_results](#backups_batch_delete_failed_results) structure is documented below.

<a name="backups_batch_delete_failed_results"></a>
The `failed_results` block supports:

* `backup_id` - Indicates the backup ID.

* `error_code` - Indicates the error code returned when a task submission exception occurs.

* `error_msg` - Indicates the error message returned when a task submission exception occurs.
