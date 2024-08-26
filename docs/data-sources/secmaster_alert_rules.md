---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_alert_rules"
description: |-
  Use this data source to get the list of SecMaster alert rules.
---

# huaweicloud_secmaster_alert_rules

Use this data source to get the list of SecMaster alert rules.

## Example Usage

```hcl
variable "workspace_id" {}
variable "rule_id" {}

data "huaweicloud_secmaster_alert_rules" "test" {
  workspace_id = var.workspace_id
  rule_id      = var.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `rule_id` - (Optional, String) Specifies the alert rule ID.

* `name` - (Optional, String) Specifies the alert rule name. Fuzzy match is supported.

* `status` - (Optional, List) Specifies the list of the status. The value can be **ENABLED** and **DISABLED**.

* `severity` - (Optional, List) Specifies the list of the severity.
  The value can be **TIPS**, **LOW**, **MEDIUM**, **HIGH** and **FATAL**.

* `pipeline_id` - (Optional, String) Specifies the pipeline ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `alert_rules` - The alert rules.

  The [alert_rules](#alert_rules_struct) structure is documented below.

<a name="alert_rules_struct"></a>
The `alert_rules` block supports:

* `id` - The alert rule ID.

* `name` - The alert rule name.

* `status` - The status.

* `severity` - The severity.

* `pipeline_id` - The pipeline ID.

* `pipeline_name` - The data pipeline name.

* `query_rule` - The query rule of the alert rule.

* `query_type` - The query type of the alert rule.

* `custom_properties` - The custom extension information.

* `event_grouping` - Whether to put events in a group.

* `created_by` - The creator.

* `updated_by` - The updater.

* `created_at` - The creation time.

* `updated_at` - The update time.

* `deleted_at` - The deletion time.

* `query_plan` - The query plan of the alert rule.

  The [query_plan](#alert_rules_query_plan_struct) structure is documented below.

* `triggers` - The alert triggering rules.

  The [triggers](#alert_rules_triggers_struct) structure is documented below.

<a name="alert_rules_query_plan_struct"></a>
The `query_plan` block supports:

* `query_interval` - The query interval.

* `query_interval_unit` - The query interval unit.

* `time_window` - The time window.

* `time_window_unit` - The time window unit.

* `execution_delay` - The execution delay in minutes.

* `overtime_interval` - The overtime interval in minutes.

<a name="alert_rules_triggers_struct"></a>
The `triggers` block supports:

* `mode` - The mode.

* `operator` - The operator.
  + **EQ**: equal,
  + **NE**: not equal,
  + **GT**: greater than,
  + **LT**: less than.

* `expression` - The expression.

* `severity` - The severity.

* `accumulated_times` - The accumulated times.
