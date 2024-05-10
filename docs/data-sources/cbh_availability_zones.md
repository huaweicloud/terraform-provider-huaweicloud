---
subcategory: "Cloud Bastion Host (CBH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbh_availability_zones"
description: |-
  Use this data source to get the list of CBH availability zones within HuaweiCloud.
---

# huaweicloud_cbh_availability_zones

Use this data source to get the list of CBH availability zones within HuaweiCloud.

## Example Usage

```hcl
variable "availability_zone_name" {}

data "huaweicloud_cbh_availability_zones" "test" {
  name = var.availability_zone_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the CBH availability zones.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the availability zone to be queried.

* `display_name` - (Optional, String) Specifies the display name of the availability zone to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `availability_zones` - All availability zones that match the filter parameters.  
  The [availability_zones](#cbh_availability_zones) structure is documented below.

<a name="cbh_availability_zones"></a>
The `availability_zones` block supports:

* `name` - The name of the availability zone.

* `region_id` - The ID of the region in which the availability zone belongs.

* `display_name` - The display name of the availability zone.

* `type` - The type of the availability zone. The valid values are as follows:
  + **Core**: Core availability zone.
  + **Dedicated**: Exclusive availability zone, only open to internal customers.

* `status` - The status of the availability zone. The value can be **Running**.
