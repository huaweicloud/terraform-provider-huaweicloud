---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_node_pool_nodes_add"
description: |-
  Use this resource to add nodes into a node pool within HuaweiCloud.
---

# huaweicloud_cce_node_pool_nodes_add

Use this resource to add nodes into a node pool within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "cluster_id" {}
variable "nodepool_id" {}
variable "server_id_1" {}
variable "server_id_2" {}

resource "huaweicloud_cce_node_pool_nodes_add" "test" {
  cluster_id  = var.cluster_id
  nodepool_id = var.nodepool_id

  node_list {
    server_id = var.server_id_1
  }
  node_list {
    server_id = var.server_id_2
  }
}

```

  ~> When the ECS instance is added into the node pool, the `image_id`, `security_group_ids` and `tags`
    will be changed. You can ignore these changes as below.

```hcl
resource "huaweicloud_compute_instance" "my_instance" {
    ...

  lifecycle {
    ignore_changes = [
      ignore_changes = [ image_id, security_group_ids, tags ]
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the CCE pool nodes add resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the cluster ID.

* `nodepool_id` - (Required, String, NonUpdatable) Specifies the node pool ID.

* `node_list` - (Required, List, NonUpdatable) Specifies the list of nodes to add into the pool.
  The [node_list](#node_list) structure is documented below.

* `remove_nodes_on_delete` - (Optional, Bool) Whether to remove nodes when delete this resource.
  If set to **false**, it will only be removed from the state. Defaults to **false**.

<a name="node_list"></a>
The `node_list` block supports:

* `server_id` - (Required, String, NonUpdatable) Specifies server ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.

* `delete` - Default is 20 minutes.
