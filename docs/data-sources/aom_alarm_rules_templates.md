---
subcategory: "Application Operations Management (AOM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_alarm_rules_templates"
description: |-
  Use this data source to get the list of AOM alarm rules templates.
---

# huaweicloud_aom_alarm_rules_templates

Use this data source to get the list of AOM alarm rules templates.

## Example Usage

```hcl
data "huaweicloud_aom_alarm_rules_templates" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the templates belong.

* `template_id` - (Optional, String) Specifies the template ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `templates` - Indicates the templates.
  The [templates](#attrblock--templates) structure is documented below.

<a name="attrblock--templates"></a>
The `templates` block supports:

* `id` - Indicates the template ID.

* `name` - Indicates the template name.

* `type` - Indicates the template type.

* `alarm_template_spec_list` - Indicates the alarm template spec list.
  The [alarm_template_spec_list](#attrblock--templates--alarm_template_spec_list) structure is documented below.

* `templating` - Indicates the variable list.
  The [templating](#attrblock--templates--templating) structure is documented below.

* `description` - Indicates the template description.

* `enterprise_project_id` - Indicates the enterprise project ID to which the template belongs.

* `created_at` - Indicates the template create time.

* `updated_at` - Indicates the template update time.

<a name="attrblock--templates--alarm_template_spec_list"></a>
The `alarm_template_spec_list` block supports:

* `alarm_notification` - Indicates the alarm notification.
  The [alarm_notification](#attrblock--templates--alarm_template_spec_list--alarm_notification) structure is documented
  below.

* `alarm_template_spec_items` - Indicates the alarm template spec items.
  The [alarm_template_spec_items](#attrblock--templates--alarm_template_spec_list--alarm_template_spec_items) structure
  is documented below.

* `related_cloud_service` - Indicates the related cloud service.

* `related_cce_clusters` - Indicates the related CCE clusters.

* `related_prometheus_instances` - Indicates the related prometheus instances.

<a name="attrblock--templates--alarm_template_spec_list--alarm_notification"></a>
The `alarm_notification` block supports:

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

<a name="attrblock--templates--alarm_template_spec_list--alarm_template_spec_items"></a>
The `alarm_template_spec_items` block supports:

* `alarm_rule_description` - Indicates the alarm rule description.

* `alarm_rule_name` - Indicates the alarm rule name.

* `alarm_rule_type` - Indicates the alarm rule type.

* `event_alarm_spec` - Indicates the event alarm spec.
  The [event_alarm_spec](#attrblock--templates--alarm_template_spec_list--alarm_template_spec_items--event_alarm_spec)
  structure is documented below.

* `metric_alarm_spec` - Indicates the metric alarm spec.
  The [metric_alarm_spec](#attrblock--templates--alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec)
  structure is documented below.

<a name="attrblock--templates--alarm_template_spec_list--alarm_template_spec_items--event_alarm_spec"></a>
The `metric_alarm_spec` block supports:

* `alarm_source` - Indicates the alarm source.

* `alarm_subtype` - Indicates the alarm subtype.

* `event_source` - Indicates the event source.

* `monitor_object_templates` - Indicates the monitor object templates.

* `monitor_objects` - Indicates the monitor objects.
  Value can be as follows:
  + **event_type**: notification type
  + **event_severity**: alarm severity
  + **event_name**: event name
  + **namespace**: namespace
  + **clusterId**: cluster ID
  + **customField**: user-defined field

* `trigger_conditions` - Indicates the trigger conditions.
  The [trigger_conditions](#attrblock--templates--alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec--trigger_conditions)
  structure is documented below.

<a name="attrblock--templates--alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec--trigger_conditions"></a>
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

<a name="attrblock--templates--alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec"></a>
The `metric_alarm_spec` block supports:

* `alarm_source` - Indicates the alarm source.

* `alarm_subtype` - Indicates the alarm subtype.

* `monitor_type` - Indicates the monitor_type.

* `alarm_tags` - Indicates the alarm tags.
  The [alarm_tags](#attrblock--templates--alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec--alarm_tags)
  structure is documented below.

* `no_data_conditions` - Indicates the no data conditions.
  The [no_data_conditions](#attrblock--templates--alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec--no_data_conditions)
  structure is documented below.

* `recovery_conditions` - Indicates the recovery conditions.
  The [recovery_conditions](#attrblock--templates--alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec--recovery_conditions)
  structure is documented below.

* `trigger_conditions` - Indicates the trigger conditions.
  The [trigger_conditions](#attrblock--templates--alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec--trigger_conditions)
  structure is documented below.

<a name="attrblock--templates--alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec--alarm_tags"></a>
The `alarm_tags` block supports:

* `auto_tags` - Indicates the automatic tag.

* `custom_annotations` - Indicates the alarm annotation.

* `custom_tags` - Indicates the custom tag.

<a name="attrblock--templates--alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec--no_data_conditions"></a>
The `no_data_conditions` block supports:

* `no_data_alert_state` - Indicates the status of the threshold rule when the data is insufficient.
  Value can be as follows:
  + **no_data**: A notification indicating insufficient data is sent.
  + **alerting**: An alarm is triggered.
  + **ok**: No exception occurs.
  + **pre_state**: Retain the previous state.

* `no_data_timeframe` - Indicates the number of periods without data.

* `notify_no_data` - Indicates whether to send a notification when data is insufficient.

<a name="attrblock--templates--alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec--recovery_conditions"></a>
The `recovery_conditions` block supports:

* `recovery_timeframe` - Indicates the number of consecutive periods for which the trigger condition is not met to clear
  an alarm.

<a name="attrblock--templates--alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec--trigger_conditions"></a>
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

<a name="attrblock--templates--templating"></a>
The `templating` block supports:

* `list` - Indicates the variable list.
  The [list](#attrblock--templates--templating--list) structure is documented below.

<a name="attrblock--templates--templating--list"></a>
The `list` block supports:

* `description` - Indicates the variable description.

* `name` - Indicates the variable name.

* `query` - Indicates the variable value.

* `type` - Indicates the variable type.
