---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_queues"
description: |-
  Use this data source to query the DLI queues within HuaweiCloud.
---

# huaweicloud_dli_queues

Use this data source to query the DLI queues within HuaweiCloud.

## Example Usage

### Filter by queue name under a specified elastic resource pool

```hcl
variable "elastic_resource_pool_name" {}
variable "queue_name" {}

data "huaweicloud_dli_queues" "test" {
  elastic_resource_pool_name = var.elastic_resource_pool_name
  queue_name                 = var.queue_name
}
```

### Filter by queue type

```hcl
data "huaweicloud_dli_queues" "test" {
  queue_type = "sql"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the queues.
  If omitted, the provider-level region will be used.

* `elastic_resource_pool_name` - (Optional, String) Specifies the name of the elastic resource pool.
  This parameter is conflict with `queue_type`, `with_privilege`, `with_charge_info` and `tags`.

* `queue_name` - (Optional, String) Specifies the name of the queue to be queried.
  This parameter is valid only when `elastic_resource_pool_name` is specified.

* `queue_type` - (Optional, String) Specifies the type of the queue to be queried. This parameter is conflict with
  `elastic_resource_pool_name`.  
  The valid values are as follows:
  + **sql**
  + **general**
  + **all**
  
* `with_privilege` - (Optional, Bool) Specifies whether to return permission information.  
  This parameter is conflict with `elastic_resource_pool_name`.

* `with_charge_info` - (Optional, Bool) Specifies whether to return charge information.  
  This parameter is conflict with `elastic_resource_pool_name`.

* `tags` - (Optional, String) Specifies the tags to filter queues.  
  This parameter is conflict with `elastic_resource_pool_name`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `queues` - The list of queues that matched filter parameters.
  The [queues](#dli_queues_attr) structure is documented below.

<a name="dli_queues_attr"></a>
The `queues` block supports:

* `name` - The name of the queue.

* `type` - The type of the queue.

* `enterprise_project_id` - The ID of the enterprise project.

* `created_at` - The creation time of the queue, in RFC3339 format.

* `owner` - The owner of the queue.

* `engine` - The engine type of the queue.

* `scaling_policies` - The scaling policies of the queue.  
  The [scaling_policies](#dli_queues_scaling_policies_attr) structure is documented below.

<a name="dli_queues_scaling_policies_attr"></a>
The `scaling_policies` block supports:

* `priority` - The priority of the scaling policy.

* `impact_start_time` - The start time of the scaling policy.

* `impact_stop_time` - The stop time of the scaling policy.

* `min_cu` - The minimum CU of the scaling policy.

* `max_cu` - The maximum CU of the scaling policy.

* `inherit_elastic_resource_pool_max_cu` - Whether to inherit the maximum CU of the elastic resource pool.
