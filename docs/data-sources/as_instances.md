---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_instances"
description: |-
  Use this data source to get AS instances within HuaweiCloud.
---

# huaweicloud_as_instances

Use this data source to get AS instances within HuaweiCloud.

## Example Usage

```hcl
variable "scaling_group_id" {}

data "huaweicloud_as_instances" "test" {
  scaling_group_id = var.scaling_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `scaling_group_id` - (Required, String) Specifies the AS group ID.

* `life_cycle_state` - (Optional, String) Specifies the instance lifecycle status in the AS group. Valid values are:
  + **INSERVICE**: The instance is enabled.
  + **PENDING**: The instance is being added to the AS group.
  + **PENDING_WAIT**: The instance is waiting to be added to the AS group.
  + **REMOVING**: The instance is being removed from the AS group.
  + **REMOVING_WAIT**: The instance is waiting to be removed from the AS group.
  + **STANDBY**: The instance is in standby state.
  + **ENTERING_STANDBY**: The instance is entering the standby state.

* `health_status` - (Optional, String) Specifies the instance health status. Valid values are:
  + **INITIALIZING**: The instance is initializing.
  + **NORMAL**: The instance is normal.
  + **ERROR**: The instance is abnormal.

* `protect_from_scaling_down` - (Optional, String) Specifies the instance protection status. Valid values are:
  + **true**: Instance protection is enabled.
  + **false**: Instance protection is disabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - The details about the instances in the AS group.
  The [instances](#ASInstances_instances) structure is documented below.

<a name="ASInstances_instances"></a>
The `instances` block supports:

* `instance_id` - Indicates the instance ID.

* `instance_name` - Indicates the instance name.

* `scaling_group_id` - Indicates the ID of the AS group to which the instance belongs.

* `scaling_group_name` - Indicates the name of the AS group to which the instance belongs.

* `life_cycle_state` - Indicates the instance lifecycle status in the AS group.

* `health_status` - Indicates the instance health status.

* `scaling_configuration_name` - Indicates the AS configuration name.

* `scaling_configuration_id` - Indicates the AS configuration ID.

* `created_at` - Indicates the time when the instance is added to the AS group. The time format complies with UTC.

* `protect_from_scaling_down` - Indicates the instance protection status.
