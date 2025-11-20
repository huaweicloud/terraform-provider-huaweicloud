---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_security_report_history_periods"
description: |-
  Use this data source to get the list of security report history periods.
---

# huaweicloud_waf_security_report_history_periods

Use this data source to get the list of security report history periods.

## Example Usage

```hcl
variable "subscription_id" {}

data "huaweicloud_waf_security_report_history_periods" "test" {
  subscription_id = var.subscription_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `subscription_id` - (Required, String) Specifies the subscription ID of the security report.
  This value can be queried through the datasource `huaweicloud_waf_security_report_subscriptions`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The security report history periods list.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `report_id` - The unique identifier for the security report.

* `subscription_id` - The subscription ID associated with the security report.

* `stat_period` - The statistical period of the historical report.
  
  The [stat_period](#stat_period_struct) structure is documented below.

<a name="stat_period_struct"></a>
The `stat_period` block supports:

* `begin_time` - The start time of the statistical period (in milliseconds since epoch).

* `end_time` - The end time of the statistical period (in milliseconds since epoch).
