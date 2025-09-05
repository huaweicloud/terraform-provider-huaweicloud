---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_supported_regions"
description: |-
  Use this data source to get the supported regions.
---

# huaweicloud_cc_supported_regions

Use this data source to get the supported regions.

## Example Usage

```hcl
data "huaweicloud_cc_supported_regions" "test"{}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `regions` - Indicates the region list.

  The [regions](#regions_struct) structure is documented below.

<a name="regions_struct"></a>
The `regions` block supports:

* `id` - Indicates the region ID。

* `name` - Indicates the region name.

* `area_id` - Indicates the geographic region.
  Cloud Connect is available in the following geographic regions:
  + **Chinese-Mainland**: Chinese mainland
  + **Asia-Pacific**: Asia Pacific
  + **Africa**
  + **Western-Latin-America**: Western Latin America
  + **Eastern-Latin-America**: Eastern Latin America
  + **Northern-Latin-America**: Northern Latin America

* `area_name` - Indicates the geographic region name.

* `used_scenes` - Indicates the cloud Connect application scenarios.
  The value can be:
  + **er**: enterprise router
  + **vpc**: VPC
  + **vgw**: virtual gateway
