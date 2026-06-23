---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_configuration_datastores"
description: |-
  Use this data source to query API that support parameter templates.
---

# huaweicloud_geminidb_configuration_datastores

Use this data source to query API that support parameter templates.

## Example Usage

```hcl
data "huaweicloud_geminidb_configuration_datastores" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `datastores` - The database API information.
  The [datastores](#datastores_struct) structure is documented below.

<a name="datastores_struct"></a>
The `datastores` block supports:

* `datastore_name` - The database API.
  + **cassandra**: GeminiDB Cassandra instance.
  + **mongodb**: GeminiDB Mongo instance.
  + **influxdb**: GeminiDB Influx instance.
  + **redis**: GeminiDB Redis instance.
  + **dynamodb**: GeminiDB DynamoDB-Compatible instance.
  + **hbase**: GeminiDB HBase instance.

* `version` - The database API version.
  + **3.11**: GeminiDB Cassandra instance 3.11.
  + **4.0**: GeminiDB Mongo instance 4.0.
  + **1.8**: GeminiDB Influx instance 1.8.
  + **5.0**: GeminiDB Redis instance 5.0.

* `mode` - The instance type.
  + **Cluster**: GeminiDB Cassandra cluster instance with classic storage, GeminiDB Influx cluster instance
  with classic storage, and proxy cluster GeminiDB Redis instance with classic storage.
  + **CloudNativeCluster**: GeminiDB Cassandra cluster instance with cloud native storage, GeminiDB Influx
  cluster (performance-enhanced) instance with cloud native storage, and GeminiDB Redis cluster instance with
  cloud native storage.
  + **ReplicaSet**: GeminiDB Mongo instance 4.0 in a replica set.
  + **EnhancedCluster**: GeminiDB Influx cluster (performance-enhanced) instance with classic storage.
  + **InfluxdbSingle**: single-node GeminiDB Influx instance.
  + **RedisCluster**: Redis Cluster GeminiDB Redis instance which classic storage.
  + **Replication**: primary/standby GeminiDB Redis instance with classic storage.
