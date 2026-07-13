---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_shard_storage"
description: |-
  Use this data source to query the shard storage information of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_shard_storage

Use this data source to query the shard storage information of a GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_shard_storage" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the shard storage information.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the GaussDB instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `group_disk_infos` - The storage information of each shard in the list.
  The [group_disk_infos](#shard_storage_group_disk_infos) structure is documented below.

<a name="shard_storage_group_disk_infos"></a>
The `group_disk_infos` block supports:

* `name` - The shard name.

* `used` - The shard storage usage.

* `total` - The shard storage size.

* `group_id` - The shard group ID.
