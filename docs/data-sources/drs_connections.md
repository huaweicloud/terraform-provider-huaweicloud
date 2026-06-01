---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_connections"
description: |-
  Use this data source to get the list of connection configurations for HuaweiCloud Data Replication Service.
---

# huaweicloud_drs_connections

Use this data source to get the list of connection configurations for HuaweiCloud Data Replication Service.

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
  The valid values are as follows:
  + **MYSQL**
  + **ORACLE**
  + **POSTGRESQL**
  + **MONGODB**

* `name` - (Optional, String) Specifies the name of the connection to filter the results.

* `inst_id` - (Optional, String) Specifies the ID of the cloud database instance to filter the results.

* `ip` - (Optional, String) Specifies the IP address of the connection to filter the results.

* `description` - (Optional, String) Specifies the description of the connection to filter the results.

* `create_time` - (Optional, String) Specifies the time range for filtering connections, separated by commas.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to filter the results.

* `sort_key` - (Optional, String) Specifies the key by which the results are sorted.

* `sort_dir` - (Optional, String) Specifies the sort direction of the results.
  The valid values are as follows:
  + **DESC**: Descending order.
  + **ASC**: Ascending order.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `connections` - The list of DRS connection configurations.

  The [connections](#connections_struct) structure is documented below.

<a name="connections_struct"></a>
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

  The [config](#config_struct) structure is documented below.

* `endpoint` - The endpoint information of the connection.

  The [endpoint](#endpoint_struct) structure is documented below.

* `vpc` - The VPC information of the connection.

  The [vpc](#vpc_struct) structure is documented below.

* `ssl` - The SSL configuration of the connection.

  The [ssl](#ssl_struct) structure is documented below.

<a name="config_struct"></a>
The `config` block supports:

* `driver_name` - The name of the database driver used by the connection.

<a name="endpoint_struct"></a>
The `endpoint` block supports:

* `id` - The ID of the endpoint.

* `endpoint_name` - The name of the endpoint.
  The valid values are as follows:
  + **ORACLE**: Self-built Oracle database.
  + **ECS_ORACLE**: Self-built Oracle database on ECS.
  + **CLOUD_GAUSSDBV5**: GaussDB distributed database.
  + **MYSQL**: Self-built MySQL database.
  + **ECS_MYSQL**: Self-built MySQL database on ECS.
  + **CLOUD_MYSQL**: RDS for MySQL.
  + **REDIS**: Self-built Redis database.
  + **ECS_REDIS**: Self-built Redis database on ECS.
  + **REDISCLUSTER**: Self-built Redis cluster database.
  + **ECS_REDISCLUSTER**: Self-built Redis cluster database on ECS.
  + **CLOUD_GAUSSDB_REDIS**: GeminiDB Redis.
  + **POSTGRESQL**: Self-built PostgreSQL database.
  + **ECS_POSTGRESQL**: Self-built PostgreSQL database on ECS.
  + **CLOUD_POSTGRESQL**: RDS for PostgreSQL.
  + **MONGODB**: Self-built MongoDB database.
  + **ECS_MONGODB**: Self-built MongoDB database on ECS.
  + **CLOUD_MONGODB**: DDS.

* `ip` - The IP address of the endpoint.

* `db_port` - The database port of the endpoint.

* `db_user` - The database user of the endpoint.

* `db_password` - The database password of the endpoint.

* `instance_id` - The ID of the database instance corresponding to the endpoint.

* `instance_name` - The name of the database instance corresponding to the endpoint.

* `db_name` - The name of the database corresponding to the endpoint.

* `source_sharding` - The sharding information of the source endpoint.

  The [source_sharding](#source_sharding_struct) structure is documented below.

<a name="source_sharding_struct"></a>
The `source_sharding` block supports:

* `id` - The ID of the sharding endpoint.

* `endpoint_name` - The name of the sharding endpoint.
  The valid values are as follows:
  + **ORACLE**: Self-built Oracle database.
  + **ECS_ORACLE**: Self-built Oracle database on ECS.
  + **CLOUD_GAUSSDBV5**: GaussDB distributed database.
  + **MYSQL**: Self-built MySQL database.
  + **ECS_MYSQL**: Self-built MySQL database on ECS.
  + **CLOUD_MYSQL**: RDS for MySQL.
  + **REDIS**: Self-built Redis database.
  + **ECS_REDIS**: Self-built Redis database on ECS.
  + **REDISCLUSTER**: Self-built Redis cluster database.
  + **ECS_REDISCLUSTER**: Self-built Redis cluster database on ECS.
  + **CLOUD_GAUSSDB_REDIS**: GeminiDB Redis.
  + **POSTGRESQL**: Self-built PostgreSQL database.
  + **ECS_POSTGRESQL**: Self-built PostgreSQL database on ECS.
  + **CLOUD_POSTGRESQL**: RDS for PostgreSQL.
  + **MONGODB**: Self-built MongoDB database.
  + **ECS_MONGODB**: Self-built MongoDB database on ECS.
  + **CLOUD_MONGODB**: DDS.

* `ip` - The IP address of the sharding endpoint.

* `db_port` - The database port of the sharding endpoint.

* `db_user` - The database user of the sharding endpoint.

* `db_password` - The database password of the sharding endpoint.

* `instance_id` - The ID of the sharding database instance.

* `instance_name` - The name of the sharding database instance.

* `db_name` - The name of the sharding database.

<a name="vpc_struct"></a>
The `vpc` block supports:

* `vpc_id` - The ID of the VPC where the connection is located.

* `subnet_id` - The ID of the subnet where the connection is located.

* `security_group_id` - The ID of the security group associated with the connection.

<a name="ssl_struct"></a>
The `ssl` block supports:

* `ssl_link` - Whether SSL is enabled for the connection.

* `ssl_cert_name` - The name of the SSL certificate used by the connection.

* `ssl_cert_key` - The key of the SSL certificate.

* `ssl_cert_check_sum` - The checksum of the SSL certificate.

* `ssl_cert_password` - The password of the SSL certificate.
