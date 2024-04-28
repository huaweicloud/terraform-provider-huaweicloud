---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_policies"
description: ""
---

# huaweicloud_as_policies

Use this data source to get a list of AS scaling policies within HuaweiCloud.

## Example Usage

```hcl
variable "scaling_group_id" {}

data "huaweicloud_as_policies" "test" {
  scaling_group_id = var.scaling_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the scaling policies.
  If omitted, the provider-level region will be used.

* `scaling_group_id` - (Required, String) Specifies the scaling group ID.

* `scaling_policy_id` - (Optional, String) Specifies the scaling policy ID.

* `scaling_policy_name` - (Optional, String) Specifies the scaling policy name.

* `scaling_policy_type` - (Optional, String) Specifies the scaling policy type.  
  The valid values are as follows:
  + **ALARM**: indicates that the scaling action is triggered by an alarm.
  + **SCHEDULED**: indicates that the scaling action is triggered as scheduled.
  + **RECURRENCE**: indicates that the scaling action is triggered periodically.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policies` - All scaling policies that match the filter parameters.  
  The [policies](#as_policies) structure is documented below.

<a name="as_policies"></a>
The `policies` block supports:

* `id` - The scaling policy ID.

* `scaling_group_id` - The scaling group ID.

* `name` - The scaling policy name.

* `status` - The scaling policy status. The value can be **INSERVICE**, **PAUSED** or **EXECUTING**.

* `type` - The scaling policy type.

* `alarm_id` - The alarm rule ID. This field is not empty while `type` is **ALARM**.

* `scheduled_policy` - The periodic or scheduled scaling policy. This field is not empty while `type` is
  **SCHEDULED** or **RECURRENCE**.  
  The [scheduled_policy](#as_scheduled_policy) structure is documented below.

* `action` - The action details of the scaling policy.  
  The [action](#as_policy_action) structure is documented below.

* `cool_down_time` - The cooling duration, in seconds.

* `created_at` - The (UTC) creation time of the event source, in RFC3339 format.

<a name="as_scheduled_policy"></a>
The `scheduled_policy` block supports:

* `launch_time` - The time when the scaling action is triggered.
  + If `type` is **SCHEDULED**, the time format is **YYYY-MM-DDThh:mmZ**.
  + If `type` is **RECURRENCE**, the time format is **hh:mm**.

* `recurrence_type` - The periodic triggering type. This field is not empty while `type` is **RECURRENCE**.
  The value can be **Daily**, **Weekly** or **Monthly**.

* `recurrence_value` - The frequency at which scaling actions are triggered.
  + When `recurrence_type` is **Daily**, this field is null, indicating daily execution.
  + When `recurrence_type` is **Weekly**, the valid value ranges from `1` to `7`, `1` represents Sunday,
    separate by commas. e.g. **1,3,5**.
  + When `recurrence_type` is **Monthly**, the valid value ranges from `1` to `31`, represent the dates of each month
    separately, separate by commas. e.g. **1,10,13,28**.

* `start_time` - The start time of the scaling action triggered periodically. The time format complies with UTC.
  The time format is **YYYY-MM-DDThh:mmZ**.

* `end_time` - The end time of the scaling action triggered periodically. The time format complies with UTC.
  The time format is **YYYY-MM-DDThh:mmZ**.

<a name="as_policy_action"></a>
The `action` block supports:

* `operation` - The operation to be performed.  
  The valid values are as follows:
  + **ADD**: add instances.
  + **REMOVE**: remove instances.
  + **SET**: set the number of instances to.

* `instance_number` - The number of instances to be operated.

* `instance_percentage` - The percentage of instances to be operated.
