---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_connection_statistics"
description: |-
  Use this data source to get the statistics of the connection number of the instance nodes.
---

# huaweicloud_dds_connection_statistics

Use this data source to get the statistics of the connection number of the instance nodes.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dds_connection_statistics" "test"{
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

* `total_connections` - The total number of connections.

* `total_inner_connections` - The total number of internal connections.

* `total_outer_connections` - The total number of external connections.

* `inner_connections` - The internal connection statistics.

  The [inner_connections](#inner_connections_struct) structure is documented below.

* `outer_connections` - The external connection statistics.

  The [outer_connections](#outer_connections_struct) structure is documented below.

<a name="inner_connections_struct"></a>
The `inner_connections` block supports:

* `client_ip` - The IP address of the client connected to this instance or node.

* `count` - The number of connections corresponding to this IP address.

<a name="outer_connections_struct"></a>
The `outer_connections` block supports:

* `client_ip` - The IP address of the client connected to this instance or node.

* `count` - The number of connections corresponding to this IP address.
