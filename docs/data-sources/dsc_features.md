---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_features"
description: |-
  Use this data source to get the list of features within HuaweiCloud.
---

# huaweicloud_dsc_features

Use this data source to get the list of features within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_features" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The list of features.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `description` - The description of the feature.

* `enabled` - Whether the feature is enabled.

* `name` - The name of the feature.
