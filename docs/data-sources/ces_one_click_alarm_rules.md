---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_one_click_alarm_rules"
description: |-
  Use this data source to get the list of CES one-click alarm rules.
---

# huaweicloud_ces_one_click_alarm_rules

Use this data source to get the list of CES one-click alarm rules.

## Example Usage

```hcl
variable "one_click_alarm_id" {}

data "huaweicloud_ces_one_click_alarm_rules" "test" {
  one_click_alarm_id = var.one_click_alarm_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `one_click_alarm_id` - (Required, String) Specifies the one-click monitoring ID for a service.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `alarms` - The alarm rule list.

  The [alarms](#alarms_struct) structure is documented below.

<a name="alarms_struct"></a>
The `alarms` block supports:

* `alarm_id` - The ID of an alarm rule.

* `name` - The alarm rule name.

* `description` - The supplementary information about an alarm rule.

* `namespace` - The metric namespace.

* `policies` - The alarm policy list.

  The [policies](#alarms_policies_struct) structure is documented below.

* `resources` - The resource list.

  The [resources](#alarms_resources_struct) structure is documented below.

* `type` - The alarm rule type.

* `enabled` - Whether to generate alarms when the alarm triggering conditions are met.

* `notification_enabled` - Whether the alarm notification is enabled.

* `alarm_notifications` - The action to be triggered by an alarm.

  The [alarm_notifications](#notifications_struct) structure is documented below.

* `ok_notifications` - The action to be triggered after an alarm is cleared.

  The [ok_notifications](#notifications_struct) structure is documented below.

* `notification_begin_time` - The time when the alarm notification was enabled.

* `notification_end_time` - The time when the alarm notification was disabled.

<a name="alarms_policies_struct"></a>
The `policies` block supports:

* `alarm_policy_id` - The alarm policy ID.

* `metric_name` - The metric name.

* `period` - How often to generate an alarm.

* `comparison_operator` - The operator of an alarm threshold.

* `filter` - The roll up method.

* `value` - The threshold.

* `unit` - The metric unit.

* `count` - The number of times that the alarm triggering conditions are met.

* `suppress_duration` - The suppression period.

* `level` - The alarm severity.

* `enabled` - Whether the one-click monitoring is enabled.

<a name="alarms_resources_struct"></a>
The `resources` block supports:

* `resource_group_id` - The resource group ID.

* `resource_group_name` - The resource group name.

* `dimensions` - The dimension information.

  The [dimensions](#resources_dimensions_struct) structure is documented below.

<a name="resources_dimensions_struct"></a>
The `dimensions` block supports:

* `name` - The name of the metric dimension.

* `value` - The value of the metric dimension.

<a name="notifications_struct"></a>
The `alarm_notifications` or `ok_notifications` block supports:

* `type` - The notification type.

* `notification_list` - The list of objects to be notified if the alarm status changes.
