---
subcategory: "ModelArts"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_node_batch_update"
description: |-
  Use this resource to batch update nodes in a resource pool within HuaweiCloud.
---

# huaweicloud_modelartsv2_node_batch_update

Use this resource to batch update nodes in a resource pool within HuaweiCloud.

-> This resource is a one-time action resource for batch update the ModelArts nodes. Deleting this resource will not
   clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Open high availability redundancy for nodes

```hcl
variable "pool_id" {}
variable "node_name" {}

resource "huaweicloud_modelartsv2_node_batch_update" "test" {
  pool_id    = var.pool_id
  node_names = [var.node_name]
  action     = "openHaRedundant"
}
```

### Batch add tags to nodes

```hcl
variable "pool_id" {}
variable "node_names" {
  type = list(string)
}

resource "huaweicloud_modelartsv2_node_batch_update" "test" {
  pool_id    = var.pool_id
  node_names = var.node_names
  action     = "createTags"

  tags {
    key   = "key1"
    value = "value1"
  }

  tags {
    key   = "key2"
    value = "value2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the resource nodes are located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `pool_id` - (Required, String, NonUpdatable) Specifies the ID of the resource pool to which the resource nodes belong.

* `node_names` - (Required, List, NonUpdatable) Specifies the name list of resource nodes to be updated.

* `action` - (Required, String, NonUpdatable) Specifies the type of the node update action.  
  The valid values are as follows:
  + **openHaRedundant**: Open high availability redundancy for nodes.
  + **closeHaRedundant**: Close high availability redundancy for nodes.
  + **createTags**: Batch add tags to nodes.
  + **deleteTags**: Batch delete tags from nodes.

* `ha_redundant_effect` - (Optional, String) Specifies the effect of the high availability redundancy. Defaults to
  **NoSchedule**.  
  The valid values are as follows:
  + **NoSchedule**: Prohibit scheduling.
  + **NoExecute**: Prohibit execution.

  -> This parameter is only used for idle **NPU** and **GPU** nodes.

* `driver` - (Optional, List) Specifies the driver version and status information of the node.  
  The [driver](#v2_node_batch_update_driver) structure is documented below.

* `tags` - (Optional, List) Specifies the list of resource tags to be operated in batch.  
  The [tags](#v2_node_batch_update_tags) structure is documented below.

<a name="v2_node_batch_update_driver"></a>
The `driver` block supports:

* `version` - (Optional, String) Specifies the version of the driver on the node.

* `update_strategy` - (Optional, String) Specifies the driver upgrade strategy of the node.

<a name="v2_node_batch_update_tags"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the key of the tag.

* `value` - (Required, String) Specifies the value of the tag.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
