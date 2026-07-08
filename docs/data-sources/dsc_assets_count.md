---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_assets_count"
description: |-
  Use this data source to get the total number of assets under the specified project within HuaweiCloud.
---

# huaweicloud_dsc_assets_count

Use this data source to get the total number of assets under the specified project within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_assets_count" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `assets_count` - The total number of assets under the project.
