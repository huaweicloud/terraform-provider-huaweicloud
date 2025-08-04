---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbo_share_types"
description: |-
  Use this data source to get the list of SFS turbo file system types and quota.
---

# huaweicloud_sfs_turbo_share_types

Use this data source to get the list of SFS turbo file system types and quota.

## Example Usage

```hcl
data "huaweicloud_sfs_turbo_share_types" "test" {
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `share_types` - The list of file system types and quotas.
  The [share_types](#share_types_struct) structure is documented below.

<a name="share_types_struct"></a>
The `share_types` block supports:

* `share_type` - The type of the SFS turbo file system.

* `scenario` - The scene of the SFS turbo file system.

* `attribution` - The attribution of the SFS turbo file system.
  The [attribution](#attribution_struct) structure is documented below.

* `support_period` - Whether the yearly/monthly billing mode is supported.

* `available_zones` - The AZs where there are the SFS turbo file system.
  The [available_zones](#available_zones_struct) structure is documented below.

* `spec_code` - The specification code of the SFS turbo file system.

* `storage_media` - The storage media of the SFS turbo file system. Possible values are **HDD**, **SDD** and **ESSD**.

* `features` - The features list supported by the instance.

<a name="attribution_struct"></a>
The `attribution` block supports:

* `capacity` - The size attributions of the SFS turbo file system.
  The [capacity](#capacity_struct) structure is documented below.

* `bandwidth` - The bandwidth attributions of the SFS turbo file system.
  The [bandwidth](#bandwidth_struct) structure is documented below.

* `iops` - The iops attributions of the SFS turbo file system.
  The [iops](#iops_struct) structure is documented below.

* `single_channel_4k_latency` - The single-channel 4K latency attributions of the SFS turbo file system.
  The [single_channel_4k_latency](#single_channel_4k_latency_struct) structure is documented below.

<a name="capacity_struct"></a>
The `capacity` block supports:

* `max` - The max capacity of the SFS turbo file system.

* `min` - The min capacity of the SFS turbo file system.

* `step` - The capacity step of the SFS turbo file system.

<a name="bandwidth_struct"></a>
The `bandwidth` block supports:

* `max` - The max bandwidth of the SFS turbo file system.

* `min` - The min bandwidth of the SFS turbo file system.

* `step` - The bandwidth step of the SFS turbo file system.

* `density` - The bandwidth density of the SFS turbo file system.

* `base` - The basic bandwidth of the SFS turbo file system.

<a name="iops_struct"></a>
The `iops` block supports:

* `max` - The max iops of the SFS turbo file system.

* `min` - The min iops of the SFS turbo file system.

<a name="single_channel_4k_latency_struct"></a>
The `single_channel_4k_latency` block supports:

* `max` - The max single-channel 4K latency of the SFS turbo file system.

* `min` - The min single-channel 4K latency of the SFS turbo file system.

<a name="available_zones_struct"></a>
The `available_zones` block supports:

* `available_zone` - The availability zone name.

* `status` - The SFS turbo file system status in this availability zone. Possible values are **normal** and **sellout**.
