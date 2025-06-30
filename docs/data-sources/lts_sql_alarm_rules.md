---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_sql_alarm_rules"
description: |-
  Use this data source to get the list of the SQL alarm rules with HuaweiCloud.
---

# huaweicloud_lts_sql_alarm_rules

Use this data source to get the list of the SQL alarm rules with HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_lts_sql_alarm_rules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `sql_alarm_rules` - The list of the SQL alarm rules.

  The [sql_alarm_rules](#sql_alarm_rules_struct) structure is documented below.

<a name="sql_alarm_rules_struct"></a>
The `sql_alarm_rules` block supports:

* `id` - The ID of the SQL alarm rule.

* `name` - The name of the SQL alarm rule.

* `sql_requests` - The request list of the SQL alarm rule.

  The [sql_requests](#sql_alarm_rules_sql_requests_struct) structure is documented below.

* `frequency` - The alarm frequency configuration list.

  The [frequency](#sql_alarm_rules_frequency_struct) structure is documented below.

* `condition_expression` - The condition expression.

* `alarm_level` - The level of the alarm.
  + **Info**
  + **Minor**
  + **Major**
  + **Critical**

* `description` - The description of the SQL alarm rule.

* `send_notifications` - Whether to send notification.

* `alarm_action_rule_name` - The name of the alarm action rule associated with the SQL alarm rule.

* `trigger_condition_count` - The count to trigger the alarm.

* `trigger_condition_frequency` - The frequency to trigger the alarm.

* `send_recovery_notifications` - Whether to send recovery notification.

* `recovery_frequency` - The frequency of recovery the alarm notification.

* `notification_frequency` - The notification frequency of the alarm, in minutes.
  + **0**
  + **5**
  + **10**
  + **15**
  + **30**
  + **60**
  + **180**
  + **360**

* `domain_id` - The ID of the domain to which the SQL alarm rule belongs.

* `topics` - The list of the SMN topics associated with the SQL alarm rule.

  The [topics](#sql_alarm_rules_topics_struct) structure is documented below.

* `template_name` - The message template name of the alarm action rule associated with the SQL alarm rule.

* `status` - The status of the SQL alarm rule.
  + **RUNNING**
  + **STOPPING**: Stopped.

* `created_at` - The creation time of the SQL alarm rule, in RFC3339 format.

* `updated_at` - The latest update of the SQL alarm rule, in RFC3339 format.

<a name="sql_alarm_rules_sql_requests_struct"></a>
The `sql_requests` block supports:

* `title` - The title of the SQL request.

* `sql` - The SQL statement.

* `log_stream_id` - The ID of the log stream.

* `log_stream_name` - The name of the log stream.

* `log_group_id` - The ID of the log group.

* `log_group_name` - The name of the log group.

* `search_time_range` - The search time range.

* `search_time_range_unit` - The unit of search time range.
  + **minute**
  + **hour**

* `is_time_range_relative` - The SQL request is relative to time range.

<a name="sql_alarm_rules_frequency_struct"></a>
The `frequency` block supports:

* `type` - The type of the frequency.
  + **CRON**
  + **HOURLY**
  + **DAILY**
  + **WEEKLY**
  + **FIXED_RATE**

* `cron_expression` - The cron expression.

* `hour_of_day` - The hour of day.

* `day_of_week` - The day of week.

* `fixed_rate_unit` - The unit of custom interval for querying alarm.
  + **minute**
  + **hour**

* `fixed_rate` - The times of custom interval for querying alarm.

<a name="sql_alarm_rules_topics_struct"></a>
The `topics` block supports:

* `name` - The name of the topic.

* `topic_urn` - The URN of the topic.

* `display_name` - The display name of the topic.

* `push_policy` - The push policy of the topic.
