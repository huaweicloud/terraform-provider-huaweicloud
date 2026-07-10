---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_protect_measure_baseline"
description: |-
  Use this data source to get the protect measure baseline within HuaweiCloud.
---

# huaweicloud_dsc_protect_measure_baseline

Use this data source to get the protect measure baseline within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_protect_measure_baseline" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `protect_measure_baseline` - The protect measure baseline list.

  The [protect_measure_baseline](#protect_measure_baseline_struct) structure is documented below.

<a name="protect_measure_baseline_struct"></a>
The `protect_measure_baseline` block supports:

* `key` - The key of the protect measure baseline map.

* `detail` - The measure requirement detail list.

  The [detail](#measure_requirement_detail_struct) structure is documented below.

<a name="measure_requirement_detail_struct"></a>
The `detail` block supports:

* `create_time` - The creation time.

* `data_type_info` - The data type detail.

  The [data_type_info](#data_type_detail_struct) structure is documented below.

* `id` - The measure ID.

* `measure_info` - The measure detail.

  The [measure_info](#measure_detail_struct) structure is documented below.

* `protect_level` - The protect level.

* `update_time` - The update time.

<a name="data_type_detail_struct"></a>
The `data_type_info` block supports:

* `category` - The data protect measure category.

* `create_time` - The creation time.

* `data_type` - The data protect type.

* `data_type_id` - The data type ID.

* `id` - The data protect measure ID.

* `is_deleted` - Whether the data protect measure is deleted.

* `life_cycle` - The life cycle type.

* `update_time` - The update time.

<a name="measure_detail_struct"></a>
The `measure_info` block supports:

* `create_time` - The creation time.

* `description` - The description.

* `id` - The measure ID.

* `is_deleted` - Whether the measure is deleted.

* `life_cycle` - The life cycle type.

* `measure_type` - The measure type.

* `name` - The measure name.

* `update_time` - The update time.
