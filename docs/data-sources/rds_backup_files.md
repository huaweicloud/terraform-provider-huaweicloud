---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_backup_files"
description: |-
  Use this data source to get the link for downloading a backup file.
---

# huaweicloud_rds_backup_files

Use this data source to get the link for downloading a backup file.

## Example Usage

```hcl
var "backup_id" {}

data "huaweicloud_rds_backup_files" "test" {
  backup_id = var.backup_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

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

* `size` - Indicates the file size in KB.

* `download_link` - Indicates the link for downloading the backup file.

* `link_expired_time` - Indicates the link expiration time.

* `database_name` - Indicates the name of the database.
