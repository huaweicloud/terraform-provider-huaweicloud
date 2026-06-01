---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_ip_num_requirement"
description: |-
  Use this data source to query the IP addresses number for creating an instance or adding nodes to an instance.
---

# huaweicloud_geminidb_ip_num_requirement

Use this data source to query the IP addresses number for creating an instance or adding nodes to an instance.

-> This data source supports the following instances: GeminiDB Cassandra, GeminiDB Mongo, GeminiDB Influx and
  GeminiDB Redis.

## Example Usage

### Query IP addresses number by instance type

```hcl
variable "node_num" {}
variable "engine_name" {}
variable "instance_mode" {}

data "huaweicloud_geminidb_ip_num_requirement" "test" {
  node_num      = var.node_num
  engine_name   = var.engine_name
  instance_mode = var.instance_mode
}
```

### Query IP addresses number by instance ID

```hcl
variable "node_num" {}
variable "instance_id" {}

data "huaweicloud_geminidb_ip_num_requirement" "test" {
  node_num    = var.node_num
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `node_num` - (Required, Int) Specifies the nodes number for creating or scaling out an instance.
  The valid value ranges from `1` to `200`.

* `engine_name` - (Optional, String) Specifies the database type.
  The valid values are as follows:
  + **cassandra**: Indicates GeminiDB Cassandra.
  + **mongodb**: Indicates GeminiDB Mongo.
  + **influxdb**: Indicates GeminiDB Influx.
  + **redis**: Indicates GeminiDB Redis.

  -> This parameter is mandatory when `instance_id` is not transferred.

* `instance_mode` - (Optional, String) Specifies the instance type.
  The valid values are as follows:
  + **Cluster**: Indicates proxy cluster GeminiDB Redis instance, cluster GeminiDB Cassandra or Influx instance
  with classic storage.
  + **CloudNativeCluster**: Indicates cluster GeminiDB Cassandra, Influx, or Redis instance with cloud native storage.
  + **RedisCluster**: Indicates Redis Cluster GeminiDB Redis instance with classic storage.
  + **Replication**: Indicates primary/standby GeminiDB Redis instance with classic storage.
  + **InfluxdbSingle**: Indicates single-node GeminiDB Influx instance with classic storage.
  + **EnhancedCluster**: Indicates GeminiDB Influx cluster (performance-enhanced) instance with classic storage.
  + **ReplicaSet**: Indicates GeminiDB Mongo instance in a replica set.

  -> This parameter is mandatory when `instance_id` is not transferred.

* `instance_id` - (Optional, String) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ip_address_count` - The number of IP addresses.
