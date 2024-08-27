---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_alert_rule_templates"
description: |-
  Use this data source to get the list of SecMaster alert rule templates.
---

# huaweicloud_secmaster_alert_rule_templates

Use this data source to get the list of SecMaster alert rule templates.

## Example Usage

```hcl
variable "workspace_id" {}
variable "severity" {}

data "huaweicloud_secmaster_alert_rule_templates" "test" {
  workspace_id = var.workspace_id
  severity     = var.severity
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `severity` - (Optional, List) Specifies the list of the severity.
  The value can be **TIPS**, **LOW**, **MEDIUM**, **HIGH** and **FATAL**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `templates` - The alert rule templates.

  The [templates](#templates_struct) structure is documented below.

<a name="templates_struct"></a>
The `templates` block supports:

* `id` - The alert rule template ID.

* `name` - The alert rule template name.

* `severity` - The severity.

* `query` - The query rule of the alert rule template.

* `query_type` - The query type of the alert rule template.

* `custom_properties` - The custom extension information.

* `event_grouping` - Whether to put events in a group.

* `data_source` - The data source.

* `version` - The version.

* `updated_at` - The update time.

* `triggers` - The alert triggering rules.

  The [triggers](#templates_triggers_struct) structure is documented below.

* `query_plan` - The query plan of the alert rule template.

  The [query_plan](#templates_query_plan_struct) structure is documented below.

<a name="templates_triggers_struct"></a>
The `triggers` block supports:

* `severity` - The severity.

* `expression` - The expression.

* `accumulated_times` - The accumulated times.

* `mode` - The mode.

* `operator` - The operator.
  + **EQ**: equal,
  + **NE**: not equal,
  + **GT**: greater than,
  + **LT**: less than.

<a name="templates_query_plan_struct"></a>
The `query_plan` block supports:

* `query_interval` - The query interval.

* `query_interval_unit` - The query interval unit.

* `time_window` - The time window.

* `time_window_unit` - The time window unit.

* `execution_delay` - The execution delay in minutes.

* `overtime_interval` - The overtime interval in minutes.
