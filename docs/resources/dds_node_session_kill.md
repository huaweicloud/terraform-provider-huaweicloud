---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_node_session_kill"
description: |-
  Manages a DDS instance node session kill resource within HuaweiCloud.
---

# huaweicloud_dds_node_session_kill

Manages a DDS instance node session kill resource within HuaweiCloud.

-> 1. This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.
  <br/>2. This resource is available for DB instances of community edition 3.4 or later.
  <br/>3. Inactive sessions cannot be terminated.

## Example Usage

```hcl
variable "node_id" {}
variable "session_id_list" {
  type = list(string)
}

resource "huaweicloud_dds_node_session_kill" "test" {
  node_id  = var.node_id
  sessions = var.session_id_list
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `node_id` - (Required, String, NonUpdatable) Specifies the node ID.
  For a cluster instance, you can select any mongos, shard, or config node.
  For a replica set instance, you can select the primary or secondary node.

* `sessions` - (Required, List, NonUpdatable) Specifies the IDs of sessions to be terminated.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
