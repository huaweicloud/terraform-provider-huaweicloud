---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_connections"
description: |-
  Use this data source to get the list of connection configurations for HuaweiCloud Data Replication Service (DRS).
---

# huaweicloud_drs_connections

Use this data source to get the list of connection configurations for HuaweiCloud Data Replication Service (DRS).

## Example Usage

```hcl
data "huaweicloud_drs_connections" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `connection_id` - (Optional, String) Specifies the ID of the DRS connection to filter the results.

* `db_type` - (Optional, String) Specifies the type of the database to filter the results.
  Valid values: mysql, oracle, postgresql, mongodb.

* `name` - (Optional, String) Specifies the name of the connection to filter the results (case-insensitive).

* `inst_id` - (Optional, String) Specifies the ID of the cloud database instance to filter the results.

* `ip` - (Optional, String) Specifies the IP address of the connection to filter the results.

* `description` - (Optional, String) Specifies the description of the connection to filter the results.

* `create_time` - (Optional, String) Specifies the time range for filtering connections, separated by commas.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to filter the results.

* `fetch_all` - (Optional, Bool) Specifies whether to ignore the offset and limit parameters and return all matching records.

* `sort_key` - (Optional, String) Specifies the key by which the results are sorted. Defaults to **created_at**.

* `sort_dir` - (Optional, String) Specifies the sort direction of the results.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `connections` - The list of DRS connection configurations.

  The [connections](#connections) structure is documented below.

<a name="connections"></a>
The `connections` block supports:

* `connection_id` - The ID of the DRS connection.

* `name` - The name of the DRS connection.

* `create_time` - The creation time of the connection, in timestamp format.

* `db_type` - The type of the database corresponding to the connection.
  The valid values are as follows:
  + **MYSQL**
  + **ORACLE**
  + **POSTGRESQL**
  + **MONGODB**

* `enterprise_project_id` - The ID of the enterprise project to which the connection belongs.

* `description` - The description of the connection.

* `config` - The driver configuration of the connection.

* `endpoint` - The endpoint information of the connection.

* `vpc` - The VPC information of the connection.

* `ssl` - The SSL configuration of the connection.

  The [config](#config) structure is documented below.

<a name="config"></a>
The `config` block supports:

* `driver_name` - The name of the database driver used by the connection.

  The [endpoint](#endpoint) structure is documented below.

<a name="endpoint"></a>
The `endpoint` block supports:

* `id` - The ID of the endpoint.

* `endpoint_name` - The name of the endpoint.
  The valid values are as follows:
  + **ORACLE**
  + **ECS_ORACLE**
  + **CLOUD_GAUSSDBV5**
  + **MYSQL**
  + **ECS_MYSQL**
  + **CLOUD_MYSQL**
  + **REDIS**
  + **ECS_REDIS**
  + **REDISCLUSTER**
  + **ECS_REDISCLUSTER**
  + **CLOUD_GAUSSDB_REDIS**
  + **POSTGRESQL**
  + **ECS_POSTGRESQL**
  + **CLOUD_POSTGRESQL**
  + **MONGODB**
  + **ECS_MONGODB**
  + **CLOUD_MONGODB**

* `ip` - The IP address of the endpoint.

* `db_port` - The database port of the endpoint.

* `db_user` - The database user of the endpoint.

* `db_password` - The database password of the endpoint (sensitive information).

* `instance_id` - The ID of the database instance corresponding to the endpoint.

* `instance_name` - The name of the database instance corresponding to the endpoint.

* `db_name` - The name of the database corresponding to the endpoint.

* `source_sharding` - The sharding information of the source endpoint.

  The [source_sharding](#source_sharding) structure is documented below.

<a name="source_sharding"></a>
The `source_sharding` block supports:

* `id` - The ID of the sharding endpoint.

* `endpoint_name` - The name of the sharding endpoint.
  The valid values are as follows:
  + **ORACLE**
  + **ECS_ORACLE**
  + **CLOUD_GAUSSDBV5**
  + **MYSQL**
  + **ECS_MYSQL**
  + **CLOUD_MYSQL**
  + **REDIS**
  + **ECS_REDIS**
  + **REDISCLUSTER**
  + **ECS_REDISCLUSTER**
  + **CLOUD_GAUSSDB_REDIS**
  + **POSTGRESQL**
  + **ECS_POSTGRESQL**
  + **CLOUD_POSTGRESQL**
  + **MONGODB**
  + **ECS_MONGODB**
  + **CLOUD_MONGODB**

* `ip` - The IP address of the sharding endpoint.

* `db_port` - The database port of the sharding endpoint.

* `db_user` - The database user of the sharding endpoint.

* `db_password` - The database password of the sharding endpoint (sensitive information).

* `instance_id` - The ID of the sharding database instance.

* `instance_name` - The name of the sharding database instance.

* `db_name` - The name of the sharding database.

  The [vpc](#vpc) structure is documented below.

<a name="vpc"></a>
The `vpc` block supports:

* `vpc_id` - The ID of the VPC where the connection is located.

* `subnet_id` - The ID of the subnet where the connection is located.

* `security_group_id` - The ID of the security group associated with the connection.

  The [ssl](#ssl) structure is documented below.

<a name="ssl"></a>
The `ssl` block supports:

* `ssl_link` - Whether SSL is enabled for the connection.

* `ssl_cert_name` - The name of the SSL certificate used by the connection.

* `ssl_cert_key` - The key of the SSL certificate (sensitive information).

* `ssl_cert_check_sum` - The checksum of the SSL certificate.

* `ssl_cert_password` - The password of the SSL certificate (sensitive information).
