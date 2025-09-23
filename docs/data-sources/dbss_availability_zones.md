---
subcategory: "Database Security Service (DBSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dbss_availability_zones"
description: |-
  Use this data source to get a list of DBSS availability zones.
---

# huaweicloud_dbss_availability_zones

Use this data source to get a list of DBSS availability zones.

## Example Usage

```hcl
data "huaweicloud_dbss_availability_zones" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `availability_zones` - The list of the availability zones.

  The [availability_zones](#availability_zones_struct) structure is documented below.

<a name="availability_zones_struct"></a>
The `availability_zones` block supports:

* `name` - The name of the availability zone.

* `number` - The number of the availability zone.

* `type` - The type of the availability zone.

* `alias` - The alias of the availability zone.

* `alias_us` - The alias in English of the availability zone.
