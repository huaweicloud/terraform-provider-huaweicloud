---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_alarmrules"
description: |-
  Use this data source to get the list of CES alarm rules.
---

# huaweicloud_ces_alarmrules

Use this data source to get the list of CES alarm rules.

## Example Usage

```hcl
data "huaweicloud_ces_alarmrules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `alarm_id` - (Optional, String) Specifies the alarm rule ID.

* `name` - (Optional, String) Specifies the name of an alarm rule.

* `namespace` - (Optional, String) Specifies the namespace of a service.

* `resource_id` - (Optional, String) Specifies the alarm resource ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `alarms` - The alarm rule list.

  The [alarms](#alarms_struct) structure is documented below.

<a name="alarms_struct"></a>
The `alarms` block supports:

* `alarm_id` - The alarm rule ID.

* `alarm_name` - The alarm rule name.

* `alarm_description` - The alarm rule description.

* `namespace` - The namespace of a service.

* `condition` - The alarm triggering condition list.

  The [condition](#alarms_condition_struct) structure is documented below.

* `resources` - The resource list.

  The [resources](#alarms_resources_struct) structure is documented below.

* `alarm_type` - The alarm rule type.

* `alarm_enabled` - Whether to generate alarms when the alarm triggering conditions are met.

* `alarm_action_enabled` - Whether to enable the action to be triggered by an alarm.

* `alarm_actions` - The action to be triggered by an alarm.

  The [alarm_actions](#notification_struct) structure is documented below.

* `ok_actions` - The action to be triggered after an alarm is cleared.

  The [ok_actions](#notification_struct) structure is documented below.

* `notification_begin_time` - The time when the alarm notification was enabled.

* `notification_end_time` - The time when the alarm notification was disabled.

* `enterprise_project_id` - The enterprise project ID.

* `alarm_template_id` - The ID of an alarm template associated with an alarm rule.

<a name="alarms_condition_struct"></a>
The `condition` block supports:

* `metric_name` - The metric name of a resource.

* `period` - The monitoring period of a metric.

* `filter` - The filter method.

* `comparison_operator` - The comparison condition of alarm thresholds.

* `value` - The alarm threshold.

* `unit` - The metric unit.

* `count` - The number of times that the alarm triggering conditions are met.

* `suppress_duration` - The interval for triggering an alarm if the alarm persists.

* `alarm_level` - The alarm severity.

<a name="alarms_resources_struct"></a>
The `resources` block supports:

* `resource_group_id` - The resource group ID.

* `resource_group_name` - The resource group name.

* `dimensions` - The dimension information.

  The [dimensions](#resources_dimensions_struct) structure is documented below.

<a name="resources_dimensions_struct"></a>
The `dimensions` block supports:

* `name` - The name of the metric dimension.

<a name="notification_struct"></a>
The `alarm_actions` or `ok_actions` blocks support:

* `type` - The type of action triggered by an alarm.

* `notification_list` - The list of objects to be notified if the alarm status changes.
