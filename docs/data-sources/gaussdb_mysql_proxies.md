---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_proxies"
description: |-
  Use this data source to get the list of GaussDB MySQL proxies.
---

# huaweicloud_gaussdb_mysql_proxies

Use this data source to get the list of GaussDB MySQL proxies.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_mysql_proxies" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB MySQL instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `proxy_list` - Indicates the list of proxies.

  The [proxy_list](#proxy_list_struct) structure is documented below.

<a name="proxy_list_struct"></a>
The `proxy_list` block supports:

* `id` - Indicates the proxy ID.

* `name` - Indicates the proxy name.

* `flavor` - Indicates the flavor of the proxy.

* `port` - Indicates the proxy port.

* `status` - Indicates the status of the proxy instance.

* `delay_threshold_in_seconds` - Indicates the delay threshold in seconds.

* `node_num` - Indicates the number of proxy nodes.

* `ram` - Indicates the memory size of the proxy.

* `connection_pool_type` - Indicates the connection pool type.

* `switch_connection_pool_type_enabled` - Indicates whether the proxy version supports session-level connection pool.

* `mode` - Indicates the proxy mode.

* `elb_vip` - Indicates the virtual IP address in ELB mode.

* `vcpus` - Indicates the number of vCPUs of the proxy.

* `transaction_split` - Indicates whether the proxy transaction splitting is enabled.

* `balance_route_mode_enabled` - Indicates whether the proxy version supports load balancing.

* `route_mode` - Indicates the routing policy of the proxy instance.

* `subnet_id` - Indicates the network ID of a subnet.

* `consistence_mode` - Indicates the consistency mode of the proxy.

* `ssl_option` - Indicates whether to enable or disable SSL.

* `new_node_auto_add_status` - Indicates whether new nodes are automatically associate with proxy.

* `new_node_weight` - Indicates the read weight of the new node.

* `address` - Indicates the address of the proxy.

* `nodes` - Indicates the node information of the proxy.

  The [nodes](#proxy_nodes_struct) structure is documented below.

* `master_node_weight` - Indicates the read weight of the master node.

  The [master_node_weight](#proxy_list_master_node_weight_struct) structure is documented below.

* `readonly_nodes_weight` - Indicates the read weight of the read-only node.

  The [readonly_nodes_weight](#proxy_list_readonly_nodes_weight_struct) structure is documented below.

<a name="proxy_nodes_struct"></a>
The `nodes` block supports:

* `id` - Indicates the proxy node ID.

* `name` - Indicates the proxy node name.

* `role` - Indicates the proxy node role.

* `az_code` - Indicates the proxy node AZ.

* `status` - Indicates the proxy node status.

* `frozen_flag` - Indicates whether the proxy node is frozen.

<a name="proxy_list_master_node_weight_struct"></a>
The `master_node_weight` block supports:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `weight` - Indicates the weight assigned to the node.

<a name="proxy_list_readonly_nodes_weight_struct"></a>
The `readonly_nodes_weight` block supports:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `weight` - Indicates the weight assigned to the node.
