---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_download_backups_authorize"
description: |-
  Manages a GaussDB to download GaussDB backups within HuaweiCloud.
---

# huaweicloud_gaussdb_download_backups_authorize

Manages a GaussDB to download GaussDB backups within HuaweiCloud.

-> This resource is a one-time action resource for authorizing users to download backups. Deleting this resource will
   not revoke the authorization, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "backup_id" {}

resource "huaweicloud_gaussdb_download_backups_authorize" "test" {
  backup_id = var.backup_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `backup_id` - (Required, String, ForceNew) Specifies the backup ID, which uniquely identifies a backup.
  The value must be in UUID format, exactly 36 characters long, and contain only letters and digits.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `bucket` - The name of the bucket where the file is stored.

* `file_paths` - The paths from which backups can be downloaded using OBS Browser+.
