---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_instance_restart"
description: |-
  Restart a DMS RocketMQ instance resource within HuaweiCloud.
---

# huaweicloud_dms_rocketmq_instance_restart

Restart a DMS RocketMQ instance resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dms_rocketmq_instance_nodes" "test" {
  instance_id = var.instance_id
}

resource "huaweicloud_dms_rocketmq_instance_restart" "test" {
  instance_id = var.instance_id

  nodes = [
    data.huaweicloud_dms_rocketmq_instance_nodes.test.nodes[0].id,
    data.huaweicloud_dms_rocketmq_instance_nodes.test.nodes[1].id
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to restart the RocketMQ instance.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RocketMQ instance to be restarted.  

* `nodes` - (Required, List, NonUpdatable) Specifies the list of node IDs to be restarted.  

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The RocketMQ instance restart resource cannot be imported.

## Timeouts

This resource does not support timeout configuration.
