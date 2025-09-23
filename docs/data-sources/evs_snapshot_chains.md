---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_snapshot_chains"
description: |-
  Use this data source to query the list of EVS snapshot chains within HuaweiCloud.
---

# huaweicloud_evs_snapshot_chains

Use this data source to query the list of EVS snapshot chains within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_evs_snapshot_chains" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `id` - (Optional, String) Specifies the snapshot chain ID.

* `volume_id` - (Optional, String) Specifies the disk ID to which the snapshot chains belong.

* `category` - (Optional, String) Specifies the category of snapshot chain.
  The valid values are **standard**, **backup** and **server_backup**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `snapshot_chains` - The snapshot chain list.
  The [snapshot_chains](#snapshot_chains_structure) structure is documented below.

<a name="snapshot_chains_structure"></a>
The `snapshot_chains` block supports:

* `id` - The snapshot chain ID.

* `availability_zone` - The AZ to which the snapshot chain belongs.

* `snapshot_count` - The number of snapshots on the snapshot chain.

* `capacity` - The total size of the snapshot chain.

* `volume_id` - The ID of disk to which the snapshot chain belongs.

* `category` - The category of snapshot chain.

* `created_at` - The time when the snapshot chain was created.

* `updated_at` - The time when the snapshot chain was updated.
