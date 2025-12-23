---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_mapreduce_cluster_node_batch_shrink"
description: |-
  Use this resource to batch shrink the nodes of the MRS cluster within HuaweiCloud.
---

# huaweicloud_mapreduce_cluster_node_batch_shrink

Use this resource to batch shrink the nodes of the MRS cluster within HuaweiCloud.

~> If you use this resource, please use `lifecycle.ignore_changes` to ignore the `node_count` parameter of the
   corresponding node in the `huaweicloud_mapreduce_cluster` resource.

-> 1. This resource only supports shrinking the nodes of the postpaid billing cluster.
   <br>2. Currently, this resource only supports shrinking the nodes of Core nodes and Task nodes, not Master nodes.
   <br>3. This resource is only a one-time action resource for batch shrinking the nodes of the cluster. Deleting
   this resource will not clear the corresponding request record, but will only remove the resource information from
   the tfstate file.

## Example Usage

### Shrink cluster nodes by count

```hcl
variable "cluster_id" {}
variable "node_group_name" {}

resource "huaweicloud_mapreduce_cluster_node_batch_shrink" "test" {
  cluster_id      = var.cluster_id
  node_group_name = var.node_group_name
  node_count      = 2
}
```

### Shrink cluster nodes by resource IDs

```hcl
variable "cluster_id" {}
variable "node_group_name" {}
variable "resource_ids" {
  type = list(string)
}

resource "huaweicloud_mapreduce_cluster_node_batch_shrink" "test" {
  cluster_id      = var.cluster_id
  node_group_name = var.node_group_name
  resource_ids    = var.resource_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the MRS cluster is located.  
  If omitted, the provider-level region will be used.  
  Changing this creates a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the cluster to be shrunk.

* `node_group_name` - (Required, String, NonUpdatable) Specifies the name of the node group to which the nodes to be
  shrunk belong.

* `node_count` - (Optional, Int, NonUpdatable) Specifies the number of nodes to be deleted from the node group.
  If specified, nodes are automatically selected for deletion according to system rules.

* `resource_ids` - (Optional, List, NonUpdatable) Specifies the ID list of resource nodes to be deleted.  
  Only nodes in the **shutdown**, **lost**, **unknown**, **detached** and **error** states can be specified for shrinkage.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
