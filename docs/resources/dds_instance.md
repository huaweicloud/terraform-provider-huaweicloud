---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_instance"
description: ""
---

# huaweicloud_dds_instance

Manages dds instance resource within HuaweiCloud.

## Example Usage: Creating a Cluster Community Edition

```hcl
variable "dds_password" {}

resource "huaweicloud_dds_instance" "instance" {
  name = "dds-instance"
  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }

  availability_zone = "{{ availability_zone }}"
  vpc_id            = "{{ vpc_id }}"
  subnet_id         = "{{ subnet_network_id }}}"
  security_group_id = "{{ security_group_id }}"
  password          = var.dds_password
  mode              = "Sharding"
  maintain_begin    = "02:00"
  maintain_end      = "03:00"

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

## Example Usage: Creating a Replica Set Community Edition

```hcl
variable "dds_password" {}

resource "huaweicloud_dds_instance" "instance" {
  name = "dds-instance"
  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }

  availability_zone = "{{ availability_zone }}"
  vpc_id            = "{{ vpc_id }}"
  subnet_id         = "{{ subnet_network_id }}}"
  security_group_id = "{{ security_group_id }}"
  password          = var.dds_password
  mode              = "ReplicaSet"
  flavor {
    type      = "replica"
    num       = 3
    storage   = "ULTRAHIGH"
    size      = 30
    spec_code = "dds.mongodb.c3.medium.4.repset"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region of the DDS instance. Changing this creates a new
  instance.

* `name` - (Required, String) Specifies the DB instance name. The DB instance name of the same type is unique in the
  same tenant.

* `datastore` - (Required, List, ForceNew) Specifies database information. The structure is described below. Changing
  this creates a new instance.

* `availability_zone` - (Required, String) Specifies the availability zone names separated by commas.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID. Changing this creates a new instance.

* `subnet_id` - (Required, String, ForceNew) Specifies the subnet Network ID. Changing this creates a new instance.

* `security_group_id` - (Required, String) Specifies the security group ID of the DDS instance.

* `password` - (Optional, String) Specifies the Administrator password of the database instance.

* `disk_encryption_id` - (Optional, String, ForceNew) Specifies the disk encryption ID of the instance. Changing this
  creates a new instance.

* `mode` - (Required, String, ForceNew) Specifies the mode of the database instance. **Sharding**, **ReplicaSet**
  are supported. Changing this creates a new instance.

* `configuration` - (Optional, List, ForceNew) Specifies the configuration information.
  The structure is described below. Changing this creates a new instance.

* `flavor` - (Required, List, ForceNew) Specifies the flavors information. The structure is described below. Changing
  this creates a new instance.

* `port` - (Optional, Int) Specifies the database access port. The valid values are range from `2100` to `9500` and
  `27017`, `27018`, `27019`. Defaults to `8635`.

* `backup_strategy` - (Optional, List) Specifies the advanced backup policy. The structure is described below.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id of the DDS instance.

* `description` - (Optional, String) Specifies the description of the DDS instance.

* `replica_set_name` - (Optional, String) Specifies the name of the replica set in the connection address.
  It must be `3` to `128` characters long and start with a letter. It is case-sensitive and can contain only letters,
  digits, and underscores (_). Default is **replica**.

* `client_network_ranges` - (Optional, List) Specifies the CIDR block where the client is located. Cross-CIDR access is
  required only when the CIDR blocks of the client and the replica set instance are different. For example, if the client
  CIDR block is 192.168.0.0/16 and the replica set instance's CIDR block is 172.16.0.0/24, add the CIDR block
  192.168.0.0/16 so that the client can access the replica set instance.
  It's only for replica set instance.

* `ssl` - (Optional, Bool) Specifies whether to enable or disable SSL. Defaults to true.

**NOTE:** The instance will be restarted in the background when switching SSL. Please operate with caution.

* `maintain_begin` - (Optional, String) Specifies begin time of the time range within which you are allowed to start a
  task that affects the running of database instances. It must be a valid value in the format of **hh:mm** in UTC+0,
  such as **02:00**, meanwhile, this time in console displays in the format of **hh:mm** in UTC+08:00, e.g. **10:00**.

* `maintain_end` - (Optional, String) Specifies end time of the time range within which you are allowed to start a
  task that affects the running of database instances. It must be a valid value in the format of **hh:mm** in UTC+0,
  such as **04:00**, meanwhile, this time in console displays in the format of **hh:mm** in UTC+08:00, e.g. **12:00**.

* `second_level_monitoring_enabled` - (Optional, Bool) Specifies whether to enable second level monitoring.

* `slow_log_desensitization` - (Optional, String) Specifies whether to enable slow original log.
  The value can be **on** or **off**.

* `balancer_status` - (Optional, String) Specifies the status of the balancer.
  The value can be **start** or **stop**. Defaults to **start**.

* `balancer_active_begin` - (Optional, String) Specifies the start time of the balancing activity time window.
  The format is **HH:MM**. It's required with `balancer_active_end`.

* `balancer_active_end` - (Optional, String) Specifies the end time of the balancing activity time window.
  The format is **HH:MM**. It's required with `balancer_active_begin`.

-> It's only for **Sharding** mode. DDS 4.0 and later DB instances do not support to set balancer configuration.
  The UTC time is used. Please convert the local time based on the time zone.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the instance.
  The valid values are as follows:
  + `prePaid`: indicates the yearly/monthly billing mode.
  + `postPaid`: indicates the pay-per-use billing mode.

  Default value is `postPaid`.
  Changing this creates a new instance.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the instance.
  Valid values are *month* and *year*. This parameter is mandatory if `charging_mode` is set to *prePaid*.
  Changing this creates a new instance.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the instance.
  If `period_unit` is set to *month*, the value ranges from 1 to 9.
  If `period_unit` is set to *year*, the value ranges from 1 to 3.
  This parameter is mandatory if `charging_mode` is set to *prePaid*.
  Changing this creates a new instance.

* `auto_renew` - (Optional, String, ForceNew) Specifies whether auto-renew is enabled.
  Valid values are `true` and `false`, defaults to `false`.
  Changing this creates a new instance.

* `tags` - (Optional, Map) The key/value pairs to associate with the DDS instance.

The `datastore` block supports:

* `type` - (Required, String, ForceNew) Specifies the DB engine. **DDS-Community** is supported.

* `version` - (Required, String, ForceNew) Specifies the DB instance version. For the Community Edition, the valid
  values are `4.0`, `4.2`, `4.4` or `5.0`.

* `storage_engine` - (Optional, String, ForceNew) Specifies the storage engine of the DB instance.
  If `version` is set to `4.0`, the value is **wiredTiger**.
  If `version` is set to `4.2`, `4.4` or `5.0`, the value is **rocksDB**.

The `configuration` block supports:

* `type` - (Required, String, ForceNew) Specifies the node type. Valid value:
  + For a Community Edition cluster instance, the value can be **mongos**, **shard** or **config**.
  + For a Community Edition replica set instance, the value is **replica**.
    Changing this creates a new instance.

* `id` - (Required, String) Specifies the ID of the template.

  -> Atfer updating the `configuration.id`, please check whether the instance needs to be restarted.

The `flavor` block supports:

* `type` - (Required, String, ForceNew) Specifies the node type. Valid value:
  + For a cluster instance, the value can be **mongos**, **shard**, or **config**.
  + For a replica set instance, the value is **replica**.

* `num` - (Required, Int) Specifies the node quantity. Valid value:
  + If the value of type is **mongos**, num indicates the number of mongos nodes in the cluster instance. Value ranges
    from `2` to `16`.
  + If the value of type is **shard**, num indicates the number of shard groups in the cluster instance. Value ranges
    from `2` to `16`.
  + If the value of type is **config**, num indicates the number of config groups in the cluster instance. Value can
    only be `1`.
  + If the value of type is **replica**, num indicates the number of replica nodes in the replica set instance. Value
    can be `3`, `5`, or `7`.

  This parameter can be updated when the value of `type` is **mongos**, **shard** or **replica**.

* `storage` - (Optional, String, ForceNew) Specifies the disk type. Valid value:
  + **ULTRAHIGH**: SSD storage.
  + **EXTREMEHIGH**: Extreme SSD storage.

  This parameter is valid for the shard and config nodes of a cluster instance and for replica set instances.

* `size` - (Optional, Int) Specifies the disk size. The value must be a multiple of `10`. The unit is GB. This parameter
  is mandatory for nodes except mongos and invalid for mongos.For a cluster instance, the storage space of a shard node
  can be `10` to `2,000` GB, and the config storage space is `20` GB. For a replica set instance, the value ranges
  from `10` to `3000` GB. This parameter can be updated when the value of `type` is shard or replica.

* `spec_code` - (Required, String) Specifies the resource specification code. In a cluster instance, multiple
  specifications need to be specified. All specifications must be of the same series, that is, general-purpose (s6),
  enhanced (c3), or enhanced II (c6). For example:
  + dds.mongodb.s6.large.4.mongos and dds.mongodb.s6.large.4.config have the same specifications.
  + dds.mongodb.s6.large.4.mongos and dds.mongodb.c3.large.4.config are not of the same specifications.

The `backup_strategy` block supports:

* `start_time` - (Required, String) Specifies the backup time window. Automated backups will be triggered during
  the backup time window. The value cannot be empty. It must be a valid value in the "hh:mm-HH:MM" format.
  The current time is in the UTC format.
  + The HH value must be 1 greater than the hh value.
  + The values from mm and MM must be the same and must be set to **00**.

* `keep_days` - (Required, Int) Specifies the number of days to retain the generated backup files. The value range is
  from 0 to 732. If this parameter is set to 0, the automated backup policy is disabled.

* `period` - (Optional, String) Specifies the backup cycle. Data will be automatically backed up on the
  selected days every week.
  + If you set the `keep_days` to 0, this parameter is no need to set.
  + If you set the `keep_days` within 6 days, set the parameter value to **1,2,3,4,5,6,7**, data is automatically
    backed up on each day every week.
  + If you set the `keep_days` between 7 and 732 days, set the parameter value to at least one day of every week.
    For example: **1**, **3,5**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the the DB instance ID.
* `db_username` - Indicates the DB Administrator name.
* `status` - Indicates the the DB instance status.
* `port` - Indicates the database port number. The port range is 2100 to 9500.
* `groups` - Indicates the instance groups information.
  The [groups](#DdsInstance_InstanceGroup) structure is documented below.
* `created_at` - Indicates the create time.
* `updated_at` - Indicates the update time.
* `time_zone` - Indicates the time zone.

<a name="DdsInstance_InstanceGroup"></a>
The `groups` block supports:

* `id` - Indicates the group ID.
* `type` - Indicates the node type.
* `name` - Indicates the group name.
* `status` - Indicates the group status.
* `size` - Indicates the disk size.
* `used` - Indicates the disk usage.
* `nodes` - Indicates the nodes info.
  The [nodes](#DdsInstance_InstanceGroupNode) structure is documented below.

<a name="DdsInstance_InstanceGroupNode"></a>
The `nodes` block supports:

* `id` - Indicates the node ID.
* `name` - Indicates the node name.
* `role` - Indicates the node role.
* `private_ip` - Indicates the private IP address of a node. This parameter is valid only for mongos nodes, replica set
  instances, and single node instances.
* `public_ip` - Indicates the EIP that has been bound on a node. This parameter is valid only for mongos nodes of
  cluster instances, primary nodes and secondary nodes of replica set instances, and single node instances.
* `status` - Indicates the node status.
* `spec_code` - Indicates the node spec code.
* `availability_zone` - Indicates the availability zone.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 60 minutes.
* `delete` - Default is 60 minutes.

## Import

DDS instance can be imported using the `id`, e.g.

```sh
terraform import huaweicloud_dds_instance.instance 9c6d6ff2cba3434293fd479571517e16in02
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `password`, `availability_zone`, `flavor`, configuration.
It is generally recommended running `terraform plan` after importing an instance.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dds_instance" "instance" {
    ...

  lifecycle {
    ignore_changes = [
      password, availability_zone, flavor, configuration,
    ]
  }
}
```
