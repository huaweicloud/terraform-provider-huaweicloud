---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_snapshot_rollback"
description: |-
  Manages an EVS snapshot rollback resource within HuaweiCloud.
---

# huaweicloud_evs_snapshot_rollback

Manages an EVS snapshot rollback resource within HuaweiCloud.

-> 1. Snapshot rollback is only supported rollback to the source volume, rollback to other specified volume is not
  supported.<br/>2. Snapshot rollback to the source volume is only allowed if the volume status is **available** or
  **error_rollbacking**.<br/>3. Snapshots with names prefixed by **autobk_snapshot_** are automatically created by the
  system when creating volume backups, please do not perform the rollback snapshot to volume operation on these
  snapshots.<br/>4. Destroying resources does not change the current operation of the snapshot rollback.

## Example Usage

```hcl
variable "volume_id" {}
variable "snapshot_id" {}

resource "huaweicloud_evs_snapshot_rollback" "test" {
  volume_id   = var.volume_id
  snapshot_id = var.snapshot_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `volume_id` - (Required, String, ForceNew) Specifies the target volume ID for snapshot rollback.
  Changing this parameter will create a new resource.

* `snapshot_id` - (Required, String, ForceNew) Specifies the ID of the snapshot.
  Changing this parameter will create a new resource.

* `name` - (Optional, String, ForceNew) Specifies the target volume name for snapshot rollback.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `snapshot_id`.
