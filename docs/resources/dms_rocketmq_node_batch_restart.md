---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_node_batch_restart"
description: |-
  Using this resource to batch restart nodes of a RocketMQ instance within HuaweiCloud.
---

# huaweicloud_dms_rocketmq_node_batch_restart

Using this resource to batch restart nodes of a RocketMQ instance within HuaweiCloud.

-> This resource is only a one-time action resource for restarting nodes of RocketMQ instance. Deleting this resource
will not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dms_rocketmq_instance_nodes" "test" {
  instance_id = var.instance_id
}

resource "huaweicloud_dms_rocketmq_node_batch_restart" "test" {
  instance_id = var.instance_id

  nodes = [
    data.huaweicloud_dms_rocketmq_instance_nodes.test.nodes[0].id,
    data.huaweicloud_dms_rocketmq_instance_nodes.test.nodes[1].id
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to restart the nodes of a RocketMQ instance.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RocketMQ instance to which nodes belong.

* `nodes` - (Required, List, NonUpdatable) Specifies the list of node IDs to be restarted.  

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
