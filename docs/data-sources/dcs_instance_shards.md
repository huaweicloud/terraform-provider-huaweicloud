---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_instance_shards"
description: |-
  Use this data source to get the information about shards and replicas.
---

# huaweicloud_dcs_instance_shards

Use this data source to get the information about shards and replicas.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dcs_instance_shards" "test" {
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

* `group_list` - Indicates the shard list.

  The [group_list](#group_list_struct) structure is documented below.

<a name="group_list_struct"></a>
The `group_list` block supports:

* `group_id` - Indicates the shard ID.

* `group_name` - Indicates the shard name.

* `replication_list` - Indicates the list of replicas in the shard.

  The [replication_list](#group_list_replication_list_struct) structure is documented below.

<a name="group_list_replication_list_struct"></a>
The `replication_list` block supports:

* `replication_role` - Indicates the role of the replica.
  The value can be:
  + **master**: master
  + **slave**: replication

* `replication_ip` - Indicates the replica IP address.

* `is_replication` - Indicates whether the replica is a newly added one.

* `replication_id` - Indicates the replica ID.

* `node_id` - Indicates the node ID.

* `status` - Indicates the replica status.

* `az_code` - Indicates the AZ where the replica is in.

* `dimensions` - Indicates the monitoring metric dimension of the replica.
  It is used to call the Cloud Eye API for querying monitoring metrics.
  + Replica monitoring is multi-dimensional. The returned array contains information about two dimensions. When querying
    monitoring data from Cloud Eye, you need to transfer parameters of multiple dimensions to obtain the metric data.
  + The first dimension is the primary dimension of the replica. The dimension name is **dcs_instance_id**, and the
    dimension value corresponds to the ID of the instance to which the replica belongs.
  + The name of the second dimension is **dcs_cluster_redis_node**, and the dimension value is the ID of the monitored
    object of the replica, which is different from the replica ID or node ID.

  The [dimensions](#replication_list_dimensions_struct) structure is documented below.

<a name="replication_list_dimensions_struct"></a>
The `dimensions` block supports:

* `name` - Indicates the monitoring dimension name.

* `value` - Indicates the dimension value.
