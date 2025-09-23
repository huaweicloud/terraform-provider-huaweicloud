---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_bandwidth_policy"
description: |-
  Manages an AS bandwidth scaling policy resource within HuaweiCloud.
---

# huaweicloud_as_bandwidth_policy

Manages an AS bandwidth scaling policy resource within HuaweiCloud.

-> AS cannot scale yearly/monthly bandwidths.

## Example Usage

### AS Recurrence Policy

```hcl
variable "scaling_policy_name" {}
variable "bandwidth_id" {}

resource "huaweicloud_as_bandwidth_policy" "test" {
  scaling_policy_name = var.scaling_policy_name
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
variable "scaling_policy_name" {}
variable "bandwidth_id" {}

resource "huaweicloud_as_bandwidth_policy" "test" {
  scaling_policy_name = var.scaling_policy_name
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
variable "scaling_policy_name" {}
variable "bandwidth_id" {}
variable "alarm_id" {}

resource "huaweicloud_as_bandwidth_policy" "test" {
  scaling_policy_name = var.scaling_policy_name
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

### AS Interval Alarm Policy

```hcl
variable "scaling_policy_name" {}
variable "bandwidth_id" {}
variable "alarm_id" {}

resource "huaweicloud_as_bandwidth_policy" "test" {
  scaling_policy_name = var.scaling_policy_name
  scaling_policy_type = "INTERVAL_ALARM"
  bandwidth_id        = var.bandwidth_id
  alarm_id            = var.alarm_id

  interval_alarm_actions {
    lower_bound = "0"
    upper_bound = "5"
    operation   = "ADD"
    size        = 1
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
  + **ALARM** (corresponding to `alarm_id`): Indicates that the scaling action is triggered by an alarm.
  + **SCHEDULED** (corresponding to `scheduled_policy`): Indicates that the scaling action is triggered as scheduled.
  + **RECURRENCE** (corresponding to `scheduled_policy`): Indicates that the scaling action is triggered periodically.
  + **INTERVAL_ALARM** (corresponding to `alarm_id`): Indicates that the scaling action is triggered by an alarm.

* `bandwidth_id` - (Required, String) Specifies the scaling bandwidth ID.

* `alarm_id` - (Optional, String) Specifies the alarm rule ID.
  This parameter is mandatory when `scaling_policy_type` is set to **ALARM** or **INTERVAL_ALARM**.

* `cool_down_time` - (Optional, Int) Specifies the cooldown period (in seconds).
  The value ranges from `0` to `86,400` and is `300` by default.

* `description` - (Optional, String) Specifies the description of the AS policy.
  The value can contain `0` to `256` characters.

* `action` - (Optional, String) Specifies identification of operation the AS bandwidth policy.
  After the AS bandwidth policy created, the status is inservice, indicates the AS bandwidth policy is enabled.
  The valid values are as follows:
  + **resume**: Indicates enable the AS bandwidth policy.
  + **pause**: Indicates disable the AS bandwidth policy.

* `scaling_policy_action` - (Optional, List) Specifies the scaling action of the AS policy.
  The [scaling_policy_action](#ASBandWidthPolicy_ScalingPolicyAction) structure is documented below.

* `scheduled_policy` - (Optional, List) Specifies the periodic or scheduled AS policy.
  This parameter is mandatory when `scaling_policy_type` is set to **SCHEDULED** or **RECURRENCE**.
  The [scheduled_policy](#ASBandWidthPolicy_ScheduledPolicy) structure is documented below.

* `interval_alarm_actions` - (Optional, List) Specifies the alarm interval of the bandwidth policy.
  The [interval_alarm_actions](#bandwidth_policy_interval_alarm) structure is documented below.

  -> 1. This parameter is valid and mandatory when `scaling_policy_type` is set to **INTERVAL_ALARM**.
  <br/>2. If the alarm rule `condition.comparison_operator` is set to **=**, the `scaling_policy_type`
    not support set to **INTERVAL_ALARM**.

<a name="ASBandWidthPolicy_ScalingPolicyAction"></a>
The `scaling_policy_action` block supports:

* `operation` - (Optional, String) Specifies the operation to be performed. The default operation is **ADD**.
  The options are as follows:
  + **ADD**: Indicates adding the bandwidth size.
  + **REDUCE**: Indicates reducing the bandwidth size.
  + **SET**: Indicates setting the bandwidth size to a specified value.

* `size` - (Optional, Int) Specifies the bandwidth (Mbit/s).
  The value is an integer from `1` to `2,000`. The default value is `1`.

* `limits` - (Optional, Int) Specifies the operation restrictions.
  If `operation` is not **SET**, this parameter takes effect and the unit is Mbit/s.
  If `operation` is set to **ADD**, this parameter indicates the maximum bandwidth allowed.
  If `operation` is set to **REDUCE**, this parameter indicates the minimum bandwidth allowed.

<a name="ASBandWidthPolicy_ScheduledPolicy"></a>
The `scheduled_policy` block supports:

* `launch_time` - (Required, String) Specifies the time when the scaling action is triggered.
  The time format complies with UTC.
  If `scaling_policy_type` is set to **SCHEDULED**, the time format is **YYYY-MM-DDThh:mmZ**.
  If `scaling_policy_type` is set to **RECURRENCE**, the time format is **hh:mm**.

* `recurrence_type` - (Optional, String) Specifies the periodic triggering type.
  This parameter is mandatory when `scaling_policy_type` is set to **RECURRENCE**.
  The valid values are as follows:
  + **Daily**: Indicates that the scaling action is triggered once a day.
  + **Weekly**: Indicates that the scaling action is triggered once a week.
  + **Monthly**: Indicates that the scaling action is triggered once a month.

* `recurrence_value` - (Optional, String) Specifies the day when a periodic scaling action is triggered.
  This parameter is mandatory when `scaling_policy_type` is set to **RECURRENCE**.
  <br/>If `recurrence_type` is set to **Daily**, the value is null, indicating that the scaling action is triggered
  once a day.
  <br/>If `recurrence_type` is set to **Weekly**, the value ranges from `1` (Sunday) to `7` (Saturday).
  The digits refer to dates in each week and separated by a comma, such as **1,3,5**.
  <br/>If `recurrence_type` is set to **Monthly**, the value ranges from `1` to `31`.
  The digits refer to the dates in each month and separated by a comma, such as **1,10,13,28**.

* `start_time` - (Optional, String) Specifies the start time of the scaling action triggered periodically.
  The time format complies with UTC. The default value is the local time.
  The time format is **YYYY-MM-DDThh:mmZ**.

* `end_time` - (Optional, String) Specifies the end time of the scaling action triggered periodically.
  The time format complies with UTC. This parameter is mandatory when scaling_policy_type is set to RECURRENCE.
  When the scaling action is triggered periodically, the end time cannot be earlier than the current and start time.
  The time format is **YYYY-MM-DDThh:mmZ**.

<a name="bandwidth_policy_interval_alarm"></a>
The `interval_alarm_actions` block supports:

* `lower_bound` - (Optional, String) Specifies the lower limit of the value range.
  The value is null by default. The minimum lower limit allowed is `-1.174271E108`.

* `upper_bound` - (Optional, String) Specifies the upper limit of the value range.
  The value is null by default. The maximum upper limit allowed is `1.174271E108`.

-> 1. If the `lower_bound` is null, the `upper_bound` must be less than or equal to `0`.
<br/>2. If the `upper_bound` is null, the `lower_bound` must be greater than or equal to `0`.
<br/>3. The `lower_bound` and the `upper_bound` cannot be both `0` at the same time.
<br/>4. The `lower_bound` and `upper_bound` can not be less than `0` when the
  alarm rule `condition.comparison_operator` is set to **>** or **>=**.
<br/>5. The `lower_bound` and `upper_bound` can be less than `0` when the
  alarm rule `condition.comparison_operator` is set to **<** or **<=**.
<br/>6. If adding multiple alarm intervals, each interval value range cannot overlap.

* `operation` - (Optional, String) Specifies the operation to be performed.
  The valid values are as follows:
  + **ADD** (default): Indicates adding the bandwidth size.
  + **REDUCE**: Indicates reducing the bandwidth size.
  + **SET**: Indicates setting the bandwidth size to a specified value.

* `size` - (Optional, Int) Specifies the operation size, unit is Mbit/s.
  The valid values from `1` to `300`, the default value is `1`.

* `limits` - (Optional, Int) Specifies the operation restrictions, unit is Mbit/s.
  The valid values from `1` to `2,000`.
  If `operation` is not **SET**, this parameter takes effect.
  If `operation` is set to **ADD**, this parameter indicates the maximum bandwidth allowed.
  If `operation` is set to **REDUCE**, this parameter indicates the minimum bandwidth allowed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `scaling_resource_type` - The scaling resource type. The value is fixed to **BANDWIDTH**.

* `status` - The AS policy status. The value can be **INSERVICE**, **PAUSED** and **EXECUTING**.

* `create_time` - The creation time of the bandwidth policy, in UTC format.

* `meta_data` - The bandwidth policy additional information.
  The [meta_data](#bandwidth_policy_meta_data_struct) structure is documented below.

<a name="bandwidth_policy_meta_data_struct"></a>
The `meta_data` block supports:

* `metadata_bandwidth_share_type` - The bandwidth sharing type in the bandwidth policy.

* `metadata_eip_id` - The EIP ID for the bandwidth in the bandwidth policy.

* `metadata_eip_address` - The EIP IP address for the bandwidth in the bandwidth policy.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 2 minutes.
* `update` - Default is 2 minutes.

## Import

The bandwidth scaling policies can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_as_bandwidth_policy.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `action`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_as_bandwidth_policy" "test" {
  ...

  lifecycle {
    ignore_changes = [
      action,
    ]
  }
}
```
