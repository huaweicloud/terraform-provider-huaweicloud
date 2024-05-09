---
subcategory: "Global Accelerator (GA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ga_availability_zones"
description: |-
  Use this data source to get the list of availability zones.
---

# huaweicloud_ga_availability_zones

Use this data source to get the list of availability zones.

## Example Usage

```hcl
variable "area" {}

data "huaweicloud_ga_availability_zones" "test" {
  area = var.area
}
```

## Argument Reference

The following arguments are supported:

* `area` - (Optional, String) The acceleration area to which the regions belong.
  The value can be one of the following:
  + **OUTOFCM**: Outside the Chinese mainland.
  + **CM**: Chinese mainland.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `regions` - The region list.

  The [regions](#regions_struct) structure is documented below.

<a name="regions_struct"></a>
The `regions` block supports:

* `region_id` - The region ID.

* `area` - The acceleration area to which the region belongs.

* `endpoint_types` - The endpoint types supported by the region.
