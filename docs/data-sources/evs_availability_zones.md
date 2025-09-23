---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_availability_zones"
description: |-
  Use this data source to query the list of EVS availability zones within HuaweiCloud.
---

# huaweicloud_evs_availability_zones

Use this data source to query the list of EVS availability zones within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_evs_availability_zones" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `availability_zones` - The list of availability zones.

  The [availability_zones](#availability_zones_struct) structure is documented below.

<a name="availability_zones_struct"></a>
The `availability_zones` block supports:

* `is_available` - Whether the availability zone is available.

* `name` - The name of availability zone.
