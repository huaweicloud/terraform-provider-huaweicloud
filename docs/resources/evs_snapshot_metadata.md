---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_snapshot_metadata"
description: |-
  Manages an EVS snapshot metadata resource within HuaweiCloud.
---

# huaweicloud_evs_snapshot_metadata

Manages an EVS snapshot metadata resource within HuaweiCloud.

## Example Usage

```hcl
variable "snapshot_id" {}

resource "huaweicloud_evs_snapshot_metadata" "test" {
  snapshot_id = var.snapshot_id

  metadata = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `snapshot_id` - (Required, String, NonUpdatable) Specifies the ID of the snapshot.

* `metadata` - (Required, Map) Specifies the user-defined metadata key-value pair.

  -> When updating the `metadata` parameter, all existing key-value pairs will be overwritten.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `snapshot_id`.
