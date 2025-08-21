---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_geolocation_detail"
description: |-
  Use this data source to query geolocation options detail.
---

# huaweicloud_waf_geolocation_detail

Use this data source to query geolocation options detail.

## Example Usage

```hcl
data "huaweicloud_waf_geolocation_detail" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `lang` - (Optional, String) Specifies the language type.
  The value can be **cn** or **en**. Defaults to **cn**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `continent` - The distribution information of country names on each continent, in JSON format.

* `geomap` - The key value represents the abbreviations of each country (except for AB and AB2, where AB indicates
  overseas and Hong Kong, Macao and Taiwan, and AB2 indicates overseas). When the key is CN, the array content inside
  is the abbreviation of each province.
  The `geomap` value in JSON format.

* `locale` - The display names of the corresponding languages for the values in `geomap`.
