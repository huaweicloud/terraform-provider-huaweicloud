---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_cluster_instance_topology"
description: |-
  Use this data source to query the topology information of a DCS cluster instance within HuaweiCloud.
---

# huaweicloud_dcs_cluster_instance_topology

Use this data source to query the topology information of a DCS cluster instance within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dcs_cluster_instance_topology" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to query the topology. If omitted, the
  provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the DCS instance.

* `start` - (Optional, Int) Specifies the offset for pagination. This parameter indicates the starting position for
  generating the list. For example, if offset is set to 3, the list starts from the 4th item. Defaults to 0.

* `limit` - (Optional, Int) Specifies the number of items per page. Defaults to 10. Value range: 1 ~ 1000.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `redis_server` - The topology information for Redis server roles.
  The [redis_server](#redis_server_struct) structure is documented below.

* `cluster_proxy` - The topology information for Proxy roles.
  The [cluster_proxy](#redis_server_struct) structure is documented below.

* `vpcendpoint` - The topology information for VPC Endpoint service roles.
  The [vpcendpoint](#redis_server_struct) structure is documented below.

* `elb` - The topology information for ELB service roles.
  The [elb](#redis_server_struct) structure is documented below.

* `cluster_lvs` - The topology information for LVS roles (legacy, Redis 3.0 only).
  The [cluster_lvs](#redis_server_struct) structure is documented below.

* `cluster_admin` - The topology information for admin roles (legacy, Redis 3.0 only).
  The [cluster_admin](#redis_server_struct) structure is documented below.

* `master` - The topology information for master roles (legacy, Redis 3.0 only).
  The [master](#redis_server_struct) structure is documented below.

<a name="redis_server_struct"></a>
The `redis_server` , `cluster_proxy` , `vpcendpoint` , `elb` , `cluster_lvs` , `cluster_admin` , `master` block
supports:

* `node_id` - The node ID.

* `node_name` - The node name.

* `ip` - The node IP.

* `port` - The node port.

* `node_type` - The node role (e.g., master, slave, proxy).

* `max_memory` - The total memory in MB.

* `used_memory` - The used memory in MB.

* `qps` - The queries per second.

* `db_num` - The number of databases.

* `relation_ip` - The associated IP.

* `relation_port` - The associated port.

* `group_id` - The shard ID.

* `status` - The node status.

* `bandwidth` - The bandwidth information.

  The [bandwidth](#redis_server_bandwidth_struct) structure is documented below.

* `dbs` - The database key information.

  The [dbs](#redis_server_dbs_struct) structure is documented below.

* `dims` - The CES monitoring information.

  The [dims](#redis_server_dims_struct) structure is documented below.

<a name="redis_server_bandwidth_struct"></a>
The `bandwidth` block supports:

* `input` - The inbound bandwidth in kbit/s.

* `output` - The outbound bandwidth in kbit/s.

<a name="redis_server_dbs_struct"></a>
The `dbs` block supports:

* `db_idx` - The database index.

* `keys` - The number of keys.

* `expires` - The number of expired keys.

* `avg_ttl` - The average TTL.

<a name="redis_server_dims_struct"></a>
The `dims` block supports:

* `dim_k` - The monitoring dimension route.

* `dim_v` - The monitoring dimension ID.
