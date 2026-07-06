---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_asset_overview"
description: |-
  Use this data source to get the asset overview information within HuaweiCloud.
---

# huaweicloud_dsc_asset_overview

Use this data source to get the asset overview information within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_asset_overview" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `asset_sensitive_num` - The total number of risk assets.

* `asset_total_num` - The total number of all assets.

* `high_level_num` - The number of high risk assets.

* `middle_level_num` - The number of medium risk assets.

* `low_level_num` - The number of low risk assets.

* `un_classed_num` - The number of unclassified assets.

* `asset_classification_list` - The list of asset classification information.

  The [asset_classification_list](#asset_classification_list_struct) structure is documented below.

<a name="asset_classification_list_struct"></a>
The `asset_classification_list` block supports:

* `color_num` - The color number corresponding to the level.

* `level_id` - The unique ID of the sensitivity level.

* `level_name` - The name of the sensitivity level.

* `sensitive_num` - The number of sensitive assets in the current level.
