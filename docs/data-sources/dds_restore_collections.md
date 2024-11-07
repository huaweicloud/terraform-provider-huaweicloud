---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_restore_collections"
description: |-
  Use this data source to get the list of DDS instance restore collections.
---

# huaweicloud_dds_restore_collections

Use this data source to get the list of DDS instance restore collections.

## Example Usage

```hcl
variable "instance_id" {}
variable "db_name" {}
variable "restore_time" {}

data "huaweicloud_dds_restore_collections" "test"{
  instance_id  = var.instance_id
  db_name      = var.db_name
  restore_time = var.restore_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `db_name` - (Required, String) Specifies the database name.

* `restore_time` - (Required, String) Specifies the restoration time point.
  The value is a UNIX timestamp, in milliseconds. The time zone is UTC.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `collections` - Indicates the list of collections.
