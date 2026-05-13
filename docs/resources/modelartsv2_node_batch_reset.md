---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_node_batch_reset"
description: |-
  Use this resource to batch reset the ModelArts nodes within HuaweiCloud.
---

# huaweicloud_modelartsv2_node_batch_reset

Use this resource to batch reset the ModelArts nodes within HuaweiCloud.

-> This resource is only a one-time action resource for batch reset the ModelArts nodes. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Reset nodes based on a percentage of nodes

```hcl
variable "pool_id" {}
variable "node_names" {
  type = list(string)
}

resource "huaweicloud_modelartsv2_node_batch_reset" "test" {
  pool_id    = var.pool_id
  node_names = var.node_names

  rolling_config {
    strategy        = "RollingByPercent"
    max_unavailable = 25
  }
}
```

### Reset nodes based on the number of nodes

```hcl
variable "pool_id" {}
variable "node_names" {
  type = list(string)
}

resource "huaweicloud_modelartsv2_node_batch_reset" "test" {
  pool_id    = var.pool_id
  node_names = var.node_names

  rolling_config {
    strategy        = "RollingByNumber"
    max_unavailable = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the resource nodes are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `pool_id` - (Required, String, NonUpdatable) Specifies the resource pool name to which the resource nodes belong.

* `node_names` - (Required, List, NonUpdatable) Specifies the name list of resource nodes to be reset.

* `rolling_config` - (Required, List, NonUpdatable) Specifies the rolling configuration for the node reset operation.  
  The [rolling_config](#modelartsv2_node_batch_reset_rolling_config) structure is documented below.

<a name="modelartsv2_node_batch_reset_rolling_config"></a>
The `rolling_config` block supports:

* `strategy` - (Required, String, NonUpdatable) Specifies the rolling strategy for the reset operation.  
  The valid values are as follows:
  + **RollingByNumber**: Sets the maximum number of nodes to reset simultaneously based on the number of nodes.
  + **RollingByPercent**: Sets the maximum number of nodes to reset simultaneously based on a percentage.

* `max_unavailable` - (Required, Int, NonUpdatable) Specifies the maximum number or percentage of nodes that can be
  reset simultaneously.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
