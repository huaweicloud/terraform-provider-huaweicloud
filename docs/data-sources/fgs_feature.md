---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_feature"
description: |-
  Use this data source to query the supported feature information of FunctionGraph within HuaweiCloud.
---

# huaweicloud_fgs_feature

Use this data source to query the supported feature information of FunctionGraph within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_fgs_feature" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the feature information is located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `feature` - The feature information, in JSON format.  
