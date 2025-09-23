---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_backup_metadata"
description: |-
  Use this data source to get backup metadata from CBR service within HuaweiCloud.
---

# huaweicloud_cbr_backup_metadata

Use this data source to get backup metadata from CBR service within HuaweiCloud.

## Example Usage

```hcl
variable "backup_id" {}

data "huaweicloud_cbr_backup_metadata" "test" {
  backup_id = var.backup_id
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the datasource.
  If omitted, the provider-level region will be used.

* `backup_id` - (Required, String) Specifies backup id.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, which is the backup ID.

* `backup_id` - The backup id.

* `backups` - The server backup informations.

* `flavor` - The server specifications.

* `floatingips` - The server floating IP address information.

* `interface` - The server API information.

* `ports` - The server port information.

* `server` - The server information .

* `volumes` - The server disk information.
