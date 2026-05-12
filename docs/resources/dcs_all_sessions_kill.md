---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_all_sessions_kill"
description: |-
  Use this resource to kill all client sessions for a specified DCS instance or node within HuaweiCloud.
---

# huaweicloud_dcs_all_sessions_kill

Use this resource to kill all client sessions for a specified DCS instance or node within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}

resource "huaweicloud_dcs_all_sessions_kill" "test" {
  instance_id    = var.instance_id
  node_id        = var.node_id
  kill_all_nodes = "false"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DCS instance.

* `node_id` - (Optional, String, NonUpdatable) Specifies the ID of the node. Required when `kill_all_nodes` is not
  specified or set to **false**.

* `kill_all_nodes` - (Optional, String, NonUpdatable) Whether to kill all client sessions of all nodes in the instance.
  Valid values are **true** and **false**.
  When set to **true**, sessions of all nodes will be killed and `node_id` should not be set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
