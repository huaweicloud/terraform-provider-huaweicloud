---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_security_report_subscription"
description: |-
  Use this data source to get the detail of a specific security report subscription.
---

# huaweicloud_waf_security_report_subscription

Use this data source to get the detail of a specific security report subscription.

## Example Usage

```hcl
variable "subscription_id" {}

data "huaweicloud_waf_security_report_subscription" "test" {
  subscription_id = var.subscription_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `subscription_id` - (Required, String) Specifies the ID of the security report subscription.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `sending_period` - The sending time period for the security report.
  The valid values are as follows:
  + **morning**：00:00~06:00
  + **noon**：06:00~12:00
  + **afternoon**：12:00~18:00
  + **evening**：18:00~24:00

* `report_name` - The security report template name.

* `report_category` - The security report type.
  The valid values are as follows:
  + **daily_report**：Indicates security daily report.
  + **weekly_report**：Indicates security weekly report.
  + **monthly_report**：Indicates security monthly report.
  + **custom_report**: Indicates custom security report.

* `topic_urn` - The URN of the SMN topic for receiving reports.

* `subscription_type` - The subscription type of the security report.
  The valid values are as follows:
  + **smn_topic**
  + **slient**
  + **message_center**

* `report_content_subscription` - The content subscription configuration of the security report.

  The [report_content_subscription](#waf_report_content_subscription) structure is documented below.

* `stat_period` - The statistical period of the security report.

  The [stat_period](#waf_security_report_subscription_stat_period) structure is documented below.

* `is_all_enterprise_project` - Whether the subscription applies to all enterprise projects.

* `enterprise_project_id` - The enterprise project ID associated with the subscription. This attribute is valid only
  when the attribute `is_all_enterprise_project` is **false**.

<a name="waf_report_content_subscription"></a>
The `report_content_subscription` block supports:

* `overview_statistics_enable` - Whether to enable overview statistics.

* `group_by_day_enable` - Whether to enable daily grouping statistics.

* `request_statistics_enable` - Whether to enable request statistics.

* `qps_statistics_enable` - Whether to enable QPS statistics.

* `bandwidth_statistics_enable` - Whether to enable bandwidth statistics.

* `response_code_statistics_enable` - Whether to enable response code statistics.

* `attack_type_distribution_enable` - Whether to enable attack type distribution statistics.

* `top_attacked_domains_enable` - Whether to enable top attacked domains statistics.

* `top_attack_source_ips_enable` - Whether to enable top attack source IPs statistics.

* `top_attacked_urls_enable` - Whether to enable top attacked URLs statistics.

* `top_attack_source_locations_enable` - Whether to enable top attack source locations statistics.

* `top_abnormal_urls_enable` - Whether to enable top abnormal URLs statistics.

<a name="waf_security_report_subscription_stat_period"></a>
The `stat_period` block supports:

* `begin_time` - The start time of the statistical period in milliseconds.

* `end_time` - The end time of the statistical period in milliseconds.
