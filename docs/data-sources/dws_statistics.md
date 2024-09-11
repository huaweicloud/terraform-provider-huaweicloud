---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_statistics"
description: |-
  Use this data source to query the resource statistics of the DWS within HuaweiCloud.
---

# huaweicloud_dws_statistics

Use this data source to query the resource statistics of the DWS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dws_statistics" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `statistics` - The list of the resource statistics.

  The [statistics](#statistics_struct) structure is documented below.

<a name="statistics_struct"></a>
The `statistics` block supports:

* `name` - The resource name.
  + **cluster.total**: The total number of DWS clusters.
  + **cluster.normal**: The number of available DWS clusters.
  + **instance.total**: The total number of nodes.
  + **instance.normal**: The number of available nodes.
  + **storage.total**: Total Capacity.

* `value` - The value of the resource.

* `unit` - The unit of the resource.
