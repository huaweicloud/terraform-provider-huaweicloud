---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_policy_v1"
sidebar_current: "docs-huaweicloud-resource-as-policy-v1"
description: |-
  Manages a V1 AS Policy resource within HuaweiCloud.
---

# huaweicloud\_as\_policy_v1

Manages a V1 AS Policy resource within HuaweiCloud.

## Example Usage

### AS Recurrence Policy

```hcl
resource "huaweicloud_as_policy_v1" "hth_aspolicy"{
  scaling_policy_name = "hth_aspolicy"
  scaling_group_id = "4579f2f5-cbe8-425a-8f32-53dcb9d9053a"
  cool_down_time = 900
  scaling_policy_type = "RECURRENCE"
  scaling_policy_action = {
    operation = "ADD"
    instance_number = 1
  }
  scheduled_policy = {
    launch_time = "07:00"
    recurrence_type = "Daily"
    start_time = "2017-11-30T12:00Z"
    end_time = "2017-12-30T12:00Z"
  }
}

```

### AS Scheduled Policy

```hcl
resource "huaweicloud_as_policy_v1" "hth_aspolicy_1"{
  scaling_policy_name = "hth_aspolicy_1"
  scaling_group_id = "4579f2f5-cbe8-425a-8f32-53dcb9d9053a"
  cool_down_time = 900
  scaling_policy_type = "SCHEDULED"
  scaling_policy_action = {
    operation = "REMOVE"
    instance_number = 1
  }
  scheduled_policy = {
    launch_time = "2017-12-22T12:00Z"
  }
}

```

Please note that the `launch_time` of the `SCHEDULED` policy cannot be earlier than the current time.

### AS Alarm Policy

```hcl
resource "huaweicloud_as_policy_v1" "hth_aspolicy_2"{
  scaling_policy_name = "hth_aspolicy_2"
  scaling_group_id = "4579f2f5-cbe8-425a-8f32-53dcb9d9053a"
  cool_down_time = 900
  scaling_policy_type = "ALARM"
  alarm_id = "37e310f5-db9d-446e-9135-c625f9c2bbfc"
  scaling_policy_action = {
    operation = "ADD"
    instance_number = 1
  }
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the AS policy. If
    omitted, the `region` argument of the provider is used. Changing this
    creates a new AS policy.

* `scaling_policy_name` - (Required) The name of the AS policy. The name can contain letters,
    digits, underscores(_), and hyphens(-),and cannot exceed 64 characters.

* `scaling_group_id` - (Required) The AS group ID. Changing this creates a new AS policy.

* `scaling_policy_type` - (Required) The AS policy type. The values can be `ALARM`, `SCHEDULED`,
    and `RECURRENCE`.

* `alarm_id` - (Optional) The alarm rule ID. This argument is mandatory
    when `scaling_policy_type` is set to `ALARM`.

* `scheduled_policy` - (Optional) The periodic or scheduled AS policy. This argument is mandatory
    when `scaling_policy_type` is set to `SCHEDULED` or `RECURRENCE`. The scheduled_policy structure
    is documented below.

* `scaling_policy_action` - (Optional) The action of the AS policy. The scaling_policy_action
    structure is documented below.

* `cool_down_time` - (Optional) The cooling duration (in seconds), and is 900 by default.

The `scheduled_policy` block supports:

* `launch_time` - (Required) The time when the scaling action is triggered. If `scaling_policy_type`
    is set to `SCHEDULED`, the time format is YYYY-MM-DDThh:mmZ. If `scaling_policy_type` is set to
    `RECURRENCE`, the time format is hh:mm.

* `recurrence_type` - (Optional) The periodic triggering type. This argument is mandatory when
    `scaling_policy_type` is set to `RECURRENCE`. The options include `Daily`, `Weekly`, and `Monthly`.

* `recurrence_value` - (Optional) The frequency at which scaling actions are triggered.

* `start_time` - (Optional) The start time of the scaling action triggered periodically.
    The time format complies with UTC. The current time is used by default. The time
    format is YYYY-MM-DDThh:mmZ.

* `end_time` - (Optional) The end time of the scaling action triggered periodically.
    The time format complies with UTC. This argument is mandatory when `scaling_policy_type`
    is set to `RECURRENCE`. The time format is YYYY-MM-DDThh:mmZ.

The `scaling_policy_action` block supports:

* `operation` - (Optional) The operation to be performed. The options include `ADD` (default), `REMOVE`,
    and `SET`.

* `instance_number` - (Optional) The number of instances to be operated. The default number is 1.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `scaling_policy_name` - See Argument Reference above.
* `scaling_policy_type` - See Argument Reference above.
* `alarm_id` - See Argument Reference above.
* `cool_down_time` - See Argument Reference above.
* `scaling_policy_action/operation` - See Argument Reference above.
* `scaling_policy_action/instance_number` - See Argument Reference above.
* `scheduled_policy/launch_time` - See Argument Reference above.
* `scheduled_policy/recurrence_type` - See Argument Reference above.
* `scheduled_policy/recurrence_value` - See Argument Reference above.
* `scheduled_policy/start_time` - See Argument Reference above.
* `scheduled_policy/end_time` - See Argument Reference above.
