---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_specified_sessions_query_and_kill"
description: |-
  Use this resource to query and kill specified sessions of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_specified_sessions_query_and_kill

Use this resource to query and kill specified sessions of a GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "component_id" {}
variable "session_id" {}

resource "huaweicloud_gaussdb_specified_sessions_query_and_kill" "test" {
  instance_id  = var.instance_id
  node_id      = var.node_id
  component_id = var.component_id
  session_ids = [var.session_id]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to kill the sessions.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GaussDB instance.

* `node_id` - (Required, String, NonUpdatable) Specifies the node ID.
  Only nodes with CN or DN (primary, standby) components are supported.

* `component_id` - (Required, String, NonUpdatable) Specifies the component ID.
  Only CN or DN (primary, standby) component IDs on the node specified by `node_id` are supported.

* `session_ids` - (Required, List, NonUpdatable) Specifies the list of session IDs to be killed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `success` - Whether the kill session request is successful.

* `session_ids` - The list of successfully killed session IDs.
