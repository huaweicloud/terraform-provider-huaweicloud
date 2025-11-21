---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_security_report_sending_records"
description: |-
  Use this data source to get the security report sending records of WAF within HuaweiCloud.
---

# huaweicloud_waf_security_report_sending_records

Use this data source to get the security report sending records of WAF within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_waf_security_report_sending_records" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `report_name` - (Optional, String) Specifies the report name.

* `report_category` - (Optional, String) Specifies the report category. Valid values are:
  + **daily_report**
  + **weekly_report**
  + **monthly_report**
  + **custom_report**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The detailed list of security report sending records.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `report_id` - The report ID, uniquely identifies the security report corresponding to this sending record.

* `subscription_id` - The subscription ID is associated with the security report subscription to which the sending
  record belongs.

* `report_name` - The report name. The name of the security report corresponding to this sending record.

* `stat_period` - The statistical period refers to the statistical time range of the report corresponding to this sent
  record.

  The [stat_period](#items_stat_period_struct) structure is documented below.

* `report_category` - The report category.

* `sending_time` - The sending time. The timestamp (in milliseconds) at which the report was actually sent.

<a name="items_stat_period_struct"></a>
The `stat_period` block supports:

* `end_time` - The end time, the end timestamp of the statistical period (in milliseconds).

* `begin_time` - The start time, the timestamp of the start of the statistical period (in milliseconds).
