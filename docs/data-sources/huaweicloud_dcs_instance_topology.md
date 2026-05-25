---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_instance_topology"
description: |-
  Use this data source to query the topology information of DCS instance.
---

# huaweicloud_dcs_instance_topology

Use this data source to query the topology information of DCS instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dcs_instance_topology" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the instance topology.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the DCS instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `redis_server` - The topology information of redis-server nodes.
  The [redis_server](#topology_info_struct) structure is documented below.

* `cluster_lvs` - The topology information of lvs nodes (deprecated, only for Redis 3.0).
  The [cluster_lvs](#topology_info_struct) structure is documented below.

* `cluster_admin` - The topology information of admin nodes (deprecated, only for Redis 3.0).
  The [cluster_admin](#topology_info_struct) structure is documented below.

* `cluster_proxy` - The topology information of proxy nodes.
  The [cluster_proxy](#topology_info_struct) structure is documented below.

* `master` - The topology information of master nodes (only for Redis 3.0).
  The [master](#topology_info_struct) structure is documented below.

* `vpc_endpoint` - The topology information of VPC Endpoint nodes (only for 4.0 and later).
  The [vpc_endpoint](#topology_info_struct) structure is documented below.

* `elb` - The topology information of ELB nodes (only for Redis 4.0 and later).
  The [elb](#topology_info_struct) structure is documented below.

<a name="topology_info_struct"></a>
The `redis_server`, `cluster_lvs`, `cluster_admin`, `cluster_proxy`, `master`, `vpc_endpoint`, `elb` block supports:

* `node_id` - The node ID.

* `node_name` - The node name.

* `ip` - The node IP address.

* `port` - The node port.

* `node_type` - The node role type. Valid values are **master**, **slave**, **proxy**.

* `max_memory` - The total memory in MB.

* `used_memory` - The used memory in MB.

* `qps` - The queries per second.

* `bandwidth` - The bandwidth information.
  The [bandwidth](#bandwidth_struct) structure is documented below.

* `db_num` - The number of databases.

* `dbs` - The keyspace information.
  The [dbs](#dbs_struct) structure is documented below.

* `relation_ip` - The associated IP address.

* `relation_port` - The associated port.

* `group_id` - The shard ID.

* `status` - The node status.

* `dims` - The CES monitoring dimension information.
  The [dims](#dims_info) structure is documented below.

<a name="bandwidth_struct"></a>
The `bandwidth` block supports:

* `input` - The upstream bandwidth in Kbit/s.

* `output` - The downlink bandwidth in Kbit/s.

<a name="dbs_struct"></a>
The `dbs` block supports:

* `db_idx` - The database index.

* `keys` - The number of node keys.

* `expires` - The number of expired node keys.

* `avg_ttl` - The average expiration time of node keys.

<a name="dims_info"></a>
The `dims` block supports:

* `dim_k` - The cloud eye monitoring dimension route.

* `dim_v` - The cloud eye monitoring dimension ID.
