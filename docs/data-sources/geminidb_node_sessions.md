---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_node_sessions"
description: |-
  Use this data source to get the list of node sessions.
---

# huaweicloud_geminidb_node_sessions

Use this data source to get the list of node sessions.

-> This data source only supports GeminiDB Redis instance.

## Example Usage

```hcl
variable "node_id" {}

data "huaweicloud_geminidb_node_sessions" "test" {
  node_id = var.node_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `node_id` - (Required, String) Specifies the node ID.

* `addr_prefix` - (Optional, String) Specifies the matching character string of a user address prefix.
  The address consists of an IP address and a port number. e.g. **192.168.1.20:8080**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `sessions` - The list of node sessions.
  The [sessions](#sessions_struct) structure is documented below.

<a name="sessions_struct"></a>
The `sessions` block supports:

* `id` - The client ID.

* `name` - The client name.

* `cmd` - The last executed command.

* `age` - The setup duration of the client connection, in seconds.

* `idle` - The idle duration of the client connection, in seconds.

* `db` - The ID of the currently accessed database.

* `addr` - The IP address and port number of the client.

* `fd` - The file descriptor for sockets.

* `sub` - The number of subscribed channels (Pub/Sub).

* `psub` - The number of subscribed channels (Pub/Sub) in batches.

* `multi` - The number of commands contained in a MULTI or EXEC transaction.
