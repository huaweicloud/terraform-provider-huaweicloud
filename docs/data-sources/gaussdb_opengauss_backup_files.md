---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_backup_files"
description: |-
  Use this data source to get the list of GaussDB OpenGauss backup files.
---

# huaweicloud_gaussdb_opengauss_backup_files

Use this data source to get the list of GaussDB OpenGauss backup files.

## Example Usage

```hcl
variable "backup_id" {}

data "huaweicloud_gaussdb_opengauss_backup_files" "test" {
  backup_id = var.backup_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `backup_id` - (Required, String) Specifies the ID of the backup.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `files` - Indicates the list of backup files.

  The [files](#files_struct) structure is documented below.

* `bucket` - Indicates the bucket name.

<a name="files_struct"></a>
The `files` block supports:

* `name` - Indicates the file name.

* `size` - Indicates the file size.

* `download_link` - Indicates the download link.

* `link_expired_time` - Indicates the link expired time in the **yyyy-mm-ddThh:mm:ssZ** format.
