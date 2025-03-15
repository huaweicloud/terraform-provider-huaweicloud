---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_alarm_histories"
description: |-
  Use this data source to get the list of CES alarm history records.
---

# huaweicloud_ces_alarm_histories

Use this data source to get the list of CES alarm history records.

## Example Usage

```hcl
data "huaweicloud_ces_alarm_histories" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `alarm_id` - (Optional, String) Specifies an alarm ID.

* `record_id` - (Optional, String) Specifies the alarm record ID.

* `name` - (Optional, String) Specifies the alarm rule name.

* `alarm_type` - (Optional, String) Specifies the alarm type.
  The valid value can be **event** (querying event alarms) or **metric** (querying metric alarms).

* `status` - (Optional, String) Specifies the alarm rule status.
  The valid value can be **ok**, **alarm** or **invalid**.

* `level` - (Optional, Int) Specifies the alarm severity.
  The valid value can be **1** (critical), **2** (major), **3** (minor) or **4** (informational).

* `namespace` - (Optional, String) Specifies the namespace of a service.

* `resource_id` - (Optional, String) Specifies the ID of a resource in an alarm rule.

* `from` - (Optional, String) Specifies the start time for querying alarm records.
  For example, **2022-02-10T10:05:46+08:00**.

* `to` - (Optional, String) Specifies the end time for querying alarm records.
  For example, **2022-02-10T10:05:47+08:00**.

* `order_by` - (Optional, String) Specifies the keyword for sorting alarms.
  The valid values are as follows:
  + **first_alarm_time**: time for generating the alarm for the first time;
  + **update_time**: alarm update time, The default value;
  + **alarm_level**: alarm severity;
  + **record_id**: primary key of the table record;

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `alarm_histories` - The alarm records.

  The [alarm_histories](#alarm_histories_struct) structure is documented below.

<a name="alarm_histories_struct"></a>
The `alarm_histories` block supports:

* `record_id` - The alarm record ID.

* `alarm_id` - The alarm rule ID.

* `name` - The alarm rule name.

* `status` - The status of an alarm record.

* `level` - The alarm severity of alarm records.

* `type` - The alarm rule type.

* `action_enabled` - Whether to send a notification.

* `begin_time` - When an alarm record is generated (UTC time).

* `end_time` - When an alarm record becomes invalid (UTC time).

* `first_alarm_time` - The UTC time when the alarm was generated for the first time.

* `last_alarm_time` - The UTC time when the alarm was generated for the last time.

* `alarm_recovery_time` - The UTC time when the alarm was cleared.

* `metric` - The metric information.

  The [metric](#alarm_histories_metric_struct) structure is documented below.

* `condition` - The conditions for triggering an alarm.

  The [condition](#alarm_histories_condition_struct) structure is documented below.

* `additional_info` - The additional field of an alarm record.

  The [additional_info](#alarm_histories_additional_info_struct) structure is documented below.

* `alarm_actions` - The action to be triggered by an alarm.

  The [alarm_actions](#alarm_histories_alarm_actions_struct) structure is documented below.

* `ok_actions` - The action to be triggered after an alarm is cleared.

  The [ok_actions](#alarm_histories_ok_actions_struct) structure is documented below.

* `data_points` - The time when the resource monitoring data is reported and the monitoring data in the alarm record.

  The [data_points](#alarm_histories_data_points_struct) structure is documented below.

<a name="alarm_histories_metric_struct"></a>
The `metric` block supports:

* `namespace` - The namespace of a service.

* `metric_name` - The metric name of a resource.

* `dimensions` - The metric dimension.

  The [dimensions](#metric_dimensions_struct) structure is documented below.

<a name="metric_dimensions_struct"></a>
The `dimensions` block supports:

* `name` - The dimension of a resource.

* `value` - The value of a resource dimension.

<a name="alarm_histories_condition_struct"></a>
The `condition` block supports:

* `period` - The rollup period of a metric, in seconds.

* `filter` - The rollup method.

* `comparison_operator` - The threshold symbol.

* `value` - The alarm threshold.

* `unit` - The data unit.

* `count` - The number of times that the alarm triggering conditions are met.

* `suppress_duration` - The alarm suppression time, in seconds.

<a name="alarm_histories_additional_info_struct"></a>
The `additional_info` block supports:

* `resource_id` - The resource ID corresponding to the alarm record.

* `resource_name` - The resource name corresponding to the alarm record.

* `event_id` - The ID of the event in the alarm record.

<a name="alarm_histories_alarm_actions_struct"></a>
The `alarm_actions` block supports:

* `type` - The notification type.

* `notification_list` - The list of objects to be notified if the alarm status changes.

<a name="alarm_histories_ok_actions_struct"></a>
The `ok_actions` block supports:

* `type` - The notification type.

* `notification_list` - The list of objects to be notified if the alarm status changes.

<a name="alarm_histories_data_points_struct"></a>
The `data_points` block supports:

* `time` - The UTC time when the resource monitoring data of the alarm record is reported.

* `value` - The resource monitoring data of the alarm record at the time point.
