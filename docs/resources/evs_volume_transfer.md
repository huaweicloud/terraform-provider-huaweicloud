---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_volume_transfer"
description: |-
  Manages an EVS volume transfer resource within HuaweiCloud.
---

# huaweicloud_evs_volume_transfer

Manages an EVS volume transfer resource within HuaweiCloud.

-> The transfer can only be created when the EVS volume is in the **available** state, other constraints that do not
   support transfer are as follows:
   <br/>1. Volumes with the prePaid billing mode do not support transfer.
   <br/>2. Frozen volumes do not support transfer.
   <br/>3. Encrypted volumes do not support transfer.
   <br/>4. Volumes with corresponding backups and snapshots do not support transfer.
   <br/>5. Volumes with backup policies do not support transfer.
   <br/>6. Volumes on DSS (Dedicated Storage Service) do not support transfer.
   <br/>7. Volumes on DESS (Dedicated Enterprise Storage Service) do not support transfer.

## Example Usage

```hcl
variable "volume_id" {}
variable "name" {}

resource "huaweicloud_evs_volume_transfer" "test" {
  volume_id = var.volume_id
  name      = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `volume_id` - (Required, String, ForceNew) Specifies the volume ID to be transferred.
  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the volume transfer record.
  Supports a maximum of `64` characters. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `auth_key` - The identity authentication key for volume transfer.

  -> The `auth_key` field is used to accept the volume transfer, after the volume transfer is successfully accepted,
     the volume transfer resource will no longer exist.

* `created_at` - The creation time of the volume transfer record, in RFC3339 format.
