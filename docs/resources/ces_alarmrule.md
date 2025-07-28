---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_alarmrule"
description: ""
---

# huaweicloud_ces_alarmrule

Manages a Cloud Eye alarm rule resource within HuaweiCloud.

## Example Usage

### Basic example

```hcl
variable "instance_id_1" {}
variable "instance_id_2" {}
variable "topic_urn" {}

resource "huaweicloud_ces_alarmrule" "test" {
  alarm_name           = "rule-test"
  alarm_action_enabled = true
  alarm_enabled        = true
  alarm_type           = "MULTI_INSTANCE"

  metric {
    namespace = "SYS.ECS"
  }

  resources {
    dimensions {
      name  = "instance_id"
      value = var.instance_id_1
    }
  }

  resources {
    dimensions {
      name  = "instance_id"
      value = var.instance_id_2
    }
  }

  condition  {
    period              = 1200
    filter              = "average"
    comparison_operator = ">"
    value               = 6.5
    unit                = "B/s"
    count               = 1
    suppress_duration   = 300
    metric_name         = "network_outgoing_bytes_rate_inband"
    alarm_level         = 4
  }

  condition  {
    period              = 3600
    filter              = "average"
    comparison_operator = ">="
    value               = 20
    unit                = "B/s"
    count               = 1
    suppress_duration   = 300
    metric_name         = "network_outgoing_bytes_rate_inband"
    alarm_level         = 4
  }

  alarm_actions {
    type              = "notification"
    notification_list = [
      var.topic_urn
    ]
  }
}
```

### Alarm rule for All instance

```hcl
variable "topic_urn" {}

resource "huaweicloud_ces_alarmrule" "test" {
  alarm_name           = "rule-test"
  alarm_action_enabled = true
  alarm_enabled        = true
  alarm_type           = "ALL_INSTANCE"

  metric {
    namespace = "AGT.ECS"
  }

  resources {
    dimensions {
      name = "instance_id"
    }

    dimensions {
      name = "mount_point"
    }
  }

  condition  {
    alarm_level         = 2
    suppress_duration   = 0
    period              = 1
    filter              = "average"
    comparison_operator = ">"
    value               = 80
    count               = 1
    metric_name         = "disk_usedPercent"
  }

  alarm_actions {
    type              = "notification"
    notification_list = [
      var.topic_urn
    ]
  }
}
```

### Alarm rule for event monitoring

```hcl
variable "topic_urn" {}

resource "huaweicloud_ces_alarmrule" "test" {
  alarm_name           = "rule-test"
  alarm_action_enabled = true
  alarm_type           = "EVENT.SYS"

  metric {
    namespace = "SYS.ECS"
  }
  
  condition  {
    metric_name         = "stopServer"
    period              = 0
    filter              = "average"
    comparison_operator = ">="
    value               = 1
    unit                = "count"
    count               = 1
    suppress_duration   = 0
    alarm_level         = 2
  }

  alarm_actions {
    type              = "notification"
    notification_list = [
      var.topic_urn
    ]
  }
}
```

### Alarm rule using the alarm template

