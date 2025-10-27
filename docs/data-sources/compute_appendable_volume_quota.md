---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_appendable_volume_quota"
description: |-
  Use this data source to get the number of disks that can be added to a yearly/monthly ECS.
---

# huaweicloud_compute_appendable_volume_quota

Use this data source to get the number of disks that can be added to a yearly/monthly ECS.

## Example Usage

```hcl
variable "server_id" {}

data "huaweicloud_compute_appendable_volume_quota" "test" {
  server_id = var.server_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `server_id` - (Required, String) Specifies the VM ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quota_count` - Indicates the quantity.

* `free_scsi` - Indicates the number of SCSI disks that can be attached.

* `free_blk` - Indicates the number of block disks that can be attached.
