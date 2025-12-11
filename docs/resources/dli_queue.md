---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_queue"
description: ""
---

# huaweicloud_dli_queue

Manages DLI Queue resource within HuaweiCloud

## Example Usage

### Create an exclusive mode queue

```hcl
variable "elastic_resource_pool_name" {}
variable "queue_name" {}

resource "huaweicloud_dli_queue" "test" {
  elastic_resource_pool_name = var.elastic_resource_pool_name
  resource_mode              = 1

  name     = var.queue_name
  cu_count = 16

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the dli queue resource. If omitted,
  the provider-level region will be used. Changing this will create a new VPC channel resource.

* `elastic_resource_pool_name` - (Required, String) Specifies the name of the elastic resource pool to which the queue
  belongs.

  ~> This parameter cannot be updated and will not trigger ForceNew(, only an error will be thrown).

* `resource_mode` - (Required, Int, ForceNew) Specifies the queue resource mode.  
  The valid value is as follows:
  + 1: indicates the exclusive resource mode.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Name of a queue. Name of a newly created resource queue. The name can contain
  only digits, lowercase letters, and underscores (\_), but cannot contain only digits or start with an underscore (_).
  Length range: `1` to `128` characters. Changing this parameter will create a new resource.

* `queue_type` - (Optional, String, ForceNew) Indicates the queue type. Changing this parameter will create a new
  resource. The options are as follows:
  + sql
  + general

  The default value is `sql`.

* `description` - (Optional, String, ForceNew) Description of a queue. Changing this parameter will create a new
  resource.

* `cu_count` - (Required, Int) Minimum number of CUs that are bound to a queue. Initial value can be `16`,
  `64`, or `256`. When scale_out or scale_in, the number must be a multiple of `16`.

* `enterprise_project_id` - (Optional, String, ForceNew) Enterprise project ID. The value 0 indicates the default
  enterprise project. Changing this parameter will create a new resource.

* `platform` - (Optional, String, ForceNew) CPU architecture of queue compute resources. Changing this parameter will
  create a new resource. The options are as follows:
  + x86_64 : default value
  + aarch64

* `feature` - (Optional, String, ForceNew)Indicates the queue feature. Changing this parameter will create a new
  resource. The options are as follows:
  + basic: basic type (default value)
  + ai: AI-enhanced (Only the SQL x86_64 dedicated queue supports this option.)

* `tags` - (Optional, Map, ForceNew) Label of a queue. Changing this parameter will create a new resource.

* `scaling_policies` - (Optional, List) Specifies the list of scaling policies of the queue associated with
  an elastic resource pool.
  This parameter is only available if `resource_mode` is set to `1`.
  If you want to use this parameter, you must ensure that there is a scaling policy with a time period from `00:00` to `24:00`.
  The [scaling_policies](#queue_scaling_policies) structure is documented below.
  
  -> After binding an elastic resource pool to a queue, the system will automatically generate a default scaling policy,
     in which the `priority` value is `1`, the `impact_start_time` value is `00:00`, the `impact_stop_time` value is `24:00`,
     and the values of `min_cu` and `max_cu` are equal to the number of CUs of the queue.
     For the default scaling policy, except for the `impact_start_time` and `impact_stop_time`, which are not allowed to
     be modified, other values can be modified according to the actual situation.

* `spark_driver` - (Optional, List) Specifies spark driver configuration of the queue.
  This parameter is only available if `queue_type` is set to `sql`.
  The [spark_driver](#queue_spark_driver) structure is documented below.

<a name="queue_scaling_policies"></a>
The `scaling_policies` block supports:

* `priority` - (Required, Int) Specifies the priority of the queue scaling policy.
  The valid value ranges from `1` to `100`. The larger value means the higher priority.

* `impact_start_time` - (Required, String) Specifies the effective time of the queue scaling policy.
  The value can be set only by hour.

* `impact_stop_time` - (Required, String) Specifies the expiration time of the queue scaling policy.
  The value can be set only by hour.

  -> The time ranges of different scaling policies in the same queue cannot overlap.
     The time range includes the start time but not the end time, e.g. `[00:00, 24:00)`.

* `min_cu` - (Required, Int) Specifies the minimum number of CUs allowed by the scaling policy.
  The number must be a multiple of `4`.

  -> The total minimum CUs of all queues in an elastic resource pool cannot be more than the minimum CUs of the pool.

* `max_cu` - (Required, Int) Specifies the maximum number of CUs allowed by the scaling policy.
  The number must be a multiple of `4`.
  
  -> The maximum CUs of any queue in an elastic resource pool cannot be more than the maximum CUs of the pool.

<a name="queue_spark_driver"></a>
The `spark_driver` block supports:

* `max_instance` - (Optional, Int) Specifies the maximum number of spark drivers that can be started on the queue.
  If the `cu_count` is `16`, the value can only be `2`.
  If The `cu_count` is greater than `16`, the minimum value is `2`, the maximum value is the number of queue CUs
  divided by `16`.

* `max_concurrent` - (Optional, Int) Specifies the maximum number of tasks that can be concurrently executed by a spark
  driver. The valid value ranges from `1` to `32`.

* `max_prefetch_instance` - (Optional, String) Specifies the maximum number of spark drivers to be pre-started on the
  queue. The minimum value is `0`. If the `cu_count` is less than `32`, the maximum value is `1`.
  If the `cu_count` is greater than or equal to `32`, the maximum value is the number of queue CUs divided by `16`.

  -> If the minimum CUs of the queue is less than `16` CUs, the `max_instance` and `max_prefetch_instance` parameters
     does not take effect.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `create_time` - Time when a queue is created.

## Timeouts

This resource provides the following timeouts configuration options:

* `update` - Default is 45 minutes.

## Import

DLI queue can be imported by `name` and `queue_type` (if omitted, the SQL type queue will be imported), separated by a
slash, e.g.

### Import a queue of the specified type (SQL type and general type)

```bash
$ terraform import huaweicloud_dli_queue.test <queue_type>/<name>
```

### Import a SQL type queue

```bash
$ terraform import huaweicloud_dli_queue.test <name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `tags`.
It is generally recommended running `terraform plan` after importing a DLI queue.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dli_queue" "test" {
  ...

  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}
```
