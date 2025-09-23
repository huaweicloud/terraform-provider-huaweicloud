---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_availability_zones"
description: |-
  Use this data source to get the availability zone list of APIG instance within HuaweiCloud.
---

# huaweicloud_apig_availability_zones

Use this data source to get the availability zone list of APIG instance within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_apig_availability_zones" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `availability_zones` - The list of the availability zones.

  The [availability_zones](#apig_data_availability_zones_struct) structure is documented below.

<a name="apig_data_availability_zones_struct"></a>
The `availability_zones` block supports:

* `id` - The ID of the availability zone.

* `name` - The name of the availability zone.

* `code` - The code of the availability zone.

* `port` - The port of the availability zone.

* `specs` - The APIG instance editions supported by the availability zone.

* `local_name` - The Chinese and English names of the availability zone.

  The [local_name](#apig_data_availability_zones_local_name_struct) structure is documented below.

<a name="apig_data_availability_zones_local_name_struct"></a>
The `local_name` block supports:

* `en_us` - The English name of the availability zone.

* `zh_cn` - The Chinese name of the availability zone.
