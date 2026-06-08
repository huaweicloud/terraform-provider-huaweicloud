---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_enlarge_fail_node_delete"
description: |-
  Manages a resource to delete GeminiDB instance enlarge failed node within HuaweiCloud.
---

# huaweicloud_geminidb_enlarge_fail_node_delete

Manages a resource to delete GeminiDB instance enlarge failed node within HuaweiCloud.

-> 1. This resource is a one-time action resource. Deleting this resource will not clear the corresponding request
  record, but will only remove the resource information from the tf state file.
  <br/>2. The instance types that support the database patch upgrade are as follows:
  **GeminiDB Cassandra**, **GeminiDB Redis**.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}

resource "huaweicloud_geminidb_enlarge_fail_node_delete" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the instance ID.

* `node_id` - (Required, String, NonUpdatable) Specifies the node ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
