---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_recycling_instances"
description: |-
  Use this data source to get the list of GeminiDB recycling instances.
---

# huaweicloud_geminidb_recycling_instances

Use this data source to get the list of GeminiDB recycling instances.

## Example Usage

```hcl
data "huaweicloud_geminidb_recycling_instances" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the list of recycling GeminiDB instances.
  The [instances](#geminidb_recycling_instances) structure is documented below.

<a name="geminidb_recycling_instances"></a>
The `instances` block supports:

* `id` - The instance ID.

* `name` - The instance name.

* `mode` - The instance mode.
  The valid values are as follows:
  + **Cluster**: GeminiDB Cassandra, GeminiDB Influx, GeminiDB Redis classic deployment mode Proxy cluster instance.
  + **CloudNativeCluster**: GeminiDB Cassandra, GeminiDB Influx, GeminiDB Redis cloud-native deployment mode cluster instance.
  + **RedisCluster**: GeminiDB Redis classic deployment mode Cluster cluster instance.
  + **Replication**: GeminiDB Redis classic deployment mode primary/standby instance.
  + **InfluxdbSingle**: GeminiDB Influx classic deployment mode single node instance.
  + **ReplicaSet**: GeminiDB Mongo replica set instance.

* `product_type` - The product type.
  Only applicable to GeminiDB Redis cloud-native deployment mode cluster.
  The valid values are as follows:
  + **Standard**: Standard type.
  + **Capacity**: Capacity type.

* `data_store` - The database information.
  The [data_store](#geminidb_recycling_instances_data_store) structure is documented below.

* `charge_type` - The billing mode.
  The valid values are as follows:
  + **prePaid**: Prepaid, i.e., yearly/monthly.
  + **postPaid**: Postpaid, i.e., pay-per-use.

* `enterprise_project_id` - The enterprise project ID.
  The value "0" indicates the default enterprise project.

* `backup_id` - The backup ID.

* `created_at` - The instance creation time.

* `deleted_at` - The instance deletion time.

* `retained_until` - The retention end time.

<a name="geminidb_recycling_instances_data_store"></a>
The `data_store` block supports:

* `type` - The database type.
  The valid values are as follows:
  + **cassandra**: GeminiDB Cassandra database instance.
  + **mongodb**: GeminiDB Mongo database instance.
  + **influxdb**: GeminiDB Influx database instance.
  + **redis**: GeminiDB Redis database instance.

* `version` - The database version.
