---
subcategory: "RGC"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_home_region"
description: |
  Use this data source to get the home region in Resource Governance Center.
---

# huaweicloud_rgc_home_region

Use this data source to get the home region in Resource Governance Center.

## Example Usage

```hcl
data "huaweicloud_rgc_home_region" "test" {}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `home_region` - The home region ID.
