---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_keyword_alarm_rules"
description: |-
  Use this data source to get the list of the keyword alarm rules with HuaweiCloud.
---

# huaweicloud_lts_keyword_alarm_rules

Use this data source to get the list of the keyword alarm rules with HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_lts_keyword_alarm_rules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `keyword_alarm_rules` - The list of the keyword alarm rules.

  The [keyword_alarm_rules](#keyword_alarm_rules_struct) structure is documented below.

<a name="keyword_alarm_rules_struct"></a>
The `keyword_alarm_rules` block supports:

* `id` - The ID of the keyword alarm rule.

* `name` - The name of the keyword alarm rule.

* `keywords_requests` - The detail of the keyword alarm rule.

  The [keywords_requests](#keyword_alarm_rules_keywords_requests_struct) structure is documented below.

* `frequency` - The configuration of the alarm query frequency.

  The [frequency](#keyword_alarm_rules_frequency_struct) structure is documented below.

* `alarm_level` - The level of the alarm.
  + **Info**
  + **Minor**
  + **Major**
  + **Critical**

* `description` - The description of the keyword alarm rule.

* `send_notifications` - Whether to send notification.

* `alarm_action_rule_name` - The name of the alarm action rule associated with the keyword alarm rule.

* `trigger_condition_count` - The count to trigger the alarm.

* `trigger_condition_frequency` - The frequency of trigger the alarm.

* `send_recovery_notifications` - Whether recovery notification is enabled.

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

* `topics` - The list of the SMN topics associated with the keyword alarm rule.

  The [topics](#keyword_alarm_rules_topics_struct) structure is documented below.

* `template_name` - The message templete name of the alarm action rule associated with the keyword alarm rule.

* `status` - The status of the keyword alarm rule.
  + **RUNNING**
  + **STOPPING**: Stopped.

* `domain_id` - The ID of the domain to which the keyword alarm rule belongs.

* `condition_expression` - The condition expression of the keyword alarm rule.

* `created_at` - The creation time of the keyword alarm rule, in RFC3339 format.

* `updated_at` - The update latest of the keyword alarm rule, in RFC3339 format.

<a name="keyword_alarm_rules_keywords_requests_struct"></a>
The `keywords_requests` block supports:

* `keywords` - The keyword in the queried logs.

* `condition` - The condition for triggering the alarm.
  + **>**
  + **<**
  + **>=**
  + **<=**

* `number` - The maximum number of logs containing the keyword that trigger the alarm.

* `log_group_id` - The ID of the log group to which the queried logs belong.

* `log_stream_id` - The ID of the log stream to which the queried logs belong.

* `search_time_range_unit` - The unit of search time range for querying the logs.

* `search_time_range` - The search time range for querying the logs.
  + **minute**
  + **hour**

* `log_stream_name` - The name of the log stream to which the queried logs belong.

* `log_group_name` - The name of the log stream to which the queried logs belong.

<a name="keyword_alarm_rules_frequency_struct"></a>
The `frequency` block supports:

* `type` - The frequency type of statistical alarm.
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

<a name="keyword_alarm_rules_topics_struct"></a>
The `topics` block supports:

* `name` - The name of the topic.

* `topic_urn` - The URN of the topic.

* `display_name` - The display name of the topic.

* `push_policy` - The push policy of the topic.
