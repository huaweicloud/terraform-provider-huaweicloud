---
subcategory: "Distributed Cache Service"
---

# huaweicloud_dcs_redis_instance

Manages a DCS redis instance resource within HuaweiCloud.
This is an alternative to `huaweicloud_dcs_instance`.

## Example Usage

### Single Instance

```hcl
variable az_code {}
variable vpc_id {}
variable subnet_id {}

resource "huaweicloud_dcs_redis_instance" "redis_single" {
  name               = "redis_name"
  engine_version     = "5.0"
  capacity           = 0.125
  resource_spec_code = "redis.ha.xu1.tiny.r2.128"
  available_zones    = [var.az]
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
  password           = "YourPassword_123"

  charging_mode = "prePaid"
  period_unit   = "month"
  auto_renew    = "true"
  period        = "1"

  whitelists {
    group_name = "group_1"
    ip_address = ["192.168.1.0/24"]
  }
}
```

### Master/Standby Instance with Backup Policy

```hcl
variable primary_az {}
variable standby_az {}
variable vpc_id {}
variable subnet_id {}

resource "huaweicloud_dcs_redis_instance" "redis_master_standby" {
  name               = "redis_name"
  engine_version     = "5.0"
  capacity           = 4
  resource_spec_code = "redis.ha.xu1.large.r2.4"
  available_zones    = [var.primary_az, var.standby_az]
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
  password           = "YourPassword_123"

  charging_mode = "prePaid"
  period_unit   = "month"
  auto_renew    = "true"
  period        = "1"

  backup_policy {
    backup_type     = "auto"
    save_days       = 5
    period_type     = "weekly"
    backup_at       = [1, 2, 3, 4, 5, 6, 7]
    begin_at        = "02:00-04:00"
  }

  whitelists {
    group_name = "group_1"
    ip_address = ["192.168.1.0/24"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the DCS redis instance resource.
  If omitted, the provider-level region will be used. Changing this creates a new DCS instance resource.

* `name` - (Required, String) Specifies the name of an instance.
  The name must be 4 to 64 characters and start with a letter.
  Only chinese, letters (case-insensitive), digits, underscores (_) ,and hyphens (-) are allowed.

* `engine_version` - (Required, String, ForceNew) Specifies the version of a cache engine.
  Currently, the valid values are `3.0`, `4.0` and `5.0`. Changing this creates a new instance.

* `capacity` - (Required, Float) The total memory of the cache, in GB.
  + **Redis4.0 and Redis5.0**: Stand-alone and active/standby type instance values: `0.125`, `0.25`, `0.5`, `1`, `2`,
    `4`, `8`, `16`, `32` and `64`.
    Cluster instance specifications support `24`, `32`, `48`, `64`, `96`, `128`, `192`, `256`, `384`, `512`, `768` and
    `1024`.
  + **Redis3.0**: Stand-alone and active/standby type instance values: `2`, `4`, `8`, `16`, `32` and `64`.
    Proxy cluster instance specifications support `64`, `128`, `256`, `512`, and `1024`.
  + **Memcached**: Stand-alone and active/standby type instance values: `2`, `4`, `8`, `16`, `32` and `64`.

* `resource_spec_code` - (Required, String) The specification code of the cache instance, which including the total
  memory, available memory, maximum number of connections allowed, maximum/assured bandwidth and reference performance.
  You can specify specification codes to create different type of Redis instances.
  It can be obtained through this data source `huaweicloud_dcs_flavors`.

* `available_zones` - (Required, List, ForceNew) The code of the AZ where the cache node resides.
  Master/Standby, Proxy Cluster, and Redis Cluster DCS instances support cross-AZ deployment.
  You can specify an AZ for the standby node. When specifying AZs for nodes, use commas (,) to separate AZs.
  Changing this creates a new instance.

* `vpc_id` - (Required, String, ForceNew) The ID of VPC which the instance belongs to.
  Changing this creates a new instance resource.

* `subnet_id` - (Required, String, ForceNew) The ID of subnet which the instance belongs to.
  Changing this creates a new instance resource.

* `security_group_id` - (Optional, String) The ID of the security group which the instance belongs to.
  This parameter is mandatory for Redis 3.0 version.

* `ip` - (Optional, String, ForceNew) The IP address of the DCS instance, which can only be the currently available IP
  address the selected subnet.
  You can specify an available IP for the Redis instance (except for the Redis Cluster type).
  If omitted, the system will automatically allocate an available IP address to the Redis instance.
  Changing this creates a new instance resource.

* `port` - (Optional, Int) Port customization, which is supported only by Redis 4.0 and Redis 5.0 instances.
  Default value is 6379.

* `password` - (Optional, String) Specifies the password of a DCS instance.
  The password of a DCS Redis instance must meet the following complexity requirements:
  + Must be a string of 8 to 32 bits in length.
  + Must contain three combinations of the following four characters: Lower case letters, uppercase letter, digital,
    Special characters include (`~!@#$%^&*()-_=+|[{}]:'",<.>/?).
  + The new password cannot be the same as the old password.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project id of the dcs instance.
  Changing this creates a new instance.

