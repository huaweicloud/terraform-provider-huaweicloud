---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_asset_quota"
description: |-
  Use this data source to get the asset quota information of the specified project within HuaweiCloud.
---

# huaweicloud_dsc_asset_quota

Use this data source to get the asset quota information of the specified project within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_asset_quota" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_size` - The total asset quota.

* `use_size` - The number of used assets.
