---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_mysql_proxies"
description: |-
  Use this data source to get the proxy list of a RDS MySQL instance.
---

# huaweicloud_rds_mysql_proxies

Use this data source to get the proxy list of a RDS MySQL instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_mysql_proxies" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of RDS MySQL instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `proxy_list` - Indicates the proxy information list of a RDS MySQL instance.

  The [proxy_list](#proxy_list_struct) structure is documented below.

* `max_proxy_num` - Indicates the maximum number of database proxies that can be enabled at the same time.

* `max_proxy_node_num` - Indicates the maximum number of proxy nodes that can be selected for a database proxy.

* `support_balance_route_mode_for_favored_version` - Indicates whether the load balancing routing mode can be set when
  the database proxy is created.

<a name="proxy_list_struct"></a>
The `proxy_list` block supports:

* `proxy` - Indicates the proxy information.

  The [proxy](#proxy_list_proxy_struct) structure is documented below.

* `master_instance` - Indicates the master instance information.

  The [master_instance](#proxy_list_master_instance_struct) structure is documented below.

* `readonly_instances` - Indicates the read-only instance information.

  The [readonly_instances](#proxy_list_readonly_instances_struct) structure is documented below.

* `proxy_security_group_check_result` - Indicates whether the security group allows access from the database proxy
  to the database.

<a name="proxy_list_proxy_struct"></a>
The `proxy` block supports:

* `id` - Indicates the proxy ID.

* `name` - Indicates the proxy name.

* `node_num` - Indicates the number of proxy nodes.

* `flavor_info` - Indicates the proxy specifications.

  The [flavor_info](#proxy_flavor_info_struct) structure is documented below.

* `proxy_mode` - Indicates the Proxy read/write Mode.
  The value can be:
  + **readwrite(default value)**: read and write.
  + **readonly**: read-only.

* `subnet_id` - Indicates the ID of the subnet to which the database proxy belongs.

* `route_mode` - Indicates the routing policy of the proxy.
  The values can be:
  + **0**: weighted load balancing.
  + **1**: load balancing (The primary node does not process read requests).
  + **2**: load balancing (The primary node processes read requests).

* `status` - Indicates the status of the proxy.

* `address` - Indicates the proxy address.

* `port` - Indicates the port number.

* `delay_threshold_in_seconds` - Indicates the delay threshold, in seconds.

* `vcpus` - Indicates the CPU size of the proxy.

* `memory` - Indicates the memory size of the proxy.

* `nodes` - Indicates the list of proxy nodes.

  The [nodes](#proxy_nodes_struct) structure is documented below.

* `mode` - Indicates the cluster mode of the proxy.
  The value can be: **Cluster**, **Ha**.

* `transaction_split` - Indicates the status of the proxy transaction splitting switch.

* `connection_pool_type` - Indicates the connection pool type.
  The value can be:
  + **CLOSED**: The connection pool is closed.
  + **SESSION**: The session-level connection pool is enabled.

* `pay_mode` - Indicates the charging mode of the proxy.
  The value can be:
  + **0**: pay-per-use billing.
  + **1**: yearly/monthly billing.

* `dns_name` - Indicates the private domain name for the read/write splitting address of the proxy.

* `seconds_level_monitor_fun_status` - Indicates the second-level monitoring status of the proxy.

* `alt_flag` - Indicates the ALT switch status.

* `force_read_only` - Indicates whether to forcibly read the route to the read-only mode.

* `ssl_option` - Indicates the SSL switch status.

* `support_balance_route_mode` - Indicates whether the proxy supports the load balancing routing mode.

* `support_proxy_ssl` - Indicates whether the database proxy supports the SSL function.

* `support_switch_connection_pool_type` - Indicates whether the proxy supports the switchover of the session
  connection pool type.

* `support_transaction_split` - Indicates whether the proxy supports transaction splitting.

<a name="proxy_flavor_info_struct"></a>
The `flavor_info` block supports:

* `group_type` - Indicates the flavor group type.

* `code` - Indicates the specification code.

<a name="proxy_nodes_struct"></a>
The `nodes` block supports:

* `id` - Indicates the proxy node ID.

* `status` - Indicates the status of the proxy node.
  The values can be:
  + **NORMAL**: The node is normal.
  + **ABNORMAL**: The node is abnormal.
  + **CREATING**: The node is being created.
  + **CREATEFAIL**: The node failed to be created.

* `role` - Indicates the role of the proxy node:
  The values can be:
  + **master**: primary node.
  + **slave**: standby node.

* `az_code` - Indicates the AZ where the proxy node is located.

* `frozen_flag` - Indicates whether the proxy node is frozen.
  The values can be:
  + **0**: unfrozen.
  + **1**: frozen.

<a name="proxy_list_master_instance_struct"></a>
The `master_instance` block supports:

* `id` - Indicates the instance ID.

* `weight` - Indicates the read weight of a the instance.

<a name="proxy_list_readonly_instances_struct"></a>
The `readonly_instances` block supports:

* `id` - Indicates the instance ID.

* `weight` - Indicates the read weight of the instance.
