---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_availability_zones"
description: |-
  Use this data source to get the list of availability zones.
---

# huaweicloud_compute_availability_zones

Use this data source to get the list of availability zones.

## Example Usage

```hcl
data "huaweicloud_compute_availability_zones" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `availability_zones` - Indicates the list of availability zones.

  The [availability_zones](#availability_zones_struct) structure is documented below.

<a name="availability_zones_struct"></a>
The `availability_zones` block supports:

* `availability_zone_id` - Indicates the ID of the availability zone.

* `type` - Indicates the type of the availability zone.

* `mode` - Indicates the mode of the availability zone.

* `category` - Indicates the category of the availability zone.

* `alias` - Indicates the mode of the availability zone.

* `public_border_group` - Indicates the public border group of the availability zone.

* `az_group_ids` - Indicates the availability zone group ids.
