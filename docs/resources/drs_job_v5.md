---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_job_v5"
description: |-
  Manages a DRS job resource within HuaweiCloud.
---

# huaweicloud_drs_job_v5

Manages a DRS job resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "source_db_password" {}
variable "target_db_password" {}
variable "source_instance_id" {}
variable "target_instance_id" {}
variable "project_id" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}

resource "huaweicloud_drs_job_v5" "test" {
  base_info {
    name                  = var.name
    job_type              = "sync"
    engine_type           = "mysql-to-mysql"
    job_direction         = "up"
    task_type             = "FULL_INCR_TRANS"
    net_type              = "eip"
    charging_mode         = "on_demand"
    enterprise_project_id = "0"
    expired_days          = "14"
    is_open_fast_clean    = false

    tags {
      key   = "tag1"
      value = "value1"
    }
  }

  source_endpoint {
    db_type       = "mysql"
    endpoint_type = "cloud"
    endpoint_role = "so"

    endpoint {
      endpoint_name = "cloud_mysql"
      ip            = "192.168.0.141"
      db_port       = "3306"
      db_user       = "root"
      db_password   = var.source_db_password
      instance_id   = var.source_instance_id
      db_name       = "user"
    }

    cloud {
      region     = "cn-north-4"
      project_id = var.project_id
    }

    ssl {
      ssl_link = false
    }
  }

  target_endpoint {
    db_type       = "mysql"
    endpoint_type = "cloud"
    endpoint_role = "ta"

    endpoint {
      endpoint_name = "cloud_mysql"
      ip            = "192.168.0.105"
      db_port       = "3306"
      db_user       = "root"
      db_password   = var.target_db_password
      instance_id   = var.target_instance_id
    }

    cloud {
      region     = "cn-north-4"
      project_id = var.project_id
      az_code    = "cn-north-4a,cn-north-4c,cn-north-4g"
    }

    vpc {
      vpc_id            = var.vpc_id
      subnet_id         = var.subnet_id
      security_group_id = var.security_group_id
    }
  }

  node_info {
    spec {
      node_type = "medium"
    }

    vpc {
      vpc_id    = var.vpc_id
      subnet_id = var.subnet_id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `base_info` - (Required, List, NonUpdatable) Specifies the basic information of the DRS job.
  The maximum number of elements is `1`.
  The [base_info](#block--base_info) structure is documented below.

* `source_endpoint` - (Required, List, NonUpdatable) Specifies the source database endpoint information.
  The [source_endpoint](#block--source_endpoint) structure is documented below.

* `target_endpoint` - (Required, List, NonUpdatable) Specifies the target database endpoint information.
  The [target_endpoint](#block--target_endpoint) structure is documented below.

* `node_info` - (Required, List, NonUpdatable) Specifies the node information of the DRS job instance.
  The maximum number of elements is `1`.
  The [node_info](#block--node_info) structure is documented below.

* `period_order` - (Optional, List, NonUpdatable) Specifies the yearly/monthly billing information.
  The maximum number of elements is `1`.
  The [period_order](#block--period_order) structure is documented below.

* `public_ip_list` - (Optional, List, NonUpdatable) Specifies the public IP information.
  The [public_ip_list](#block--public_ip_list) structure is documented below.

<a name="block--base_info"></a>
The `base_info` block supports:

* `name` - (Optional, String, NonUpdatable) Specifies the job name. The name consists of `4` to `50` characters,
  starting with a letter. Only letters, digits, underscores (_) and hyphens (-) are allowed.

* `job_type` - (Optional, String, NonUpdatable) Specifies the job scenario. The options are as follows:
  + **migration**: Real-time migration.
  + **sync**: Real-time synchronization.
  + **cloudDataGuard**: Real-time disaster recovery.

* `multi_write` - (Optional, Bool, NonUpdatable) Specifies whether the disaster recovery type is dual-primary
  disaster recovery. This parameter is mandatory when `job_type` is **cloudDataGuard**. If the disaster recovery
  type is dual-primary, set this parameter to **true**; otherwise, set it to **false**.

* `engine_type` - (Optional, String, NonUpdatable) Specifies the engine type. The options are as follows:
  + **oracle-to-gaussdbv5**: Oracle to GaussDB distributed.
  + **mysql-to-mysql**: MySQL to MySQL.
  + **redis-to-gaussredis**: Redis to GeminiDB Redis.
  + **rediscluster-to-gaussredis**: Redis cluster to GeminiDB Redis.

* `job_direction` - (Optional, String, NonUpdatable) Specifies the migration direction. The options are as follows:
  + **up**: Ingress. In the disaster recovery scenario, the local cloud is the standby.
  + **down**: Egress. In the disaster recovery scenario, the local cloud is the primary.
  + **non-dbs**: Self-built.

* `task_type` - (Optional, String, NonUpdatable) Specifies the migration mode. The options are as follows:
  + **FULL_TRANS**: Full migration.
  + **FULL_INCR_TRANS**: Full + incremental migration.
  + **INCR_TRANS**: Incremental migration.

* `net_type` - (Optional, String, NonUpdatable) Specifies the network type. The options are as follows:
  + **eip**: Public network.
  + **vpc**: VPC network. The VPC network is not supported in the disaster recovery scenario.
  + **vpn**: VPN or dedicated line network.

* `charging_mode` - (Optional, String, NonUpdatable) Specifies the charging mode. The default value is **on_demand**.
  The options are as follows:
  + **period**: Yearly/Monthly billing.
  + **on_demand**: Pay-per-use billing.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  Defaults to **0**, which indicates the default enterprise project.

* `description` - (Optional, String, NonUpdatable) Specifies the job description. The description can contain up to
  `256` characters and cannot contain the following special characters: !<&gt;&lt;&gt;&'"`.

* `start_time` - (Optional, String, NonUpdatable) Specifies the scheduled start time of the job.

* `expired_days` - (Optional, String, NonUpdatable) Specifies the number of days after which a job in abnormal
  status will automatically end. The value ranges from `14` to `100`. Defaults to **14**.

* `tags` - (Optional, List, NonUpdatable) Specifies the tag information. A maximum of `20` tags can be added.
  The [tags](#block--tags) structure is documented below.

* `is_open_fast_clean` - (Optional, Bool, NonUpdatable) Specifies whether to enable fast cleanup of Binlog for
  RDS for MySQL and MariaDB. Defaults to **false**.

<a name="block--tags"></a>
The `tags` block supports:

* `key` - (Optional, String, NonUpdatable) Specifies the tag key. The maximum length is `36` characters.
  Only letters, digits, underscores (_), hyphens (-) and Chinese characters are allowed.

* `value` - (Optional, String, NonUpdatable) Specifies the tag value. The maximum length is `43` characters.
  Only letters, digits, underscores (_), hyphens (-) and Chinese characters are allowed.

<a name="block--source_endpoint"></a>
The `source_endpoint` and `target_endpoint` blocks support:

* `db_type` - (Required, String, NonUpdatable) Specifies the database type. The options are as follows:
  + **oracle**: Oracle.
  + **gaussdbv5**: GaussDB distributed.
  + **redis**: Redis.
  + **rediscluster**: Redis cluster.
  + **gaussredis**: GeminiDB Redis.

* `endpoint_type` - (Required, String, NonUpdatable) Specifies the database instance type. The options are as follows:
  + **offline**: Self-built database.
  + **ecs**: Huawei Cloud ECS self-built database.
  + **cloud**: Huawei Cloud database.

* `endpoint_role` - (Required, String, NonUpdatable) Specifies the database instance role. The options are as follows:
  + **so**: Source database.
  + **ta**: Target database.

* `endpoint` - (Required, List, NonUpdatable) Specifies the database basic information.
  The maximum number of elements is `1`.
  The [endpoint](#block--endpoint) structure is documented below.

* `cloud` - (Optional, List, NonUpdatable) Specifies the region and project information of the database instance.
  The maximum number of elements is `1`.
  The [cloud](#block--cloud) structure is documented below.

* `vpc` - (Optional, List, NonUpdatable) Specifies the VPC, subnet, and security group information of the database
  instance. This parameter is mandatory when the database instance type is **ecs**.
  The maximum number of elements is `1`.
  The [vpc](#block--endpoint_vpc) structure is documented below.

* `config` - (Optional, List, NonUpdatable) Specifies the database configuration information.
  The maximum number of elements is `1`.
  The [config](#block--config) structure is documented below.

* `ssl` - (Optional, List, NonUpdatable) Specifies the SSL certificate information of the database.
  The maximum number of elements is `1`.
  The [ssl](#block--ssl) structure is documented below.

* `customized_dns` - (Optional, List, NonUpdatable) Specifies the custom DNS information.
  The maximum number of elements is `1`.
  The [customized_dns](#block--customized_dns) structure is documented below.

<a name="block--endpoint"></a>
The `endpoint` block supports:

* `endpoint_name` - (Required, String, NonUpdatable) Specifies the database scenario type. The options are as follows:
  + **oracle**: Self-built Oracle database.
  + **ecs_oracle**: Huawei Cloud ECS self-built Oracle database.
  + **cloud_gaussdbv5**: Huawei Cloud GaussDB distributed.
  + **mysql**: Self-built MySQL database.
  + **ecs_mysql**: Huawei Cloud ECS self-built MySQL database.
  + **cloud_mysql**: Huawei Cloud RDS for MySQL.
  + **redis**: Self-built Redis.
  + **ecs_redis**: Huawei Cloud ECS self-built Redis.
  + **rediscluster**: Self-built Redis cluster.
  + **ecs_rediscluster**: Huawei Cloud ECS self-built Redis cluster.
  + **cloud_gaussdb_redis**: Huawei Cloud GeminiDB Redis.

* `db_user` - (Required, String, NonUpdatable) Specifies the database username.

* `db_password` - (Required, String, NonUpdatable) Specifies the database password.

* `id` - (Optional, String, NonUpdatable) Specifies the database information ID.

* `ip` - (Optional, String, NonUpdatable) Specifies the database IP address. For Redis cluster, fill in the IP
  addresses and ports of all shards in the source Redis cluster, separated by colons (:) between IP and port,
  and separated by commas (,) between multiple entries. A maximum of `32` IP addresses or domain names are supported.

* `db_port` - (Optional, String, NonUpdatable) Specifies the database port. The value ranges from `1` to `65535`.

* `instance_id` - (Optional, String, NonUpdatable) Specifies the Huawei Cloud database instance ID.

* `instance_name` - (Optional, String, NonUpdatable) Specifies the Huawei Cloud database instance name.

* `db_name` - (Optional, String, NonUpdatable) Specifies the database name. For example, for Oracle:
  **serviceName.orcl**.

* `source_sharding` - (Optional, String, NonUpdatable) Specifies the physical source database information in
  JSON string format.

<a name="block--cloud"></a>
The `cloud` block supports:

* `region` - (Required, String, NonUpdatable) Specifies the region ID. When the database instance type is **ecs**
  or **cloud**, this parameter is mandatory. If the region has sub-projects, the region ID is composed of the
  region project ID and sub-project ID, joined by an underscore (_), for example: **cn-north-4_abc**.

* `project_id` - (Required, String, NonUpdatable) Specifies the project ID in the region.

* `az_code` - (Optional, String, NonUpdatable) Specifies the availability zone (AZ) where the database is located.

<a name="block--endpoint_vpc"></a>
The `vpc` block (in source_endpoint/target_endpoint) supports:

* `vpc_id` - (Required, String, NonUpdatable) Specifies the VPC ID where the database instance is located.

* `subnet_id` - (Required, String, NonUpdatable) Specifies the subnet ID where the database instance is located.

* `security_group_id` - (Optional, String, NonUpdatable) Specifies the security group ID where the database
  instance is located.

<a name="block--config"></a>
The `config` block supports:

* `is_target_readonly` - (Optional, Bool, NonUpdatable) Specifies whether to set the target instance to read-only.
  This is effective for MySQL migration and disaster recovery when `job_direction` is **up**. Defaults to **true**.

* `node_num` - (Optional, Int, NonUpdatable) Specifies the number of subtasks for connecting to the source Redis
  cluster in the Redis cluster to GeminiDB Redis migration scenario. The value ranges from `1` to `16`, and must
  not exceed the number of shards in the source Redis cluster. Defaults to **0**.

<a name="block--ssl"></a>
The `ssl` block supports:

* `ssl_link` - (Optional, Bool, NonUpdatable) Specifies whether to use SSL secure connection. If the database
  has SSL enabled, set this parameter to **true**.

* `ssl_cert_name` - (Optional, String, NonUpdatable) Specifies the SSL certificate name.

* `ssl_cert_key` - (Optional, String, NonUpdatable) Specifies the SSL certificate content, encrypted with base64.

* `ssl_cert_check_sum` - (Optional, String, NonUpdatable) Specifies the checksum value of the SSL certificate content.
  This is required for source database secure connections.

* `ssl_cert_password` - (Optional, String, NonUpdatable) Specifies the SSL certificate password. This is mandatory
  when the certificate file suffix is **.p12**.

<a name="block--customized_dns"></a>
The `customized_dns` block supports:

* `is_set_dns` - (Required, Bool, NonUpdatable) Specifies whether to set custom DNS.

* `set_dns_action` - (Required, String, NonUpdatable) Specifies the action for setting custom DNS. The options are
  as follows:
  + **add**: Add custom DNS IP.
  + **keep**: Keep custom DNS IP.
  + **update**: Update custom DNS IP (takes effect when the DNS IP changes).
  + **recover**: Restore the system default DNS IP (may cause domain name resolution failure, use with caution).

* `dns_ip` - (Required, String, NonUpdatable) Specifies the custom DNS IP. The maximum length is `15` characters.

<a name="block--node_info"></a>
The `node_info` block supports:

* `spec` - (Required, List, NonUpdatable) Specifies the specification information of the job instance.
  The maximum number of elements is `1`.
  The [spec](#block--spec) structure is documented below.

* `vpc` - (Optional, List, NonUpdatable) Specifies the VPC information of the job instance. This is mandatory for
  self-built jobs.
  The maximum number of elements is `1`.
  The [vpc](#block--node_info_vpc) structure is documented below.

* `base_info` - (Optional, List, NonUpdatable) Specifies the basic information of the job instance. This is mandatory
  for self-built jobs.
  The maximum number of elements is `1`.
  The [base_info](#block--node_info_base_info) structure is documented below.

<a name="block--spec"></a>
The `spec` block supports:

* `node_type` - (Required, String, NonUpdatable) Specifies the instance specification code. The options are as follows:
  + **micro**: Extra small.
  + **small**: Small.
  + **medium**: Medium.
  + **high**: Large.

<a name="block--node_info_vpc"></a>
The `vpc` block (in node_info) supports:

* `vpc_id` - (Required, String, NonUpdatable) Specifies the VPC ID where the job instance is located.

* `subnet_id` - (Required, String, NonUpdatable) Specifies the subnet ID where the job instance is located.

* `custom_node_ip` - (Optional, String, NonUpdatable) Specifies the IP address of the job instance. Multiple IP
  addresses are separated by commas (,). Only IPv4 addresses are supported. Example: **192.168.0.10,192.168.0.11**.

* `security_group_id` - (Optional, String, NonUpdatable) Specifies the security group ID where the job instance
  is located.

<a name="block--node_info_base_info"></a>
The `base_info` block (in node_info) supports:

* `instance_type` - (Required, String, NonUpdatable) Specifies the instance type. The options are as follows:
  + **single**: Single.
  + **ha**: Primary/standby.

* `arch` - (Required, String, NonUpdatable) Specifies the CPU architecture. The options are as follows:
  + **x86**: x86 architecture.
  + **arm**: ARM architecture.

* `availability_zone` - (Required, String, NonUpdatable) Specifies the availability zone ID. For instances that
  are not single, you need to specify availability zones for all nodes, separated by commas (,).
  Example: **cn-north-4a** for single instance, **cn-north-4a,cn-north-4b** for primary/standby instance.

* `status` - (Optional, String, NonUpdatable) Specifies the status.

* `role` - (Optional, String, NonUpdatable) Specifies the primary/standby role of the job.

<a name="block--period_order"></a>
The `period_order` block supports:

* `period_type` - (Required, Int, NonUpdatable) Specifies the subscription period type. The options are as follows:
  + **2**: Monthly.
  + **3**: Yearly.

* `period_num` - (Required, Int, NonUpdatable) Specifies the subscription period number. When `period_type` is **2**,
  the value indicates the number of months. When `period_type` is **3**, the value indicates the number of years.

* `is_auto_renew` - (Optional, Int, NonUpdatable) Specifies whether to auto-renew. The options are as follows:
  + **0**: No (default, manual payment required).
  + **1**: Yes (automatic payment).

<a name="block--public_ip_list"></a>
The `public_ip_list` block supports:

* `id` - (Required, String, NonUpdatable) Specifies the public IP ID.

* `public_ip` - (Required, String, NonUpdatable) Specifies the public IP address.

* `type` - (Required, String, NonUpdatable) Specifies the type of the public IP bound to the job. For
  primary/standby jobs, **master** indicates the primary and **slave** indicates the standby. For other types,
  it is fixed to **master**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
