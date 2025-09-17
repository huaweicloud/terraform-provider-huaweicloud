---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_features"
description: |-
  Use this data source to query the features of CBR within HuaweiCloud.
---

# huaweicloud_cbr_features

Use this data source to query the features of CBR within HuaweiCloud.

-> The API used by this datasource is currently in public beta and is temporarily unavailable in some regions.

## Example Usage

```hcl
data "huaweicloud_cbr_features" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `feature_value` - The feature values in JSON format. This contains all the feature information provided by the CBR service.
