---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_backup_download_links"
description: |-
  Use this data source to get the list of DDS instance backup download links.
---

# huaweicloud_dds_backup_download_links

Use this data source to get the list of DDS instance backup download links.

## Example Usage

```hcl
variable "instance_id" {}
variable "backup_id" {}

data "huaweicloud_dds_backup_download_links" "test"{
  instance_id = var.instance_id
  backup_id   = var.backup_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `backup_id` - (Required, String) Specifies the backup ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `files` - Indicates the list of backup files.

  The [files](#files_struct) structure is documented below.

* `bucket` - Indicates the name of the bucket where the file is located.

<a name="files_struct"></a>
The `files` block supports:

* `name` - Indicates the file name.

* `download_link` - Indicates the link for downloading the backup file.

* `size` - Indicates the file size in KB.

* `link_expired_time` - Indicates the link expiration time. The format is **yyyy-mm-ddThh:mm:ssZ**.
