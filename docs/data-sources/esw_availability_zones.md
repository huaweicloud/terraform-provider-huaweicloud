---
subcategory: "Enterprise Switch (ESW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_esw_availability_zones"
description: |-
  Use this data source to get the list of ESW available zones.
---

# huaweicloud_esw_availability_zones

Use this data source to get the list of ESW available zones.

## Example Usage

```hcl
data "huaweicloud_esw_availability_zones" "test" {}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the available zones. If omitted, the provider-level region
  will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `available_zones` - Indicates the list of available zones.
