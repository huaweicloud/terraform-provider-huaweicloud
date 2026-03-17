---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_mapreduce_scaling_policy_v2"
description: |-
  Manages an auto scaling policy of an MRS cluster within HuaweiCloud.
---

# huaweicloud_mapreduce_scaling_policy_v2

Manages an auto scaling policy of an MRS cluster within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "node_group_name" {}

resource "huaweicloud_mapreduce_scaling_policy_v2" "test" {
  cluster_id         = var.cluster_id
  node_group_name    = var.node_group_name
  resource_pool_name = "default"

  auto_scaling_policy {
    auto_scaling_enable = true
    min_capacity        = 0
    max_capacity        = 5

    resources_plans {
      period_type    = "daily"
      start_time     = "06:00"
      end_time       = "20:00"
      min_capacity   = 0
      max_capacity   = 2
      effective_days = ["MONDAY"]
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
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the scaling policy is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the cluster.

* `node_group_name` - (Required, String, NonUpdatable) Specifies the name of the node group.

* `resource_pool_name` - (Required, String, NonUpdatable) Specifies the name of the resource pool.  
  When the cluster version does not support elastic scaling by the specified resource pool, this parameter
  must be set to **default**.

* `auto_scaling_policy` - (Required, List) Specifies the configurations of the auto scaling policy.  
  The [auto_scaling_policy](#scaling_policy_v2_auto_scaling_policy) structure is documented below.

<a name="scaling_policy_v2_auto_scaling_policy"></a>
The `auto_scaling_policy` block supports:

* `auto_scaling_enable` - (Optional, Bool) Specifies whether to enable the auto scaling policy.  
  Defaults to **false**.

* `min_capacity` - (Optional, Int) Specifies the minimum number of nodes in the node group.  
  Defaults to `0`.  
  The valid value ranges from `0` to `500`.

* `max_capacity` - (Optional, Int) Specifies the maximum number of nodes in the node group.  
  Defaults to `0`.  
  The valid value ranges from `0` to `500`.

* `resources_plans` - (Optional, List) Specifies the list of resource plans.  
  The [resources_plans](#scaling_policy_v2_resources_plans) structure is documented below.

* `rules` - (Optional, List) Specifies the list of auto scaling rules.  
  The [rules](#scaling_policy_v2_rule) structure is documented below.

-> At least one of the `resources_plans` and `rules` parameters must be specified.

* `tags` - (Optional, Map) Specifies the key/value pairs associated with auto scaling policy.  
  The maximum number of tags is `20`.

<a name="scaling_policy_v2_resources_plans"></a>
The `resources_plans` block supports:

* `period_type` - (Required, String) Specifies the period type of the resource plan.  
  The valid values are as follows:
  + **daily**

* `start_time` - (Required, String) Specifies the start time of the resource plan, in `HH:mm` format.

* `end_time` - (Required, String) Specifies the end time of the resource plan, in `HH:mm` format.  
  The interval between `start_time` and `end_time` must be greater than or equal to `30` minutes.

* `min_capacity` - (Optional, Int) Specifies the minimum number of retained nodes for the node group in the
  resource plan.  
  Defaults to `0`.  
  The valid value ranges from `0` to `500`.

* `max_capacity` - (Optional, Int) Specifies the maximum number of retained nodes for the node group in the
  resource plan.
  Defaults to `0`.  
  The valid value ranges from `0` to `500`.

* `effective_days` - (Optional, List) Specifies the effective day list of the resource plan.  
  The valid values are as follows:
  + **MONDAY**
  + **TUESDAY**
  + **WEDNESDAY**
  + **THURSDAY**
  + **FRIDAY**
  + **SATURDAY**
  + **SUNDAY**

  If not specified, it means the resource plan will take effect on all days.

<a name="scaling_policy_v2_rule"></a>
The `rules` block supports:

* `name` - (Required, String) Specifies the name of the rule.  
  The maximum length is `64` characters.  
  Only letters, digits, hyphens (-), and underscores (_) are allowed.  
  The rule name must be unique in the same node group.

* `adjustment_type` - (Required, String) Specifies the adjustment type of the rule.  
  The valid values are as follows:
  + **scale_out**
  + **scale_in**

* `cool_down_minutes` - (Required, Int) Specifies the cool down time of the cluster after the scaling rule
  is triggered, in minutes.  
  The valid value ranges from `5` to `10,080`.

* `scaling_adjustment` - (Required, Int) Specifies the number of adjusted nodes in one scaling action.  
  The valid value ranges from `1` to `100`.

* `trigger` - (Required, List) Specifies the trigger condition list of the rule.  
  The [trigger](#scaling_policy_v2_rule_trigger) structure is documented below.

* `description` - (Optional, String) Specifies the description of the rule.  
  Only letters, digits, hyphens (-), and underscores (_) are allowed.  

<a name="scaling_policy_v2_rule_trigger"></a>
The `trigger` block supports:

* `metric_name` - (Required, String) Specifies the name of the metric.  
  For details about metric names, see [documentation](https://support.huaweicloud.com/intl/en-us/bestpractice-mrs/mrs_05_0132.html#mrs_05_0132__en-us_topic_0000001168444324_table15133845184415).

* `metric_value` - (Required, String) Specifies the threshold value of the metric.  
  The parameter value can only be an integer or number with two decimal places.

* `evaluation_periods` - (Required, Int) Specifies the number of consecutive periods that meet the metric threshold.  
  The interval between the consecutive periods is `5` minutes.
  The valid value ranges from `1` to `200`.

* `comparison_operator` - (Optional, String) Specifies the comparison operator of the metric judgment logic.  
  The valid values are as follows:
  + **LT**
  + **GT**
  + **LTOE**
  + **GTOE**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The resource can be imported using the `cluster_id`, `node_group_name` and `resource_pool_name`, separated by
slashes, e.g.

```bash
$ terraform import huaweicloud_mapreduce_scaling_policy_v2.test <cluster_id>/<node_group_name>/<resource_pool_name>
```
