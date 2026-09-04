---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_node_batch_bind"
description: |-
  Manages a ModelArts resource pool node batch bind action resource within HuaweiCloud.
---

# huaweicloud_modelartsv2_node_batch_bind

Manages a ModelArts resource pool node batch bind action resource. This resource is used to batch bind nodes
to a logical sub-pool within a dedicated resource pool.

-> This resource is a one-time action resource. Deleting this resource will not clear the corresponding
   request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Bind nodes to a logical sub-pool

```hcl
variable "pool_id" {}
variable "quota_name" {}

resource "huaweicloud_modelartsv2_node_batch_bind" "test" {
  pool_id = var.pool_id

  nodes {
    name       = "os-node-created-8888g"
    quota_name = var.quota_name
  }

  drain = true
}
```

### Unbind nodes from any logical sub-pool

```hcl
variable "pool_id" {}

resource "huaweicloud_modelartsv2_node_batch_bind" "test" {
  pool_id = var.pool_id

  nodes {
    name = "os-node-created-8888g"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the resource pool is located.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `pool_id` - (Required, String, ForceNew) Specifies the ID of the resource pool to which the nodes belong.
  Changing this parameter will create a new resource.

* `nodes` - (Required, List, ForceNew) Specifies the list of nodes to be bound to the logical sub-pool.
  Changing this parameter will create a new resource.
  The [nodes](#NodeBatchBindNodes) structure is documented below.

* `drain` - (Optional, Bool, ForceNew) Specifies whether to drain the nodes during the bind operation.
  Changing this parameter will create a new resource.

<a name="NodeBatchBindNodes"></a>
The `nodes` block supports:

* `name` - (Required, String) Specifies the name of the node to be bound.

* `quota_name` - (Optional, String) Specifies the ID of the logical sub-pool to which the node is bound.
  If left empty, the node will be unbound from any logical sub-pool.
