---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_instance_sessions"
description: |-
  Use this data source to get the list of instance sessions.
---

# huaweicloud_geminidb_instance_sessions

Use this data source to get the list of instance sessions.

-> This data source only supports GeminiDB Redis instances.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_geminidb_instance_sessions" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `node_id` - (Optional, String) Specifies the node ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `node_sessions` - The list of instance sessions.
  The [node_sessions](#instacne_sessions_struct) structure is documented below.

<a name="instacne_sessions_struct"></a>
The `node_sessions` block supports:

* `node_id` - The node ID.

* `sessions` - The list of node sessions.
  The [sessions](#sessions_struct) structure is documented below.

<a name="sessions_struct"></a>
The `sessions` block supports:

* `id` - The session ID.

* `name` - The connection name.

* `cmd` - The last executed command.

* `age` - The connection duration, in seconds.

* `idle` - The idle duration, in seconds.

* `db` - The ID of the database that is being used by the client.

* `addr` - The IP address and port of the client.

* `fd` - The file descriptor for sockets.

* `sub` - The number of subscribed channels.

* `psub` - The number of subscribed modes.

* `multi` - The number of commands executed in a transaction.
