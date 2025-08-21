---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_bundle"
description: |-
  Use this data source to query user bundle information.
---

# huaweicloud_waf_bundle

Use this data source to query user bundle information.

## Example Usage

```hcl
data "huaweicloud_waf_bundle" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `name` - The cloud mode bundle name.
  The valid values are as follows:
  + **None**: Indicates None.
  + **BASIC**: Indicates cloud mode getting started edition (yearly/monthly).
  + **Professional**: Indicates cloud mode standard edition (yearly/monthly).
  + **Enterprise**: Indicates cloud mode professional edition (yearly/monthly).
  + **Ultimate**: Indicates cloud mode enterprise edition (yearly/monthly).
  + **cloud.waf.postpaid**: Indicates cloud mode (pay-per-use).

* `type` - The cloud mode bundle type.
  The valid values are as follows:
  + `-2`: Indicates frozen.
  + `-1`: Indicates None.
  + `1`: Indicates cloud mode getting started edition (yearly/monthly).
  + `2`: Indicates cloud mode standard edition (yearly/monthly).
  + `3`: Indicates cloud mode professional edition (yearly/monthly).
  + `4`: Indicates cloud mode enterprise edition (yearly/monthly).
  + `22`: Indicates cloud mode (pay-per-use).

* `host` - The cloud mode supports domain quota information, in JSON format.

* `premium_name` - The dedicated mode bundle name.
  The valid values are as follows:
  + **None**: Indicates None.
  + **Instance.professional**: Indicates dedicated mode version specification is WI-100.
  + **Instance.enterprise**: Indicates dedicated mode version specification is WI-500.

* `premium_type` - The dedicated mode bundle type.
  The valid values are as follows:
  + `-2`: Indicates frozen.
  + `-1`: Indicates None.
  + `12`: Indicates dedicated mode version specification is WI-100.
  + `13`: Indicates dedicated mode version specification is WI-500.

* `premium_host` - The dedicated mode supports domain quota information, in JSON format.

* `options` - The policy related information.

* `rule` - The rule quota related information.

* `upgrade` - The different versions supports rule information.

* `feature` - The features information.
