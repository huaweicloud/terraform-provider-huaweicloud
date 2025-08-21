---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_security_report_subscriptions"
description: |-
  Use this data source to get the list of security report subscriptions.
---

# huaweicloud_waf_security_report_subscriptions

Use this data source to get the list of security report subscriptions.

## Example Usage

```hcl
data "huaweicloud_waf_security_report_subscriptions" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `report_name` - (Optional, String) Specifies the security report template name.

* `report_category` - (Optional, String) Specifies the security report type.
  The valid values are as follows:
  + **daily_report**：Indicates security daily report.
  + **weekly_report**：Indicates security weekly report.
  + **monthly_report**：Indicates security monthly report.
  + **custom_report**: Indicates custom security report.

* `report_status` - (Optional, String) Specifies the security report status.
  The valid values are as follows:
  + **opened**
  + **closed**

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The security report subscription list.

  The [items](items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `subscription_id` - The subscription ID.

* `report_id` - The security report ID.

* `report_name` - The security report template name.

* `report_category` - The security report type.

* `report_status` - The security report status.

* `sending_period` - The send time period.
  The valid values are as follows:
  + **morning** (00:00~06:00)
  + **noon** (06:00~12:00)
  + **afternoon** (12:00~18:00)
  + **evening** (18:00~24:00)

* `is_all_enterprise_project` - Whether the security report belongs to all enterprise project.

* `enterprise_project_id` - The enterprise project ID.

* `template_eps_id` - The enterprise project ID to which the security report belongs.

* `is_report_created` - The security report generation status.

* `latest_create_time` - The security report generation time.
