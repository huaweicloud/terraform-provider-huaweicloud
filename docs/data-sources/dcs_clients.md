---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_clients"
description: |-
  Use this data source to query the list of clients for a specified DCS instance within HuaweiCloud.
---

# huaweicloud_dcs_clients

Use this data source to query the list of clients for a specified DCS instance within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "node_id" {}

data "huaweicloud_dcs_clients" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
}
```

### Query clients with filters

```hcl
variable "instance_id" {}
variable "node_id" {}

data "huaweicloud_dcs_clients" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
  addr        = "127.0.0.1"
  sort        = "age"
  order       = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to query the DCS clients. If omitted, the
  provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DCS instance.

* `node_id` - (Required, String) Specifies the ID of the node.
  For read/write splitting or proxy cluster instances, use the proxy node ID.
  For single-node, master-standby, or cluster instances, use the data node ID.

* `addr` - (Optional, String) Specifies the client address to filter by.

* `sort` - (Optional, String) Specifies the field to sort the client list by.
  Valid values are any field name in the `clients` block.

* `order` - (Optional, String) Specifies the sorting order.
  Valid values are **asc** (ascending) and **desc** (descending).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `clients` - The list of clients.

  The [clients](#clients_struct) structure is documented below.

<a name="clients_struct"></a>
The `clients` block supports:

* `id` - The client ID.

* `addr` - The address and port of the client.

* `fd` - The file descriptor used within the container.

* `name` - The name of the client.

* `cmd` - The last command executed by the client.

* `age` - The duration of the connection in seconds.

* `idle` - The idle time in seconds.

* `db` - The database ID currently being used by the client.

* `flags` - The client flags.

* `sub` - The number of subscribed channels.

* `psub` - The number of subscribed patterns.

* `multi` - The number of commands executed in a transaction.

* `qbuf` - The length of the query buffer in bytes.

* `qbuf_free` - The remaining space in the query buffer in bytes.

* `obl` - The length of the output buffer in bytes.

* `oll` - The number of objects in the output list.

* `omem` - The total memory used by the output buffer and list.

* `events` - The file operation events (r: read, w: write).

* `network` - The network type used by the client.

* `peer` - The client address and port.

* `user` - The client user.
