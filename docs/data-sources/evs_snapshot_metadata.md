---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_snapshot_metadata"
description: |-
  Use this data source to get the EVS snapshot metadata within HuaweiCloud.
---

# huaweicloud_evs_snapshot_metadata

Use this data source to get the EVS snapshot metadata within HuaweiCloud.

## Example Usage

```hcl
variable "snapshot_id" {}

data "huaweicloud_evs_snapshot_metadata" "test" {
  snapshot_id = var.snapshot_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `snapshot_id` - (Required, String) Specifies the snapshot ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `metadata` - The user-defined metadata key-value pair.
