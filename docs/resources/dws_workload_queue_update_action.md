---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_workload_queue_update_action"
description: |-
  Use this resource to update the resource configuration of workload queue within HuaweiCloud.
---

# huaweicloud_dws_workload_queue_update_action

Use this resource to update the resource configuration of workload queue within HuaweiCloud.

-> This resource performs a one-time for updating the configuration of workload queue. Deleting this resource will
   not revert the configuration on the cluster, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "cluster_id" {}
variable "queue_name" {}

variable "configuration_list" {
  type = list(object({
    resource_name  = string
    resource_value = number
  }))
}

resource "huaweicloud_dws_workload_queue_update_action" "test" {
  cluster_id = var.cluster_id
  name       = var.queue_name

  dynamic "configuration" {
    for_each = var.configuration_list

    content {
      resource_name  = configuration.value.resource_name
      resource_value = configuration.value.resource_value
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the workload queue is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the cluster to which the workload queue belongs.

* `name` - (Required, String, NonUpdatable) Specifies the name of the workload queue to be updated.

* `logical_cluster_name` - (Optional, String, NonUpdatable) Specifies the name of the logical cluster to
  which the workload queue belongs.

* `configuration` - (Required, List, NonUpdatable) Specifies the list of workload queue resource items to be updated.
  The [configuration](#dws_workload_queue_configuration) structure is documented below.

  -> The configuration requires the simultaneous declaration of memory, tablespace, activestatements, cpu_limit,
     and cpu_share.  
     Missing any of these fields will cause the parameter validation to fail.

<a name="dws_workload_queue_configuration"></a>
The `configuration` block supports:

* `resource_name` - (Required, String, NonUpdatable) Specifies the resource attribute name.  
  The valid values are as follows:
  + **memory**
  + **tablespace**
  + **activestatements**
  + **cpu_limit**
  + **cpu_share**

* `resource_value` - (Required, Int, NonUpdatable) Specifies the resource attribute value.

* `value_unit` - (Optional, String, NonUpdatable) Specifies the unit of the resource attribute.

* `resource_description` - (Optional, String, NonUpdatable) Specifies the description of the resource attribute.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
