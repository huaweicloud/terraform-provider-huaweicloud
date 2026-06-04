---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_connection"
description: |-
  Manages DRS connection resource within HuaweiCloud.
---

# huaweicloud_drs_connection

Manages DRS connection resource within HuaweiCloud.

## Example Usage

### Create a DRS connection to MySQL

```hcl
variable "name" {}
variable "db_user" {}
variable "db_password" {}

resource "huaweicloud_drs_connection" "test" {
  name        = var.name
  db_type     = "mysql"
  description = "Test DRS connection for MySQL"

  endpoint {
    endpoint_name = "mysql"
    ip            = "192.168.0.100"
    db_port       = "3306"
    db_user       = var.db_user
    db_password   = var.db_password
  }

  ssl {
    ssl_link = false
  }

  lifecycle {
    ignore_changes = [
      endpoint.0.db_password
    ]
  }
}
```

### Create a DRS connection with RDS instance

```hcl
variable "name" {}
variable "db_user" {}
variable "db_password" {}
variable "rds_id" {}
variable "vpc_id" {}
variable "subnet_id" {}

resource "huaweicloud_drs_connection" "test" {
  name        = var.name
  db_type     = "mysql"
  description = "Test DRS connection with RDS instance"

  endpoint {
    endpoint_name = "cloud_mysql"
    instance_id   = var.rds_id
    db_port       = "3306"
    db_user       = var.db_user
    db_password   = var.db_password
  }

  vpc {
    vpc_id    = var.vpc_id
    subnet_id = var.subnet_id
  }

  ssl {
    ssl_link = false
  }

  lifecycle {
    ignore_changes = [
      endpoint.0.db_password,
    ]
  }
}
```

### Create a DRS connection for MongoDb

