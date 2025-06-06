---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_alarms"
description: |-
  Use this data source to get the alarm list within HuaweiCloud.
---

# huaweicloud_lts_alarms

Use this data source to get the alarm list within HuaweiCloud.

## Query active alarms in the last 30 minutes

```hcl
data "huaweicloud_lts_alarms" "filter_by_time_range" {
  type                 = "active_alert"
  whether_custom_field = false
  time_range           = 30
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `type` - (Required, String) Specifies the type of the alarm to be queried.  
  The valid values are as follows:
  + **active_alert**: Query active alarms.
  + **history_alert**: Query historical alarms.

* `whether_custom_field` - (Optional, Bool) Specifies whether to customize the query time range, defaults to **false**.

* `time_range` - (Optional, String) Specifies the time range of the alarm to be queried, in minutes.  
  This parameter is required when `whether_custom_field` set to **false**.

* `search` - (Optional, String) Specifies the keyword search criteria.

* `alarm_level_ids` - (Optional, List) Specifies the list of alarm levels.  
  The valid values are as follows:
  + **Critical**
  + **Major**
  + **Minor**
  + **Info**

* `start_time` - (Optional, Int) Specifies the start time of a customized time segment, in milliseconds.  
  This parameter is required when `whether_custom_field` set to **true**.

* `end_time` - (Optional, Int) Specifies the end time of a customized time segment, in milliseconds.  
  This parameter is required when `whether_custom_field` set to **true**.

* `sort` - (Optional, List) Specifies the sort criteria of the queried alarms.  
  The [sort](#data_alarms_sort) structure is documented below.

* `step` - (Optional, Int) Specifies the step of the query, in milliseconds.

<a name="data_alarms_sort"></a>
The `sort` block supports:

* `order` - (Required, String) Specifies the sort mode of the alarm.  
  The valid values are as follows:
  + **asc**
  + **desc**

* `order_by` - (Required, List) Specifies the fields to be sorted.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `alarms` - The list of the queried alarms.  
  The [alarms](#data_alarms) structure is documented below.

<a name="data_alarms"></a>
The `alarms` block supports:

* `id` - The ID of the alarm.

* `type` - The type of the alarm.
  + **keywords**
  + **SQL**

* `timeout` - The time when the alarm is automatically cleared, in milliseconds.

* `arrives_at` - The time when the alarm arrives, in milliseconds.

* `ends_at` - The time when the alarm is cleared, in milliseconds.

* `starts_at` - The time when the alarm is generated, in milliseconds.

* `annotations` - The details of the alarm.  
  The [annotations](#data_alarms_annotations) structure is documented below.

* `metadata` - The metadata of the alarm.  
  The [metadata](#data_alarms_metadata) structure is documented below.

<a name="data_alarms_annotations"></a>
The `annotations` block supports:

* `type` - The type of the alarm rule.

* `message` - The detail information of the alarm.

* `log_info` - The log information of the alarm.

* `current_value` - The current value of the alarm.

* `old_annotations` - The raw data of the alarm detail.

* `alarm_action_rule_name` - The name of the alarm action rule.

* `alarm_rule_alias` - The alias of the alarm rule.

* `alarm_rule_url` - The URL of the alarm rule.

* `alarm_status` - The status of the alarm trigger.

* `condition_expression` - The condition expression of the alarm trigger.

* `condition_expression_with_value` - The condition of the alarm trigger.

* `notification_frequency` - The notification frequency of the alarm.

* `recovery_policy` - Whether the alarm is recovered.

* `frequency` - The frequency of the alarm.

<a name="data_alarms_metadata"></a>
The `metadata` block supports:

* `event_id` - The ID of the alarm rule.

* `event_name` - The name of the alarm rule.

* `event_type` - The mode of the alarm.

* `event_severity` - The level of the alarm.

* `resource_provider` - The source of the alarm.

* `lts_alarm_type` - The type of the alarm rule.

* `resource_id` - The ID of the resource.

* `resource_type` - The type of the resource.

* `log_group_name` - The original name of the log group.

* `log_stream_name` - The original name of the log stream.

* `event_subtype` - The type of the alarm.
