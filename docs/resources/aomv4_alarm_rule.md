---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aomv4_alarm_rule"
description: |-
  Manages an AOM cloud alarm rule resource within HuaweiCloud.
---

# huaweicloud_aomv4_alarm_rule

Manages an AOM cloud alarm rule resource within HuaweiCloud.

## Example Usage

### Create an event alarm rule with all system event

```hcl
variable "alarm_rule_name" {}
variable "action_rule_id" {}

resource "huaweicloud_aomv4_alarm_rule" "test" {
  name        = var.alarm_rule_name
  type        = "event"
  description = "test"
  enable      = true

  alarm_notifications {
    notification_type         = "direct"
    notification_enable       = true
    bind_notification_rule_id = var.action_rule_id
  }

  event_alarm_spec {
    event_source = "CCE"
    alarm_source = "systemEvent"

    trigger_conditions {
      trigger_type = "immediately"

      thresholds = {
        "Critical" = 2
      }
    }
  }
}
```

### Create an event alarm rule with specific system event

```hcl
variable "alarm_rule_name" {}
variable "action_rule_id" {}

resource "huaweicloud_aomv4_alarm_rule" "test" {
  name        = var.alarm_rule_name
  type        = "event"
  description = "test"
  enable      = true

  alarm_notifications {
    notification_type         = "direct"
    notification_enable       = true
    bind_notification_rule_id = var.action_rule_id
  }

  event_alarm_spec {
    event_source = "CCE"
    alarm_source = "systemEvent"

    monitor_objects = [
      {
        event_name = "扩容节点超时##ScaleUpTimedOut;数据卷扩容失败##VolumeResizeFailed"
      },
    ]

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
```

### Create a metric alarm rule with metrics

