---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_mysql_proxy"
description: |-
  Manages RDS mysql proxy resource within HuaweiCloud.
---

# huaweicloud_rds_mysql_proxy

Manages RDS mysql proxy resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "replica_node_id_1" {}
variable "replica_node_id_2" {}

resource "huaweicloud_rds_mysql_proxy" "test" {
  instance_id = var.instance_id
  flavor      = "rds.proxy.large.2"
  node_num    = 3
  route_mode  = 0

  master_node_weight {
    id     = var.instance_id
    weight = 10
  }

  readonly_nodes_weight {
    id     = var.replica_node_id_1
    weight = 20
  }

  readonly_nodes_weight {
    id     = var.replica_node_id_2
    weight = 30
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the RDS MySQL proxy resource. If omitted,
  the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS MySQL instance.

* `flavor` - (Required, String, NonUpdatable) Specifies the flavor of the proxy.
  + When the site supports the proxy in primary/standby mode, this parameter does not take effect.

* `node_num` - (Required, Int, NonUpdatable) Specifies the node number of the proxy nodes.
  + When the site supports the proxy in primary/standby mode, set this parameter to **2**.
  + When the site supports the proxy in cluster mode, the minimum value of this parameter is **2**.

* `proxy_name` - (Optional, String, NonUpdatable) Specifies the name of the proxy. The name must start with a letter and
  consist of **4** to **64** characters. Only letters, digits, hyphens (-), underscores (_), and periods (.) are allowed.

* `proxy_mode` - (Optional, String, NonUpdatable) Specifies the read/write mode of the proxy. Value options:
  + **readwrite(default value)**: read and write.
  + **readonly**: read-only.

* `subnet_id` - (Optional, String, NonUpdatable) Specifies the network ID of a subnet.

* `route_mode` - (Optional, Int) Specifies the routing policy of the proxy. Value options:
  + **0**: weighted load balancing.
  + **1**: load balancing (The primary node does not process read requests).
  + **2**: load balancing (The primary node processes read requests).

* `master_node_weight` - (Optional, List) Specifies the read weight of the master node.
  The [master_node_weight](#node_weight_struct) structure is documented below.

* `readonly_nodes_weight` - (Optional, List) Specifies the read weight of the read-only node.
  The [readonly_nodes_weight](#node_weight_struct) structure is documented below.

<a name="node_weight_struct"></a>
The `master_node_weight` and `readonly_nodes_weight` block supports:

* `id` - (Required, String) Specifies the ID of the node.

* `weight` - (Required, Int) Specifies the weight assigned to the node.
  + If `route_mode` is `0`, the value is `0` to `1,000`.
  + If `route_mode` is `1`, the value for the primary node is `0` and the value for read replicas is `0` or `1`.
  + If `route_mode` is `2`, the value for the primary node is `1` and the value for read replicas is `0` or `1`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID.

* `status` - Indicates the status of the proxy.

* `address` - Indicates the read/write splitting address of the proxy.

* `port` - Indicates the port number.

* `delay_threshold_in_seconds` - Indicates the delay threshold, in seconds.

* `vcpus` - Indicates the vCPUs of the proxy.

* `memory` - Indicates the memory size of the proxy.

* `nodes` - Indicates the list of proxy nodes.
  The [nodes](#nodes_struct) structure is documented below.

* `mode` - Indicates the proxy mode. The value can be: **Cluster**, **Ha**.

* `flavor_group_type` - Indicates the CPU architecture. The value can be: **X86**, **ARM**.

* `transaction_split` - Indicates the status of transaction splitting for the proxy.

* `connection_pool_type` - Indicates the connection pool type. The value can be:
  + **CLOSED**: The connection pool is closed.
  + **SESSION**: The session-level connection pool is enabled.

* `pay_mode` - Indicates the billing mode of the proxy. The value can be:
  + **0**: pay-per-use billing.
  + **1**: yearly/monthly billing.

* `dns_name` - Indicates the private domain name for the read/write splitting address of the proxy.

* `seconds_level_monitor_fun_status` - Indicates the status of monitoring by seconds of the proxy. The value can be:
  **on**, **off**.

* `alt_flag` - Indicates the ALT switch status.

* `force_read_only` - Indicates whether to forcibly route read requests to read replicas.

* `ssl_option` - Indicates the SSL switch status.

* `support_balance_route_mode` - Indicates whether load balancing can be enabled for the proxy.

* `support_proxy_ssl` - Indicates whether SSL can be enabled for the proxy.

* `support_switch_connection_pool_type` - Indicates whether the session connection pool type can be changed for the proxy.

* `support_transaction_split` - Indicates whether transaction splitting can be enabled for the proxy.

<a name="nodes_struct"></a>
The `nodes` block supports:

* `id` - Indicates the proxy node ID.

* `status` - Indicates the proxy node status. The values can be:
  + **NORMAL**: The node is normal.
  + **ABNORMAL**: The node is abnormal.
  + **CREATING**: The node is being created.
  + **CREATEFAIL**: The node failed to be created.

* `role` - Indicates the role of the proxy node. The values can be:
  + **master**: primary node.
  + **slave**: standby node.

* `az_code` - Indicates the AZ where the proxy node is located.

* `frozen_flag` - Indicates whether the proxy node is frozen. The values can be:
  + **0**: unfrozen.
  + **1**: frozen.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 10 minutes.

## Import

The RDS MySQL proxy can be imported using the `instance_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_rds_mysql_proxy.test <instance_id>/<id>
```
