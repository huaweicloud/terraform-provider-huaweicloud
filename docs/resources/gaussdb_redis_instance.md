---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_redis_instance"
description: ""
---

# huaweicloud_gaussdb_redis_instance

GeminiDB Redis instance management within HuaweiCould.

## Example Usage

### create a geminidb redis instance with tags

```hcl
resource "huaweicloud_gaussdb_redis_instance" "test" {
  name              = "gaussdb_redis_instance_1"
  password          = var.password
  flavor            = "geminidb.redis.xlarge.4"
  volume_size       = 16
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
  availability_zone = var.availability_zone

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### create a geminidb redis instance with backup strategy

```hcl
resource "huaweicloud_gaussdb_redis_instance" "test" {
  name              = "gaussdb_redis_instance_1"
  password          = var.password
  flavor            = "geminidb.redis.xlarge.4"
  volume_size       = 16
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
  availability_zone = var.availability_zone

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the Redis instance resource.
  See [Region and Endpoints](https://developer.huaweicloud.com/intl/en-us/endpoint?GaussDB%20NoSQL) for more detail. If
  omitted, the provider-level region will be used. Changing this creates a new Redis instance resource.

* `availability_zone` - (Required, String, ForceNew) Specifies the AZ name.
  See [Region and Endpoints](https://developer.huaweicloud.com/intl/en-us/endpoint?GaussDB%20NoSQL) for more detail.
  For a three-AZ deployment instance, use commas (,) to separate the AZs, for example, `cn-north-4a,cn-north-4b,cn-north-4c`.
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the instance name, which can be the same as an existing instance name. The value
  must be `4` to `64` characters in length and start with a letter. It is case-sensitive and can contain only letters,
  digits, hyphens (-), and underscores (_). Chinese characters must be in UTF-8 or Unicode format.

* `flavor` - (Required, String) Specifies the instance specifications. For details,
  see [DB Instance Specifications](https://support.huaweicloud.com/intl/en-us/redisug-nosql/nosql_05_0059.html).

* `node_num` - (Optional, Int) Specifies the number of nodes, ranges from `2` to `12`. Defaults to `3`.

* `volume_size` - (Required, Int) Specifies the storage space in GB. For a GaussDB for Redis instance, the minimum and
  maximum storage space depends on the flavor and nodes_num. For details,
  see [DB Instance Specifications](https://support.huaweicloud.com/intl/en-us/redisug-nosql/nosql_05_0059.html)

* `password` - (Required, String) Specifies the database password. The value must be `8` to `32` characters in length,
  including uppercase and lowercase letters, digits, and special characters, such as ~!@#%^*-_=+? You are advised to
  enter a strong password to improve security, preventing security risks such as brute force cracking.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID. Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the network ID of a subnet. Changing this parameter will create a
  new resource.

* `mode` - (Optional, String, ForceNew) Specifies the instance type. Value options: **Cluster**, **Replication**.
  Defaults to **Cluster**.

* `availability_zone_detail` - (Optional, List, ForceNew) Specifies Multi-AZ details of the active/standby instance.
  The system ignores this parameter if single-AZ deployment is selected. Currently, only GeminiDB Redis instances are
  supported.

  The [availability_zone_detail](#availability_zone_detail_struct) structure is documented below.

* `security_group_id` - (Optional, String) Specifies the security group ID. Required if the selected subnet doesn't
  enable network ACL.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID, Only valid for users who
  have enabled the enterprise multi-project service.

* `force_import` - (Optional, Bool) If specified, try to import the instance instead of creating if the name already
  existed.

* `datastore` - (Optional, List, ForceNew) Specifies the database information. Structure is documented below. Changing
  this parameter will create a new resource.

* `port` - (Optional, Int, ForceNew) Specifies the port number for accessing the instance. You can specify a port number
  based on your requirements. The port number ranges from `1,024` to `65,535`, excluding `2,180`, `2,887`, `3,887`,
  `6,377`, `6,378`, `6,380`, `8,018`, `8,079`, `8,091`, `8,479`, `8,484`, `8,999`, `12,017`, `12,333`, and `50,069`.
  Defaults to `6,379`. If you want to use this instance for dual-active DR, set the port to `8,635`.

  Changing this parameter will create a new resource.

* `backup_strategy` - (Optional, List) Specifies the advanced backup policy. Structure is documented below. Do nothing
  in update method if change this parameter.

* `ssl` - (Optional, Bool) Specifies whether SSL is enabled. Defaults to **false**.

* `tags` - (Optional, Map) The key/value pairs to associate with the instance.

* `charging_mode` - (Optional, String) Specifies the charging mode of the GaussDB for Redis instance. Valid values are
  **prePaid** and **postPaid**, defaults to **postPaid**. Do nothing in update method if change this parameter.

* `period_unit` - (Optional, String) Specifies the charging period unit of the GaussDB for Redis instance. Valid values
  are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**. Do nothing in update
  method if change this parameter.

* `period` - (Optional, Int) Specifies the charging period of the GaussDB for Redis instance.  
  If `period_unit` is set to **month** , the value ranges from `1` to `9`.  
  If `period_unit` is set to **year**, the value ranges from `1` to `3`.  
  This parameter is mandatory if `charging_mode` is set to **prePaid**.  
  Changing this will do nothing.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled. Valid values are "true" and "false".

<a name="availability_zone_detail_struct"></a>
The `availability_zone_detail` block supports:

* `primary_availability_zone` - (Required, String, ForceNew) Specifies the primary AZ, it must be a single AZ and be
  different from the standby AZ. Changing this parameter will create a new resource.

* `secondary_availability_zone` - (Required, String, ForceNew) Specifies the standby AZ, it must be a single AZ and be
  different from the primary AZ. Changing this parameter will create a new resource.

The `datastore` block supports:

* `engine` - (Required, String, ForceNew) Specifies the database engine. Only "redis" is supported now.
  Changing this parameter will create a new resource.

* `version` - (Required, String, ForceNew) Specifies the database version. Only "5.0" is supported now.
  Changing this parameter will create a new resource.

* `storage_engine` - (Required, String, ForceNew) Specifies the storage engine. Only "rocksDB" is supported now.
  Changing this parameter will create a new resource.

The `backup_strategy` block supports:

* `start_time` - (Required, String) Specifies the backup time window. Automated backups will be triggered during the
  backup time window. It must be a valid value in the "hh:mm-HH:MM" format. The current time is in the UTC format. The
  HH value must be `1` greater than the hh value. The values of mm and MM must be the same and must be set to 00. Example
  value: **08:00-09:00**, **03:00-04:00**.

* `keep_days` - (Optional, Int) Specifies the number of days to retain the generated backup files. The value ranges from
  0 to 35. If this parameter is set to `0`, the automated backup policy is not set. If this parameter is not transferred,
  the automated backup policy is enabled by default. Backup files are stored for seven days by default.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the DB instance ID.
* `status` - Indicates the DB instance status.
* `db_user_name` - Indicates the default username.
* `nodes` - Indicates the instance nodes information. Structure is documented below.
* `private_ips` - Indicates the IP address list of the db.
* `lb_ip_address` - Indicates the LB IP address of the db.
* `lb_port` - Indicates the LB port of the db.

The `nodes` block contains:

* `id` - Indicates the node ID.
* `name` - Indicates the node name.
* `status` - Indicates the node status.
* `support_reduce` - Indicates whether the node support reduce or not.
* `public_ip` - Indicates the public IP address of a node.
* `private_ip` - Indicates the private IP address of a node.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 120 minutes.
* `delete` - Default is 30 minutes.

## Import

GaussDB Redis instance can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_gaussdb_redis_instance.instance_1 d54b21f037ed447aad4bfd20927711c6in12
```
