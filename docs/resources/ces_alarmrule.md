---
subcategory: "Cloud Eye"
---

# huaweicloud_ces_alarmrule

Manages a Cloud Eye alarm rule resource within HuaweiCloud.

## Example Usage

### Basic example

```hcl
resource "huaweicloud_ces_alarmrule" "alarm_rule" {
  alarm_name = "alarm_rule"

  metric {
    namespace   = "SYS.ECS"
    metric_name = "network_outgoing_bytes_rate_inband"
    dimensions {
      name  = "instance_id"
      value = var.webserver_instance_id
    }
  }
  condition {
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 6
    unit                = "B/s"
    count               = 1
  }
  alarm_actions {
    type              = "notification"
    notification_list = [var.smn_topic_id]
  }
}
```

## Alarm rule for event monitoring

```hcl
resource "huaweicloud_ces_alarmrule" "alarm_rule" {
  alarm_name           = "alarm_rule"
  alarm_action_enabled = true
  alarm_type           = "EVENT.SYS"

  metric {
    namespace   = "SYS.ECS"
    metric_name = "stopServer"
  }
  
  condition  {
    period              = 0
    filter              = "average"
    comparison_operator = ">="
    value               = 1
    unit                = "count"
    count               = 1
  }

  alarm_actions {
    type              = "notification"
    notification_list = [var.smn_topic_id]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the alarm rule resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `alarm_name` - (Required, String) Specifies the name of an alarm rule. The value can be a string of 1 to 128
  characters that can consist of letters, digits, underscores (_), hyphens (-) and chinese characters.

* `metric` - (Required, List, ForceNew) Specifies the alarm metrics. The structure is described below. Changing this
  creates a new resource.

* `condition` - (Required, List) Specifies the alarm triggering condition. The structure is described below.

* `alarm_description` - (Optional, String) The value can be a string of 0 to 256 characters.

* `alarm_enabled` - (Optional, Bool) Specifies whether to enable the alarm. The default value is true.

* `alarm_level` - (Optional, Int) Specifies the alarm severity. The value can be 1, 2, 3 or 4,
  which indicates *critical*, *major*, *minor*, and *informational*, respectively.
  The default value is 2.

* `alarm_type` - (Optional, String) Specifies the alarm type. The value can be **EVENT.SYS**, **EVENT.CUSTOM**
  or **MULTI_INSTANCE**. Defaults to **MULTI_INSTANCE**.

* `alarm_actions` - (Optional, List, ForceNew) Specifies the action triggered by an alarm. The structure is described
  below. Changing this creates a new resource.

* `ok_actions` - (Optional, List, ForceNew) Specifies the action triggered by the clearing of an alarm. The structure is
  described below. Changing this creates a new resource.

* `alarm_action_enabled` - (Optional, Bool) Specifies whether to enable the action to be triggered by an alarm. The
  default value is true.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id of the alarm rule. Changing
  this creates a new resource.

-> **Note** If alarm_action_enabled is set to true, either alarm_actions or ok_actions cannot be empty. If alarm_actions
and ok_actions coexist, their corresponding notification_list must be of the **same value**.

The `metric` block supports:

* `namespace` - (Required, String, ForceNew) Specifies the namespace in **service.item** format. **service** and **item**
  each must be a string that starts with a letter and contains only letters, digits, and underscores (_).
  Changing this creates a new resource.
  For details, see [Services Interconnected with Cloud Eye](https://support.huaweicloud.com/intl/en-us/api-ces/ces_03_0059.html).

* `metric_name` - (Required, String, ForceNew) Specifies the metric name. The value can be a string of 1 to 64
  characters that must start with a letter and contain only letters, digits, and underscores (_).
  Changing this creates a new resource.
  For details, see [Services Interconnected with Cloud Eye](https://support.huaweicloud.com/intl/en-us/api-ces/ces_03_0059.html).

* `dimensions` - (Optional, List, ForceNew) Specifies the list of metric dimensions. The structure is described below.
  Changing this creates a new resource.

The `dimensions` block supports:

* `name` - (Required, String, ForceNew) Specifies the dimension name. The value can be a string of 1 to 32 characters
  that must start with a letter and contain only letters, digits, underscores (_), and hyphens (-).
  Changing this creates a new resource.

* `value` - (Required, String, ForceNew) Specifies the dimension value. The value can be a string of 1 to 64 characters
  that must start with a letter or a number and contain only letters, digits, underscores (_), and hyphens (-).
  Changing this creates a new resource.

The `condition` block supports:

* `period` - (Required, Int) Specifies the alarm checking period in seconds. The value can be 0, 1, 300, 1200, 3600, 14400,
  and 86400.

  Note: If period is set to 1, the raw metric data is used to determine whether to generate an alarm. When the value of
  `alarm_type` is **EVENT.SYS** or **EVENT.CUSTOM**, period can be set to 0.

* `filter` - (Required, String) Specifies the data rollup methods. The value can be max, min, average, sum, and vaiance.

* `comparison_operator` - (Required, String) Specifies the comparison condition of alarm thresholds. The value can be >,
  =, <, >=, or <=.

* `value` - (Required, Float) Specifies the alarm threshold. The value ranges from 0 to Number of
  1.7976931348623157e+108.

* `count` - (Required, Int) Specifies the number of consecutive occurrence times. The value ranges from 1 to 5.

* `unit` - (Optional, String, ForceNew) Specifies the data unit. Changing this creates a new resource.
  For details, see [Services Interconnected with Cloud Eye](https://support.huaweicloud.com/intl/en-us/api-ces/ces_03_0059.html).

* `suppress_duration` - (Optional, Int, ForceNew) Specifies the interval for triggering an alarm if the alarm persists.
  Changing this creates a new resource. Possible values are as follows:
  + **0**: Cloud Eye triggers the alarm only once;
  + **300**: Cloud Eye triggers the alarm every 5 minutes;
  + **600**: Cloud Eye triggers the alarm every 10 minutes;
  + **900**: Cloud Eye triggers the alarm every 15 minutes;
  + **1800**: Cloud Eye triggers the alarm every 30 minutes;
  + **3600**: Cloud Eye triggers the alarm every hour;
  + **10800**: Cloud Eye triggers the alarm every 3 hours;
  + **21600**: Cloud Eye triggers the alarm every 6 hours;
  + **43200**: Cloud Eye triggers the alarm every 12 hour;
  + **86400**: Cloud Eye triggers the alarm every day.

  The default value is **0**.

the `alarm_actions` block supports:

* `type` - (Optional, String) specifies the type of action triggered by an alarm. the
  value can be *notification* or *autoscaling*.
    + notification: indicates that a notification will be sent to the user.
    + autoscaling: indicates that a scaling action will be triggered.

* `notification_list` - (Required, List) specifies the list of objects to be notified if the alarm status changes, the
  maximum length is 5. If `type` is set to *notification*, the value of notification_list cannot be empty. If `type` is
  set to *autoscaling*, the value of notification_list must be **[]**
  and the value of namespace must be *SYS.AS*.

  Note: to enable the *autoscaling* alarm rules take effect, you must bind scaling policies.

the `ok_actions` block supports:

* `type` - (Optional, String) specifies the type of action triggered by an alarm. the value is notification.
  notification: indicates that a notification will be sent to the user.

* `notification_list` - (Required, List) specifies the list of objects to be notified if the alarm status changes, the
  maximum length is 5.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the alarm rule ID.

* `alarm_state` - Indicates the alarm status. The value can be:
  + ok: The alarm status is normal;
  + alarm: An alarm is generated;
  + insufficient_data: The required data is insufficient.

* `update_time` - Indicates the time when the alarm status changed. The value is a UNIX timestamp and the unit is ms.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `update` - Default is 10 minute.
* `delete` - Default is 5 minute.

## Import

CES alarm rules can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_ces_alarmrule.alarm_rule al1619578509719Ga0X1RGWv
```
