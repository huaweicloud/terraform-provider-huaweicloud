---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_collector_channel_operation"
description: |-
  Manages a SecMaster collector channel operation resource within HuaweiCloud.
---

# huaweicloud_secmaster_collector_channel_operation

Manages a SecMaster collector channel operation resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not change the current status.

## Example Usage

```hcl
variable "workspace_id" {}
variable "channel_id" {}
variable "action" {}

resource "huaweicloud_secmaster_collector_channel_operation" "test" {
  workspace_id = var.workspace_id
  channel_id   = var.channel_id
  action       = var.action
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the workspace ID.

* `channel_id` - (Required, String, NonUpdatable) Specifies the collector channel ID.

* `action` - (Required, String, NonUpdatable) Specifies the operation action of the collector channel.
  The valid values are as follows:
  + **START**: Start.
  + **STOP**: Stop.
  + **REMOVE**: Remove.
  + **RESTART**: Restart.
  + **REFRESH**: Refresh.
  + **INSTALL**: Install.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
