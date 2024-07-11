---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_snapshots"
description: |-
  Use this data source to get the list of DWS snapshots within HuaweiCloud.
---

# huaweicloud_dws_snapshots

Use this data source to get the list of DWS snapshots within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dws_snapshots" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `snapshots` - All snapshots that match the filter parameters.

  The [snapshots](#snapshots_struct) structure is documented below.

<a name="snapshots_struct"></a>
The `snapshots` block supports:

* `id` - The ID of the snapshot.

* `name` - The name of the snapshot.

* `cluster_id` - The cluster ID corresponding of the snapshot.

* `type` - The type of the snapshot.
  + **MANUAL**
  + **AUTOMATED**

* `size` - The size of the snapshot, in GB.

* `status` - The current status of the snapshot.
  + **AVAILABLE**
  + **UNAVAILABLE**

* `description` - The description of the snapshot.

* `created_at` - The creation time of the snapshot, in RFC3339 format.

* `finished_at` - The completion time of the snapshot, in RFC3339 format.