```hcl
variable "topic_urn" {}
variable "alarm_template_id" {}
variable "instance_id" {}

resource "huaweicloud_ces_alarmrule" "test" {
  alarm_name           = "rule-test"
  alarm_enabled        = true
  alarm_action_enabled = true
  alarm_type           = "MULTI_INSTANCE"
  alarm_template_id    = var.alarm_template_id

  metric {
    namespace = "SYS.ECS"
  }

  resources {
    dimensions {
      name  = "instance_id"
      value = var.instance_id
    }
  }

  alarm_actions {
    type              = "notification"
    notification_list = [
      var.topic_urn
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the alarm rule resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `alarm_name` - (Required, String) Specifies the name of an alarm rule. The value can be a string of `1` to `128`
  characters that can consist of English letters, Chinese characters, digits, underscores (_), hyphens (-).

* `metric` - (Required, List, ForceNew) Specifies the alarm metrics. The structure is described below. Changing this
  creates a new resource.

* `alarm_template_id` - (Optional, String, ForceNew) Specifies the ID of the alarm template.
  When using `alarm_template_id`, the fields `alarm_name`, `alarm_description`, `alarm_action_enabled`, `alarm_actions`
  and `ok_actions` cannot be updated.
  Changing this creates a new resource.

* `condition` - (Optional, List) Specifies the alarm triggering condition.
  The [condition](#Condition) structure is documented below.

* `resources` - (Optional, List) Specifies the list of the resources to add into the alarm rule.
  The structure is described below.

* `alarm_description` - (Optional, String) The value can be a string of 0 to 256 characters.

* `alarm_enabled` - (Optional, Bool) Specifies whether to enable the alarm. The default value is true.

* `alarm_type` - (Optional, String) Specifies the alarm type. The value can be **EVENT.SYS**, **EVENT.CUSTOM**,
  **MULTI_INSTANCE** and **ALL_INSTANCE**. Defaults to **MULTI_INSTANCE**.

* `alarm_actions` - (Optional, List) Specifies the action triggered by an alarm. The structure is described
  below.

* `ok_actions` - (Optional, List) Specifies the action triggered by the clearing of an alarm. The structure is
  described below.

* `alarm_action_enabled` - (Optional, Bool) Specifies whether to enable the action to be triggered by an alarm. The
  default value is true.

* `notification_begin_time` - (Optional, String) Specifies the alarm notification start time, for example: **05:30**.

* `notification_end_time` - (Optional, String) Specifies the alarm notification stop time, for example: **22:10**.

* `effective_timezone` - (Optional, String) Specifies the time zone, for example: **GMT-08:00**, **GMT+08:00** or
  **GMT+0:00**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the alarm rule.

-> **Note** If alarm_action_enabled is set to true, either alarm_actions or ok_actions cannot be empty. If alarm_actions
and ok_actions coexist, their corresponding notification_list must be of the **same value**.

The `metric` block supports:

* `namespace` - (Required, String, ForceNew) Specifies the namespace in **service.item** format. **service** and **item**
  each must be a string that starts with a letter and contains only letters, digits, and underscores (_).
  Changing this creates a new resource.
  For details, see [Services Interconnected with Cloud Eye](https://support.huaweicloud.com/intl/en-us/api-ces/ces_03_0059.html).

The `resources` block supports:

* `dimensions` - (Optional, List) Specifies the list of metric dimensions. The structure is described below.

The `dimensions` block supports:

* `name` - (Required, String) Specifies the dimension name. The value can be a string of 1 to 32 characters
  that must start with a letter and contain only letters, digits, underscores (_), and hyphens (-).

* `value` - (Optional, String) Specifies the dimension value. The value can be a string of 1 to 64 characters
  that must start with a letter or a number and contain only letters, digits, underscores (_), and hyphens (-).

<a name="Condition"></a>
The `condition` block supports:

* `period` - (Required, Int) Specifies the alarm checking period in seconds. The value can be 0, 1, 300, 1200, 3600, 14400,
  and 86400.

  Note: If period is set to 1, the raw metric data is used to determine whether to generate an alarm. When the value of
  `alarm_type` is **EVENT.SYS** or **EVENT.CUSTOM**, period can be set to 0.

* `filter` - (Required, String) Specifies the data rollup methods. The value can be max, min, average, sum, and variance.

* `comparison_operator` - (Required, String) Specifies the comparison condition of alarm thresholds. The value can be >,
  =, <, >=, or <=.

* `value` - (Required, Float) Specifies the alarm threshold. The value ranges from 0 to Number of
  1.7976931348623157e+108.

* `count` - (Required, Int) Specifies the number of consecutive occurrence times. The value ranges from 1 to 5.

* `unit` - (Optional, String) Specifies the data unit.
  For details, see [Services Interconnected with Cloud Eye](https://support.huaweicloud.com/intl/en-us/api-ces/ces_03_0059.html).

* `suppress_duration` - (Optional, Int) Specifies the interval for triggering an alarm if the alarm persists.
  Possible values are as follows:
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

  The default value is `0`.

* `metric_name` - (Required, String) Specifies the metric name of the condition. The value can be a string of
  1 to 64 characters that must start with a letter and contain only letters, digits, and underscores (_).
  For details, see [Services Interconnected with Cloud Eye](https://support.huaweicloud.com/intl/en-us/api-ces/ces_03_0059.html).

* `alarm_level` - (Optional, Int) Specifies the alarm severity of the condition. The value can be 1, 2, 3 or 4,
  which indicates *critical*, *major*, *minor*, and *informational*, respectively.
  The default value is 2.

the `alarm_actions` block supports:

* `type` - (Required, String) Specifies the type of action triggered by an alarm. the
  value can be *notification* or *autoscaling*.
    + notification: indicates that a notification will be sent to the user.
    + autoscaling: indicates that a scaling action will be triggered.

* `notification_list` - (Required, List) specifies the list of objects to be notified if the alarm status changes, the
  maximum length is 5. If `type` is set to *notification*, the value of notification_list cannot be empty. If `type` is
  set to *autoscaling*, the value of notification_list must be **[]**
  and the value of namespace must be *SYS.AS*.

  Note: to enable the *autoscaling* alarm rules take effect, you must bind scaling policies.

the `ok_actions` block supports:

* `type` - (Required, String) Specifies the type of action triggered by an alarm. the value is notification.
  notification: indicates that a notification will be sent to the user.

* `notification_list` - (Required, List) specifies the list of objects to be notified if the alarm status changes, the
  maximum length is 5.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the alarm rule ID.

* `alarm_state` - Indicates the alarm status. The value can be:
  + ok: The alarm status is normal;
  + alarm: An alarm is generated;
  + insufficient_data: The required data is insufficient.

* `update_time` - Indicates the time when the alarm status changed. The value is a UNIX timestamp and the unit is ms.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 5 minutes.

## Import

CES alarm rules can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ces_alarmrule.alarm_rule al1619578509719Ga0X1RGWv
```
