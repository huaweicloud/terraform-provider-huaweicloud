---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_read_replica_instance"
description: ""
---

# huaweicloud_rds_read_replica_instance

Manage RDS Read Replica Instance resource.

## Example Usage

### Create a Rds read replica instance

```hcl
variable "primary_instance_id" {}
variable "security_group_id" {}
variable "availability_zone" {}

resource "huaweicloud_rds_read_replica_instance" "replica_instance" {
  name                = "test_rds_readonly_instance"
  flavor              = "rds.mysql.x1.large.2.rr"
  primary_instance_id = var.primary_instance_id
  availability_zone   = [var.availability_zone]
  security_group_id   = var.security_group_id

  db {
    port = "8888"
  }
  
  volume {
    type              = "CLOUDSSD"
    size              = 50
    limit_size        = 200
    trigger_threshold = 10
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the rds read replica instance resource. If
  omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.
  Currently, read replicas can be created **only** in the same region as that of the primary DB instance.

* `availability_zone` - (Required, String, ForceNew) Specifies the AZ name. Changing this parameter will create a new
  resource.

* `name` - (Required, String, ForceNew) Specifies the DB instance name. The DB instance name of the same type must be
  unique for the same tenant. The value must be 4 to 64 characters in length and start with a letter. It is
  case-sensitive and can contain only letters, digits, hyphens (-), and underscores (_). Changing this parameter will
  create a new resource.

* `primary_instance_id` - (Required, String, ForceNew) Specifies the DB instance ID, which is used to create a read
  replica. Changing this parameter will create a new resource.

* `volume` - (Required, List, ForceNew) Specifies the volume information. The [volume](#Rds_volume) structure is
  documented below. Changing this parameter will create a new resource.

* `db` - (Optional, List, ForceNew) Specifies the database information. The [db](#Rds_db) structure is documented below.
  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id of the read replica instance.

* `flavor` - (Required, String) Specifies the specification code.

* `security_group_id` - (Optional, String) Specifies the security group which the read replica instance belongs to.

* `fixed_ip` - (Optional, String) Specifies an intranet floating IP address of read replica instance.

* `ssl_enable` - (Optional, Bool) Specifies whether to enable the SSL for read replica instance.

* `parameters` - (Optional, List) Specify an array of one or more parameters to be set to the read replica instance
  after launched. You can check on console to see which parameters supported. The [parameters](#Rds_parameters)
  structure is documented below.

* `description` - (Optional, String) Specifies the description of the instance. The value consists of 0 to 64
  characters, including letters, digits, periods (.), underscores (_), and hyphens (-).

* `maintain_begin` - (Optional, String) Specifies the time at which the maintenance time window starts, for example, **22:00**.

* `maintain_end` - (Optional, String) Specifies the time at which the maintenance time window ends, for example, **01:00**.

-> **Note** For RDS for MySQL and RDS for PostgreSQL databases, the maintenance begin time and end time must be on the
  hour, and the interval between them must be one to four hours.<br>
  For RDS for SQL Server databases, the interval between the maintenance begin time and end time must be four hours.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the read replica instance. Valid values
  are **prePaid** and **postPaid**, defaults to **postPaid**. Changing this creates a new resource.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the read replica instance. Valid
  values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**. Changing this
  creates a new resource.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the read replica instance. If `period_unit` is
  set to **month**, the value ranges from 1 to 9. If `period_unit` is set to **year**, the value ranges from 1 to 3.
  This parameter is mandatory if `charging_mode` is set to **prePaid**. Changing this creates a new resource.

* `auto_renew` - (Optional, String) Specifies whether auto-renew is enabled. Valid values are **true** and **false**.

* `tags` - (Optional, Map) A mapping of tags to assign to the RDS read replica instance. Each tag is represented by one
  key-value pair.

<a name="Rds_db"></a>
The `db` block supports:

* `port` - (Optional, Int) Specifies the database port.
  + The MySQL database port ranges from `1,024` to `65,535` (excluding `12,017` and `33,071`, which are occupied by
    the RDS system and cannot be used). The default value is `3,306`.
  + The PostgreSQL database port ranges from `2,100` to `9,500`. The default value is `5,432`.
  + The Microsoft SQL Server database port can be `1,433` or ranges from `2,100` to `9,500`, excluding `5,355` and
    `5,985`. The default value is `1,433`.

<a name="Rds_volume"></a>
The `volume` block supports:

* `type` - (Required, String, ForceNew) Specifies the volume type. It must same with the type of the primary instance.
  Its value can be any of the following and is case-sensitive:
  + **ULTRAHIGH**: SSD storage.
  + **LOCALSSD**: local SSD storage.
  + **CLOUDSSD**: cloud SSD storage. This storage type is supported only with general-purpose and dedicated DB
    instances.
  + **ESSD**: extreme SSD storage.

  Changing this parameter will create a new resource.

* `size` - (Optional, Int) Specifies the volume size. Its value range is from `40` GB to `4,000` GB. The value must
  be a multiple of 10 and greater than the original size.

* `limit_size` - (Optional, Int) Specifies the upper limit of automatic expansion of storage, in GB.

* `trigger_threshold` - (Optional, Int) Specifies the threshold to trigger automatic expansion.  
  If the available storage drops to this threshold or `10` GB, the automatic expansion is triggered.  
  The valid values are as follows:
  + **10**
  + **15**
  + **20**

<a name="Rds_parameters"></a>
The `parameters` block supports:

* `name` - (Required, String) Specifies the parameter name. Some of them needs the instance to be restarted
  to take effect.

* `value` - (Required, String) Specifies the parameter value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the instance ID.

* `status` - Indicates the instance status.

* `type` - Indicates the type of the read replica instance. The value can be **Single**, **Ha**, **Replica**,
  **Enterprise**.

* `db/type` - Indicates the DB engine. The value can be **MySQL**, **PostgreSQL**, **SQLServer**, **MariaDB**.

* `db/version` - Indicates the database version.

* `db/user_name` - Indicates the default username of database.

* `volume/disk_encryption_id` - Indicates the key ID for disk encryption. It is same with the primary instance.

* `private_ips` - Indicates the private IP address list.

* `public_ips` - Indicates the public IP address list.

* `subnet_id` - Indicates the subnet id.

* `vpc_id` - Indicates the VPC ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

RDS read replica instance can be imported by `id`, e.g.

```bash
$ terraform import huaweicloud_rds_read_replica_instance.replica_instance <id>
```
