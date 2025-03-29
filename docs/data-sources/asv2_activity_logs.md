---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_asv2_activity_logs"
description: |-
  Use this data source to get a list of AS (V2 version) scaling activity logs within HuaweiCloud.
---

# huaweicloud_asv2_activity_logs

Use this data source to get a list of AS (V2 version) scaling activity logs within HuaweiCloud.

## Example Usage

```hcl
variable "group_id" {}

data "huaweicloud_asv2_activity_logs" "test" {
  scaling_group_id = var.group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `scaling_group_id` - (Required, String) Specifies the AS group ID.

* `log_id` - (Optional, String) Specifies the scaling action log ID.

* `start_time` - (Optional, String) Specifies the start time that complies with UTC for querying scaling action logs.
  The format of the start time is yyyy-MM-ddThh:mm:ssZ.

* `end_time` - (Optional, String) Specifies the end time that complies with UTC for querying scaling action logs.
  The format of the end time is yyyy-MM-ddThh:mm:ssZ.

* `type` - (Optional, String) Specifies the types of the scaling actions to be queried. Different types are separated by
  commas (,). Valid values are:
  + **NORMAL**: Indicates a common scaling action.
  + **MANUAL_REMOVE**: Indicates manually removing instances from an AS group.
  + **MANUAL_DELETE**: Indicates manually removing and deleting instances from an AS group.
  + **MANUAL_ADD**: Indicates manually adding instances to an AS group.
  + **ELB_CHECK_DELETE**: Indicates that instances are removed from an AS group and deleted based on the ELB health check
    result.
  + **AUDIT_CHECK_DELETE**: Indicates that instances are removed from an AS group and deleted based on the audit.
  + **DIFF**: Indicates that the number of expected instances is different from the actual number of instances.
  + **MODIFY_ELB**: Indicates the load balancer migration.
  + **ENTER_STANDBY**: Indicates setting instances to standby mode.
  + **EXIT_STANDBY**: Indicates canceling standby mode for instances.

* `status` - (Optional, String) Specifies the status of the scaling action. Valid values are:
  + **SUCCESS**: The scaling action has been performed.
  + **FAIL**: Performing the scaling action failed.
  + **DOING**: The scaling action is being performed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `scaling_activity_log` - The scaling action logs.

  The [scaling_activity_log](#scaling_activity_log_struct) structure is documented below.

<a name="scaling_activity_log_struct"></a>
The `scaling_activity_log` block supports:

* `id` - The scaling action log ID.

* `scaling_value` - The number of added or deleted instances during the scaling.

* `desire_value` - The expected number of instances for the scaling action.

* `type` - The type of the scaling action.

* `status` - The status of the scaling action. Valid values are:
  + **SUCCESS**: The scaling action has been performed.
  + **FAIL**: Performing the scaling action failed.
  + **DOING**: The scaling action is being performed.

* `description` - The description of the scaling action.

* `start_time` - The start time of the scaling action. The time format must comply with UTC.

* `end_time` - The end time of the scaling action. The time format must comply with UTC.

* `instance_value` - The number of instances in the AS group.

* `instance_removed_list` - The names of the ECSs that are removed from the AS group in a scaling action.

  The [instance_removed_list](#scaling_activity_log_scaling_instance_struct) structure is documented below.

* `instance_standby_list` - The ECSs that are set to standby mode or for which standby mode is canceled in a scaling action.

  The [instance_standby_list](#scaling_activity_log_scaling_instance_struct) structure is documented below.

* `instance_failed_list` - The ECSs for which a scaling action fails.

  The [instance_failed_list](#scaling_activity_log_scaling_instance_struct) structure is documented below.

* `instance_deleted_list` - The names of the ECSs that are removed from the AS group and deleted in a scaling action.

  The [instance_deleted_list](#scaling_activity_log_scaling_instance_struct) structure is documented below.

* `instance_added_list` - The names of the ECSs that are added to the AS group in a scaling action.

  The [instance_added_list](#scaling_activity_log_scaling_instance_struct) structure is documented below.

* `lb_bind_failed_list` - The load balancers that failed to be bound to the AS group.

  The [lb_bind_failed_list](#scaling_activity_log_modify_lb_struct) structure is documented below.

* `lb_unbind_failed_list` - The load balancers that failed to be unbound from the AS group.

  The [lb_unbind_failed_list](#scaling_activity_log_modify_lb_struct) structure is documented below.

* `lb_bind_success_list` - The load balancers that are bound to the AS group.

  The [lb_bind_success_list](#scaling_activity_log_modify_lb_struct) structure is documented below.

* `lb_unbind_success_list` - The load balancers that are unbound from the AS group.

  The [lb_unbind_success_list](#scaling_activity_log_modify_lb_struct) structure is documented below.

<a name="scaling_activity_log_scaling_instance_struct"></a>
The `instance_removed_list`, `instance_standby_list`, `instance_failed_list`, `instance_deleted_list`, and
`instance_added_list` block supports:

* `instance_name` - The ECS name.

* `instance_id` - The ECS ID.

* `failed_reason` - The cause of the instance scaling failure.

* `failed_details` - The details of the instance scaling failure.

* `instance_config` - The information about instance configurations.

<a name="scaling_activity_log_modify_lb_struct"></a>
The `lb_bind_failed_list`, `lb_unbind_failed_list`, `lb_bind_success_list`, and `lb_unbind_success_list` block supports:

* `lbaas_listener` - The information about an enhanced load balancer.

  The [lbaas_listener](#lbaas_listener_struct) structure is documented below.

* `listener` - The information about a classic load balancer.

* `failed_reason` - The cause of a load balancer migration failure.

* `failed_details` - The details of a load balancer migration failure.

<a name="lbaas_listener_struct"></a>
The `lbaas_listener` block supports:

* `listener_id` - The listener ID.

* `pool_id` - The backend ECS group ID.

* `protocol_port` - The backend protocol port, which is the port on which a backend ECS listens for traffic.

* `weight` - The weight, which determines the portion of requests a backend ECS processes when being compared to other
  backend ECSs added to the same listener.
