---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_mapreduce_scaling_policy"
description: ""
---

# huaweicloud_mapreduce_scaling_policy

Manages a scaling policy of MapReduce cluster within HuaweiCloud.  

## Example Usage

```hcl
variable "mrs_cluster_id" {}
variable "task_node_group_name" {}
variable "script_uri" {}

resource "huaweicloud_mapreduce_scaling_policy" "test" {
  cluster_id = var.mrs_cluster_id
  node_group = var.task_node_group_name

  auto_scaling_enable = true
  min_capacity        = 4
  max_capacity        = 10

  resources_plans {
    period_type  = "daily"
    start_time   = "01:00"
    end_time     = "03:00"
    min_capacity = 5
    max_capacity = 10
  }

  rules {
    name               = "default-expand-1"
    adjustment_type    = "scale_out"
    cool_down_minutes  = 20
    scaling_adjustment = 1
    trigger {
      metric_name         = "YARNAppRunning"
      metric_value        = "75"
      comparison_operator = "GT"
      evaluation_periods  = 1
    }
  }

  rules {
    name               = "default-shrink-1"
    adjustment_type    = "scale_in"
    cool_down_minutes  = 20
    scaling_adjustment = 1
    trigger {
      metric_name         = "YARNAppRunning"
      metric_value        = "25"
      comparison_operator = "LT"
      evaluation_periods  = 1
    }
  }

  exec_scripts {
    name          = "script_1"
    uri           = var.script_uri
    parameters    = ""
    nodes         = [var.task_node_group_name]
    active_master = false
    fail_action   = "continue"
    action_stage  = "before_scale_out"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, ForceNew) The MRS cluster ID to which the auto scaling policy applies.

  Changing this parameter will create a new resource.

* `node_group` - (Required, String) Name of the node group to which the auto scaling policy applies.  
  Currently, only Task nodes support auto scaling rules.

* `auto_scaling_enable` - (Required, Bool) Whether to enable the auto scaling policy.  

* `min_capacity` - (Required, Int) Minimum number of nodes in the node group. Value range: 0 to 500.  

* `max_capacity` - (Required, Int) Maximum number of nodes in the node group. Value range: 0 to 500.  

* `resources_plans` - (Optional, List) The list of resources plans.  
The [resources_plans](#ScalingPolicy_ResourcesPlan) structure is documented below.

* `rules` - (Optional, List) The list of auto scaling rules.  
  When auto scaling is enabled, either `resources_plans` or `rules` must be configured.
  The [rules](#ScalingPolicy_Rule) structure is documented below.

* `exec_scripts` - (Optional, List) The list of custom scaling automation scripts.  
  When auto scaling is enabled, either `resources_plans` or `rules` must be configured.
  The [exec_scripts](#ScalingPolicy_ExecScript) structure is documented below.

<a name="ScalingPolicy_ResourcesPlan"></a>
The `resources_plans` block supports:

* `period_type` - (Required, String) Cycle type of a resource plan.  
  Currently, only the following cycle type is supported: **daily**.  

* `start_time` - (Required, String) The start time of a resource plan.  
  The value is in the format of **hour:minute**, indicating that the time ranges from 00:00 to 23:59.

* `end_time` - (Required, String) End time of a resource plan.  
  The value is in the format of **hour:minute**.
  The interval between end_time and start_time must be greater than or equal to 30 minutes.

* `min_capacity` - (Required, Int) Minimum number of the preserved nodes in a node group in a resource plan.
   Value range: 0 to 500.  

* `max_capacity` - (Required, Int) Maximum number of the preserved nodes in a node group in a resource plan.
   Value range: 0 to 500.  

<a name="ScalingPolicy_Rule"></a>
The `rules` block supports:

* `name` - (Required, String) Name of an auto scaling rule.  
  The name can contain only 1 to 64 characters. Only letters, digits, hyphens (-), and underscores (_) are allowed.
  Rule names must be unique in a node group.

* `adjustment_type` - (Required, String) Auto scaling rule adjustment type.  
  The following options are supported:  
    + **scale_out**: cluster scale-out.
    + **scale_in**: cluster scale-in.

* `cool_down_minutes` - (Required, Int) Cluster cooling time after an auto scaling rule is triggered,
  when no auto scaling operation is performed.  
  The unit is minute. Value range: 0 to 10,080. One week is equal to 10,080 minutes.

* `scaling_adjustment` - (Required, Int) Number of nodes that can be adjusted once. Value range: 1 to 100.  

* `trigger` - (Required, List) Condition for triggering a rule.  
  The [trigger](#ScalingPolicy_Trigger) structure is documented below.

* `description` - (Optional, String) Description about an auto scaling rule.  
  It contains a maximum of 1,024 characters.

<a name="ScalingPolicy_Trigger"></a>
The `trigger` block supports:

* `metric_name` - (Required, String) Metric name.  
  This triggering condition makes a judgment according to the value of the metric.  
  A metric name contains a maximum of 64 characters.  
  For details about metric names, see [Configuring Auto Scaling for an MRS Cluster](https://support.huaweicloud.com/intl/en-us/qs-mrs/mrs_09_0005.html).

* `metric_value` - (Required, String) Metric threshold to trigger a rule.  
  The parameter value can only be an integer or number with two decimal places.
  The value type and range must correspond to the metric_name.

* `comparison_operator` - (Optional, String) Metric judgment logic operator.  
  The following options are supported:
    + **LT**: less than.
    + **GT**: greater than.
    + **LTOE**: less than or equal to.
    + **GTOE**: greater than or equal to.

* `evaluation_periods` - (Required, Int) Number of consecutive five-minute periods,
  during which a metric threshold is reached.  
  Value range: 1 to 288.

<a name="ScalingPolicy_ExecScript"></a>
The `exec_scripts` block supports:

* `name` - (Required, String) Name of a custom automation script.  
  The name can contain only 1 to 64 characters. Only letters, digits, hyphens (-), and underscores (_) are allowed.
  Script names must be unique in a node group.

* `uri` - (Required, String) Path of a custom automation script.  
  Set this parameter to an OBS bucket path or a local VM path.  
  OBS bucket path: Enter a script path manually. for example, s3a://XXX/scale.sh.  
  Local VM path: Enter a script path. The script path must start with a slash (/) and end with .sh.

* `parameters` - (Optional, String) Parameters of a custom automation script.  
  Multiple parameters are separated by space.  
  The following predefined system parameters can be transferred:  
    + **${mrs_scale_node_num}**: Number of the nodes to be added or removed.
    + **${mrs_scale_type}**: Scaling type. The value can be **scale_out** or **scale_in**.
    + **${mrs_scale_node_hostnames}**: Host names of the nodes to be added or removed.
    + **${mrs_scale_node_ips}**: IP addresses of the nodes to be added or removed.
    + **${mrs_scale_rule_name}**: Name of the rule that triggers auto scaling.

  Other user-defined parameters are used in the same way as those of common shell scripts. Parameters are separated by space.

* `nodes` - (Required, List) Type of a node where the custom automation script is executed.  
  The node type can be **Master**, **Core**, or **Task**.

* `active_master` - (Optional, Bool) Whether the custom automation script runs only on the active Master node.  
  The default value is **false**, indicating that the custom automation script can run on all Master nodes.

* `action_stage` - (Required, String) Time when a script is executed.  
  The following options are supported:  
    + **before_scale_out**: before scale-out.
    + **before_scale_in**: before scale-in.
    + **after_scale_out**: after scale-out.
    + **after_scale_in**: after scale-in.

* `fail_action` - (Required, String) Whether to continue to execute subsequent scripts and create a cluster after
  the custom automation script fails to be executed.  
  The following options are supported:  
    + **continue**: Continue to execute subsequent scripts.
    + **errorout**: Stop the action.  
  
  -> You are advised to set this parameter to **continue** in the commissioning phase so that the cluster
     can continue to be installed and started no matter whether the custom automation script is executed successfully.  
     The scale-in operation cannot be undone. Therefore, `fail_action` must be set to **continue** for the
     scripts that are executed after scale-in.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The scaling policy of MapReduce cluster can be imported using `cluster_id`, `node_group`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_mapreduce_scaling_policy.test <cluster_id>/<node_group>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `exec_scripts`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition
should be updated to align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_mapreduce_scaling_policy" "test" {
    ...

  lifecycle {
    ignore_changes = [
      exec_scripts,
    ]
  }
}
```
