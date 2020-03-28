---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_snapshot"
sidebar_current: "docs-huaweicloud-resource-evs-snapshot"
description: |-
  Provides an EVS snapshot resource.
---

# huaweicloud_evs_snapshot

Provides an EVS snapshot resource.
 
# Example Usage

 ```hcl
resource "huaweicloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  description = "test volume"
  size = 40
}

resource "huaweicloud_evs_snapshot" "snapshot_1" {
  volume_id = huaweicloud_blockstorage_volume_v2.volume_1.id
  name = "snapshot-001"
  description = "Daily backup"
}
 ```

# Argument Reference

The following arguments are supported:

* `volume_id` - (Required) The id of the snapshot's source disk. Changing the parameter creates a new snapshot.

* `name` - (Required) The name of the snapshot. The value can contain a maximum of 255 bytes.

* `description` - (Optional) The description of the snapshot. The value can contain a maximum of 255 bytes.

* `force` - (Optional) Specifies the flag for forcibly creating a snapshot. Default to false.

# Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the snapshot.

* `status` - The status of the snapshot.

* `size` - The size of the snapshot in GB.

 
# Import

EVS snapshot can be imported using the `snapshot id`, e.g.

```
 $ terraform import huaweicloud_evs_snapshot.snapshot_1 3a11b255-3bb6-46f3-91e4-3338baa92dd6
```