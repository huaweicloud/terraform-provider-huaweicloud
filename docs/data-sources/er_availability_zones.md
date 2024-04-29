---
subcategory: "Enterprise Router (ER)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_er_availability_zones"
description: ""
---

# huaweicloud_er_availability_zones

Use this data source to query availability zones where ER instances can be created within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_er_availability_zones" "all" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are supported:

* `id` - The data source ID.

* `names` - The names of availability zone.
