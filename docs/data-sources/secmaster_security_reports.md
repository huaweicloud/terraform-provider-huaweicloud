---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_security_reports"
description: |-
  Use this data source to get the list of SecMaster security reports.
---

# huaweicloud_secmaster_security_reports

Use this data source to get the list of SecMaster security reports.

## Example Usage

```hcl
variable "workspace_id" {}
variable "report_period" {}
variable "status" {}

data "huaweicloud_secmaster_security_reports" "test" {
  workspace_id  = var.workspace_id
  report_period = var.report_period
  status        = var.status
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID to which the reports belong.

* `report_period` - (Required, String) Specifies the report period.
  The value can be **weekly**, **daily**, **monthly** or **annual**.

* `status` - (Required, String) Specifies the report status.
  The value can be **enable** or **disable**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `reports` - The list of report details.

  The [reports](#reports_struct) structure is documented below.

<a name="reports_struct"></a>
The `reports` block supports:

* `id` - The report ID.

* `report_name` - The report name.

* `report_period` - The report period.

* `report_range` - The data range.

  The [report_range](#report_range_struct) structure is documented below.

* `language` - The language.

* `notification_task` - The notification task ID.

* `layout_id` - The layout ID.

* `status` - The report status.

* `is_generated` - Whether the report has been generated.

* `report_rule_infos` - The report rules.

  The [report_rule_infos](#report_rule_infos_struct) structure is documented below.

<a name="report_range_struct"></a>
The `report_range` block supports:

* `start` - The start time.

* `end` - The end time.

<a name="report_rule_infos_struct"></a>
The `report_rule_infos` block supports:

* `id` - The report sending rule ID.

* `project_id` - The tenant ID.

* `workspace_id` - The workspace ID.

* `cycle` - The data cycle.

* `rule` - The cron expression.

* `start_time` - The report data cycle start time.

* `end_time` - The report data cycle end time.

* `email_title` - The email title.

* `email_to` - The recipient email.

* `email_content` - The email content.

* `report_file_type` - The report type.
