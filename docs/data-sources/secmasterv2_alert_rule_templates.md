---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmasterv2_alert_rule_templates"
description: |-
  Use this data source to get the list of alert rule templates (V2).
---

# huaweicloud_secmasterv2_alert_rule_templates

Use this data source to get the list of alert rule templates (V2).

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmasterv2_alert_rule_templates" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `template_name` - (Optional, String) Specifies the alert rule template name.

* `severity` - (Optional, String) Specifies the alert level.
  The value can be **TIPS**, **LOW**, **MEDIUM**, **HIGH** or **FATAL**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The list of alert rule templates.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `template_id` - The alert rule template ID.

* `template_name` - The alert rule template name.

* `accumulated_times` - The cumulative number of times.

* `alert_description` - The alert description.

* `alert_name` - The query type of the alert rule template.

* `alert_remediation` - The custom extension information.

* `alert_type` - The alert type.

* `create_by` - The user ID.

* `create_time` - The creation time.

* `custom_properties` - The custom properties.

* `description` - The alert rule template description.

* `event_grouping` - Whether to put events in a group.

* `job_mode` - The mode corresponding to alert rule template.

* `process_status` - The process status.
  The value can be **COMPLETED**, **CREATING**, **UPDATING**, **ENABLING**, **DISABLING**, **DELETING**,
  **CREATE_FAILED**, **UPDATE_FAILED**, **ENABLE_FAILED**, **DISABLE_FAILED**, **DELETE_FAILED** or **RECOVERING**.

* `query` - The query statement.

* `query_type` - The query type.
  The value can be **SQL** or **CBSL**.

* `schedule` - The alert rule schedule.

  The [schedule](#templates_schedule_struct) structure is documented below.

* `severity` - The alert level.

* `simulation` - Whether simulation alert.

* `status` - The alert rule template status.

* `suppression` - Whether alert suppression.

* `table_name` - The table name.

* `triggers` - The alert trigger rules information.

  The [triggers](#templates_triggers_struct) structure is documented below.

* `update_by` - The user ID which update the alert rule template.

* `update_time` - The update time.

* `update_time_by_user` - The update time which the user update the alert rule template.

<a name="templates_schedule_struct"></a>
The `schedule` block supports:

* `delay_interval` - The delay interval.

* `frequency_interval` - The scheduling interval.

* `frequency_unit` - The scheduling interval unit.
  The value can be **MINUTE**, **HOUR** or **DAY**.

* `overtime_interval` - The overtime interval.

* `period_interval` - The time window interval.

* `period_unit` - The time window unit.
  The value can be **MINUTE**, **HOUR** or **DAY**.

<a name="templates_triggers_struct"></a>
The `triggers` block supports:

* `accumulated_times` - The cumulative number of times.

* `expression` - The expression.

* `job_id` - The  job ID.

* `mode` - The mode.

* `operator` - The operator type.
  The valid values are as follows:
  + **EQ**: equal,
  + **NE**: not equal,
  + **GT**: greater than,
  + **LT**: less than.

* `severity` - The alert level.
