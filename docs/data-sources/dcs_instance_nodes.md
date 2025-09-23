---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_instance_nodes"
description: |-
  Use this data source to get the nodes of a specified instance.
---

# huaweicloud_dcs_instance_nodes

Use this data source to get the nodes of a specified instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dcs_instance_nodes" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `nodes` - Indicates the node information.

  The [nodes](#nodes_struct) structure is documented below.

<a name="nodes_struct"></a>
The `nodes` block supports:

* `logical_node_id` - Indicates the logical node ID.

* `name` - Indicates the node name.

* `status` - Indicates the node status.
  The value can be:
  + **Creating**
  + **Active**
  + **Inactive**
  + **Deleting**
  + **AddSharding**

* `az_code` - Indicates the AZ code.

* `node_role` - Indicates the node role.
  The value can be:
  + **redis-server**: server node
  + **redis-proxy**: proxy node

* `node_type` - Indicates  the node master/standby role.
  The value can be:
  + **master**: master node
  + **slave**: standby node
  + **proxy**: node of a Proxy Cluster instance

* `node_ip` - Indicates the node IP.

* `node_port` - Indicates the node port.

* `node_id` - Indicates the node ID.

* `priority_weight` - Indicates  the replica promotion priority.
  The value can be:
  + Priority ranges from **0** to **100** in descending order.
  + **0** indicates that the replica will never be automatically promoted.
  + **1** indicates the highest priority.
  + **100** indicates the lowest priority.

* `is_access` - Indicates whether the IP address of the node can be directly accessed.

* `group_id` - Indicates the instance shard ID.

* `group_name` - Indicates the instance shard name.

* `is_remove_ip` - Indicates whether the IP address is removed from the read-only domain name.

* `replication_id` - Indicates the instance replica ID.

* `dimensions` - Indicates the monitoring metric dimension of the replica.
  It is used to call the Cloud Eye API for querying monitoring metrics.
  + Replica monitoring is multi-dimensional. The returned array contains information about two dimensions. When querying
    monitoring data from Cloud Eye, transfer parameters of multiple dimensions to obtain the metric data.
  + The first dimension is the primary dimension of the replica. The dimension name is dcs_instance_id, and the dimension
    value is the ID of the instance where the Indicates the replica is located.
  + The name of the second dimension is dcs_cluster_redis_node, and the dimension value is the ID of the monitored object
    of the replica, which is different from the replica ID or node ID.

  The [dimensions](#nodes_dimensions_struct) structure is documented below.

<a name="nodes_dimensions_struct"></a>
The `dimensions` block supports:

* `name` - Indicates the monitoring dimension name.
  The value can be:
  + **dcs_instance_id**: instance dimension
  + **dcs_cluster_redis_node**: data node dimension

* `value` - Indicates the dimension value.
