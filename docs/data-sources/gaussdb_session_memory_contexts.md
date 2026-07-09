---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_session_memory_contexts"
description: |-
  Use this data source to query the session memory context list of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_session_memory_contexts

Use this data source to query the session memory context list of a GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "session_id" {}

data "huaweicloud_gaussdb_session_memory_contexts" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
  session_id  = var.session_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the session memory contexts.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `node_id` - (Required, String) Specifies the node ID. Only nodes containing CN or DN (primary/standby) components are
  supported.

* `session_id` - (Required, String) Specifies the session ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `memory_context_info` - The list of session memory context information.
  The [memory_context_info](#memory_context_info) structure is documented below.

<a name="memory_context_info"></a>
The `memory_context_info` block supports:

* `context_name` - The memory context name.

* `amount` - The number of session contexts.

* `size` - The total size of the session memory context.