```hcl
variable "alarm_rule_name" {}

resource "huaweicloud_aomv4_alarm_rule" "test" {
  name             = var.alarm_rule_name
  type             = "metric"
  prom_instance_id = "0"
  enable           = true

  alarm_notifications {
    notification_type = "direct"
  }

  metric_alarm_spec {
    monitor_type = "all_metric"

    recovery_conditions {
      recovery_timeframe = 1
    }

    trigger_conditions {
      metric_query_mode       = "PROM"
      metric_statistic_method = "single"
      metric_name             = "duration"
      promql                  = "label_replace(avg_over_time(duration{}[59999ms]),\"__name__\",\"duration\",\"\",\"\")"
      trigger_times           = "3"
      trigger_type            = "FIXED_RATE"
      aggregation_window      = "1m"
      trigger_interval        = "30s"
      aggregation_type        = "average"
      operator                = ">"

      thresholds = {
        "Critical" = "1"
      }
    }

    alarm_tags {
      custom_tags = ["key=value"]
    }

    no_data_conditions {
      notify_no_data      = true
      no_data_timeframe   = 1
      no_data_alert_state = "no_data"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the alarm rule name.
  Changing this creates a new resource.

* `type` - (Required, String, ForceNew) Specifies the alarm rule type.
  Valid values are as follows:
  + **metric**: metric alarm rule
  + **event**: event alarm rule

  Changing this creates a new resource.

* `alarm_notifications` - (Required, List) Specifies the alarm notification module.
  The [alarm_notifications](#block--alarm_notifications) structure is documented below.

* `description` - (Optional, String) Specifies the alarm rule description.

* `enable` - (Optional, Bool) Specifies whether to enable the alarm rule. Defaults to **false**.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the alarm rule belongs.
  Changing this creates a new resource.

* `event_alarm_spec` - (Optional, List) Specifies the structure of an event alarm rule.
  It's required if `type` is **event**.
  The [event_alarm_spec](#block--event_alarm_spec) structure is documented below.

* `metric_alarm_spec` - (Optional, List) Specifies the structure of a metric alarm rule.
  It's required if `type` is **metric**.
  The [metric_alarm_spec](#block--metric_alarm_spec) structure is documented below.

* `prom_instance_id` - (Optional, String) Specifies the prometheus instance ID to which the metric alarm rule belongs.
  It's required if `type` is **metric**.

<a name="block--alarm_notifications"></a>
The `alarm_notifications` block supports:

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

<a name="block--event_alarm_spec"></a>
The `event_alarm_spec` block supports:

* `alarm_source` - (Required, String, ForceNew) Specifies the alarm rule source.
  Valid values are **systemEvent** and **customEvent**.
  Changing this creates a new resource.

* `event_source` - (Required, String) Specifies the alarm source.
  Valid values are **RDS**, **EVS**, **CCE**, **LTS** and **AOM**.

* `trigger_conditions` - (Optional, List) Specifies the trigger conditions.
  The [trigger_conditions](#block--event_alarm_spec--trigger_conditions) structure is documented below.

* `monitor_objects` - (Optional, List) Specifies the monitored objects. It's an array of map objects.
  Key-value pair, key can be as follows:
  + **event_type**: notification type
  + **event_severity**: alarm severity
  + **event_name**: event name
  + **namespace**: namespace
  + **clusterId**: cluster ID
  + **customField**: user-defined field

<a name="block--event_alarm_spec--trigger_conditions"></a>
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

<a name="block--metric_alarm_spec"></a>
The `metric_alarm_spec` block supports:

* `monitor_type` - (Required, String, ForceNew) Specifies the monitoring type.
  Valid values are as follows:
  + **all_metric**: Select metrics from all metrics.
  + **promql**: Select metrics using PromQL.

  Changing this creates a new resource.

* `recovery_conditions` - (Required, List) Specifies the alarm clearance condition.
  The [recovery_conditions](#block--metric_alarm_spec--recovery_conditions) structure is documented below.

* `trigger_conditions` - (Required, List) Specifies the trigger conditions.
  The [trigger_conditions](#block--metric_alarm_spec--trigger_conditions) structure is documented below.

* `alarm_tags` - (Optional, List) Specifies the alarm tags.
  The [alarm_tags](#block--metric_alarm_spec--alarm_tags) structure is documented below.

* `no_data_conditions` - (Optional, List) Specifies the action taken for insufficient data.
  The [no_data_conditions](#block--metric_alarm_spec--no_data_conditions) structure is documented below.

* `monitor_objects` - (Optional, List) Specifies the monitored objects. It's an array of Map objects.

<a name="block--metric_alarm_spec--recovery_conditions"></a>
The `recovery_conditions` block supports:

* `recovery_timeframe` - (Optional, Int) Specifies the number of consecutive periods for which the trigger condition is
  not met to clear an alarm.

<a name="block--metric_alarm_spec--trigger_conditions"></a>
The `trigger_conditions` block supports:

* `metric_name` - (Required, String) Specifies the metric name.

* `metric_query_mode` - (Required, String) Specifies the metric query mode.
  Valid values are as follows:
  + **AOM**: native AOM
  + **PROM**: AOM prometheus
  + **NATIVE_PROM**: native prometheus

* `promql` - (Required, String) Specifies the prometheus statement.

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

* `promql_expr` - (Optional, String) Specifies the prometheus statement template.

* `promql_for` - (Optional, String) Specifies the native prometheus monitoring duration.

* `query_param` - (Optional, String) Specifies the query parameters.

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

<a name="block--metric_alarm_spec--alarm_tags"></a>
The `alarm_tags` block supports:

* `auto_tags` - (Optional, List) Specifies the automatic tag.

* `custom_annotations` - (Optional, List) Specifies the custom tag.

* `custom_tags` - (Optional, List) Specifies the alarm annotation.

<a name="block--metric_alarm_spec--no_data_conditions"></a>
The `no_data_conditions` block supports:

* `no_data_alert_state` - (Optional, String) Specifies the status of the threshold rule when the data is insufficient.
  Valid values are as follows:
  + **no_data**: A notification indicating insufficient data is sent.
  + **alerting**: An alarm is triggered.
  + **ok**: No exception occurs.
  + **pre_state**: Retain the previous state.

* `no_data_timeframe` - (Optional, Int) Specifies the number of periods without data.

* `notify_no_data` - (Optional, Bool) Specifies whether to send a notification when data is insufficient.
  Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is same as the `name`.

* `alarm_rule_id` - Indicates the alarm rule ID.

* `status` - Indicates the alarm status.
  Value can be as follows:
  + **OK**: normal
  + **alarm**: threshold-crossing
  + **Effective**: in use
  + **Invalid**: not in use

* `created_at` - Indicates the time when an alarm rule was created.

* `updated_at` - Indicates the time when an alarm rule was updated.

## Import

Alarm rule resources can be imported using `name`, e.g.

```bash
$ terraform import huaweicloud_aomv4_alarm_rule.test <name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `metric_alarm_spec.0.trigger_conditions`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the rule, or the resource definition should be updated to
align with the rule. Also you can ignore changes as below.

```hcl
resource "huaweicloud_aomv4_alarm_rule" "test" {
    ...

  lifecycle {
    ignore_changes = [
      metric_alarm_spec.0.trigger_conditions,
    ]
  }
}
```
