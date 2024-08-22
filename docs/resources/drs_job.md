---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_job"
description: ""
---

# huaweicloud_drs_job

Manages DRS job resource within HuaweiCloud.

## Example Usage

### Create a DRS job to migrate data to the HuaweiCloud RDS database

```hcl
variable "name" {}
variable "source_db_ip" {}
variable "source_db_port" {}
variable "source_db_user" {}
variable "source_db_password" {}
variable "source_db_port" {}
variable "destination_db_password" {}

resource "huaweicloud_rds_instance" "mysql" {
  ...
}

resource "huaweicloud_drs_job" "test" {
  name           = var.name
  type           = "migration"
  engine_type    = "mysql"
  direction      = "up"
  net_type       = "eip"
  migration_type = "FULL_INCR_TRANS"
  description    = "terraform demo"

  source_db {
    engine_type = "mysql"
    ip          = var.source_db_ip
    port        = var.source_db_port
    user        = var.source_db_user
    password    = var.source_db_password
    ssl_link    = false
  }

  destination_db {
    region      = huaweicloud_rds_instance.mysql.region
    ip          = huaweicloud_rds_instance.mysql.fixed_ip
    port        = 3306
    engine_type = "mysql"
    user        = "root"
    password    = var.destination_db_password
    instance_id = huaweicloud_rds_instance.mysql.id
    subnet_id   = huaweicloud_rds_instance.mysql.subnet_id
  }

  lifecycle {
    ignore_changes = [
      source_db.0.password, destination_db.0.password,
    ]
  }
}
```

### Create a DRS job to synchronize database level data to the HuaweiCloud RDS database and net type is VPC

```hcl
variable "name" {}
variable "source_db_ip" {}
variable "source_db_port" {}
variable "source_db_user" {}
variable "source_db_password" {}
variable "source_db_port" {}
variable "destination_db_password" {}
variable "database_name" {}
variable "source_db_vpc_id" {}
variable "source_db_subnet_id" {}

resource "huaweicloud_rds_instance" "mysql" {
  ...
}

resource "huaweicloud_drs_job" "test" {
  name           = var.name
  type           = "sync"
  engine_type    = "mysql"
  direction      = "up"
  net_type       = "vpc"
  migration_type = "FULL_INCR_TRANS"
  description    = "terraform demo"

  source_db {
    engine_type = "mysql"
    ip          = var.source_db_ip
    port        = var.source_db_port
    user        = var.source_db_user
    password    = var.source_db_password
    ssl_link    = false
    vpc_id      = var.source_db_vpc_id
    subnet_id   = var.source_db_subnet_id
  }

  destination_db {
    region      = huaweicloud_rds_instance.mysql.region
    ip          = huaweicloud_rds_instance.mysql.fixed_ip
    port        = 3306
    engine_type = "mysql"
    user        = "root"
    password    = var.destination_db_password
    instance_id = huaweicloud_rds_instance.mysql.id
    subnet_id   = huaweicloud_rds_instance.mysql.subnet_id
  }

  databases = [var.database_name]

  policy_config {
    filter_ddl_policy = "drop_database"
    conflict_policy   = "overwrite"
    index_trans       = true
  }

  lifecycle {
    ignore_changes = [
      source_db.0.password, destination_db.0.password,
    ]
  }
}
```

### Create a DRS job to synchronize table level data to the HuaweiCloud RDS database

