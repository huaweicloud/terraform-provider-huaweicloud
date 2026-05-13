---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_sessions_kill"
description: |-
  Use this resource to kill specific client sessions for a specified DCS instance within HuaweiCloud.
---

# huaweicloud_dcs_sessions_kill

Use this resource to kill specific client sessions for a specified DCS instance within HuaweiCloud.

## Example Usage

### Kill specific client sessions

```hcl
variable "instance_id" {}
variable "node_id" {}

resource "huaweicloud_dcs_kill_clients" "test" {
  instance_id  = var.instance_id
  node_id      = var.node_id
  client_addrs = ["127.0.0.1:6379", "192.168.0.10:54321"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DCS instance.

* `node_id` - (Optional, String, NonUpdatable) Specifies the ID of the node.

* `client_addrs` - (Optional, List, NonUpdatable) Specifies the list of client addresses to be killed. The value is the
  address and port of the peer in the session information.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
