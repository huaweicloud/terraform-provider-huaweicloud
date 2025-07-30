---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_resource_pool_node_batch_resize"
description: |-
  Use this resource to batch adjust step size of the hyper instance nodes under the resource pool within HuaweiCloud.
---

# huaweicloud_modelarts_resource_pool_node_batch_resize

Use this resource to batch adjust step size of the hyper instance nodes under the resource pool within HuaweiCloud.

-> This resource is only a one-time operation resource for batch adjust the step size of resource pool nodes.
   Deleting this resource will not clear the corresponding request records, only removing the resource information
   from the tfstate file.

## Increase the step size of prepaid nodes from 1 to 2

```hcl
variable "resource_pool_name" {}
variable "node_batch_uids" {
  type    = list(string)
}
variable "source_node_pool_configuration" {
  type = object({
    node_pool = string
    flavor    = string
  })
}
variable "target_node_pool_configuration" {
  type = object({
    node_pool = string
    flavor    = string
  })
}

resource "huaweicloud_modelarts_resource_pool_node_batch_resize" "test" {
  resource_pool_name = var.resource_pool_name

  dynamic "nodes" {
    for_each = var.node_batch_uids

    content {
      batch_uid = nodes.value
    }
  }

  source {
    node_pool = var.source_node_pool_configuration.node_pool
    flavor    = var.source_node_pool_configuration.flavor

    creating_step {
      type = "hyperinstance"
      step = 1
    }
  }

  target {
    node_pool = var.target_node_pool_configuration.node_pool
    flavor    = var.target_node_pool_configuration.flavor

    creating_step {
      type = "hyperinstance"
      step = 2
    }
  }

  billing = jsonencode({
    autoPay = "1"
  })
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the resource pool is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.
  
* `resource_pool_name` - (Required, String, NonUpdatable) Specifies the resource pool name to which the resource nodes belong.

* `nodes` - (Required, List, NonUpdatable) Specifies the list of nodes to be scaled.  
  The [nodes](#resource_pool_node_batch_resize_nodes) structure is documented below.

* `source` - (Required, List, NonUpdatable) Specifies the configuration of the source node pool to which the node
  to be scaled belongs.  
  The [source](#resource_pool_node_batch_resize_node_pool_configuration) structure is documented below.

* `target` - (Required, List, NonUpdatable) Specifies the configuration of the target node pool to which the node
  to be scaled belongs.  
  The [target](#resource_pool_node_batch_resize_node_pool_configuration) structure is documented below.

* `billing` - (Optional, String, NonUpdatable) Specifies whether to automatically pay, in JSON format.  
  This parameter is **required** only when upgrading specification of the nodes and cannot be set when downgrading.

<a name="resource_pool_node_batch_resize_nodes"></a>
The `nodes` block supports:

* `batch_uid` - (Required, String, NonUpdatable) Specifies the batch UID of the node.

* `delete_node_names` - (Optional, List, NonUpdatable) Specifies the list of nodes to be deleted.  
  This parameter is **required** only when downgrading specification of the nodes and cannot be set when upgrading.

<a name="resource_pool_node_batch_resize_node_pool_configuration"></a>
The `source` and `target` block supports:

* `flavor` - (Required, String, NonUpdatable) Specifies the flavor of the node pool.

* `node_pool` - (Required, String, NonUpdatable) Specifies the name of the node pool.

* `creating_step` - (Required, List, NonUpdatable) Specifies the creating step of the node pool.  
  The [creating_step](#resource_pool_node_batch_resize_source_creating_step) structure is documented below.

<a name="resource_pool_node_batch_resize_source_creating_step"></a>
The `creating_step` block supports:

* `step` - (Required, Int, NonUpdatable) Specifies the step number of the nodes.

* `type` - (Required, String, NonUpdatable) Specifies the type of the nodes.  
  The valid values are as follows:
  + **hyperinstance**

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 90 minutes.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `server_ids` - The list of service IDs corresponding to the currently upgraded specification nodes.
