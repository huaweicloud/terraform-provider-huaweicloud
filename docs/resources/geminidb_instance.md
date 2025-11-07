---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_instance"
description: |-
  Manages a GeminiDB instance resource within HuaweiCloud.
---

# huaweicloud_geminidb_instance

Manages a GeminiDB instance resource within HuaweiCloud.

## Example Usage

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}
variable "availability_zone" {}
variable "configuration_id" {}
variable "enterprise_project_id" {}

resource "huaweicloud_geminidb_instance" "test" {
  name                  = "test_name"
  availability_zone     = var.availability_zone
  vpc_id                = var.vpc_id
  subnet_id             = var.subnet_id
  security_group_id     = var.security_group_id
  password              = "test_password"
  mode                  = "Cluster"
  configuration_id      = var.configuration_id
  enterprise_project_id = var.enterprise_project_id
  port                  = 8888
  ssl_option            = "on"

  datastore {
    type           = "redis"
    version        = "5.0"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "5"
    size      = "16"
    storage   = "ULTRAHIGH"
    spec_code = "geminidb.redis.large.2"
  }

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the rds instance resource. If omitted, the
  provider-level region will be used. Changing this creates a new rds instance resource.

* `name` - (Required, String) Specifies the instance name, which can be the same as an existing instance name. The value
  must be **4** to **64** characters in length and start with a letter. It is case-sensitive and can contain only letters,
  digits, hyphens (-), and underscores (_). Chinese characters must be in UTF-8 or Unicode format.

* `datastore` - (Required, List, NonUpdatable) Specifies the database information.
  The [datastore](#datastore_struct) structure is documented below.

* `availability_zone` - (Required, String, NonUpdatable) Specifies the AZ ID.

* `vpc_id` - (Required, String, NonUpdatable) Specifies the VPC ID.

* `subnet_id` - (Required, String, NonUpdatable) Specifies the subnet ID.

* `security_group_id` - (Required, String) Specifies the security group ID.

* `password` - (Required, String) Specifies the database password.

* `mode` - (Required, String, NonUpdatable) Specifies the instance type. Value options:
  + **Cluster**: GeminiDB Cassandra cluster instance with classic storage.
  + **CloudNativeCluster**: GeminiDB Cassandra cluster instance with cloud native storage.
  + **ReplicaSet**: GeminiDB Mongo instance 4.0 in a replica set.
  + **Cluster**: GeminiDB Influx cluster instance which classic storage.
  + **CloudNativeCluster**: GeminiDB Influx cluster (performance-enhanced) instance with cloud native storage.
  + **EnhancedCluster**: GeminiDB Influx cluster (performance-enhanced) instance with classic storage.
  + **InfluxdbSingle**: single-node GeminiDB Influx instance.
  + **Cluster**: proxy cluster GeminiDB Redis instance with classic storage.
  + **CloudNativeCluster**: GeminiDB Redis cluster instance with cloud native storage.
  + **RedisCluster**: Redis Cluster GeminiDB Redis instance which classic storage.
  + **Replication**: primary/standby GeminiDB Redis instance with classic storage.

* `flavor` - (Required, List, NonUpdatable) Specifies the instance specifications.
  The [flavor](#flavor_struct) structure is documented below.

* `product_type` - (Optional, String, NonUpdatable) Specifies the product type. This parameter is mandatory when you
  create a GeminiDB Redis cluster instance with cloud native storage. Value options: **Capacity**.

* `configuration_id` - (Optional, String) Specifies the parameter template ID.

* `backup_strategy` - (Optional, List) Specifies the advanced backup policy.
  The [backup_strategy](#backup_strategy_struct) structure is documented below.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `ssl_option` - (Optional, String) Specifies whether SSL is enabled. Value options:
  + **on**: indicating that SSL is enabled by default.
  + **off**: indicating that SSL is disabled by default.

* `dedicated_resource_id` - (Optional, String, NonUpdatable) Specifies the dedicated resource ID. This parameter is
  available only after a dedicated resource pool is enabled.

* `port` - (Optional, Int) Specifies the database port. Currently, only **GeminiDB Redis** instances support
  user-defined ports. If you do not specify a port number, port **6379** is used by default when you create a **GeminiDB
  Redis** instance. If you want to use this instance for dual-active DR, set the port to **8635**. Values: the value
  ranges from **1024** to **65535**. The disabled ports are **2180**, **2887**, **3887**, **6377**, **6378**, **6380**,
  **8018**, **8079**, **8091**, **8479**, **8484**, **8999**, **9864**, **9866**, **9867**, **12017**, **12333** and
  **50069**.

* `availability_zone_detail` - (Optional, List, NonUpdatable) Specifies Multi-AZ details of the active/standby instance.
  Currently, only **GeminiDB Redis** instances are supported. The system ignores this parameter if single-AZ deployment
  is selected.
  The [availability_zone_detail](#availability_zone_detail_struct) structure is documented below.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the instance.

* `delete_node_list` - (Optional, List) Specifies the ID of the nodes to be deleted. Make sure that the node
  can be deleted. This parameter can not be specified when add nodes. When delete nodes, if this parameter is specified,
  then the nodes specified by this parameter will be deleted, and id this parameter is not transferred, the nodes will be
  deleted randomly.

* `charging_mode` - (Optional, String, NonUpdatable) Specifies the charging mode of the instance. Valid values are
  **prePaid** and **postPaid**, defaults to **postPaid**.

* `period_unit` - (Optional, String, NonUpdatable) Specifies the charging period unit of the instance. Valid values are
  **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.

* `period` - (Optional, Int, NonUpdatable) Specifies the charging period of the instance.  
  If `period_unit` is set to **month** , the value ranges from `1` to `9`.  
  If `period_unit` is set to **year**, the value ranges from `1` to `3`.  
  This parameter is mandatory if `charging_mode` is set to **prePaid**.

* `auto_renew` - (Optional, String) Specifies whether auto-renew is enabled. Valid values are **true** and **false**.

<a name="datastore_struct"></a>
The `datastore` block supports:

* `type` - (Required, String, NonUpdatable) Specifies the database type. Value options:
  + **cassandra**: GeminiDB Cassandra instances will be created.
  + **mongodb**: GeminiDB Mongo instances will be created.
  + **influxdb**: GeminiDB Influx instances will be created.
  + **redis**: GeminiDB Redis instances will be created.
  + **dynamodb**: GeminiDB DynamoDB-Compatible instances will be created.
  + **hbase**: GeminiDB HBase instances will be created.

* `version` - (Required, String, NonUpdatable) Specifies the database version. Value options:
  + **3.11**: GeminiDB Cassandra instance 3.11.
  + **4.0**: GeminiDB Mongo instance 4.0.
  + **1.8**: GeminiDB Influx cluster instance 1.8 with classic storage.
  + **1.8**: GeminiDB Influx cluster (performance-enhanced) instance 1.8 with classic storage.
  + **1.7**: GeminiDB Influx cluster (performance-enhanced) instance 1.7 with cloud native storage.
  + **5.0**: GeminiDB Redis instance 5.0.
  + For a GeminiDB DynamoDB-Compatible instance, the value is empty string.
  + For a GeminiDB HBase instance, the value is empty string.

* `storage_engine` - (Required, String, NonUpdatable) Specifies the storage engine. Value options:
  + **rocksDB**: A GeminiDB Cassandra instance supports RocksDB.
  + **rocksDB**: A GeminiDB Mongo instance supports RocksDB.
  + **rocksDB**: A GeminiDB Influx instance supports RocksDB.
  + **rocksDB**: A GeminiDB Redis instance supports RocksDB.

<a name="flavor_struct"></a>
The `flavor` block supports:

* `num` - (Required, Int) Specifies the node quantity.
  + Each GeminiDB Cassandra instance can contain 3 to 60 nodes.
  + Each GeminiDB Mongo replica set 4.0 can contain 3 nodes.
  + Each GeminiDB Influx cluster instance can contain 3 to 16 nodes.
  + Each GeminiDB Influx single-node instance can contain 1 node.
  + Each GeminiDB Redis instance can contain 3 to 12 nodes.

* `size` - (Required, Int) Specifies the disk size. For GeminiDB Cassandra, GeminiDB Mongo, and GeminiDB Influx instances,
  the minimum storage space is **100 GB**, and the maximum limit depends on instance specifications. The maximum and
  minimum storage space of a GeminiDB Redis instance depends on node quantity and specifications of the instance.
  + For details about GeminiDB Cassandra instances, see
    [Instance Specifications](https://support.huaweicloud.com/intl/en-us/cassandraug-nosql/nosql_05_0017.html).
  + For details about GeminiDB Mongo instances, see
    [Instance Specifications](https://support.huaweicloud.com/intl/en-us/mongoug-nosql/nosql_05_0029.html).
  + For details about GeminiDB Influx instances, see
    [Instance Specifications](https://support.huaweicloud.com/intl/en-us/influxug-nosql/nosql_05_0045.html).
  + For details about GeminiDB Redis instances, see
    [Instance Specifications](https://support.huaweicloud.com/intl/en-us/redisug-nosql/nosql_05_0059.html).

* `storage` - (Required, String, NonUpdatable) Specifies the disk type. Value options: **ULTRAHIGH**.

* `spec_code` - (Required, String) Specifies the resource specification code.

<a name="backup_strategy_struct"></a>
The `backup_strategy` block supports:

* `start_time` - (Required, String) Specifies the backup time window. Automated backups will be triggered during the
  backup time window. It must be a valid value in the **hh:mm-HH:MM** format. The current time is in the UTC format. The
  **HH** value must be **1** greater than the **hh** value. The values of **mm** and **MM** must be the same and must be
  set to **00**. Example value: **08:00-09:00**, **03:00-04:00**.

* `keep_days` - (Optional, Int) Specifies the number of days to retain the generated backup files. The value ranges from
  **0** to **35**. If this parameter is set to **0**, the automated backup policy is not set. If this parameter is not
  transferred, the automated backup policy is enabled by default. Backup files are stored for seven days by default.

<a name="availability_zone_detail_struct"></a>
The `availability_zone_detail` block supports:

* `primary_availability_zone` - (Required, String, NonUpdatable) Specifies the primary AZ, which must be a single AZ and
  different from the standby AZ.

* `secondary_availability_zone` - (Required, String, NonUpdatable) Specifies the standby AZ, which must be a single AZ
  and be different from the primary AZ.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the DB instance ID.

* `datastore` - Specifies the database information.
  The [datastore](#datastore_struct) structure is documented below.

* `status` - Indicates the DB instance status. The value can be:
  + **normal**: An instance is available.
  + **abnormal**: An instance is abnormal.
  + **creating**: An instance is being created.
  + **frozen**: An instance is frozen.
  + **data_disk_full**: An instance disk is full.
  + **createfail**: An instance failed to be created.
  + **enlargefail**: Instance nodes failed to be added.

* `db_user_name` - Indicates the default username.

* `groups` - Indicates the group information.
  The [groups](#groups_struct) structure is documented below.

* `time_zone` - Indicates the time zone.

* `actions` - Indicates the operation that is executed on the instance.

* `lb_ip_address` - Indicates the load balancer IP address. This parameter is returned only when a load balancer IP
  address is assigned.

* `lb_port` - Indicates the load balancer port. This parameter is returned only when a load balancer IP address is assigned.

* `dual_active_info` - Indicates the active-active instance information.
  The [dual_active_info](#dual_active_info_struct) structure is documented below.

* `created` - Indicates the instance creation time.

* `updated` - Indicates the time when an instance is updated.

<a name="datastore_struct"></a>
The `datastore` block supports:

* `patch_available` - Indicates whether the current instance can be patched. The value can be:
  + **true**: A database can be upgraded through a patching API.
  + **false**: A database cannot be upgraded through a patching API.

* `whole_version` - Indicates the whole version of a GeminiDB instance.

<a name="groups_struct"></a>
The `groups` block supports:

* `id` - Indicates the group ID.

* `status` - Indicates the group status. The value can be:
  + **normal**: A group is normal.
  + **abnormal**: A group is abnormal.
  + **creating**: A group is being created.
  + **createfail**: A group failed to be created.
  + **deleted**: A group is deleted.
  + **resizefailed**: The group specifications failed to be modified.
  + **enlargefail**: A group failed to be scaled out.

* `volume` - Indicates the volume information.
  The [volume](#volume_struct) structure is documented below.

* `nodes` - Indicates the node information.
  The [nodes](#nodes_struct) structure is documented below.

<a name="volume_struct"></a>
The `volume` block supports:

* `size` - Indicates the storage (GB).

* `used` - Indicates the used storage (GB).

<a name="nodes_struct"></a>
The `nodes` block supports:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `status` - Indicates the node status. The value can be:
  + **normal**: A node is normal.
  + **abnormal**: A node is abnormal.
  + **creating**: A node is being created.
  + **createfail**: A node failed to be created.
  + **deleted**: A node is deleted.
  + **resizefailed**: Node specifications failed to be changed.
  + **enlargefail**: Nodes failed to be added.

* `role` - Indicates the node role. This parameter is available only for GeminiDB MongoAPI replica set instances.

* `subnet_id` - Indicates the ID of the subnet where the instance node is deployed.

* `private_ip` - Indicates the private IP address of a node. This parameter value is available after an ECS is created.
  Otherwise, the value is "".

* `public_ip` - Indicates the public IP address of a node. This parameter is valid only for nodes bound with EIPs.

* `spec_code` - Indicates the resource specification code.

* `availability_zone` - Indicates the AZ.

* `support_reduce` - Indicates whether instance nodes can be deleted.

<a name="dual_active_info_struct"></a>
The `dual_active_info` block supports:

* `role` - Indicates the active-active role.

* `status` - Indicates the active-active status.

* `destination_instance_id` - Indicates the ID of the peer instance in the active-active pair.

* `destination_region` - Indicates the peer region of the active-active pair.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 120 minutes.
* `delete` - Default is 30 minutes.

## Import

The GeminiDB instance can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_geminidb_instance.test <id>
```

Note that the imported state may not be identical to your resource definition, due to the attribute missing from the
API response. The missing attribute is: `password`, `flavor.0.storage`, `ssl_option`, `delete_node_list`, `auto_renew`,
`period` and `period_unit`. It is generally recommended running `terraform plan` after importing a GeminiDB instance.
You can then decide if changes should be applied to the GeminiDB instance or the resource definition should be updated
to align with the GeminiDB instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_geminidb_instance" "test" {
  ...

  lifecycle {
    ignore_changes = [
      password, flavor.0.storage, ssl_option, delete_node_list, auto_renew, period, period_unit
    ]
  }
}
```
