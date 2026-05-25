---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_instance_restart"
description: |-
  Manages a resource to restart GeminiDB instance or node within HuaweiCloud.
---

# huaweicloud_geminidb_instance_restart

Manages a resource to restart GeminiDB instance or node within HuaweiCloud.

-> 1. This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.
  <br/>2. The instance types that support the instance restart operation are as follows:
  **GeminiDB Cassandra**, **GeminiDB Mongo**, **GeminiDB Influx** and **GeminiDB Redis**.
  <br/>3. The instance type that supports the node restart operation are as follows:
  GeminiDB Redis instance with cloud native storage.

## Example Usage

### Restart instance

```hcl
variable "instance_id" {}

resource "huaweicloud_geminidb_instance_restart" "test" {
  instance_id = var.instance_id
}
```

### Restart instance node

```hcl
variable "instance_id" {}
variable "node_id" {}

resource "huaweicloud_geminidb_instance_restart" "test" {
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

* `node_id` - (Optional, String, NonUpdatable) Specifies the ID of a node to be restarted.
  -> 1. Only GeminiDB Redis instances with cloud native storage can be restarted based on the node ID.
  <br/>2. If this parameter is not specified, the entire instance is restarted.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
