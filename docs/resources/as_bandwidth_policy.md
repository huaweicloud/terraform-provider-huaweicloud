---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_bandwidth_policy"
description: ""
---

# huaweicloud_as_bandwidth_policy

Manages an AS bandwidth scaling policy resource within HuaweiCloud.

-> AS cannot scale yearly/monthly bandwidths.

## Example Usage

### AS Recurrence Policy

```hcl
variable "bandwidth_id" {}

resource "huaweicloud_as_bandwidth_policy" "bw_policy" {
  scaling_policy_name = "bw_policy"
  scaling_policy_type = "RECURRENCE"
  bandwidth_id        = var.bandwidth_id
  cool_down_time      = 600

  scaling_policy_action {
    operation = "ADD"
    size      = 1
  }
  scheduled_policy {
    launch_time      = "07:00"
    recurrence_type  = "Weekly"
    recurrence_value = "1,3,5"
    start_time       = "2022-09-30T12:00Z"
    end_time         = "2022-12-30T12:00Z"
  }
}
```

### AS Scheduled Policy

```hcl
variable "bandwidth_id" {}

resource "huaweicloud_as_bandwidth_policy" "bw_policy" {
  scaling_policy_name = "bw_policy"
  scaling_policy_type = "SCHEDULED"
  bandwidth_id        = var.bandwidth_id
  cool_down_time      = 600

  scaling_policy_action {
    operation = "ADD"
    size      = 1
  }
  scheduled_policy {
    launch_time = "2022-09-30T12:00Z"
  }
}
```

### AS Alarm Policy

```hcl
variable "bandwidth_id" {}
variable "alarm_id" {}

resource "huaweicloud_as_bandwidth_policy" "test" {
  scaling_policy_name = "bw_policy"
  scaling_policy_type = "ALARM"
  bandwidth_id        = var.bandwidth_id
  alarm_id            = var.alarm_id

  scaling_policy_action {
    operation = "ADD"
    size      = 1
    limits    = 300
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `scaling_policy_name` - (Required, String) Specifies the AS policy name.
  The name contains only letters, digits, underscores (_), and hyphens (-), and cannot exceed 64 characters.

* `scaling_policy_type` - (Required, String) Specifies the AS policy type. The options are as follows:
  - **ALARM** (corresponding to `alarm_id`): indicates that the scaling action is triggered by an alarm.
  - **SCHEDULED** (corresponding to `scheduled_policy`): indicates that the scaling action is triggered as scheduled.
  - **RECURRENCE** (corresponding to `scheduled_policy`): indicates that the scaling action is triggered periodically.

* `bandwidth_id` - (Required, String) Specifies the scaling bandwidth ID.

* `alarm_id` - (Optional, String) Specifies the alarm rule ID.
  This parameter is mandatory when `scaling_policy_type` is set to ALARM.

* `cool_down_time` - (Optional, Int) Specifies the cooldown period (in seconds).
  The value ranges from 0 to 86400 and is 300 by default.

* `description` - (Optional, String) Specifies the description of the AS policy.
  The value can contain 0 to 256 characters.

* `scaling_policy_action` - (Optional, List) Specifies the scaling action of the AS policy.
  The [object](#ASBandWidthPolicy_ScalingPolicyAction) structure is documented below.

* `scheduled_policy` - (Optional, List) Specifies the periodic or scheduled AS policy.
  This parameter is mandatory when `scaling_policy_type` is set to SCHEDULED or RECURRENCE.
  The [object](#ASBandWidthPolicy_ScheduledPolicy) structure is documented below.

<a name="ASBandWidthPolicy_ScalingPolicyAction"></a>
The `scaling_policy_action` block supports:

* `operation` - (Optional, String) Specifies the operation to be performed. The default operation is ADD.
  The options are as follows:
  - **ADD**: indicates adding the bandwidth size.
  - **REDUCE**: indicates reducing the bandwidth size.
  - **SET**: indicates setting the bandwidth size to a specified value.

* `size` - (Optional, Int) Specifies the bandwidth (Mbit/s).
  The value is an integer from 1 to 2000. The default value is 1.

* `limits` - (Optional, Int) Specifies the operation restrictions.
  - If operation is not SET, this parameter takes effect and the unit is Mbit/s.
  - If operation is set to ADD, this parameter indicates the maximum bandwidth allowed.
  - If operation is set to REDUCE, this parameter indicates the minimum bandwidth allowed.

<a name="ASBandWidthPolicy_ScheduledPolicy"></a>
The `scheduled_policy` block supports:

* `launch_time` - (Required, String) Specifies the time when the scaling action is triggered.
  The time format complies with UTC.
  - If scaling_policy_type is set to SCHEDULED, the time format is YYYY-MM-DDThh:mmZ.
  - If scaling_policy_type is set to RECURRENCE, the time format is hh:mm.

* `recurrence_type` - (Optional, String) Specifies the periodic triggering type.
  This parameter is mandatory when scaling_policy_type is set to RECURRENCE. The options are as follows:
  - **Daily**: indicates that the scaling action is triggered once a day.
  - **Weekly**: indicates that the scaling action is triggered once a week.
  - **Monthly**: indicates that the scaling action is triggered once a month.

* `recurrence_value` - (Optional, String) Specifies the day when a periodic scaling action is triggered.
  This parameter is mandatory when scaling_policy_type is set to RECURRENCE.
  - If recurrence_type is set to Daily, the value is null, indicating that the scaling action is triggered once a day.
  - If recurrence_type is set to Weekly, the value ranges from 1 (Sunday) to 7 (Saturday).
    The digits refer to dates in each week and separated by a comma, such as 1,3,5.
  - If recurrence_type is set to Monthly, the value ranges from 1 to 31.
    The digits refer to the dates in each month and separated by a comma, such as 1,10,13,28.

* `start_time` - (Optional, String) Specifies the start time of the scaling action triggered periodically.
  The time format complies with UTC. The default value is the local time.
  The time format is YYYY-MM-DDThh:mmZ.

* `end_time` - (Optional, String) Specifies the end time of the scaling action triggered periodically.
  The time format complies with UTC. This parameter is mandatory when scaling_policy_type is set to RECURRENCE.
  When the scaling action is triggered periodically, the end time cannot be earlier than the current and start time.
  The time format is YYYY-MM-DDThh:mmZ.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `scaling_resource_type` - The scaling resource type. The value is fixed to **BANDWIDTH**.

* `status` - The AS policy status. The value can be **INSERVICE**, **PAUSED** and **EXECUTING**.

## Import

The bandwidth scaling policies can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_as_bandwidth_policy.test 0ce123456a00f2591fabc00385ff1234
```
