---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_node_status_change"
description: |-
  Manages a DCS node status change resource within HuaweiCloud.
---

# huaweicloud_dcs_node_status_change

Manages a DCS node status change resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_ids" {}

resource "huaweicloud_dcs_node_status_change" "test" {
  instance_id = var.instance_id
  node_ids    = var.node_ids
  action      = "stop"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. This parameter is non-updatable.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DCS instance.

* `node_ids` - (Optional, List of String, NonUpdatable) Specifies the list of node IDs to start or stop.
  If this parameter is not configured, all nodes of the instance will be started or stopped by default.

* `action` - (Required, String, NonUpdatable) Specifies the operation to perform on the instance nodes.
  The valid values are as follows:
  + **start**: Start the nodes.
  + **stop**: Stop the nodes.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
