---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_datasource_connections"
description: ""
---

# huaweicloud_dli_datasource_connections

Use this data source to get the list of DLI enhanced datasource connections within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}

data "huaweicloud_dli_datasource_connections" "test" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the datasource connection.

* `tags` - (Optional, Map) Specifies the key/value pairs of the datasource connections.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `connections` - The list of the datasource connections.  
  The [connections](#datasource_connections) structure is documented below.

<a name="datasource_connections"></a>
The `connections` block supports:

* `id` - The ID of the datasource connection.

* `created_at` - The creation time of the datasource connection.

* `hosts` - List of the user-defined hosts information.
  The [hosts](#datasource_connections_hosts) structure is documented below.

* `is_privis` - Whether the data source connection has setted project permissions. The valid values are as follows:
  + **false**: Indicates that the project permission is granted.
  + **true**: Indicates that the project permission is not granted.

* `name` - The name of the datasource connection.

* `queues` - List of queues associated with the datasource connection.
  The [queues](#datasource_connections_queues) structure is documented below.

* `elastic_resource_pools` - List of resource pools associated with the datasource connection.
  The [elastic_resource_pools](#datasource_connections_resource_pools) structure is documented below.

* `routes` - List of routes.
  The [routes](#datasource_connections_routes) structure is documented below.

* `status` - The current status of the datasource connection.

* `subnet_id` - The subnet ID associated with the datasource connection.

* `vpc_id` - The VPC ID associated with the datasource connection.

<a name="datasource_connections_hosts"></a>
The `hosts` block supports:

* `name` - The user-defined host name.

* `ip` - IPv4 address of the host.

<a name="datasource_connections_queues"></a>
The `queues` block supports:

* `id` - The peer ID of the datasource connection.

* `name` - The queue name.

* `status` - The status of the peering connection to which the datasource connection belongs.

* `error_msg` - Error message when the peering connection status is falied.

* `updated_at` - The latest update time of the queue.

<a name="datasource_connections_resource_pools"></a>
The `elastic_resource_pools` block supports:

* `id` - The peer ID of the datasource connection.

* `name` - The name of the elastic resource pool.

* `status` - The status of the peering connection to which the datasource connection belongs.

* `error_msg` - Error message when the Peering connection status status is falied.

* `updated_at` - The latest update time of the elastic resource pool.

<a name="datasource_connections_routes"></a>
The `routes` block supports:

* `name` - The route name.

* `cidr` - The CIDR of the route.

* `created_at` - The creation time of the route.
