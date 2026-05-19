---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_node_session_kill"
description: |-
  Manages a resource to kill node session within HuaweiCloud.
---

# huaweicloud_geminidb_node_session_kill

Manages a resource to kill node session within HuaweiCloud.

-> 1. This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.
  <br/>2. This resource only supports GeminiDB Redis instance.

## Example Usage

```hcl
variable "node_id" {}
variable "session_id_list" {
  type = list(string)
}

resource "huaweicloud_geminidb_node_session_kill" "test" {
  node_id     = var.node_id
  is _all     = false
  session_ids = var.session_id_list
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `node_id` - (Required, String, NonUpdatable) Specifies the node ID.

* `is_all` - (Required, Bool, NonUpdatable) Specifies whether all sessions are closed.
  The valid values are as follows:
  + **true**: All sessions are closed.
  + **false**: Not all sessions are closed.

* `session_ids` - (Optional, List, NonUpdatable) Specifies the ID of the session to be closed.
  When the parameter `is_all` is set to **false**, this parameter is required.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
