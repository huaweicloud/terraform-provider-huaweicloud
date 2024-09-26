---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_snapshot_statistics"
description: |-
  Use this data source to query the snapshot statistics under the specified DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_cluster_snapshot_statistics

Use this data source to query the snapshot statistics under the specified DWS cluster within HuaweiCloud.

## Example Usage

```hcl
variable "dws_cluster_id" {}

data "huaweicloud_dws_cluster_snapshot_statistics" "test" {
  cluster_id = var.dws_cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specified the ID of the DWS cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `statistics` - The list of the snapshot statistics.

  The [statistics](#snapshot_statistics_struct) structure is documented below.

<a name="snapshot_statistics_struct"></a>
The `statistics` block supports:

* `name` - The name of the resource statistic.
  + **storage.free**: The free capacity available for the snapshots.
  + **storage.paid**: The paid capacity by the snapshots.
  + **storage.used**: The capacity used by the snapshots.

* `value` - The value of the resource statistic.

* `unit` - The unit of the resource statistic.
