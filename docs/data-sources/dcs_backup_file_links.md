---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_backup_file_links"
description: |-
  Use this data source to get the list of DCS backup file links.
---

# huaweicloud_dcs_backup_file_links

Use this data source to get the list of DCS backup file links.

## Example Usage

```hcl
var "instance_id" {}
var "backup_id" {}

data "huaweicloud_dcs_backups" "test" {
  instance_id = var.instance_id
  backup_id   = var.backup_id
  expiration  = 3600
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the DCS instance.

* `backup_id` - (Required, String) Specifies the ID of the DCS instance backup.

* `expiration` - (Required, Int) Specifies the expiration time of the download link.
  The value ranges from **300** to **86400** seconds.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `bucket_name` - Indicates the name of the OBS bucket where the backup file is stored.

* `file_path` - Indicates the path of the backup file in the OBS bucket.

* `links` - Indicates the list of backup file download links.
  The [links](#links_struct) structure is documented below.

<a id="links_struct"></a>
The `links` block supports:

* `file_name` - Indicates the name of the backup file.

* `link` - Indicates the download link of the backup file.
