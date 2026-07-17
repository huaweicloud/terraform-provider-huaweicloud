---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_feature_gates"
description: |-
  Use this data source to query the CCI feature gates within HuaweiCloud.
---

# huaweicloud_cciv2_feature_gates

Use this data source to query the CCI feature gates within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cciv2_feature_gates" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the CCI feature gates.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `feature_gates` - The list of feature gates.
  The [feature_gates](#cciv2_feature_gates_feature_gates) structure is documented below.

<a name="cciv2_feature_gates_feature_gates"></a>
The `feature_gates` block supports:

* `deprecated` - Whether the feature will be deprecated.

* `description` - The description of the feature.

* `feature` - The feature name.

* `type` - The type of the feature value. Valid values are **boolean**, **int**, **list**, and **string**.

* `value` - The feature value.
