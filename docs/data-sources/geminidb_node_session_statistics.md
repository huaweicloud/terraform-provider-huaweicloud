---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_node_session_statistics"
description: |-
  Use this data source to get the list of node session statistics.
---

# huaweicloud_geminidb_node_session_statistics

Use this data source to get the list of node session statistics.

-> This data source only supports GeminiDB Redis instance.

## Example Usage

```hcl
variable "node_id" {}

data "huaweicloud_geminidb_node_session_statistics" "test" {
  node_id = var.node_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `node_id` - (Required, String) Specifies the node ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `top_source_ips` - The top 10 clients with the most connections.
  The [top_source_ips](#top_source_ips_struct) structure is documented below.

* `top_dbs` - The top 10 databases with the most connections.
  The [top_dbs](#top_dbs_struct) structure is documented below.

* `total_connection_count` - The total client connections.

* `active_connection_count` - The number of active client connections.

<a name="top_source_ips_struct"></a>
The `top_source_ips` block supports:

* `client_ip` - The client IP address.

* `connection_count` - The number of client connections.

<a name="top_dbs_struct"></a>
The `top_dbs` block supports:

* `db` - The database ID.

* `connection_count` - The number of client connections.
