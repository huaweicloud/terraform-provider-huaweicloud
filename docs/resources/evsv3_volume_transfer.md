---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evsv3_volume_transfer"
description: |-
  Manages an EVS volume transfer (V3) resource within HuaweiCloud.
---

# huaweicloud_evsv3_volume_transfer

Manages an EVS volume transfer (V3) resource within HuaweiCloud.

-> The transfer can only be created when the EVS volume is in the **available** state, other constraints that do not
   support transfer are as follows:
   <br/>1. Volumes with the prePaid billing mode do not support transfer.
   <br/>2. Frozen volumes do not support transfer.
   <br/>3. Encrypted volumes do not support transfer.
   <br/>4. Volumes with corresponding backups and snapshots do not support transfer.
   <br/>5. Volumes with backup policies do not support transfer.
   <br/>6. Volumes on DSS (Dedicated Storage Service) do not support transfer.
   <br/>7. Volumes on DESS (Dedicated Enterprise Storage Service) do not support transfer.
   <br/>8. EVS system disk does not support transfer.

## Example Usage

```hcl
variable "volume_id" {}
variable "name" {}

resource "huaweicloud_evsv3_volume_transfer" "test" {
  volume_id = var.volume_id
  name      = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `volume_id` - (Required, String, NonUpdatable) Specifies the volume ID to be transferred.

* `name` - (Required, String, NonUpdatable) Specifies the name of the volume transfer record.
  Supports a maximum of `64` characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `auth_key` - The identity authentication key for volume transfer.

  -> The `auth_key` field is used to accept the volume transfer, after the volume transfer is successfully accepted,
     the volume transfer resource will no longer exist.

* `created_at` - The creation time of the volume transfer record, in RFC3339 format.

* `links` - The links to the cloud disk transfer record.
  The [links](#links_struct) structure is documented below.

<a name="links_struct"></a>
The `links` block supports:

* `href` - The corresponding shortcut link.

* `rel` - The shortcut link marker name.

## Import

Volumes can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_evsv3_volume_transfer.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `auth_key`.
It is generally recommended running terraform plan after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_evsv3_volume_transfer" "test" {
    ...

  lifecycle {
    ignore_changes = [
      auth_key, 
    ]
  }
}
```
