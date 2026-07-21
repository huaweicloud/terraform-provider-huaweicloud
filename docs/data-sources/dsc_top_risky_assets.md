---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_top_risky_assets"
description: |-
  Use this data source to get the top five risky assets within HuaweiCloud.
---

# huaweicloud_dsc_top_risky_assets

Use this data source to get the top five risky assets within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_top_risky_assets" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `risk_asset_list` - The list of the top five assets with the highest risk score in the current project.

  The [risk_asset_list](#risk_asset_list_struct) structure is documented below.

<a name="risk_asset_list_struct"></a>
The `risk_asset_list` block supports:

* `asset_id` - The unique ID of the asset.

* `asset_name` - The name of the asset.

* `asset_type` - The major category type of the asset.

* `data_source` - The data source type of the asset.

* `deducted_point` - The deducted point value of the security risk.
