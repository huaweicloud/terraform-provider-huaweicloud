---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_security_reports"
description: |-
  Use this data source to get the security report list of HSS within HuaweiCloud.
---

# huaweicloud_hss_security_reports

Use this data source to get the security report list of HSS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_security_reports" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `report_category` - (Optional, String) Specifies the report category.  
  The valid values are as follows:
  + **daily_report**: Daily security report.
  + **weekly_report**: Weekly security report.
  + **monthly_report**: Monthly security report.
  + **custom_report**: Custom report.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The security report list.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `report_id` - The report ID.

* `report_sub_id` - The report sub ID.

* `default_report` - Whether it is the default report. Default reports cannot be deleted.

* `latest_create_time` - The latest creation time in milliseconds. A null value indicates that the report has not been
  generated yet.

* `report_name` - The report name.

* `report_category` - The report category.  
  The valid values are as follows:
  + **daily_report**: Daily security report.
  + **weekly_report**: Weekly security report.
  + **monthly_report**: Monthly security report.
  + **custom_report**: Custom report.

* `report_status` - The report status.  
  The valid values are as follows:
  + **opened**: Opened.
  + **closed**: Closed.

* `report_create_time` - The report creation time in milliseconds.

* `sending_period` - The report sending period.  
  The valid values are as follows:
  + **morning**: From `0:00` to `6:00`.
  + **noon**: From `6:00` to `12:00`.
  + **afternoon**: From `12:00` to `18:00`.
  + **evening**: From `18:00` to `24:00`.