```hcl
variable "name" {}
variable "db_password" {}
variable "db_name" {}

resource "huaweicloud_drs_connection" "test" {
  name        = var.name
  db_type     = "mongodb"
  description = "Test DRS connection for MongoDb"

  endpoint {
    endpoint_name = "mongodb"
    ip            = "192.168.0.1:8080"
    db_user       = "mog"
    db_password   = var.db_password
    db_name       = var.db_name

    source_sharding {
      endpoint_name = "mongodb"
      ip            = "192.168.0.1:8000"
      db_user       = "mog"
      db_password   = var.db_password
      db_name       = var.db_name
    }

    source_sharding {
      endpoint_name = "mongodb"
      ip            = "192.168.0.2:8000"
      db_user       = "mog"
      db_password   = var.db_password
      db_name       = var.db_name
    }
  }

  ssl {
    ssl_link = false
  }

  lifecycle {
    ignore_changes = [
      endpoint.0.db_password,
      endpoint.0.source_sharding.0.db_password,
      endpoint.0.source_sharding.0.endpoint_name,
      endpoint.0.source_sharding.1.db_password,
      endpoint.0.source_sharding.1.endpoint_name,
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the connection name.
  The name valid length is `4` to `50` characters, can contain letters, digits, hyphens (-) and underscores (_). Special
  characters are not allowed.

* `db_type` - (Required, String) Specifies the database type.
  Valid values are: **mysql**, **postgresql**, **mongodb**, **oracle**.

* `description` - (Optional, String) Specifies the description of the connection. The description cannot exceed
  255 characters.

* `endpoint` - (Required, List) Specifies the endpoint configuration of the database.
  The [endpoint](#endpoint_struct) structure is documented below.

* `config` - (Optional, List) Specifies the connection configuration items. The configuration items vary depending on
  the connection type.
  The [config](#config_struct) structure is documented below.

* `vpc` - (Optional, List) Specifies the VPC, subnet, security group and other information where the database instance
  is located.
  The [vpc](#vpc_struct) structure is documented below.

* `ssl` - (Optional, List) Specifies the database SSL certificate information.
  The [ssl](#ssl_struct) structure is documented below.

* `cloud` - (Optional, List) Specifies the region, project and other information where the database instance is
  located.
  The [cloud](#cloud_struct) structure is documented below.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

<a name="endpoint_struct"></a>
The `endpoint` block supports:

* `endpoint_name` - (Required, String) Specifies the database scenario type.
  The valid values are as follows:
  + **oracle**: On-premises self-built Oracle database.
  + **ecs_oracle**: Huawei Cloud ECS self-built Oracle database.
  + **cloud_gaussdbv5**: Huawei Cloud GaussDB distributed database.
  + **mysql**: Third-party cloud/on-premises self-built MySQL database.
  + **ecs_mysql**: Huawei Cloud ECS self-built MySQL database.
  + **cloud_mysql**: Huawei Cloud RDS for MySQL database.
  + **redis**: On-premises self-built Redis database.
  + **ecs_redis**: Huawei Cloud ECS self-built Redis database.
  + **rediscluster**: On-premises self-built Redis cluster database.
  + **ecs_rediscluster**: Huawei Cloud ECS self-built Redis cluster database.
  + **cloud_gaussdb_redis**: Huawei Cloud GeminiDB Redis database.
  + **postgresql**: On-premises self-built PostgreSQL database.
  + **ecs_postgresql**: Huawei Cloud ECS self-built PostgreSQL database.
  + **cloud_postgresql**: Huawei Cloud RDS for PostgreSQL database.
  + **mongodb**: On-premises self-built MongoDB database.
  + **ecs_mongodb**: Huawei Cloud ECS self-built MongoDB database.
  + **cloud_mongodb**: Huawei Cloud DDS (Document Database Service).

* `db_user` - (Required, String) Specifies the database username.

* `db_password` - (Required, String) Specifies the database password.

* `id` - (Optional, String) Specifies the database information ID.

* `ip` - (Optional, String) Specifies the database IP address. Constraints:

  + For self-built MongoDB databases, concatenate the IP address and port with a colon (:), separate multiple values
    with commas (,), and support up to 3 IP addresses or domain names.
  + For DDS instances, concatenate the IP address and port with a colon (:), and separate multiple IP-port pairs with
    commas (,).
  + For Redis clusters, fill in the IP addresses and corresponding ports of all shards in the source Redis cluster.
    Concatenate the IP address and port with a colon (:), separate multiple IP-port pairs with commas (,), and it is
    recommended to use the slave node IP addresses. Support up to 32 IP addresses or domain names, separated by commas.

  Examples:
  + MySQL: `ip`
  + MongoDB: `ip:port,ip:port,ip:port`
  + DDS: `ip:port,ip:port`
  + Redis cluster: `ip:port,ip:port`

* `db_port` - (Optional, String) Specifies the database port. The value ranges from `1` to `65535`.

* `instance_id` - (Optional, String) Specifies the Huawei Cloud database instance ID.

* `instance_name` - (Optional, String) Specifies the Huawei Cloud database instance name.

* `db_name` - (Optional, String) Specifies the database name. For example:
  + Oracle: `serviceName.orcl`

* `source_sharding` - (Optional, List) Specifies the physical source database information.
  The [source_sharding](#source_sharding_struct) structure is documented below.

<a name="source_sharding_struct"></a>
The `source_sharding` block supports:

* `endpoint_name` - (Required, String) Specifies the database scenario type.
  The valid values are as follows:
  + **oracle**: On-premises self-built Oracle database.
  + **ecs_oracle**: Huawei Cloud ECS self-built Oracle database.
  + **cloud_gaussdbv5**: Huawei Cloud GaussDB distributed database.
  + **mysql**: Third-party cloud/on-premises self-built MySQL database.
  + **ecs_mysql**: Huawei Cloud ECS self-built MySQL database.
  + **cloud_mysql**: Huawei Cloud RDS for MySQL database.
  + **redis**: On-premises self-built Redis database.
  + **ecs_redis**: Huawei Cloud ECS self-built Redis database.
  + **rediscluster**: On-premises self-built Redis cluster database.
  + **ecs_rediscluster**: Huawei Cloud ECS self-built Redis cluster database.
  + **cloud_gaussdb_redis**: Huawei Cloud GeminiDB Redis database.
  + **postgresql**: On-premises self-built PostgreSQL database.
  + **ecs_postgresql**: Huawei Cloud ECS self-built PostgreSQL database.
  + **cloud_postgresql**: Huawei Cloud RDS for PostgreSQL database.
  + **mongodb**: On-premises self-built MongoDB database.
  + **ecs_mongodb**: Huawei Cloud ECS self-built MongoDB database.
  + **cloud_mongodb**: Huawei Cloud DDS (Document Database Service).

* `db_user` - (Required, String) Specifies the database username.

* `db_password` - (Required, String) Specifies the database password.

* `id` - (Optional, String) Specifies the database information ID.

* `ip` - (Optional, String) Specifies the database IP address. Refer to `ip` in the endpoint block for constraints and
  examples.

* `db_port` - (Optional, String) Specifies the database port.

* `instance_id` - (Optional, String) Specifies the Huawei Cloud database instance ID.

* `instance_name` - (Optional, String) Specifies the Huawei Cloud database instance name.

* `db_name` - (Optional, String) Specifies the database name.

<a name="config_struct"></a>
The `config` block supports:

* `driver_name` - (Optional, String) Specifies the driver name.

<a name="vpc_struct"></a>
The `vpc` block supports:

* `vpc_id` - (Required, String) Specifies the VPC ID where the database instance is located.

* `subnet_id` - (Required, String) Specifies the subnet ID where the database instance is located.

* `security_group_id` - (Optional, String) Specifies the security group ID where the database instance is located.

<a name="ssl_struct"></a>
The `ssl` block supports:

* `ssl_link` - (Optional, Bool) Specifies whether to enable SSL secure connection. Set to **true** if the database has
  SSL enabled.

* `ssl_cert_name` - (Optional, String) Specifies the SSL certificate name.

* `ssl_cert_key` - (Optional, String) Specifies the SSL certificate content encrypted with Base64.

* `ssl_cert_check_sum` - (Optional, String) Specifies the checksum value of the SSL certificate content.

* `ssl_cert_password` - (Optional, String) Specifies the SSL certificate password. This parameter is required when the
  certificate file suffix is **.p12**.

<a name="cloud_struct"></a>
The `cloud` block supports:

* `region` - (Optional, String) Specifies the region ID. This parameter is required when the database instance type is
  **ecs** (Huawei Cloud ECS self-built database) or **cloud** (Huawei Cloud database). For details, see Regions and
  Endpoints.

  Note: When there are sub-projects under the region, the Region ID is concatenated by the region project ID and
  sub-project ID with an underscore (_), for example: `cn-north-4_abc`.

* `project_id` - (Optional, String) Specifies the project ID of the tenant in a region.

* `az_code` - (Optional, String) Specifies the availability zone (AZ) name where the database is located.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also `connection_id`).

* `create_time` - The creation time of the connection.

## Import

The DRS connection can be imported by `id`. e.g.

```bash
$ terraform import huaweicloud_drs_connection.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `endpoint.0.db_password`,
`endpoint.0.source_sharding.*.db_password`, `endpoint.0.source_sharding.*.endpoint_name`,`ssl.0.ssl_cert_key`,
`ssl.0.ssl_cert_password`,`cloud`.
It is generally recommended running **terraform plan** after importing a connection. You can then
decide if changes should be applied to the connection, or the resource definition should be updated to align with the
connection. Also you can ignore changes as below.

```hcl
resource "huaweicloud_drs_connection" "test" {
  ...

lifecycle {
  ignore_changes = [
    endpoint.0.db_password,
    endpoint.0.source_sharding.*.db_password,
    endpoint.0.source_sharding.*.endpoint_name,
    ssl.0.ssl_cert_key,
    ssl.0.ssl_cert_password,
    cloud
  ]
}
}
```
