---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_instance_node_delete"
description: |-
  Manages a DDS instance node delete resource within HuaweiCloud.
---

# huaweicloud_dds_instance_node_delete

Manages a DDS instance node delete resource within HuaweiCloud.

-> 1. This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.
  <br/>2. This resource is available only to replica set instances.
  <br/>3. For a 7-node replica set instance, `2` or `4` standby nodes can be deleted.
  <br/>4. For a 5-node replica set instance, `2` standby nodes can be deleted.
  <br/>5. The standby node of a 3-node replica set instance cannot be deleted.
  <br/>6. Nodes cannot be deleted from instances that have abnormal nodes.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_list" {
  type = list(string)
}

resource "huaweicloud_dds_instance_node_delete" "test" {
  instance_id = var.instance_id
  node_list   = var.node_list
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DDS instance.

* `num` - (Optional, Int, NonUpdatable) Specifies the number of nodes to be deleted.

* `node_list` - (Optional, List, NonUpdatable) Specifies the ID list of nodes to be deleted.

-> 1. Either `num` or `node_list` must be set.
  <br/>2. If both `num` and `node_list` are specified, the value of `node_list` takes effect.
  <br/>3. The role of the node to be deleted cannot be Primary or Hidden.
  <br/>4. If there is a multi-AZ instance, ensure that at least one node is deployed in each AZ
  after specified nodes are deleted.

## Attribute Reference

* `id` - Indicates the resource ID. It's same as the instance ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
