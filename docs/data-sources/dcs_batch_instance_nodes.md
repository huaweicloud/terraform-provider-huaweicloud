---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_batch_instance_nodes"
description: |-
  Use this data source to query the node information, number of valid instances, and number of nodes of all instances.
---

# huaweicloud_dcs_batch_instance_nodes

Use this data source to query the node information, number of valid instances, and number of nodes of all instances.

## Example Usage

```hcl
data "huaweicloud_dcs_batch_instance_nodes" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the DCS instances are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - The instance list.  
  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `instance_id` - The instance ID.

* `node_count` - The total number of nodes in the current instance.

* `nodes` - The list of node details.  
  The [nodes](#nodes_struct) structure is documented below.

<a name="nodes_struct"></a>
The `nodes` block supports:

* `logical_node_id` - The logical node ID.

* `name` - The node name.

* `status` - The node status.  
  Valid values are:
  + **Creating**: Creating.
  + **Active**: Running.
  + **Inactive**: Fault.
  + **Deleting**: Deleting.
  + **AddSharding**: Adding sharding.

* `az_code` - The availability zone code.

* `node_role` - The node role.

* `node_type` - The node type.

* `node_ip` - The node IP address.

* `node_port` - The node port.

* `node_id` - The node ID.

* `priority_weight` - The replica promotion priority.

* `is_access` - Whether the IP address of the node can be directly accessed.

* `group_id` - The instance shard ID.

* `group_name` - The instance shard name.

* `is_remove_ip` - Whether the IP address is removed from the read-only domain name.

* `replication_id` - The Instance replica ID.

* `dimensions` - The monitoring metric dimension of the replica used to call the Cloud Eye API for querying monitoring
  metrics.
  + Replica monitoring is multi-dimensional. The returned array contains information about two dimensions. When querying
    monitoring data from Cloud Eye, transfer parameters of multiple dimensions to obtain the metric data.
  + The first dimension is the primary dimension of the replica. The dimension name is dcs_instance_id, and the dimension
    value is the ID of the instance where the replica is located.
  + The name of the second dimension is dcs_cluster_redis_node, and the dimension value is the ID of the monitored
    object of the replica, which is different from the replica ID or node ID.

  The [dimensions](#dimensions_struct) structure is documented below.

<a name="dimensions_struct"></a>
The `dimensions` block supports:

* `name` - The monitoring dimension name. The value can be:
  + **dcs_instance_id**: instance dimension
  + **dcs_cluster_redis_node**: data node dimension

* `value` - The dimension value.
