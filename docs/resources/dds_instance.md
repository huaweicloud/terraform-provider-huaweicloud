---
subcategory: "Document Database Service (DDS)"
---

# huaweicloud\_dds\_instance

Manages dds instance resource within HuaweiCloud
This is an alternativeto `huaweicloud_dds_instance_v3`

## Example Usage: Creating a Cluster Community Edition

```hcl
resource "huaweicloud_dds_instance" "instance" {
  name = "dds-instance"
  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  availability_zone = "{{ availability_zone }}"
  vpc_id            = "{{ vpc_id }}"
  subnet_id         = "{{ subnet_network_id }}}"
  security_group_id = "{{ security_group_id }}"
  password          = "Test@123"
  mode              = "Sharding"
  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.c3.medium.4.mongos"
  }
  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c3.medium.4.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c3.large.2.config"
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
  }
}
```

## Example Usage: Creating a Replica Set
```hcl
resource "huaweicloud_dds_instance" "instance" {
  name = "dds-instance"
  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  availability_zone = "{{ availability_zone }}"
  vpc_id            = "{{ vpc_id }}"
  subnet_id         = "{{ subnet_network_id }}}"
  security_group_id = "{{ security_group_id }}"
  password          = "Test@123"
  mode              = "ReplicaSet"
  flavor {
    type      = "replica"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 30
    spec_code = "dds.mongodb.c3.medium.4.repset"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) Specifies the region of the DDS instance. Changing this creates
	a new instance.

* `name` - (Required) Specifies the DB instance name. The DB instance name of the same
	type is unique in the same tenant.

* `datastore` - (Required) Specifies database information. The structure is described
	below. Changing this creates a new instance.

* `availability_zone` - (Required) Specifies the ID of the availability zone. Changing
	this creates a new instance.

* `vpc_id` - (Required) Specifies the VPC ID. Changing this creates a new instance.

* `subnet_id` - (Required) Specifies the subnet Network ID. Changing this creates a new instance.

* `security_group_id` - (Required) Specifies the security group ID of the DDS instance.

* `password` - (Required) Specifies the Administrator password of the database instance.

* `disk_encryption_id` - (Required) Specifies the disk encryption ID of the instance.
	Changing this creates a new instance.

* `mode` - (Required) Specifies the mode of the database instance. Changing this creates
	a new instance.

* `flavor` - (Required) Specifies the flavors information. The structure is described below.
	Changing this creates a new instance.

* `backup_strategy` - (Optional) Specifies the advanced backup policy. The structure is
	described below. Changing this creates a new instance.

* `ssl` - (Optional) Specifies whether to enable or disable SSL. Defaults to true.

**Note:** The instance will be restarted in the background when switching SSL. Please operate with caution.

The `datastore` block supports:

* `type` - (Required) Specifies the DB engine. 'DDS-Community' and 'DDS-Enhanced' are supported.

* `version` - (Required) Specifies the DB instance version. For the Community Edition,
  the valid values are 3.2, 3.4, or 4.0. For the Enhanced Edition, only 3.4 is supported now.

* `storage_engine` - (Optional) Specifies the storage engine of the DB instance. 
  DDS Community Edition supports wiredTiger engine, and the Enhanced Edition supports rocksDB engine.

The `flavor` block supports:

* `type` - (Required) Specifies the node type. Valid value:
  * For a Community Edition cluster instance, the value can be mongos, shard, or config.
  * For an Enhanced Edition cluster instance, the value is shard.
  * For a Community Edition replica set instance, the value is replica.
  * For a Community Edition single node instance, the value is single.

* `num` - (Required) Specifies the node quantity. Valid value:
	* In a Community Edition cluster instance,the number of mongos ranges from 2 to 16.
  * In a Community Edition cluster instance,the number of shards ranges from 2 to 16.
  * In an Enhanced Edition cluster instance, the number of shards ranges from 2 to 12.
	* config: the value is 1.
	* replica: the value is 1.
  * single: The value is 1.

* `storage` - (Optional) Specifies the disk type. Valid value: ULTRAHIGH which indicates the type SSD.

* `size` - (Optional) Specifies the disk size. The value must be a multiple of 10. The unit is GB.
  This parameter is mandatory for nodes except mongos and invalid for mongos.

* `spec_code` - (Required) Specifies the resource specification code. In a cluster instance,
  multiple specifications need to be specified. All specifications must be of the same series,
  that is, general-purpose (s6), enhanced (c3), or enhanced II (c6). For example:
  * dds.mongodb.s6.large.4.mongos and dds.mongodb.s6.large.4.config have the same specifications.
  * dds.mongodb.s6.large.4.mongos and dds.mongodb.c3.large.4.config are not of the same specifications.

The `backup_strategy ` block supports:

* `start_time` - (Required) Specifies the backup time window. Automated backups will be triggered
	during the backup time window. The value cannot be empty. It must be a valid value in the
	"hh:mm-HH:MM" format. The current time is in the UTC format.
	* The HH value must be 1 greater than the hh value.
	* The values from mm and MM must be the same and must be set to any of the following 00, 15, 30, or 45.

* `keep_days` - (Optional) Specifies the number of days to retain the generated backup files. The
	value range is from 0 to 732.
	* If this parameter is set to 0, the automated backup policy is not set.
	* If this parameter is not transferred, the automated backup policy is enabled by default.
    Backup files are stored for seven days by default.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `datastore` - See Argument Reference above.
* `availability_zone` - See Argument Reference above.
* `vpc_id` - See Argument Reference above.
* `subnet_id` - See Argument Reference above.
* `security_group_id` - See Argument Reference above.
* `password` - See Argument Reference above.
* `disk_encryption_id` - See Argument Reference above.
* `ssl` - See Argument Reference above.
* `mode` - See Argument Reference above.
* `flavor` - See Argument Reference above.
* `backup_strategy` - See Argument Reference above.
* `db_username` - Indicates the DB Administator name.
* `status` - Indicates the the DB instance status.
* `port` - Indicates the database port number. The port range is 2100 to 9500.
* `nodes` - Indicates the instance nodes information. Structure is documented below.

The `nodes` block contains:

  - `id` - Indicates the node ID.
  - `name` - Indicates the node name.
  - `role` - Indicates the node role.
  - `type` - Indicates the node type.
  - `private_ip` - Indicates the private IP address of a node. This parameter is valid only for
     mongos nodes, replica set instances, and single node instances.
  - `public_ip` - Indicates the EIP that has been bound on a node. This parameter is valid only for
     mongos nodes of cluster instances, primary nodes and secondary nodes of replica set instances,
     and single node instances.
  - `status` - Indicates the node status.
