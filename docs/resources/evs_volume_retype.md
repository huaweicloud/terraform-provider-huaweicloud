---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_volume_retype"
description: |-
  Manages an EVS volume retype resource within HuaweiCloud.
---

# huaweicloud_evs_volume_retype

Manages an EVS volume retype resource within HuaweiCloud.

-> 1. The current resource is a one-time resource, and destroying this resource will not recover the volume type,
but will only remove the resource information from the tfstate file.<br>2. The `new_type` parameters of this resource
will affect other resource with `volume_type` parameter, such as `huaweicloud_evs_volume`. You can handle the changes
in the affected resource by `lifecycle.ignore_changes`.

## Example Usage

```hcl
variable "volume_id" {}
variable "snapshot_id" {}

resource "huaweicloud_evs_volume_retype" "test" {
  volume_id   = var.volume_id
  snapshot_id = var.snapshot_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `volume_id` - (Required, String, NonUpdatable) Specifies the target volume ID for snapshot rollback.
  Changing this parameter will create a new resource.

* `new_type` - (Required, String, NonUpdatable) Specifies the cloud disk type which change to. Possible values are:
  + **SAS**: High I/O type.
  + **SSD**: Ultra-high I/O type.
  + **GPSSD**: General purpose SSD type.
  + **ESSD**: Extreme SSD type.
  + **GPSSD2**: General purpose SSD V2 type.
  + **ESSD2**: Extreme SSD V2 type.

  -> The field has the following restrictions:
  <br/>1. When the specified cloud disk type does not exist in the availability_zone, the cloud disk type change
  fails.
  <br/>2. When the original type is SAS, it can be changed to any of the other types mentioned above.
  <br/>3. When the original type includes SSD, it can be retyped to other types including SSD, but cannot be
  retyped to SAS.

* `is_auto_pay` - (Optional, String, NonUpdatable) Specifies whether to pay immediately. This parameter is valid only
  when chargingMode is set to prePaid. Possible values are:
  + **true**: An order is immediately paid from the account balance.
  + **false**: An order is not paid immediately after being created.

* `iops` - (Optional, String, NonUpdatable) Specifies the new maximum IOPS of the disk. This parameter is supported
  only for general purpose SSD V2 and extreme SSD V2 disks.

* `throughput` - (Optional, String, NonUpdatable) Specifies the new maximum throughput of the disk, in the unit of
  MiB/s. This parameter is supported only for general purpose SSD V2 disks.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