```hcl
variable "name" {}
variable "source_db_ip" {}
variable "source_db_port" {}
variable "source_db_user" {}
variable "source_db_password" {}
variable "source_db_port" {}
variable "destination_db_password" {}
variable "database_name" {}
variable "table_name" {}

resource "huaweicloud_rds_instance" "mysql" {
  ...
}

resource "huaweicloud_drs_job" "test" {
  name           = var.name
  type           = "sync"
  engine_type    = "mysql"
  direction      = "up"
  net_type       = "eip"
  migration_type = "FULL_INCR_TRANS"
  description    = "terraform demo"

  source_db {
    engine_type = "mysql"
    ip          = var.source_db_ip
    port        = var.source_db_port
    user        = var.source_db_user
    password    = var.source_db_password
    ssl_link    = false
  }

  destination_db {
    region      = huaweicloud_rds_instance.mysql.region
    ip          = huaweicloud_rds_instance.mysql.fixed_ip
    port        = 3306
    engine_type = "mysql"
    user        = "root"
    password    = var.destination_db_password
    instance_id = huaweicloud_rds_instance.mysql.id
    subnet_id   = huaweicloud_rds_instance.mysql.subnet_id
  }

  tables {
    database    = var.database_name
    table_names = [var.table_name]
  }

  lifecycle {
    ignore_changes = [
      source_db.0.password, destination_db.0.password,
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the job name. The name consists of 4 to 50 characters, starting with
 a letter. Only letters, digits, underscores (\_) and hyphens (-) are allowed.

* `type` - (Required, String, ForceNew) Specifies the job type. Changing this parameter will create a new
 resource. The options are as follows:
  + **migration**: Online Migration.
  + **sync**: Data Synchronization.
  + **cloudDataGuard**: Disaster Recovery.

* `engine_type` - (Required, String, ForceNew) Specifies the migration engine type.
 Changing this parameter will create a new resource. The options are as follows:
  + **mysql**:  MySQL migration, MySQL synchronization use.
  + **mongodb**: Mongodb migration use.
  + **cloudDataGuard-mysql**: Disaster recovery use.
  + **gaussdbv5**: GaussDB (for openGauss) synchronization use.
  + **mysql-to-kafka**: Synchronization from MySQL to Kafka use.
  + **taurus-to-kafka**: Synchronization from GaussDB(for MySQL) to Kafka use.
  + **gaussdbv5ha-to-kafka**: Synchronization from GaussDB primary/standby to Kafka use.
  + **postgresql**: Synchronization from PostgreSQL to PostgreSQL use.

  For details, see [engine types](https://support.huaweicloud.com/intl/en-us/api-drs/drs_api_0159.html)

* `direction` - (Required, String, ForceNew) Specifies the direction of data flow.
  Changing this parameter will create a new resource. The options are as follows:
  + **up**: To the cloud. The destination database must be a database in the current cloud.
  + **down**: Out of the cloud. The source database must be a database in the current cloud.
  + **non-dbs**: self-built database.
  
* `source_db` - (Required, List, ForceNew) Specifies the source database configuration.
  The [db_info](#block--db_info) structure of the `source_db` is documented below.
  Changing this parameter will create a new resource.

* `destination_db` - (Required, List, ForceNew) Specifies the destination database configuration.
  The [db_info](#block--db_info) structure of the `destination_db` is documented below.
  Changing this parameter will create a new resource.

* `net_type` - (Optional, String, ForceNew) Specifies the network type.
  Changing this parameter will create a new resource. The default value is **eip**. The options are as follows:
  + **eip**: suitable for migration from an on-premises or other cloud database to a destination cloud database.
   An EIP will be automatically bound to the replication instance and released after the replication task is complete.
  + **vpc**: suitable for migration from one cloud database to another.
  + **vpn**: suitable for migration from an on-premises self-built database to a destination cloud database,
   or from one cloud database to another in a different region.

* `migration_type` - (Optional, String, ForceNew) Specifies migration type.
  Changing this parameter will create a new resource. The default value is **FULL_INCR_TRANS**. The options are as follows:
  + **FULL_TRANS**: Full migration. Suitable for scenarios where services can be interrupted. It migrates all database
   objects and data, in a non-system database, to a destination database at a time.
  + **INCR_TRANS**: Incremental migration. Suitable for migration from an on-premises self-built database to a
   destination cloud database, or from one cloud database to another in a different region.
  + **FULL_INCR_TRANS**:  Full+Incremental migration. This allows to migrate data with minimal downtime. After a full
   migration initializes the destination database, an incremental migration parses logs to ensure data consistency
   between the source and destination databases.

* `migrate_definer` - (Optional, Bool, ForceNew) Specifies whether to migrate the definers of all source database
  objects to the `user` of `destination_db`. The default value is **true**.
  Changing this parameter will create a new resource.

* `limit_speed` - (Optional, List, ForceNew) Specifies the migration speed by setting a time period.
  The default is no speed limit. The maximum length is 3. The [limit_speed](#block--limit_speed) structure is documented
  below. Changing this parameter will create a new resource.

* `policy_config` - (Optional, List, ForceNew) Specifies the policy information used to configure migration and
  synchronization policies. The [policy_config](#block--policy_config) structure is documented below.
  Changing this parameter will create a new resource.

* `multi_write` - (Optional, Bool, ForceNew) Specifies whether to enable multi write. It is mandatory when `type`
  is **cloudDataGuard**. When the disaster recovery type is dual-active disaster recovery, set `multi_write` to **true**,
  otherwise to **false**. The default value is **false**. Changing this parameter will create a new resource.

* `expired_days` - (Optional, Int, ForceNew) Specifies how many days after the task is abnormal, it will automatically
  end. The value ranges from 14 to 100. the default value is **14**. Changing this parameter will create a new resource.

* `start_time` - (Optional, String, ForceNew) Specifies the time to start the job. The time format is a time stamp
  accurating to milliseconds, e.g. **1684466549755**, which indicates **2023-05-19 11:22:29.755**.
  Start immediately by default. Changing this parameter will create a new resource.

* `destination_db_readnoly` - (Optional, Bool, ForceNew) Specifies the destination DB instance as read-only helps
  ensure the migration is successful. Once the migration is complete, the DB instance automatically changes to
  Read/Write. Changing this parameter will create a new resource.

  -> This parameter is valid only when MySQL migration and DR and `direction` is set to **up**. The default value is **true**,
  you need to manually set this parameter to **false** in other application scenarios except MySQL migration and DR scenarios.

* `description` - (Optional, String) Specifies the description of the job, which contain a
  maximum of 256 characters, and certain special characters (including !<>&'"\\) are not allowed.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id.
  Changing this parameter will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the DRS job.

* `force_destroy` - (Optional, Bool) Specifies whether to forcibly destroy the job even if it is running.
  The default value is **false**.

* `action` - (Optional, String) Specifies the action of job. The options are as follows:
  + **stop**: Stop the job. Available when job status is **FULL_TRANSFER_STARTED**, **FULL_TRANSFER_COMPLETE** or
    **INCRE_TRANSFER_STARTED**.
  + **restart**: Continue the job. Available when job status is **PAUSING**.
  + **reset**: Retry the job. Available when job status is **FULL_TRANSFER_FAILED** or **INCRE_TRANSFER_FAILED**.

  -> It will only take effect when **updating** a job.

* `pause_mode` - (Optional, String) Specifies the stop type of job. It's valid when `action` is **stop**.
  Default value is **target**. The options are as follows:
  + **target**: Stop playback.
  + **all**: Stop log capture and playback.

* `is_sync_re_edit` - (Optional, Bool) Specifies whether to start the sync re-edit job. It's valid when `action` is **restart**.

* `databases` - (Optional, List)  Specifies the list of the databases which the job migrates or synchronizes. Means to
  transfer database level data. This parameter conflicts with `tables`.

* `tables` - (Optional, List)  Specifies the list of the tables which the job migrates or synchronizes. Means to transfer
  table level data. This parameter conflicts with `databases`.
  The [tables](#block--tables) structure is documented below.

  ->   1. `databases` and `tables` will only take effect when `type` is **migration** or **sync**.
  <br/>2. When `type` is **migration**, they are not allowed to **update**, if they are empty, means to migrate all objects.
  <br/>3. When `type` is **sync**, exactly one data level of `databases` and `tables` must be specified. It's **not allowed**
       to transfer the data level to another. Only when `status` is **INCRE_TRANSFER_STARTED** or **INCRE_TRANSFER_FAILED**,
       **update** will take effect.
  <br/>4. It's only for synchronization from **MySQL** to **MySQL**, migration from **Redis** to **GeminiDB Redis**,
       migration from cluster **Redis** to **GeminiDB Redis**, and synchronization from **Oracle** to **GaussDB Distributed**.

* `master_az` - (Optional, String, ForceNew) Specifies the AZ where the primary task is located.

* `slave_az` - (Optional, String, ForceNew) Specifies the AZ where the standby task is located.
  
  -> It takes effect when both `master_az` and `slave_az` are specified. Only MySQL and gaussdbv5ha-to-kafka scenarios
  are supported.

* `alarm_notify` - (Optional, List, ForceNew)  Specifies the information body for setting task exception notification.
  The [alarm_notify](#block--alarm_notify) structure is documented below.

<a name="block--db_info"></a>
The `db_info` block supports:

* `engine_type` - (Required, String, ForceNew) Specifies the engine type of database. Changing this parameter will
  create a new resource. The options are as follows: **mysql**, **mongodb**, **gaussdbv5**, **taurus**, **gaussdbv5ha**,
  **kafka**, **postgresql**.

* `ip` - (Required, String, ForceNew) Specifies the IP of database. Changing this parameter will create a new resource.

* `port` - (Required, Int, ForceNew) Specifies the port of database. Changing this parameter will create a new resource.

* `user` - (Required, String, ForceNew) Specifies the user name of database.
  Changing this parameter will create a new resource.

* `password` - (Required, String, ForceNew) Specifies the password of database.
  Changing this parameter will create a new resource.

* `instance_id` - (Optional, String, ForceNew) Specifies the instance id of database when it is a RDS database.
  Changing this parameter will create a new resource.

* `vpc_id` - (Optional, String, ForceNew) Specifies vpc ID of database.
  Changing this parameter will create a new resource.

* `subnet_id` - (Optional, String, ForceNew) Specifies subnet ID of database when it is a RDS database.
  It is mandatory when `direction` is **down**. Changing this parameter will create a new resource.

  -> When `net_type` is **vpc**, if `direction` is **up**, `source_db.vpc_id` and `source_db.subnet_id` is mandatory, if
  `direction` is **down**, `destination_db.vpc_id` and `destination_db.subnet_id` is mandatory.

* `region` - (Optional, String, ForceNew) Specifies the region which the database belongs when it is a RDS database.
  Changing this parameter will create a new resource.

* `name` - (Optional, String, ForceNew) Specifies the name of database.
  Changing this parameter will create a new resource.

* `ssl_enabled` - (Optional, Bool, ForceNew) Specifies whether to enable SSL connection.
  Changing this parameter will create a new resource.

* `ssl_cert_key` - (Optional, String, ForceNew) Specifies the SSL certificate content, encrypted with base64.
  It is mandatory when `ssl_enabled` is **true**. Changing this parameter will create a new resource.

* `ssl_cert_name` - (Optional, String, ForceNew) Specifies SSL certificate name.
  It is mandatory when `ssl_enabled` is **true**. Changing this parameter will create a new resource.

* `ssl_cert_check_sum` - (Optional, String, ForceNew) Specifies the checksum of SSL certificate content.
  It is mandatory when `ssl_enabled` is **true**. Changing this parameter will create a new resource.

* `ssl_cert_password` - (Optional, String, ForceNew) Specifies SSL certificate password. It is mandatory when
  `ssl_enabled` is **true** and the certificate file suffix is **.p12**. Changing this parameter will create a new resource.

<a name="block--limit_speed"></a>
The `limit_speed` block supports:

* `speed` - (Required, String, ForceNew) Specifies the transmission speed, the value range is 1 to 9999, unit: **MB/s**.
  Changing this parameter will create a new resource.

* `start_time` - (Required, String, ForceNew) Specifies the time to start speed limit, this time is UTC time. The start
  time is the whole hour, if there is a minute, it will be ignored, the format is **hh:mm**, and the hour number
  is two digits, for example: 01:00. Changing this parameter will create a new resource.

* `end_time` - (Required, String, ForceNew) Specifies the time to end speed limit, this time is UTC time. The input must
  end at 59 minutes, the format is **hh:mm**, for example: 15:59. Changing this parameter will create a new resource.

<a name="block--policy_config"></a>
The `policy_config` block supports:

* `filter_ddl_policy` - (Optional, String, ForceNew) Specifies the DDL filtering policy. Valid value is **drop_database**.
  For MySQL synchronization, this parameter can only be set to **drop_database**.
  Changing this parameter will create a new resource.

* `conflict_policy` - (Optional, String, ForceNew) Specifies the incremental conflict policy.

  Valid values are as fallows:
  + **ignore**: Ignore the conflict. The system will ignore the conflicting data and continue the subsequent
    synchronization process.
  + **stop**: Report an error. The synchronization task will be stopped and fail. You can view the details in
    synchronization logs.
  + **overwrite**: Overwrite the existing data with the synchronized data. Conflicting data will be overwritten.

  Changing this parameter will create a new resource.

* `index_trans` - (Optional, Bool, ForceNew) Specifies the object synchronization scope, indicating whether to
  synchronize normal indexes. If it's **true**, all indexes will be synchronized, otherwise, only primary key or unique
  indexes are synchronized. Changing this parameter will create a new resource.

<a name="block--tables"></a>
The `tables` block supports:

* `database` - (Required, String) Specifies the name of database to which the tables belong.

* `table_names` - (Required, List) Specifies the names of table which belong to a same datebase.

<a name="block--alarm_notify"></a>
The `alarm_notify` block supports:

* `topic_urn` - (Required, String, ForceNew) Specifies the SMN topic URN which is subscribed.

* `delay_time` - (Optional, Int, ForceNew) Specifies the Delay threshold between the source and destination database,
  in seconds. Value ranges from **1** to **3600**. Default is **0** and no notifications will be sent to recipient. If
  the delay exceeds a specified value and lasts for 6 minutes, DRS will notify specified recipients. This option is
  available only for **full+incremental** tasks.

* `rpo_delay` - (Optional, Int, ForceNew) Specifies the RPO delay threshold, in seconds. Value ranges from **1** to **3600**.
  Default is **0** and no notifications will be sent to recipient. If the RPO delay between the service database and
  the DRS instance exceeds a specified value and lasts for 6 minutes, DRS will notify specified recipients.

* `rto_delay` - (Optional, Int, ForceNew) Specifies the RTO delay threshold, in seconds. Value ranges from **1** to **3600**.
  Default is **0** and no notifications will be sent to recipient. If the RTO delay between the DRS instance and the
  DR database exceeds a specified value and lasts for 6 minutes, DRS will notify specified recipients.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `order_id` - The order ID which will return if `charging_mode` is **prePaid**.

* `master_job_id` - The master job ID which will return if job is dual-AZ.

* `slave_job_id` - The slave job ID which will return if job is dual-AZ.

* `created_at` - Create time. The format is ISO8601:YYYY-MM-DDThh:mm:ssZ.

* `updated_at` - Update time. The format is ISO8601:YYYY-MM-DDThh:mm:ssZ.

* `vpc_id` - The VPC ID to which the DRS instance belongs.

* `subnet_id` - The subnet ID to which the DRS instance belongs.

* `security_group_id` - The security group ID to which the DRS instance belongs.

* `source_db` - The source database configuration.
  The [db_info](#attrblock--db_info) structure of the `source_db` is documented below.

* `destination_db` - The destination database configuration.
  The [db_info](#attrblock--db_info) structure of the `destination_db` is documented below.

* `status` - Status.

* `progress` - Progress.

* `public_ip` - Public IP.

* `private_ip` - Private IP.

<a name="attrblock--db_info"></a>
The `db_info` block supports:

* `security_group_id` - The security group ID to which the databese instance belongs.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

* `update` - Default is 10 minutes.

* `delete` - Default is 10 minutes.

## Import

The DRS job can be imported by `id`. e.g.

```bash
$ terraform import huaweicloud_drs_job.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `enterprise_project_id`, `force_destroy`,
`source_db.0.password` and `destination_db.0.password`, `action`, `is_sync_re_edit`, `pause_mode`, `auto_renew`,
`alarm_notify.0.topic_urn`. It is generally recommended running **terraform plan** after importing a job. You can then
decide if changes should be applied to the job, or the resource definition should be updated to align with the job. Also
you can ignore changes as below.

```hcl
resource "huaweicloud_drs_job" "test" {
    ...

  lifecycle {
    ignore_changes = [
      source_db.0.password, destination_db.0.password, action,
    ]
  }
}
```
