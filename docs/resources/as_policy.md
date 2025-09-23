---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_policy"
description: |-
  Manages an AS policy resource within HuaweiCloud.
---

# huaweicloud_as_policy

Manages an AS policy resource within HuaweiCloud.

## Example Usage

### AS Recurrence Policy

```hcl
variable "as_group_id" {}

resource "huaweicloud_as_policy" "my_aspolicy" {
  scaling_policy_name = "my_aspolicy"
  scaling_policy_type = "RECURRENCE"
  scaling_group_id    = var.as_group_id

  scaling_policy_action {
    operation       = "ADD"
    instance_number = 1
  }
  scheduled_policy {
    launch_time     = "07:00"
    recurrence_type = "Daily"
    start_time      = "2022-11-30T12:00Z"
    end_time        = "2022-12-30T12:00Z"
  }
}
```

### AS Scheduled Policy

```hcl
variable "as_group_id" {}

resource "huaweicloud_as_policy" "my_aspolicy_1" {
  scaling_policy_name = "my_aspolicy_1"
  scaling_policy_type = "SCHEDULED"
  scaling_group_id    = var.as_group_id

  scaling_policy_action {
    operation       = "REMOVE"
    instance_number = 1
  }
  scheduled_policy {
    launch_time = "2022-12-22T12:00Z"
  }
}
```

### AS Alarm Policy

```hcl
variable "as_group_id" {}

resource "huaweicloud_ces_alarmrule" "alarm_rule" {
  alarm_name = "as_alarm_rule"

  metric {
    namespace   = "SYS.AS"
    metric_name = "cpu_util"
    dimensions {
      name  = "AutoScalingGroup"
      value = var.as_group_id
    }
  }
  condition {
    period              = 300
    filter              = "average"
    comparison_operator = ">="
    value               = 60
    unit                = "%"
    count               = 1
  }
  alarm_actions {
    type              = "autoscaling"
    notification_list = []
  }
}

resource "huaweicloud_as_policy" "my_aspolicy_2" {
  scaling_policy_name = "my_aspolicy_2"
  scaling_policy_type = "ALARM"
  scaling_group_id    = var.as_group_id
  alarm_id            = huaweicloud_ces_alarmrule.alarm_rule.id
  cool_down_time      = 900

  scaling_policy_action {
    operation       = "ADD"
    instance_number = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the AS policy. If omitted, the
  provider-level region will be used. Changing this creates a new AS policy.

* `scaling_policy_name` - (Required, String) Specifies the name of the AS policy. The name contains only letters, digits,
  underscores(_), and hyphens(-), and cannot exceed 64 characters.

* `scaling_group_id` - (Required, String, ForceNew) Specifies the AS group ID. Changing this creates a new AS policy.

* `scaling_policy_type` - (Required, String) Specifies the AS policy type. The value can be `ALARM`, `SCHEDULED` or `RECURRENCE`.
  + **ALARM**: indicates that the scaling action is triggered by an alarm.
  + **SCHEDULED**: indicates that the scaling action is triggered as scheduled.
  + **RECURRENCE**: indicates that the scaling action is triggered periodically.

* `action` - (Optional, String) Specifies the operation for the AS policy.
  The default value is **resume**. The valid values are as follows:
  + **resume**: Enables the AS policy.
  + **pause**: Disables the AS policy.

* `alarm_id` - (Optional, String) Specifies the alarm rule ID. This parameter is mandatory when `scaling_policy_type`
  is set to `ALARM`. You can create an alarm rule with
  [huaweicloud_ces_alarmrule](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs/resources/ces_alarmrule).

* `scheduled_policy` - (Optional, List) Specifies the periodic or scheduled AS policy.
  This parameter is mandatory when `scaling_policy_type` is set to `SCHEDULED` or `RECURRENCE`.
  The [object](#scheduled_policy_object) structure is documented below.

* `scaling_policy_action` - (Optional, List) Specifies the action of the AS policy.
  The [object](#scaling_policy_action_object) structure is documented below.

* `cool_down_time` - (Optional, Int) Specifies the cooling duration (in seconds).
  The value ranges from 0 to 86400 and is 300 by default.

<a name="scheduled_policy_object"></a>
The `scheduled_policy` block supports:

* `launch_time` - (Required, String) Specifies the time when the scaling action is triggered.
  + If `scaling_policy_type` is set to `SCHEDULED`, the time format is **YYYY-MM-DDThh:mmZ**.
  + If `scaling_policy_type` is set to `RECURRENCE`, the time format is **hh:mm**.

  -> the `launch_time` of the `SCHEDULED` policy cannot be earlier than the current time.

* `recurrence_type` - (Optional, String) Specifies the periodic triggering type. This argument is mandatory when
  `scaling_policy_type` is set to `RECURRENCE`. The options include `Daily`, `Weekly`, and `Monthly`.

* `recurrence_value` - (Optional, String) Specifies the frequency at which scaling actions are triggered.

* `start_time` - (Optional, String) Specifies the start time of the scaling action triggered periodically. The time format
  complies with UTC. The current time is used by default. The time format is YYYY-MM-DDThh:mmZ.

* `end_time` - (Optional, String) Specifies the end time of the scaling action triggered periodically. The time format complies
  with UTC. This argument is mandatory when `scaling_policy_type`
  is set to `RECURRENCE`. The time format is YYYY-MM-DDThh:mmZ.

<a name="scaling_policy_action_object"></a>
The `scaling_policy_action` block supports:

* `operation` - (Optional, String) Specifies the operation to be performed. The options include `ADD` (default), `REMOVE`,
  and `SET`.

* `instance_number` - (Optional, Int) Specifies the number of instances to be operated.

* `instance_percentage` - (Optional, Int) Specifies the percentage of instances to be operated.

-> At most one of `instance_number` and `instance_percentage` can be set. When neither `instance_number` nor
  `instance_percentage` is specified, the number of operation instances is **1**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `status` - The AS policy status. The value can be *INSERVICE*, *PAUSED* or *EXECUTING*.

* `create_time` - The creation time of the AS policy, in UTC format.

## Import

AS policies can be imported by their `id`, e.g.

```bash
$ terraform import huaweicloud_as_policy.test <id>
```
