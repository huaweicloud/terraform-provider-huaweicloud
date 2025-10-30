---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_alarm_rules_template"
description: |-
  Manages an AOM cloud alarm rules template resource within HuaweiCloud.
---
# huaweicloud_aom_alarm_rules_template

Manages an AOM cloud alarm rules template resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "action_rule_id" {}

resource "huaweicloud_aom_alarm_rules_template" "test" {
  name        = var.name
  type        = "statics"
  description = "test"

  alarm_template_spec_list {
    related_cloud_service = "CCEFromProm"

    alarm_notifications {
      notification_type         = "direct"
      notification_enable       = true
      bind_notification_rule_id = var.action_rule_id
    }

    alarm_template_spec_items {
      alarm_rule_name = "cce_event"
      alarm_rule_type = "event"

      event_alarm_spec {
        event_source = "CCE"

        monitor_objects = [
          {
            event_name = "扩容节点超时##ScaleUpTimedOut;数据卷扩容失败##VolumeResizeFailed"
          },
        ]

        monitor_object_templates = ["clusterId"]

        trigger_conditions {
          trigger_type = "immediately"
          event_name   = "扩容节点超时##ScaleUpTimedOut"

          thresholds = {
            "Critical" = 2
          }
        }

        trigger_conditions {
          trigger_type       = "accumulative"
          event_name         = "数据卷扩容失败##VolumeResizeFailed"
          aggregation_window = 300
          frequency          = "600"
          operator           = ">="

          thresholds = {
            "Info" = 5
          }
        }
      }
    }

    alarm_template_spec_items {
      alarm_rule_name = "cce_metric"
      alarm_rule_type = "metric"

      metric_alarm_spec {
        monitor_type = "promql"

        recovery_conditions {
          recovery_timeframe = 1
        }

        trigger_conditions {
          metric_query_mode  = "NATIVE_PROM"
          promql             = "increase(kube_pod_container_status_restarts_total[5m]) > 3"
          trigger_times      = "3"
          trigger_type       = "FIXED_RATE"
          aggregation_window = "1m"
          trigger_interval   = "30s"
          aggregation_type   = "average"
          operator           = ">"
          promql_for         = "1m"

          thresholds = {
            "Critical" = "1"
          }
        }
      }
    }
  }

  alarm_template_spec_list {
    related_cloud_service = "DRS"

    alarm_notifications {
      notification_type = "direct"
    }

    alarm_template_spec_items {
      alarm_rule_name = "drs"
      alarm_rule_type = "metric"

      metric_alarm_spec {
        monitor_type = "resource"

        recovery_conditions {
          recovery_timeframe = 1
        }

        trigger_conditions {
          metric_query_mode       = "PROM"
          metric_name             = "huaweicloud_sys_drs_cpu_util"
          promql                  = "label_replace(avg_over_time(duration{}[59999ms]),\"__name__\",\"duration\",\"\",\"\")"
          trigger_times           = "3"
          trigger_type            = "FIXED_RATE"
          aggregation_window      = "1m"
          trigger_interval        = "30s"
          aggregation_type        = "average"
          operator                = ">"
          metric_statistic_method = "single"

          thresholds = {
            "Critical" = "1"
          }
        }
      }
    }
  }

  templating {
    list {
      name  = "key"
      type  = "constant"
      query = "value"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the template name.
  Changing this creates a new resource.

* `type` - (Required, String) Specifies the template type.
  Valid values are as follows:
  + **statics**: Static alarm template.
  + **dynamic**: Dynamic alarm template.

* `alarm_template_spec_list` - (Required, List) Specifies the alarm template spec list.
  The [alarm_template_spec_list](#alarm_template_spec_list) structure is documented below.

* `templating` - (Optional, List) Specifies the variable list.
  The [templating](#templating) structure is documented below.

* `description` - (Optional, String) Specifies the description of template.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the template belongs.
  Changing this creates a new resource.

<a name="alarm_template_spec_list"></a>
The `alarm_template_spec_list` block supports:

* `related_cloud_service` - (Optional, String) Specifies the related cloud service of the alarm rules.

* `related_cce_clusters` - (Optional, List) Specifies the related cce clusters of the alarm rules.

* `related_prometheus_instances` - (Optional, List) Specifies the related prometheus instances of the alarm rules.

* `alarm_notification` - (Optional, List) Specifies the alarm notification.
  The [alarm_notification](#alarm_template_spec_list--alarm_notification) structure is documented below.

* `alarm_template_spec_items` - (Optional, List) Specifies the alarm template spec items.
  The [alarm_template_spec_items](#alarm_template_spec_list--alarm_template_spec_items) structure is documented below.

<a name="alarm_template_spec_list--alarm_notification"></a>
The `alarm_notification` block supports:

* `notification_type` - (Required, String) Specifies the notification type.
  Valid values are as follows:
  + **direct**: Direct alarm reporting.
  + **alarm_policy**: Alarm reporting after noise reduction.

* `notification_enable` - (Optional, Bool) Specifies whether to enable an alarm action rule. Defaults to **false**.
  If the `notification_type` is **direct**, set this parameter to **true**.

* `bind_notification_rule_id` - (Optional, String) Specifies the alarm action rule ID.
  It's required if `notification_enable` is **true**.

* `route_group_enable` - (Optional, Bool) Specifies whether to enable the grouping rule. Defaults to **false**.
  If the `notification_type` is **alarm_policy**, set this parameter to **true**.

* `route_group_rule` - (Optional, String) Specifies the grouping rule name.
  It's required if `route_group_enable` is **true**.

* `notify_resolved` - (Optional, Bool) Specifies whether to send a notification when an alarm is cleared.
  Defaults to **false**.

* `notify_triggered` - (Optional, Bool) Specifies whether to send a notification when an alarm is triggered.
  Defaults to **false**.

* `notify_frequency` - (Optional, String) Specifies the notification frequency.
  If the `notification_type` is **alarm_policy**, set this parameter to **-1**.
  If the `notification_type` is **direct**, set this parameter to any of the following:
  + **0**: alarm sent only once
  + **300**: every 5 minutes
  + **600**: every 10 minutes
  + **900**: every 15 minutes
  + **1800**: every 30 minutes
  + **3600**: every hour
  + **10800**: every 3 hours
  + **21600**: every 6 hours
  + **43200**: every 12 hours
  + **86400**: every day

<a name="alarm_template_spec_list--alarm_template_spec_items"></a>
The `alarm_template_spec_items` block supports:

* `alarm_rule_name` - (Required, String) Specifies the alarm rule name.

* `alarm_rule_type` - (Required, String) Specifies the alarm rule type.
  Valid values are as follows:
  + **metric**: metric alarm rule
  + **event**: event alarm rule

* `alarm_rule_description` - (Optional, String) Specifies the alarm rule description.

* `event_alarm_spec` - (Optional, List) Specifies the event alarm spec.
  The [event_alarm_spec](#alarm_template_spec_list--alarm_template_spec_items--event_alarm_spec) structure is documented
  below.

* `metric_alarm_spec` - (Optional, List) Specifies the metric alarm spec.
  The [metric_alarm_spec](#alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec) structure is
  documented below.

<a name="alarm_template_spec_list--alarm_template_spec_items--event_alarm_spec"></a>
The `event_alarm_spec` block supports:

* `alarm_source` - (Optional, String)  Specifies the alarm rule source.

* `alarm_subtype` - (Optional, String) Specifies the alarm subtype.

* `event_source` - (Optional, String) Specifies the alarm source.

* `monitor_object_templates` - (Optional, List) Specifies the monitor object templates.

* `monitor_objects` - (Optional, List) Specifies the monitored objects. It's an array of map objects.
  Key-value pair, key can be as follows:
  + **event_type**: notification type
  + **event_severity**: alarm severity
  + **event_name**: event name
  + **namespace**: namespace
  + **clusterId**: cluster ID
  + **customField**: user-defined field

* `trigger_conditions` - (Optional, List) Specifies the trigger conditions.
  The [trigger_conditions](#alarm_template_spec_list--alarm_template_spec_items--event_alarm_spec--trigger_conditions)
  structure is documented below.

<a name="alarm_template_spec_list--alarm_template_spec_items--event_alarm_spec--trigger_conditions"></a>
The `trigger_conditions` block supports:

* `trigger_type` - (Required, String) Specifies the trigger mode.
  Valid values are as follows:
  + **immediately**: An alarm is triggered immediately if the alarm condition is met.
  + **accumulative**: An alarm is triggered if the alarm condition is met for a specified number of times.

* `event_name` - (Optional, String) Specifies the event name.

* `thresholds` - (Optional, Map) Specifies the thresholds. Key-value pair. The key indicates the alarm severity while
  the value indicates the number of accumulated trigger times. Leave this parameter empty if `trigger_type` is set
  to **immediately**.

* `aggregation_window` - (Optional, Int) Specifies the statistical period, in seconds. For example, 3600 indicates one
  hour. Leave this parameter empty if `trigger_type` is set to **immediately**.

* `frequency` - (Optional, String) Specifies the event alarm notification frequency. Leave this parameter empty if
  `trigger_type` is set to **immediately**. Valid values are as follows:
  + **0**: alarm sent only once
  + **300**: every 5 minutes
  + **600**: every 10 minutes
  + **900**: every 15 minutes
  + **1800**: every 30 minutes
  + **3600**: every hour
  + **10800**: every 3 hours
  + **21600**: every 6 hours
  + **43200**: every 12 hours
  + **86400**: every day

* `operator` - (Optional, String) Specifies the operator. Options: >, <, =, >=, and <=. Leave this parameter empty if
  `trigger_type` is set to **immediately**.

<a name="alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec"></a>
The `metric_alarm_spec` block supports:

* `alarm_source` - (Optional, String) Specifies the alarm source.

* `alarm_subtype` - (Optional, String) Specifies the alarm subtype.

* `alarm_tags` - (Optional, List) Specifies the alarm tags.
  The [alarm_tags](#alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec--alarm_tags) structure
  is documented below.

* `monitor_type` - (Optional, String) Specifies the monitor type.

* `no_data_conditions` - (Optional, List) Specifies the no data conditions.
  The [no_data_conditions](#alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec--no_data_conditions)
  structure is documented below.

* `recovery_conditions` - (Optional, List) Specifies the recovery conditions.
  The [recovery_conditions](#alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec--recovery_conditions)
  structure is documented below.

* `trigger_conditions` - (Optional, List) Specifies the trigger conditions.
  The [trigger_conditions](#alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec--trigger_conditions)
  structure is documented below.

<a name="alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec--alarm_tags"></a>
The `trigger_conditions` block supports:

* `auto_tags` - (Optional, List) Specifies the automatic tag.

* `custom_annotations` - (Optional, List) Specifies the custom tag.

* `custom_tags` - (Optional, List) Specifies the alarm annotation.

<a name="alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec--no_data_conditions"></a>
The `trigger_conditions` block supports:

* `no_data_alert_state` - (Optional, String) Specifies the status of the threshold rule when the data is insufficient.
  Valid values are as follows:
  + **no_data**: A notification indicating insufficient data is sent.
  + **alerting**: An alarm is triggered.
  + **ok**: No exception occurs.
  + **pre_state**: Retain the previous state.

* `no_data_timeframe` - (Optional, Int) Specifies the number of periods without data.

* `notify_no_data` - (Optional, Bool) Specifies whether to send a notification when data is insufficient.
  Defaults to **false**.

<a name="alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec--recovery_conditions"></a>
The `trigger_conditions` block supports:

* `recovery_timeframe` - (Optional, Int) Specifies the number of consecutive periods for which the trigger condition is
  not met to clear an alarm.

<a name="alarm_template_spec_list--alarm_template_spec_items--metric_alarm_spec--trigger_conditions"></a>
The `trigger_conditions` block supports:

* `metric_query_mode` - (Required, String) Specifies the metric query mode.
  Valid values are as follows:
  + **AOM**: native AOM
  + **PROM**: AOM prometheus
  + **NATIVE_PROM**: native prometheus

* `metric_name` - (Optional, String) Specifies the metric name.

* `promql` - (Optional, String) Specifies the prometheus statement.

* `aggregate_type` - (Optional, String) Specifies the aggregation mode.
  Valid values are **by**, **avg**, **max**, **min** and **sum**.

* `aggregation_type` - (Optional, String) Specifies the statistical mode.
  Valid values are **average**, **minimum**, **maximum**, **sum** and **sampleCount**.

* `aggregation_window` - (Optional, String) Specifies the statistical period.
  Valid values are **15s**, **30s**, **1m**, **5m**, **15m** and **1h**.

* `aom_monitor_level` - (Optional, String) Specifies the monitoring layer.

* `expression` - (Optional, String) Specifies the expression of a combined operation.

* `metric_labels` - (Optional, List) Specifies the metric dimension.

* `metric_namespace` - (Optional, String) Specifies the metric namespace.

* `metric_unit` - (Optional, String) Specifies the metric unit.

* `mix_promql` - (Optional, String) Specifies the promQL of a combined operation.

* `metric_statistic_method` - (Optional, String) Specifies the metric statistics method to be used when you set
  Configuration Mode to Select from all metrics during alarm rule setting.
  Valid values are as follows:
  + **single**: single metric
  + **mix**: multi-metric combined operations

* `operator` - (Optional, String) Specifies the operator. Options: >, <, =, >=, and <=.

* `promql_expr` - (Optional, List) Specifies the prometheus statement template.

* `promql_for` - (Optional, String) Specifies the native prometheus monitoring duration.

* `query_match` - (Optional, String) Specifies the query filter criteria.

* `thresholds` - (Optional, Map) Specifies the thresholds. Key-value pair. The key indicates the alarm severity while
  the value indicates the alarm threshold.

* `trigger_type` - (Optional, String) Specifies the trigger type.
  Valid values are as follows:
  + **FIXED_RATE**: fixed interval
  + **HOURLY**: every hour
  + **DAILY**: every day
  + **WEEKLY**: every week
  + **CRON**: Cron expression

* `trigger_interval` - (Optional, String) Specifies the check interval.
  Valid values are as follows:
  + If `trigger_type` is set to **HOURLY**, set this parameter to empty.
  + If `trigger_type` is set to **DAILY**, set 00:00–23:00. Example: **03:00**.
  + If `trigger_type` is set to **WEEKLY**, select a day in a week and then select 00:00–23:00.
    Example: **1 03:00** indicates 03:00 on every Monday.
  + If `trigger_type` is set to **CRON**, specify a standard cron expression.
  + If `trigger_type` is set to **FIXED_RATE**, select 15s, 30s, 1–59 min, or 1–24 h.
    Example: **15s**, **30s**, **1min**, or **1h**.

* `trigger_times` - (Optional, String) Specifies the number of consecutive periods.

<a name="templating"></a>
The `templating` block supports:

* `list` - (Required, List) Specifies the
  The [list](#templating--list) structure is documented below.

<a name="templating--list"></a>
The `list` block supports:

* `name` - (Required, String) Specifies the name.

* `description` - (Optional, String) Specifies the description.

* `query` - (Optional, String) Specifies the query.

* `type` - (Optional, String) Specifies the type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The create time of the template.

* `updated_at` - The update time of the template.

## Import

The template can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_aom_alarm_rules_template.test <id>
```
