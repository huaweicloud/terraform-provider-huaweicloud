---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_workload_plan_stage"
description: |-
  Manages a GaussDB(DWS) workload plan stage resource within HuaweiCloud.
---

# huaweicloud_dws_workload_plan_stage

Manages a GaussDB(DWS) workload plan stage resource within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "plan_id" {}
variable "stage_name" {}
variable "pool_name" {}

resource "huaweicloud_dws_workload_plan_stage" "test" {
  cluster_id = var.cluster_id
  plan_id    = var.plan_id
  name       = var.stage_nme
  start_time = "01:00:00"
  end_time   = "00:00:00"

  queues {
    name = var.pool_name
  
    configuration {
      resource_name  = "cpu"
      resource_value = 1
    }
    configuration {
      resource_name  = "cpu_limit"
      resource_value = 0
    }
    configuration {
      resource_name  = "memory"
      resource_value = 0
    }
    configuration {
      resource_name  = "concurrency"
      resource_value = 10
    }
    configuration {
      resource_name  = "shortQueryConcurrencyNum"
      resource_value = -1
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the ID of the cluster to which the workload plan belongs.
  Changing this creates a new resource.

* `plan_id` - (Required, String, ForceNew) Specifies the ID of the plan to which the workload plan stage belongs.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the workload plan stage, which must be unique
  and contains `3` to `28` characters, composed only of lowercase letters, numbers, or underscores (_),
  and must start with a lowercase letter. Changing this creates a new resource.

* `start_time` - (Required, String, ForceNew) Specifies the start time of the workload plan.
  The time format is **hh:mm:ss**. Changing this creates a new resource.

* `end_time` - (Required, String, ForceNew) Specifies the end time of the workload plan.
  The time format is **hh:mm:ss**. Changing this creates a new resource.

* `queues` - (Required, List, ForceNew) Specifies the workload queue configurations.
  Changing this creates a new resource.
  The [queues](#block_queues) structure is documented below.

* `month` - (Optional, String, ForceNew) Specifies the execution months of the workload plan. The valid value ranges
  from `1` to `12`, separate by commas, e.g. **1,3,5**. Default to all months.
  Changing this creates a new resource.

* `day` - (Optional, String, ForceNew) Specifies the execution days of the workload plan. The valid value ranges
  from `1` to `31`, separate by commas, e.g. **1,13,25**. Default to all days.
  Changing this creates a new resource.

<a name="block_queues"></a>
The `queues` block supports:

* `configuration` - (Required, List, ForceNew) Specifies the configuration information for workload queue.
  Changing this creates a new resource.
  The [configuration](#block_queues_configuration) structure is documented below.

* `name` - (Required, String, ForceNew) Specifies the name of workload queue which the workload plan stage running.
  Changing this creates a new resource.

<a name="block_queues_configuration"></a>
The `configuration` block supports:

* `resource_name` - (Required, String, ForceNew) Specifies the resource name to be configured for the workload queue.
  Changing this creates a new resource. Value options:  
  + **cpu**: Percentage of CPU time that can be used by users associated with the current workload queue to execute jobs.
  + **cpu_limit**: Maximum percentage of CPU cores used by a database user in a workload queue.
  + **memory**: Percentage of the memory that can be used by a workload queue.
  + **concurrency**: Maximum number of concurrent queries in a workload queue.
  + **shortQueryConcurrencyNum**: Maximum number of concurrent short queries in a workload queue.

* `resource_value` - (Required, Int, ForceNew) Specifies the value of the resource attribute for the workload queue.
  Changing this creates a new resource.  
  When `resource_name` is **cpu**, the value is an integer ranging from `1` to `99`.  
  When `resource_name` is **cpu_limit**, the value is an integer ranging from `0` to `100`. `0` indicates no limit.  
  When `resource_name` is **memory**, `0` indicates no limit.  
  When `resource_name` is **shortQueryConcurrencyNum**, `-1` indicates no limit.  

* `value_unit` - (Optional, String, ForceNew) Specifies the value unit of the resource attribute for the workload queue.
  Changing this creates a new resource.

* `resource_description` - (Optional, String, ForceNew) Specifies the description of the resource attribute for
  the workload queue. Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID.

## Import

The workload plan stage can be imported using `cluster_id`, `plan_id` and `name`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_dws_workload_plan_stage.test <cluster_id>/<plan_id>/<name>
```
