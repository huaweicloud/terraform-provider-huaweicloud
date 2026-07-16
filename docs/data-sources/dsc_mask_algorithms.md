---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_mask_algorithms"
description: |-
  Use this data source to get the mask algorithm list within HuaweiCloud.
---

# huaweicloud_dsc_mask_algorithms

Use this data source to get the mask algorithm list within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_mask_algorithms" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the algorithm name used to filter results.

* `type` - (Optional, String) Specifies the algorithm type used to filter results.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `algorithms` - The mask algorithm list.

  The [algorithms](#algorithms_struct) structure is documented below.

<a name="algorithms_struct"></a>
The `algorithms` block supports:

* `algorithm` - The mask algorithm identifier.

* `algorithm_id` - The mask algorithm ID.

* `algorithm_name` - The mask algorithm name.

* `algorithm_name_en` - The mask algorithm English name.

* `category` - The mask algorithm category.

* `data` - The data content processed by the mask algorithm.

* `parameter` - The configuration parameters of the mask algorithm.

* `processed_data` - The processed data content.
