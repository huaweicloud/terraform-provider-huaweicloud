---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_alert_rule_template_detail"
description: |-
  Use this data source to query a specific alert rule template detail.
---

# huaweicloud_secmaster_alert_rule_template_detail

Use this data source to query a specific alert rule template detail.

## Example Usage

```hcl
variable "workspace_id" {}
variable "template_id" {}

data "huaweicloud_secmaster_alert_rule_template_detail" "test" {
  workspace_id = var.workspace_id
  template_id  = var.template_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `template_id` - (Required, String) Specifies the alert rule template ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rule_template_id` - The alert rule template ID.

* `template_name` - The alert rule template name.

* `data_source` - The data source.

* `version` - The version.

* `query` - The query statement.

* `query_type` - The query type.

* `severity` - The severity.
  The value can be **TIPS**, **LOW**, **MEDIUM**, **HIGH** or **FATAL**.

* `custom_properties` - The custom properties.

* `event_grouping` - The event grouping.

* `update_time` - The update time.

* `schedule` - The schedule rule.

  The [schedule](#schedule_struct) structure is documented below.

* `triggers` - The triggers information.

  The [triggers](#triggers_struct) structure is documented below.

<a name="schedule_struct"></a>
The `schedule` block supports:

* `frequency_interval` - The scheduling interval.

* `frequency_unit` - The scheduling interval time unit.
  The value can be **MINUTE**, **HOUR** or **DAY**.

* `period_interval` - The time window interval.

* `period_unit` - The time window unit.
  The value can be **MINUTE**, **HOUR** or **DAY**.

* `delay_interval` - The delay interval.

* `overtime_interval` - The overtime interval.

<a name="triggers_struct"></a>
The `triggers` block supports:

* `mode` - The mode.

* `operator` - The operator.
  The valid values are as follows:
  + **EQ**: equal,
  + **NE**: not equal,
  + **GT**: greater than,
  + **LT**: less than.

* `expression` - The expression.

* `severity` - The severity.

* `accumulated_times` - The cumulative number of times.