* `description` - (Optional, String) Specifies the description of an instance.
  It is a string that contains a maximum of 1,024 characters.

* `whitelist_enable` - (Optional, Bool) Enable or disable the IP address whitelists. Defaults to true.
  If the whitelist is disabled, all IP addresses connected to the VPC can access the instance.

* `whitelists` - (Optional, List) Specifies the IP addresses which can access the instance.
  This parameter is valid for Redis 4.0 and 5.0 versions. The structure is described below.

* `maintain_begin` - (Optional, String) Time at which the maintenance time window starts.
  The valid values are `22:00:00`, `02:00:00`, `06:00:00`, `10:00:00`, `14:00:00` and `18:00:00`.
  Default value is `02:00:00`.
  + The start time and end time of a maintenance time window must indicate the time segment of a supported maintenance
    time window.
  + Parameters `maintain_begin` and `maintain_end` must be set in pairs.
  + If parameter maintain_begin is left blank, parameter maintain_end is also blank.
    In this case, the system automatically allocates the default start time 02:00:00.

* `maintain_end` - (Optional, String) Time at which the maintenance time window ends.
  The valid values are `22:00:00`, `02:00:00`, `06:00:00`, `10:00:00`, `14:00:00` and `18:00:00`.
  Default value is `06:00:00`.
  + The start time and end time of a maintenance time window must indicate the time segment of a supported maintenance
    time window.
  + The end time is four hours later than the start time.
    For example, if the start time is 22:00:00, the end time is 02:00:00.
  + Parameters `maintain_begin` and `maintain_end` must be set in pairs.
  + If parameter maintain_end is left blank, parameter maintain_begin is also blank.
    In this case, the system automatically allocates the default end time 06:00:00.

* `backup_policy` - (Optional, List) Specifies the backup configuration to be used with the instance.
  The structure is described below.

* `rename_commands` - (Optional, Map, ForceNew) Critical command renaming, which is supported only by Redis 4.0 and
  Redis 5.0 instances but not by Redis 3.0 instance.
  The valid commands that can be renamed are: *command*, *keys*, *flushdb*, *flushall* and *hgetall*.

* `tags` - (Optional, Map) The key/value pairs to associate with the redis instance.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the redis instance.
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

* `auto_renew` - (Optional, String, ForceNew) Specifies whether auto renew is enabled.
  Valid values are `true` and `false`, defaults to `false`.
  Changing this creates a new instance.

The `whitelists` block supports:

* `group_name` - (Required, String) Specifies the name of IP address group.

* `ip_address` - (Required, List) Specifies the list of IP address or CIDR which can be whitelisted for an instance.
  The maximum of ip address is 20.

The `backup_policy` block supports:

* `backup_type` - (Optional, String) Backup type. Default value is `auto`. The valid values are as follows:
  + `auto`: automatic backup.
  + `manual`: manual backup.

* `save_days` - (Required, Int) Retention time. Unit: day, the value ranges from 1 to 7.

* `period_type` - (Optional, String) Interval at which backup is performed. Default value is `weekly`.
  Currently, only weekly backup is supported.

* `backup_at` - (Required, List) Day in a week on which backup starts, the value ranges from 1 to 7.
  Where: 1 indicates Monday; 7 indicates Sunday.

* `begin_at` - (Required, String) Time at which backup starts.
  Format: `hh24:00-hh24:00`, "00:00-01:00" indicates that backup starts at 00:00:00.

* `timezone_offset` - (Optional, String) Time zone in which backup is performed.
  The value ranges from GMT–12:00 to GMT+12:00.
  If omitted, the current time zone of the DCS-Server VM is used by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in UUID format.

* `status` - Cache instance status. The valid values are as follows:
  + `RUNNING`: The instance is running properly.
    Only instances in the Running state can provide in-memory cache service.
  + `ERRORv`: The instance is not running properly.
  + `RESTARTING`: The instance is being restarted.
  + `FROZEN`: The instance has been frozen due to low balance.
    You can unfreeze the instance by recharging your account in My Order.
  + `EXTENDING`: The instance is being scaled up.
  + `RESTORING`: The instance data is being restored.
  + `FLUSHING`: The DCS instance is being cleared.

* `domain_name` - Domain name of the instance. Usually, we use domain name and port to connect to the Redis.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minute.
* `update` - Default is 20 minute.
* `delete` - Default is 15 minute.

## Import

DCS redis instance can be imported using the `id`, e.g.

```sh
terraform import huaweicloud_dcs_redis_instance.redis_single b7a3b725-77c7-4662-af43-ad071b93707a
```
