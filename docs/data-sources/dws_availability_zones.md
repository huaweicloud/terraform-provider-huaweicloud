---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_availability_zones"
description: |-
  Use this data source to get the list of DWS availability zones within HuaweiCloud.
---

# huaweicloud_dws_availability_zones

Use this data source to get the list of DWS availability zones within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dws_availability_zones" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `availability_zones` - All availability zones that match the filter parameters.

  The [availability_zones](#availability_zones_struct) structure is documented below.

<a name="availability_zones_struct"></a>
The `availability_zones` block supports:

* `name` - The name of the availability zone.

* `display_name` - The display name of the availability zone.

* `status` - The current status of the availability zone.
  + **available**
  + **unavailable**

* `public_border_group` - The availability zone group.
