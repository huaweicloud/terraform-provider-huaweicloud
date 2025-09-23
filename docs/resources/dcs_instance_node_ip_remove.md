---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_instance_node_ip_remove"
description: |-
  Manages a DCS instance node IP remove resource within HuaweiCloud.
---

# huaweicloud_dcs_instance_node_ip_remove

Manages a DCS instance node IP remove resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "group_id" {}
variable "node_id" {}

resource "huaweicloud_dcs_instance_node_ip_remove" "test"{
  instance_id = var.instance_id
  group_id    = var.group_id
  node_id     = var.node_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DCS instance.

* `group_id` - (Required, String, NonUpdatable) Specifies the ID of the shard.

* `node_id` - (Required, String, NonUpdatable) Specifies the ID of the node that is promoted to master.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is formatted `<instance_id>/<group_id>/<node_id>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
