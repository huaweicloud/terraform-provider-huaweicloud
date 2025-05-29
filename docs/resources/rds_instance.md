---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_instance"
description: ""
---

# huaweicloud_rds_instance

Manage RDS instance resource within HuaweiCloud.

## Example Usage

### create a single db instance

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "secgroup_id" {}
variable "availability_zone" {}
variable "postgreSQL_password" {}

resource "huaweicloud_rds_instance" "instance" {
  name              = "terraform_test_rds_instance"
  flavor            = "rds.pg.n1.large.2"
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
  availability_zone = [var.availability_zone]

  db {
    type     = "PostgreSQL"
    version  = "12"
    password = var.postgreSQL_password
  }

  volume {
    type = "ULTRAHIGH"
    size = 100
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
}
```

### create a primary/standby db instance

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "secgroup_id" {}
variable "availability_zone1" {}
variable "availability_zone2" {}
variable "postgreSQL_password" {}

resource "huaweicloud_rds_instance" "instance" {
  name                = "terraform_test_rds_instance"
  flavor              = "rds.pg.n1.large.2.ha"
  ha_replication_mode = "async"
  vpc_id              = var.vpc_id
  subnet_id           = var.subnet_id
  security_group_id   = var.secgroup_id
  availability_zone   = [
    var.availability_zone1,
    var.availability_zone2,
  ]

  db {
    type     = "PostgreSQL"
    version  = "12"
    password = var.postgreSQL_password
  }
  volume {
    type = "ULTRAHIGH"
    size = 100
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
}
```

### create a single db instance with encrypted volume

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "secgroup_id" {}
variable "availability_zone" {}
variable "kms_id" {}
variable "postgreSQL_password" {}

resource "huaweicloud_rds_instance" "instance" {
  name              = "terraform_test_rds_instance"
  flavor            = "rds.pg.n1.large.2"
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
  availability_zone = [var.availability_zone]

  db {
    type     = "PostgreSQL"
    version  = "12"
    password = var.postgreSQL_password
  }
  volume {
    type               = "ULTRAHIGH"
    size               = 100
    disk_encryption_id = var.kms_id
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
}
```

### create db instance with customized parameters

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "secgroup_id" {}
variable "availability_zone" {}
variable "postgreSQL_password" {}

resource "huaweicloud_rds_instance" "instance" {
  name              = "terraform_test_rds_instance"
  flavor            = "rds.pg.n1.large.2"
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
  availability_zone = [var.availability_zone]

  db {
    type     = "PostgreSQL"
    version  = "12"
    password = var.postgreSQL_password
  }

  volume {
    type = "ULTRAHIGH"
    size = 100
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  parameters {
    name  = "div_precision_increment"
    value = "12"
  }

  parameters {
    name  = "connect_timeout"
    value = "13"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the rds instance resource. If omitted, the
  provider-level region will be used. Changing this creates a new rds instance resource.

* `availability_zone` - (Required, List) Specifies the list of AZ name.

* `name` - (Required, String) Specifies the DB instance name. The DB instance name of the same type must be unique for
  the same tenant. The value must be 4 to 64 characters in length and start with a letter. It is case-sensitive and can
  contain only letters, digits, hyphens (-), and underscores (_).

* `flavor` - (Required, String) Specifies the specification code.

  -> **NOTE:** Services will be interrupted for 5 to 10 minutes when you change RDS instance flavor.If this parameter is
  changed, a temporary instance will be generated. This temporary instance will occupy the association of the VPC
  security group and cannot be deleted for 12 hours.

* `db` - (Required, List, ForceNew) Specifies the database information. Structure is documented below. Changing this
  parameter will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID. Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the network id of a subnet. Changing this parameter will create a
  new resource.

* `security_group_id` - (Required, String) Specifies the security group which the RDS DB instance belongs to.

* `volume` - (Required, List) Specifies the volume information. Structure is documented below.

* `restore` - (Optional, List, ForceNew) Specifies the restoration information. It only supported restore to postpaid
  instance. Structure is documented below. Changing this parameter will create a new resource.

* `fixed_ip` - (Optional, String) Specifies an intranet floating IP address of RDS DB instance.

* `backup_strategy` - (Optional, List) Specifies the advanced backup policy. Structure is documented below.

* `ha_replication_mode` - (Optional, String) Specifies the replication mode for the standby DB instance.
  + For MySQL, the value is **async** or **semisync**.
  + For PostgreSQL, the value is **async** or **sync**.
  + For Microsoft SQL Server, the value is **sync**.
  + For MariaDB, the value is **async** or **semisync**.

  -> **NOTE:** **async** indicates the asynchronous replication mode. **semisync** indicates the semi-synchronous
  replication mode. **sync** indicates the synchronous replication mode.

* `lower_case_table_names` - (Optional, String, ForceNew) Specifies the case-sensitive state of the database table name,
  the default value is "1". Changing this parameter will create a new resource.
    + 0: Table names are stored as fixed and table names are case-sensitive.
    + 1: Table names will be stored in lower case and table names are not case-sensitive.

* `param_group_id` - (Optional, String) Specifies the parameter group ID.

* `collation` - (Optional, String) Specifies the Character Set, only available to Microsoft SQL Server DB instances.

* `time_zone` - (Optional, String, ForceNew) Specifies the UTC time zone. For MySQL and PostgreSQL Chinese mainland site
  and international site use UTC by default. The value ranges from UTC-12:00 to UTC+12:00 at the full hour. For
  Microsoft SQL Server international site use UTC by default and Chinese mainland site use China Standard Time. The time
  zone is expressed as a character string, refer to
  [HuaweiCloud Document](https://support.huaweicloud.com/intl/en-us/api-rds/rds_01_0002.html#rds_01_0002__table613473883617).

* `switch_strategy` - (Optional, String) Specifies the database switchover policy.
  + **reliability**: reliability first.
  + **availability**: availability first.
  
  Defaults to **reliability**.

* `charging_mode` - (Optional, String) Specifies the charging mode of the RDS DB instance. Valid values are **prePaid**
  and **postPaid**, defaults to **postPaid**.

* `period_unit` - (Optional, String) Specifies the charging period unit of the RDS DB instance. Valid values are **month**
  and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.

* `period` - (Optional, Int) Specifies the charging period of the RDS DB instance. If `period_unit` is set to **month**,
  the value ranges from `1` to `9`. If `period_unit` is set to **year**, the value ranges from `1` to `3`. This parameter
  is mandatory if `charging_mode` is set to **prePaid**.

* `auto_renew` - (Optional, String) Specifies whether auto-renew is enabled. Valid values are "true" and "false".

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id of the RDS instance.

* `ssl_enable` - (Optional, Bool) Specifies whether to enable the SSL for MySQL database.

* `description` - (Optional, String) Specifies the description of the instance. The value consists of 0 to 64
  characters, including letters, digits, periods (.), underscores (_), and hyphens (-).

* `dss_pool_id` - (Optional, String) Specifies the exclusive storage ID for Dec users. It is different for each az
  configuration. When creating an instance for Dec users, it is needed to be specified for all nodes of the instance
  and separated by commas if database instance type is not standalone or read-only.

* `maintain_begin` - (Optional, String) Specifies the time at which the maintenance time window starts, for example, **22:00**.

* `maintain_end` - (Optional, String) Specifies the time at which the maintenance time window ends, for example, **01:00**.

-> **Note** For RDS for MySQL and RDS for PostgreSQL databases, the maintenance begin time and end time must be on the
  hour, and the interval between them must be one to four hours.<br>
  For RDS for SQL Server databases, the interval between the maintenance begin time and end time must be four hours.

* `tags` - (Optional, Map) A mapping of tags to assign to the RDS instance. Each tag is represented by one key-value
  pair.

* `parameters` - (Optional, List) Specify an array of one or more parameters to be set to the RDS instance after
  launched. You can check on console to see which parameters supported. Structure is documented below.

* `binlog_retention_hours` - (Optional, Int) Specify the binlog retention period in hours. This parameter applies only to
  MySQL Server databases. Value range: `0` to `168` (7x24).

* `msdtc_hosts` - (Optional, List) Specify the host information for MSDTC.
  The [msdtc_hosts](#RdsInstance_MsdtcHosts) structure is documented below.

  -> **NOTE:** Only adding MSDTC hosts is supported, deletion is not allowed.

* `power_action` - (Optional, String) Specifies the power action to be done for the instance.
  Value options: **ON**, **OFF** and **REBOOT**.

  -> **NOTE:** The `power_action` is a one-time action.

* `tde_enabled` - (Optional, Bool) Specifies whether enable TDE for the instance.

  -> **NOTE:** TDE cannot be disabled after being enabled.

* `rotate_day` - (Optional, Int) Specifies the rotation days of TDE rotation.

* `secret_id` - (Optional, String) Specifies the key ID of TDE rotation.

* `secret_name` - (Optional, String) Specifies the key name of TDE rotation.

* `secret_version` - (Optional, String) Specifies the key version of TDE rotation.

  -> **NOTE:** `rotate_day`, `secret_id`, `secret_name` and `secret_version` will only take effect when `tde_enabled`
  is **true**.

* `read_write_permissions` - (Optional, String) Specifies the read write permissions of the instance. Valid values:
  + **readwrite**: read write permissions.
  + **readonly**: readonly permissions.

* `seconds_level_monitoring_enabled` - (Optional, Bool) Specifies whether to enable seconds level monitoring.

* `seconds_level_monitoring_interval` - (Optional, Int) Specifies the seconds level monitoring interval. Valid values:
  `1`, `5`. It is mandatory when `seconds_level_monitoring_enabled` is **true**.

* `minor_version_auto_upgrade_enabled` - (Optional, Bool) Specifies whether to enable minor version auto upgrade.

* `private_dns_name_prefix` - (Optional, String) Specifies the prefix of the private domain name. The value contains
  `8` to `64` characters. Only uppercase letters, lowercase letters, and digits are allowed.

* `slow_log_show_original_status` - (Optional, String) Specifies the slow log show original status of the instance.
  Only **MySQL** and **PostgreSQL** are supported. Value options: **on**, **off**.

The `db` block supports:

* `type` - (Required, String, ForceNew) Specifies the DB engine. Available value are **MySQL**, **PostgreSQL**,
  **SQLServer** and **MariaDB**. Changing this parameter will create a new resource.

* `version` - (Required, String, ForceNew) Specifies the database version. Changing this parameter will create a new
  resource. Available values detailed in
  [DB Engines and Versions](https://support.huaweicloud.com/intl/en-us/productdesc-rds/en-us_topic_0043898356.html).

* `password` - (Optional, String) Specifies the database password. The value should contain 8 to 32 characters,
  including uppercase and lowercase letters, digits, and the following special characters: ~!@#%^*-_=+? You are advised
  to enter a strong password to improve security, preventing security risks such as brute force cracking.

* `port` - (Optional, Int) Specifies the database port.
  + The MySQL database port ranges from 1024 to 65535 (excluding 12017 and 33071, which are occupied by the RDS system
      and cannot be used). The default value is 3306.
  + The PostgreSQL database port ranges from 2100 to 9500. The default value is 5432.
  + The Microsoft SQL Server database port can be 1433 or ranges from 2100 to 9500, excluding 5355 and 5985. The
      default value is 1433.
  + The MariaDB database port ranges from 1024 to 65535 (excluding 12017 and 33071, which are occupied by the RDS system
      and cannot be used). The default value is 3306.

The `volume` block supports:

* `size` - (Required, Int) Specifies the volume size. Its value range is from 40 GB to 4000 GB. The value must be a
  multiple of 10 and greater than the original size.

* `type` - (Required, String, ForceNew) Specifies the volume type. Its value can be any of the following and is
  case-sensitive:
  + **ULTRAHIGH**: SSD storage.
  + **LOCALSSD**: local SSD storage.
  + **CLOUDSSD**: cloud SSD storage. This storage type is supported only with general-purpose and dedicated DB
    instances.
  + **ESSD**: extreme SSD storage.

  Changing this parameter will create a new resource. For details about volume types, see
  [DB Instance Storage Types](https://support.huaweicloud.com/intl/en-us/productdesc-rds/rds_01_0020.html).

* `disk_encryption_id` - (Optional, String, ForceNew) Specifies the key ID for disk encryption.
  Changing this parameter will create a new resource.

* `limit_size` - (Optional, Int) Specifies the upper limit of automatic expansion of storage, in GB.

* `trigger_threshold` - (Optional, Int) Specifies the threshold to trigger automatic expansion.  
  If the available storage drops to this threshold or `10` GB, the automatic expansion is triggered.  
  The valid values are as follows:
  + **10**
  + **15**
  + **20**

The `restore` block supports:

* `instance_id` - (Required, String, ForceNew) Specifies the source DB instance ID. Changing this parameter will create
  a new resource.

* `backup_id` - (Required, String, ForceNew) Specifies the ID of the backup used to restore data. Changing this
  parameter will create a new resource.

* `database_name` - (Optional, Map, ForceNew) Specifies the database to be restored. This parameter applies only to
  Microsoft SQL Server databases. Changing this parameter will create a new resource.

The `backup_strategy` block supports:

* `keep_days` - (Required, Int) Specifies the retention days for specific backup files. The value range is from 0 to 732.

  -> **NOTE:** Primary/standby DB instances of Microsoft SQL Server do not support disabling the automated backup
  policy.

* `start_time` - (Required, String) Specifies the backup time window. Automated backups will be triggered during the
  backup time window. It must be a valid value in the **hh:mm-HH:MM**
  format. The current time is in the UTC format. The HH value must be 1 greater than the hh value. The values of mm and
  MM must be the same and must be set to any of the following: 00, 15, 30, or 45. Example value: 08:15-09:15 23:00-00:
  00.

* `period` - (Optional, String) Specifies the backup cycle. Automatic backups will be performed on the specified days of
  the week, except when disabling the automatic backup policy. The value range is a comma-separated number, where each
  number represents a day of the week. For example, a value of 1,2,3,4 would set the backup cycle to Monday, Tuesday,
  Wednesday, and Thursday. The default value is 1,2,3,4,5,6,7.

The `parameters` block supports:

* `name` - (Required, String) Specifies the parameter name. Some of them needs the instance to be restarted
  to take effect.

* `value` - (Required, String) Specifies the parameter value.

<a name="RdsInstance_MsdtcHosts"></a>
The `msdtc_hosts` block supports:

* `ip` - (Required, String) Specifies the host IP address.

* `host_name` - (Required, String) Specifies the host name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the DB instance ID.

* `status` - Indicates the DB instance status.

* `db/user_name` - Indicates the default username of database.

* `created` - Indicates the creation time.

* `nodes` - Indicates the instance nodes information. Structure is documented below.

* `private_ips` - Indicates the private IP address list. It is a blank string until an ECS is created.

* `private_dns_names` - Indicates the private domain name list of the DB instance.

* `public_ips` - Indicates the public IP address list.

* `msdtc_hosts` - Indicates the host information for MSDTC.
  The [msdtc_hosts](#RdsInstance_MsdtcHostsResp) structure is documented below.

The `nodes` block contains:

* `availability_zone` - Indicates the AZ.

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `role` - Indicates the node type. The value can be master or slave, indicating the primary node or standby node
  respectively.

* `status` - Indicates the node status.

<a name="RdsInstance_MsdtcHostsResp"></a>
The `msdtc_hosts` block supports:

* `id` - Indicates the host ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

RDS instance can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rds_instance.instance_1 <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `db`, `restore`,`param_group_id`,
`power_action`, `availability_zone`, `read_write_permissions`, `rotate_day`, `secret_id`, `secret_name`, `secret_version`,
`dss_pool_id`, `lower_case_table_names`, `slow_log_show_original_status`, `charging_mode`, `period_unit`, `period`,
`auto_renew`, `auto_pay`. It is generally recommended running `terraform plan` after importing a RDS instance. You can
then decide if changes should be applied to the instance, or the resource definition should be updated to align with the
instance. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_rds_instance" "instance_1" {
  ...

  lifecycle {
    ignore_changes = [
      "db", "restore", "param_group_id", "power_action", "availability_zone", "read_write_permissions", "rotate_day",
      "secret_id", "secret_name", "secret_version", "dss_pool_id", "lower_case_table_names", "slow_log_show_original_status",
      "charging_mode", "period_unit", "period", "auto_renew", "auto_pay",
    ]
  }
}
```
