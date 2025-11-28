---
subcategory: "Enterprise Switch (ESW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_esw_flavors"
description: |-
  Use this data source to get the list of ESW flavors.
---

# huaweicloud_esw_flavors

Use this data source to get the list of ESW flavors.

## Example Usage

```hcl
data "huaweicloud_esw_flavors" "test" {}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the flavors. If omitted, the provider-level region will be
  used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - Indicates the list of flavors.
  The [flavors](#flavors_struct) structure is documented below.

<a name="flavors_struct"></a>
The `flavors` block supports:

* `id` - Indicates the ID of the flavor.

* `name` - Indicates the name of the flavor.

* `connections` - Indicates the number of layer 2 connections that can be supported.

* `bandwidth` - Indicates the maximum bandwidth that the instance can handle.

* `pps` - Indicates the maximum number of packets an instance can handle.

* `available_zones` - Indicates the list of available zones.
