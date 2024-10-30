---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_instance"
description: ""
---

# huaweicloud_gaussdb_mysql_instance

GaussDB mysql instance management within HuaweiCould.

## Example Usage

### create a basic instance

```hcl
resource "huaweicloud_gaussdb_mysql_instance" "instance_1" {
  name              = "gaussdb_instance_1"
  password          = var.password
  flavor            = "gaussdb.mysql.4xlarge.x86.4"
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
}
```

### create a gaussdb mysql instance with backup strategy

```hcl
resource "huaweicloud_gaussdb_mysql_instance" "instance_1" {
  name              = "gaussdb_instance_1"
  password          = var.password
  flavor            = "gaussdb.mysql.4xlarge.x86.4"
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 7
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the GaussDB mysql instance resource. If omitted,
  the provider-level region will be used. Changing this creates a new instance resource.

* `name` - (Required, String) Specifies the instance name, which can be the same as an existing instance name.
  The value must be `4` to `64` characters in length and start with a letter.
  It is case-sensitive and can contain only letters, digits, hyphens (-), and underscores (_).

* `flavor` - (Required, String) Specifies the instance specifications. Please use
  `gaussdb_mysql_flavors` data source to fetch the available flavors.

* `password` - (Required, String) Specifies the database password. The value must be `8` to `32` characters in length,
  including uppercase and lowercase letters, digits, and special characters, such as ~!@#%^*-_=+? You are advised to
  enter a strong password to improve security, preventing security risks such as brute force cracking.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID. Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the network ID of a subnet. Changing this parameter will create a
  new resource.

* `dedicated_resource_id` - (Optional, String, ForceNew) Specifies the dedicated resource ID. Changing this parameter
  will create a new resource.

* `dedicated_resource_name` - (Optional, String, ForceNew) Specifies the dedicated resource name. Changing this parameter
  will create a new resource.

* `security_group_id` - (Optional, String) Specifies the security group ID. Required if the selected subnet doesn't
  enable network ACL.

* `port` - (Optional, Int) Specifies the database port.

* `ssl_option` - (Optional, String) Specifies whether to enable SSL. Value options:
  + **true**: SSL is enabled.
  + **false**: SSL is disabled.

  Defaults to **true**.

* `slow_log_show_original_switch` - (Optional, Bool) Specifies the slow log show original switch of the instance.

* `description` - (Optional, String) Specifies the description of the instance.

* `private_write_ip` - (Optional, String) Specifies the private IP address of the DB instance.

* `configuration_id` - (Optional, String) Specifies the configuration ID.

* `private_dns_name_prefix` - (Optional, String) Specifies the prefix of the private domain name. The value contains
  `8` to `63` characters. Only uppercase letters, lowercase letters, and digits are allowed.

* `maintain_begin` - (Optional, String) Specifies the start time for a maintenance window, for example, **22:00**.

* `maintain_end` - (Optional, String) Specifies the end time for a maintenance window, for example, **01:00**.

-> **Note** The start time and end time of a maintenance window must be on the hour, and the interval between them at
  most four hours.

* `seconds_level_monitoring_enabled` - (Optional, Bool) Specifies whether to enable seconds level monitoring.

* `seconds_level_monitoring_period` - (Optional, Int) Specifies the seconds level collection period.
  + This parameter is valid only when `seconds_level_monitoring_enabled` is set to **true**.
  + This parameter can not be specified when `seconds_level_monitoring_enabled` is set to **false**.
  + Value options:
    - **1**: The collection period is 1s.
    - **5** (default value): The collection period is 5s.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id. Required if EPS enabled.

* `table_name_case_sensitivity` - (Optional, Bool) Whether the kernel table name is case sensitive. The value can
  be `true` (case sensitive) and `false` (case insensitive). Defaults to `false`. This parameter only works during
  creation.

* `read_replicas` - (Optional, Int) Specifies the count of read replicas. Defaults to `1`.

* `time_zone` - (Optional, String, ForceNew) Specifies the time zone. Defaults to "UTC+08:00". Changing this parameter
  will create a new resource.

* `availability_zone_mode` - (Optional, String, ForceNew) Specifies the availability zone mode: "single" or "multi".
  Defaults to "single". Changing this parameter will create a new resource.

* `master_availability_zone` - (Optional, String, ForceNew) Specifies the availability zone where the master node
  resides. The parameter is required in multi availability zone mode. Changing this parameter will create a new
  resource.

* `charging_mode` - (Optional, String) Specifies the charging mode of the instance. Valid values are **prePaid**
  and **postPaid**, defaults to **postPaid**. Changing this will do nothing.

* `period_unit` - (Optional, String) Specifies the charging period unit of the instance.
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.
  Changing this will do nothing.

* `period` - (Optional, Int) Specifies the charging period of the instance.  
  If `period_unit` is set to **month** , the value ranges from `1` to `9`.  
  If `period_unit` is set to **year**, the value ranges from `1` to `3`.  
  This parameter is mandatory if `charging_mode` is set to **prePaid**.  
  Changing this will do nothing.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled.
  Valid values are **true** and **false**.

* `datastore` - (Optional, List, ForceNew) Specifies the database information. Structure is documented below. Changing
  this parameter will create a new resource.

* `backup_strategy` - (Optional, List) Specifies the advanced backup policy. Structure is documented below.

* `parameters` - (Optional, List) Specifies an array of one or more parameters to be set to the instance after launched.
  The [parameters](#parameters_struct) structure is documented below.

* `auto_scaling` - (Optional, List) Specifies the auto-scaling policies.
  The [auto_scaling](#auto_scaling_struct) structure is documented below.

* `force_import` - (Optional, Bool) If specified, try to import the instance instead of creating if the name already
  existed.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the GaussDB Mysql instance.

* `volume_size` - (Optional, Int) Specifies the volume size of the instance. The new storage space must be greater than
  the current storage and must be a multiple of `10` GB. Only valid when in prePaid mode.

The `datastore` block supports:

* `engine` - (Required, String, ForceNew) Specifies the database engine. Only "gaussdb-mysql" is supported now.
  Changing this parameter will create a new resource.

* `version` - (Required, String, ForceNew) Specifies the database version. Only "8.0" is supported now.
  Changing this parameter will create a new resource.

The `backup_strategy` block supports:

* `start_time` - (Required, String) Specifies the backup time window. Automated backups will be triggered during the
  backup time window. It must be a valid value in the "hh:mm-HH:MM" format. The current time is in the UTC format. The
  HH value must be 1 greater than the hh value. The values of mm and MM must be the same and must be set to 00. Example
  value: **08:00-09:00**, **03:00-04:00**.

* `keep_days` - (Optional, Int) Specifies the number of days to retain the generated backup files.  
  The value ranges from `0` to `35`. If this parameter is set to `0`, the automated backup policy is not set.
  If this parameter is not transferred, the automated backup policy is enabled by default.
  Backup files are stored for seven days by default.

* `audit_log_enabled` - (Optional, Bool) Specifies whether audit log is enabled. The default value is `false`.

* `sql_filter_enabled` - (Optional, Bool) Specifies whether sql filter is enabled. The default value is `false`.

* `encryption_status` - (Optional, String) Specifies whether to enable or disable encrypted backup. Value options:
  + **ON**: enabled
  + **OFF**: disabled

* `encryption_type` - (Optional, String) Specifies the encryption type. Currently, only **kms (case-insensitive)** is
  supported. It is mandatory when `encryption_status` is set to **ON**.

* `kms_key_id` - (Optional, String) Specifies the KMS ID. It is mandatory when `encryption_status` is set to **ON**.

<a name="parameters_struct"></a>
The `parameters` block supports:

* `name` - (Required, String) Specifies the name of the parameter.

* `value` - (Required, String) Specifies the value of the parameter.

<a name="auto_scaling_struct"></a>
The `auto_scaling` block supports:

* `status` - (Required, String) Specifies whether auto-scaling is enabled. Value options:
  + **ON**: enabled.
  + **OFF**: disabled.

* `scaling_strategy` - (Required, List) Specifies the auto-scaling policy.
  The [scaling_strategy](#scaling_strategy_struct) structure is documented below.

* `monitor_cycle` - (Optional, Int) Specifies the observation period, in seconds. During the entire observation period,
  if the average CPU usage is greater than or equal to the preset value, a scale-up is triggered. It is mandatory when
  `status` is set to **ON**. Value options: **300**, **600**, **900** or **1800**.

* `silence_cycle` - (Optional, Int) Specifies the silent period, in seconds. It indicates the minimum interval between
  two auto scale-up operations or two scale-down operations. It is mandatory when `status` is set to **ON**. Value
  options: **300**,  **600**, **1800**, **3600**, **7200**, **10800**, **86400** or **604800**.

* `enlarge_threshold` - (Optional, Int) Specifies the average CPU usage (%). It is mandatory when `status` is set to
  **ON**. Value options: **50–100**.

* `max_flavor` - (Optional, String) Specifies the maximum specifications. It is mandatory when the instance specifications
  are automatically scaled up or down.

* `reduce_enabled` - (Optional, Bool) Specifies whether auto-down is enabled. It is mandatory when `status` is set to
  **ON**. Value options:
  + **true**: enabled.
  + **false**: disabled.

* `max_read_only_count` - (Optional, Int) Specifies the maximum number of read replicas. It is mandatory when read
  replicas are automatically added or deleted.

* `read_only_weight` - (Optional, Int) Specifies the read weights of read replicas. It is mandatory when read replicas
  are automatically added or deleted.

<a name="scaling_strategy_struct"></a>
The `scaling_strategy` block supports:

* `flavor_switch` - (Required, String) Specifies whether instance specifications can be automatically scaled up or down.
  Value options:
  + **ON**: Yes
  + **OFF**: No

* `read_only_switch` - (Required, String) Specifies whether read replicas can be automatically added or deleted. To use
  this function, ensure that there is only one proxy instance.
  Value options:
  + **ON**: Yes
  + **OFF**: No

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the DB instance ID.

* `status` - Indicates the DB instance status.

* `mode` - Indicates the instance mode.

* `db_user_name` - Indicates the default username.

* `private_dns_name` - Indicates the private domain name.

* `upgrade_flag` - Indicates whether the version can be upgraded.

* `current_version` - Indicates the current database version.

* `current_kernel_version` - Indicates the current database kernel version.

* `created_at` - Indicates the creation time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `updated_at` - Indicates the Update time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `nodes` - Indicates the instance nodes information.
  The [nodes](#nodes_struct) structure is documented below.

* `auto_scaling` - Indicates the auto-scaling policies.
  The [auto_scaling](#auto_scaling_struct) structure is documented below.

<a name="nodes_struct"></a>
The `nodes` block contains:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `type` - Indicates the node type: master or slave.

* `status` - Indicates the node status.

* `private_read_ip` - Indicates the private IP address of a node.

* `availability_zone` - Indicates the availability zone where the node resides.

<a name="auto_scaling_struct"></a>
The `auto_scaling` block supports:

* `id` - Indicates the ID of an auto-scaling policy.

* `min_flavor` - Indicates the minimum specifications.

* `silence_start_at` - Indicates the start time of the silent period.

* `min_read_only_count` - Indicates the minimum number of read replicas.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 60 minutes.
* `delete` - Default is 30 minutes.

## Import

GaussDB instance can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_gaussdb_mysql_instance.test <id>
```

Note that the imported state may not be identical to your resource definition, due to the attribute missing from the
API response. The missing attribute is: `table_name_case_sensitivity`, `enterprise_project_id`, `password`, `ssl_option`,
`encryption_type`, `kms_key_id` and `parameters`. It is generally recommended running `terraform plan` after importing
a GaussDB MySQL instance. You can then decide if changes should be applied to the GaussDB MySQL instance, or the resource
definition should be updated to align with the GaussDB MySQL instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_gaussdb_mysql_instance" "test" {
  ...

  lifecycle {
    ignore_changes = [
      new_node_weight, proxy_mode, readonly_nodes_weight, parameters, ssl_option, encryption_type, kms_key_id
    ]
  }
}
```
