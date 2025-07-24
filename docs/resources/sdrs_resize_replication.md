---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_resize_replication"
description: |-
  Using this resource to resize a replication pair's disk in SDRS within HuaweiCloud.
---

# huaweicloud_sdrs_resize_replication

Using this resource to resize a replication pair's disk in SDRS within HuaweiCloud.

-> This is a one-time action resource to resize a replication pair's disk. Deleting this resource will
not change the current configuration, but will only remove the resource information from the tfstate file.

-> Before using this resource, please note the following restrictions:
<br/>1. The status of the replication pair must be **available** or **protected** or **error-extending**.
<br/>2. The status of the cloud disks must be **available** or **in-use**.
<br/>3. If the cloud disks are pay-per-use, they can be resized directly.
<br/>4. If the cloud disks are prepaid, they cannot be resized directly. You need to delete the replication pair first,
resize the cloud disks, and then create a new replication pair.
<br/>5. Running this resource may cause unexpected change to the `size` field of the `huaweicloud_evs_volume`.
Please using `lifecycle` to control the `size` field of the `huaweicloud_evs_volume`.

## Example Usage

```hcl
variable "replication_id" {}
variable "new_size" {}

resource "huaweicloud_sdrs_resize_replication" "test" {
  replication_id = var.replication_id
  new_size       = var.new_size
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `replication_id` - (Required, String, NonUpdatable) Specifies the ID of the replication pair to resize.

* `new_size` - (Required, Integer, NonUpdatable) Specifies the new size of the replication pair's disk in GB.
  Must be greater than the current size.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the resource.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
