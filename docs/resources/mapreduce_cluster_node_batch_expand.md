---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_mapreduce_cluster_node_batch_expand"
description: |-
  Use this resource to batch expand the nodes of the MRS cluster within HuaweiCloud.
---

# huaweicloud_mapreduce_cluster_node_batch_expand

Use this resource to batch expand the nodes of the MRS cluster within HuaweiCloud.

~> If you use this resource, please use `lifecycle.ignore_changes` to ignore the `node_count` parameter of the
   corresponding node in the `huaweicloud_mapreduce_cluster` resource.

-> This resource is only a one-time action resource for batch expanding the nodes of the cluster. Deleting this
   resource will not clear the corresponding request record, but will only remove the resource information from the
   tfstate file.

## Example Usage

```hcl
variable "cluster_id" {}
variable "node_group_name" {}

resource "huaweicloud_mapreduce_cluster_node_batch_expand" "test" {
  cluster_id      = var.cluster_id
  node_group_name = var.node_group_name
  node_count      = 2
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the MRS cluster is located.  
  If omitted, the provider-level region will be used.  
  Changing this creates a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies ID of the cluster to which the nodes
   to be expanded belong.

* `node_group_name` - (Required, String, NonUpdatable) Specifies the name of the node group to which the nodes to be
  expanded belong.

* `node_count` - (Required, Int, NonUpdatable) Specifies the number of nodes to be expanded.

* `skip_bootstrap_scripts` - (Optional, Bool, NonUpdatable) Specifies whether to skip bootstrap scripts.  
  Defaults to **true**.
  + **true**: Skip bootstrap scripts.
  + **false**: Do not skip bootstrap scripts.

* `scale_without_start` - (Optional, Bool, NonUpdatable) Specifies whether to start the components on the node after it
  has been expanded.  
  Defaults to **false**.
  + **true**: Do not start the components on the node after it has been expanded.
  + **false**: Start the components on the node after it has been expanded.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
