---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_quotas"
description: |-
  Use this data source to query the GeminiDB resource quotas.
---

# huaweicloud_geminidb_quotas

Use this data source to query the GeminiDB resource quotas.

## Example Usage

```hcl
data "huaweicloud_geminidb_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `datastore_type` - (Optional, String) Specifies the database type.
  The valid values are as follows:
  + **cassandra**: Indicates GeminiDB Cassandra.
  + **mongodb**: Indicates GeminiDB Mongo.
  + **influxdb**: Indicates GeminiDB Influx.
  + **redis**: Indicates GeminiDB Redis.

* `mode` - (Optional, String) Specifies the instance type.
  The valid values are as follows:
  + **Cluster**: Indicates proxy cluster GeminiDB Redis instance, cluster GeminiDB Cassandra or Influx instance
  with classic storage.
  + **CloudNativeCluster**: Indicates cluster GeminiDB Cassandra, Influx, or Redis instance with cloud native storage.
  + **RedisCluster**: Indicates Redis Cluster GeminiDB Redis instance with classic storage.
  + **Replication**: Indicates primary/standby GeminiDB Redis instance with classic storage.
  + **InfluxdbSingle**: Indicates single-node GeminiDB Influx instance with classic storage.
  + **EnhancedCluster**: Indicates GeminiDB Influx cluster (performance-enhanced) instance with classic storage.
  + **ReplicaSet**: Indicates GeminiDB Mongo instance in a replica set.

  -> If `datastore_type` is not transferred, this parameter is automatically ignored.
    This parameter is mandatory when `datastore_type` is transferred.

* `product_type` - (Optional, String) Specifies the product type.
  The valid values are as follows:
  + **Capacity**
  + **Standard**
  + **Performance**

  -> This parameter is mandatory when you query a GeminiDB Redis cluster instance with cloud native storage.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The GeminiDB quotas.
  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `type` - The quota resource type.

* `quota` - The current quota.
  If this parameter is set to `0`, no quantity limit is set for resources.

* `used` - The number of used resources.
