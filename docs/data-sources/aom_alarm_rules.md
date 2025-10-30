---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_alarm_rules"
description: |-
  Use this data source to get the list of AOM alarm rules.
---

# huaweicloud_aom_alarm_rules

Use this data source to get the list of AOM alarm rules.

## Example Usage

```hcl
data "huaweicloud_aom_alarm_rules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `alarm_rule_name` - (Optional, String) Specifies the alarm rule name.

* `alarm_rule_status` - (Optional, String) Specifies the alarm rule status.
  Valid values are:
  + **OK**: normal
  + **alarm**: threshold-crossing
  + **Effective**: in use
  + **Invalid**: not in use

* `alarm_rule_type` - (Optional, String) Specifies the alarm rule type.
  Valid values are:
  + **metric**: metric alarm rule
  + **event**: event alarm rule

* `bind_notification_rule_id` - (Optional, String) Specifies the name of the bound alarm action rule.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the rule belongs.

* `event_severity` - (Optional, String) Specifies the alarm severity.
  Valid values are **Critical**, **Major**, **Minor** and **Info**.

* `event_source` - (Optional, String) Specifies the source of an event alarm rule.
  Valid values are **RDS**, **EVS**, **CCE**, **LTS** and **AOM**.

* `prom_instance_id` - (Optional, String) Specifies the prometheus instance ID.

* `related_cce_clusters` - (Optional, String) Specifies the related CCE cluster ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `alarm_rules` - Indicates the alarm rules list.
  The [alarm_rules](#attrblock--alarm_rules) structure is documented below.

<a name="attrblock--alarm_rules"></a>
The `alarm_rules` block supports:

* `name` - Indicates the alarm rule name.

* `alarm_rule_id` - Indicates the alarm rule ID.

* `enable` - Indicates whether the alarm rule enabled.

* `status` - Indicates the alarm rule status.

* `type` - Indicates the alarm rule type.

* `metric_alarm_spec` - Indicates the metric alarm rule.
  The [metric_alarm_spec](#attrblock--alarm_rules--metric_alarm_spec) structure is documented below.

* `event_alarm_spec` - Indicates the event alarm rule.
  The [event_alarm_spec](#attrblock--alarm_rules--event_alarm_spec) structure is documented below.

* `alarm_notifications` - Indicates the alarm notification.
  The [alarm_notifications](#attrblock--alarm_rules--alarm_notifications) structure is documented below.

* `prom_instance_id` - Indicates the prometheus instance ID.

* `description` - Indicates the alarm rule description.

* `enterprise_project_id` - Indicates the enterprise project ID to which the rule belongs.

* `created_at` - Indicates the time when an alarm rule was created.

* `updated_at` - Indicates the time when an alarm rule was updated.

<a name="attrblock--alarm_rules--alarm_notifications"></a>
The `alarm_notifications` block supports:

* `bind_notification_rule_id` - Indicates the alarm action rule ID.

* `notification_enable` - Indicates whether the alarm action rule enabled.

* `notification_type` - Indicates the notification type.
  Value can be as follows:
  + **direct**: direct alarm reporting
  + **alarm_policy**: alarm reporting after noise reduction

* `notify_frequency` - Indicates the notification frequency.

* `notify_resolved` - Indicates whether the notification enabled when an alarm is cleared.

* `notify_triggered` - Indicates whether the notification enabled when an alarm is triggered.

* `route_group_enable` - Indicates whether the grouping rule enabled.

* `route_group_rule` - Indicates the grouping rule name.

<a name="attrblock--alarm_rules--event_alarm_spec"></a>
The `event_alarm_spec` block supports:

* `alarm_source` - Indicates the alarm rule source.
  Value can be **systemEvent** and **customEvent**.

* `event_source` - Indicates the alarm source.

* `monitor_objects` - Indicates the list of monitored objects. Key-value pair.
  Value can be as follows:
  + **event_type**: notification type
  + **event_severity**: alarm severity
  + **event_name**: event name
  + **namespace**: namespace
  + **clusterId**: cluster ID
  + **customField**: user-defined field

* `trigger_conditions` - Indicates the trigger conditions.
  The [trigger_conditions](#attrblock--alarm_rules--event_alarm_spec--trigger_conditions) structure is documented below.

<a name="attrblock--alarm_rules--event_alarm_spec--trigger_conditions"></a>
The `trigger_conditions` block supports:

* `aggregation_window` - Indicates the statistical period, in seconds.

* `event_name` - Indicates the event name.

* `frequency` - Indicates the event alarm notification frequency.

* `operator` - Indicates the operator.

* `thresholds` - Key-value pair. The key indicates the alarm severity while the value indicates the alarm threshold.

* `trigger_type` - Indicates the trigger mode.
  Value can be as follows:
  + **immediately:** An alarm is triggered immediately if the alarm condition is met.
  + **accumulative**: An alarm is triggered if the alarm condition is met for a specified number of times.

<a name="attrblock--alarm_rules--metric_alarm_spec"></a>
The `metric_alarm_spec` block supports:

* `alarm_tags` - Indicates the alarm tags.
  The [alarm_tags](#attrblock--alarm_rules--metric_alarm_spec--alarm_tags) structure is documented below.

* `monitor_objects` - Indicates the list of monitored objects.

* `monitor_type` - Indicates the monitoring type.
  Value can be as follows:
  + **all_metric**: Select metrics from all metrics.
  + **promql**: Select metrics using PromQL.

* `no_data_conditions` - Indicates the action taken for insufficient data.
  The [no_data_conditions](#attrblock--alarm_rules--metric_alarm_spec--no_data_conditions) structure is documented below.

* `recovery_conditions` - Indicates the alarm clearance condition.
  The [recovery_conditions](#attrblock--alarm_rules--metric_alarm_spec--recovery_conditions) structure is documented below.

* `trigger_conditions` - Indicates the trigger conditions.
  The [trigger_conditions](#attrblock--alarm_rules--metric_alarm_spec--trigger_conditions) structure is documented below.

<a name="attrblock--alarm_rules--metric_alarm_spec--alarm_tags"></a>
The `alarm_tags` block supports:

* `auto_tags` - Indicates the automatic tag.

* `custom_annotations` - Indicates the alarm annotation.

* `custom_tags` - Indicates the custom tag.

<a name="attrblock--alarm_rules--metric_alarm_spec--no_data_conditions"></a>
The `no_data_conditions` block supports:

* `no_data_alert_state` - Indicates the status of the threshold rule when the data is insufficient.
  Value can be as follows:
  + **no_data**: A notification indicating insufficient data is sent.
  + **alerting**: An alarm is triggered.
  + **ok**: No exception occurs.
  + **pre_state**: Retain the previous state.

* `no_data_timeframe` - Indicates the number of periods without data.

* `notify_no_data` - Indicates whether to send a notification when data is insufficient.

<a name="attrblock--alarm_rules--metric_alarm_spec--recovery_conditions"></a>
The `recovery_conditions` block supports:

* `recovery_timeframe` - Indicates the number of consecutive periods for which the trigger condition is not met to clear
  an alarm.

<a name="attrblock--alarm_rules--metric_alarm_spec--trigger_conditions"></a>
The `trigger_conditions` block supports:

* `aggregate_type` - Indicates the aggregation mode.
  Value can be **by**, **avg**, **max**, **min** and **sum**.

* `aggregation_type` - Indicates the statistical mode.
  Value can be **average**, **minimum**, **maximum**, **sum** and **sampleCount**.

* `aggregation_window` - Indicates the statistical period.

* `aom_monitor_level` - Indicates the monitoring layer.

* `expression` - Indicates the expression of a combined operation.

* `metric_labels` - Indicates the metric dimension.

* `metric_name` - Indicates the metric name.

* `metric_namespace` - Indicates the metric namespace.

* `metric_query_mode` - Indicates the metric query mode.
  Value can be as follows:
  + **AOM**: native AOM
  + **PROM**: AOM Prometheus
  + **NATIVE_PROM**: native Prometheus

* `metric_statistic_method` - Indicates the metric statistics method to be used.
  Value can be as follows:
  + **single**: single metric
  + **mix**: multi-metric combined operations

* `metric_unit` - Indicates the metric unit.

* `mix_promql` - Indicates the promQL of a combined operation.

* `operator` - Indicates the operator.

* `promql` - Indicates the prometheus statement.

* `promql_expr` - Indicates the prometheus statement template.

* `promql_for` - Indicates the native Prometheus monitoring duration.

* `query_match` - Indicates the query filter criteria.

* `query_param` - Indicates the query parameters.

* `thresholds` - Key-value pair. The key indicates the alarm severity while the value indicates the alarm threshold.

* `trigger_interval` - Indicates the check interval.

* `trigger_times` - Indicates the number of consecutive periods.

* `trigger_type` - Indicates the trigger type.
  Value can be as follows:
  + **FIXED_RATE**: fixed interval
  + **HOURLY**: every hour
  + **DAILY**: every day
  + **WEEKLY**: every week
  + **CRON**: Cron expression
