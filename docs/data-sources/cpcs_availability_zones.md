---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpcs_availability_zones"
description: |-
  Use this data source to get the list of availability zones for creating password service cludters.
---

# huaweicloud_cpcs_availability_zones

Use this data source to get the list of availability zones for creating password service cludters.

-> Currently, this data source is valid only in cn-north-9 region.

## Example Usage

```hcl
data "huaweicloud_cpcs_availability_zones" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `availability_zone` - The list of availability zones.
  The [availability_zone](#availability_zones_struct) structure is documented below.

<a name="availability_zones_struct"></a>
The `availability_zone` block supports:

* `id` - The availability zone ID.

* `display_name` - The display name.

* `locales` - The availability zone name.
  The [locales](#locales_struct) structure is documented below.

* `type` - The availability zone type.

* `region_id` - The region ID.

* `status` - The availability zone status.

<a name="locales_struct"></a>
The `locales` block supports:

* `en_us` - The English name.

* `zh_cn` - The Chinese name.
