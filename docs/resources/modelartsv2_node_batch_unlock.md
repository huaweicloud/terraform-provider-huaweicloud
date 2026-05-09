---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_node_batch_unlock"
description: |-
  Use this resource to batch unlock nodes in a resource pool within HuaweiCloud.
---

# huaweicloud_modelartsv2_node_batch_unlock

Use this resource to batch unlock nodes in a resource pool within HuaweiCloud.

-> This resource is only a one-time action resource for batch unlocking the ModelArts nodes. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "resource_pool_id" {}
variable "node_names" {
  type = list(string)
}

resource "huaweicloud_modelartsv2_node_batch_unlock" "test" {
  pool_id    = var.resource_pool_id
  node_names = var.node_names
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the resource nodes are located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `pool_id` - (Required, String, NonUpdatable) Specifies the resource pool ID to which the resource nodes belong.

* `node_names` - (Required, List, NonUpdatable) Specifies the name list of resource nodes to be unlocked.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
