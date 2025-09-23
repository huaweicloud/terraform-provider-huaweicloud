---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_geo_blockings"
description: |-
  Use this data source to get the list of Live geo blockings within HuaweiCloud.
---

# huaweicloud_live_geo_blockings

Use this data source to get the list of Live geo blockings within HuaweiCloud.

## Example Usage

```hcl
variable "domain_name" {}

data "huaweicloud_live_geo_blockings" "test" {
  domain_name = var.domain_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `domain_name` - (Required, String) Specifies the streaming domain name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `apps` - The list of the application.

  The [apps](#apps_struct) structure is documented below.

<a name="apps_struct"></a>
The `apps` block supports:

* `app_name` - The application name.

* `area_whitelist` - The restricted area list, an empty list indicates no restrictions.
  Except for China, codes for other regions are capitalized with `2` letters.
  Some valid values are as follows:
  + **CN-IN**: Chinese Mainland.
  + **CN-HK**: Hong Kong, China.
  + **CN-MO**: Macao, China.
  + **CN-TW**: Taiwan, China.
  + **BR**: Brazil.
