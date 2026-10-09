---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_node_batch_unsubscribe"
description: |-
  Use this resource to batch unsubscribe the ModelArts nodes within HuaweiCloud.
---

# huaweicloud_modelartsv2_node_batch_unsubscribe

Use this resource to batch delete the ModelArts nodes within HuaweiCloud.

-> This resource is only a one-time action resource for batch unsubscribe the ModelArts nodes. Deleting this resource
   will not clear the corresponding request record, but will only remove the resource information from the tfstate file.

~> This resource can only be used to delete nodes with prePaid billing.

## Example Usage

### Batch delete the A2 nodes or A3 hyper nodes

```hcl
variable "resource_pool_name" {}
variable "node_ids" {
  type = list(string)
}

resource "huaweicloud_modelartsv2_node_batch_unsubscribe" "test" {
  resource_pool_name = var.resource_pool_name
  node_ids           = var.node_ids
}
```

### Batch delete the A5 hyper nodes

```hcl
variable "resource_pool_name" {}
variable "node_ids" {
  type = list(string)
}
variable "batch_uids" {
  type = list(string)
}

resource "huaweicloud_modelartsv2_node_batch_unsubscribe" "test" {
  resource_pool_name = var.resource_pool_name
  node_ids           = var.node_ids
  resource_ids       = var.batch_uids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the resource nodes are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `resource_pool_name` - (Required, String, NonUpdatable) Specifies the resource pool name to which the resource nodes
  belong.

* `node_ids` - (Required, List, NonUpdatable) Specifies the ID list of resource nodes that used to refresh the
  deleting status.

* `resource_ids` - (Optional, List, NonUpdatable) Specifies the ID list of resources to be deleted.  
  Optional for the A2 nodes and A3 hyper nodes, the elements of this list are also node IDs.  
  Required for the A5 hyper nodes, the elements of this list are batch IDs (defined in the `metadata.labels`), and only
  one entry is required for these batch IDs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 45 minutes.
