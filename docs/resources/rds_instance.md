---
subcategory: "Relational Database Service (RDS)"
---

# huaweicloud\_rds\_instance

Manage RDS instance resource
This is an alternative to `huaweicloud_rds_instance_v3`

## Example Usage

### create a single db instance

```hcl
resource "huaweicloud_networking_secgroup" "secgroup" {
  name        = "terraform_test_security_group"
  description = "terraform security group acceptance test"
}

resource "huaweicloud_rds_instance" "instance" {
  availability_zone = ["{{ availability_zone }}"]
  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "9.5"
    port     = "8635"
  }
  name              = "terraform_test_rds_instance"
  security_group_id = huaweicloud_networking_secgroup.secgroup.id
  subnet_id         = "{{ subnet_id }}"
  vpc_id            = "{{ vpc_id }}"
  volume {
    type = "ULTRAHIGH"
    size = 100
  }
  flavor = "rds.pg.c2.medium"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
}
```

### create a primary/standby db instance

```hcl
resource "huaweicloud_networking_secgroup" "secgroup" {
  name        = "terraform_test_security_group"
  description = "terraform security group acceptance test"
}

resource "huaweicloud_rds_instance" "instance" {
  availability_zone = ["{{ availability_zone_1 }}", "{{ availability_zone_2 }}"]
  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "9.5"
    port     = "8635"
  }
  name              = "terraform_test_rds_instance"
  security_group_id = huaweicloud_networking_secgroup.secgroup.id
  subnet_id         = "{{ subnet_id }}"
  vpc_id            = "{{ vpc_id }}"
  volume {
    type = "ULTRAHIGH"
    size = 100
  }
  flavor              = "rds.pg.s1.medium.ha"
  ha_replication_mode = "async"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
}
```

### create a single db instance with encrypted volume

```hcl
resource "huaweicloud_kms_key" "key" {
  key_alias       = "key_1"
  key_description = "first test key"
  is_enabled      = true
}

resource "huaweicloud_networking_secgroup" "secgroup" {
  name        = "terraform_test_security_group"
  description = "terraform security group acceptance test"
}

resource "huaweicloud_rds_instance" "instance" {
  availability_zone = ["{{ availability_zone }}"]
  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "9.5"
    port     = "8635"
  }
  name              = "terraform_test_rds_instance"
  security_group_id = huaweicloud_networking_secgroup.secgroup.id
  subnet_id         = "{{ subnet_id }}"
  vpc_id            = "{{ vpc_id }}"
  volume {
    disk_encryption_id = huaweicloud_kms_key.key.id
    type               = "ULTRAHIGH"
    size               = 100
  }
  flavor = "rds.pg.c2.medium"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the rds instance resource. If omitted, the provider-level region will be used. Changing this creates a new rds instance resource.

* `availability_zone` -
  (Required)
  Specifies the AZ name. Changing this parameter will create a new resource.

* `db` -
  (Required)
  Specifies the database information. Structure is documented below. Changing this parameter will create a new resource.

* `flavor` -
  (Required)
  Specifies the specification code.

* `name` -
  (Required)
  Specifies the DB instance name. The DB instance name of the same type
  must be unique for the same tenant. The value must be 4 to 64
  characters in length and start with a letter. It is case-sensitive
  and can contain only letters, digits, hyphens (-), and underscores
  (_).  Changing this parameter will create a new resource.

* `security_group_id` -
  (Required)
  Specifies the security group which the RDS DB instance belongs to.
  Changing this parameter will create a new resource.

* `vpc_id` -
  (Required)
  Specifies the VPC ID. Changing this parameter will create a new resource.

* `subnet_id` -
  (Required)
  Specifies the network id of a subnet. Changing this parameter will create a new resource.

* `volume` -
  (Required)
  Specifies the volume information. Structure is documented below.

* `backup_strategy` -
  (Optional)
  Specifies the advanced backup policy. Structure is documented below.

* `ha_replication_mode` -
  (Optional)
  Specifies the replication mode for the standby DB instance. For MySQL, the value
  is async or semisync. For PostgreSQL, the value is async or sync. For
  Microsoft SQL Server, the value is sync. NOTE: async indicates the
  asynchronous replication mode. semisync indicates the
  semi-synchronous replication mode. sync indicates the synchronous
  replication mode.  Changing this parameter will create a new resource.

* `param_group_id` -
  (Optional)
  Specifies the parameter group ID. Changing this parameter will create a new resource.

* `enterprise_project_id` - 
  (Optional) 
  The enterprise project id of the RDS instance. Changing this creates a new RDS instance.

* `tags` - (Optional) A mapping of tags to assign to the RDS instance.
  Each tag is represented by one key-value pair.

The `db` block supports:

* `password` -
  (Required)
  Specifies the database password. The value cannot be
  empty and should contain 8 to 32 characters, including uppercase
  and lowercase letters, digits, and the following special
  characters: ~!@#%^*-_=+? You are advised to enter a strong
  password to improve security, preventing security risks such as
  brute force cracking.  Changing this parameter will create a new resource.

* `port` -
  (Optional)
  Specifies the database port information. The MySQL database port
  ranges from 1024 to 65535 (excluding 12017 and 33071, which are
  occupied by the RDS system and cannot be used). The PostgreSQL
  database port ranges from 2100 to 9500. The Microsoft SQL Server
  database port can be 1433 or ranges from 2100 to 9500, excluding
  5355 and 5985. If this parameter is not set, the default value is
  as follows: For MySQL, the default value is 3306. For PostgreSQL,
  the default value is 5432. For Microsoft SQL Server, the default
  value is 1433.  Changing this parameter will create a new resource.

* `type` -
  (Required)
  Specifies the DB engine. Value: MySQL, PostgreSQL, SQLServer. Changing this parameter will create a new resource.

* `version` -
  (Required)
  Specifies the database version. Changing this parameter will create a new resource.
  Available value for attributes:

type | version
---- | ---
MySQL| 5.6 <br>5.7  <br>8.0 <br>*8.0 is available only for users with the required permission.* <br>*You can contact customer service to apply for the permission.*
PostgreSQL | 9.5 <br> 9.6 <br>10 <br>11
SQLServer| 2008_R2_EE <br>2008_R2_WEB <br>2012_SE <br>2014_SE <br>2016_SE <br>2017_SE <br>2012_EE <br>2014_EE <br>2016_EE <br>2017_EE <br>2012_WEB <br>2014_WEB <br>2016_WEB <br>2017_WEB

The `volume` block supports:

* `disk_encryption_id` -
  (Optional)
  Specifies the key ID for disk encryption. Changing this parameter will create a new resource.

* `size` -
  (Required)
  Specifies the volume size. Its value range is from 40 GB to 4000
  GB. The value must be a multiple of 10. Changing this resize the volume.

* `type` -
  (Required)
  Specifies the volume type. Its value can be any of the following
  and is case-sensitive: ULTRAHIGH: indicates the SSD type.
  ULTRAHIGHPRO: indicates the ultra-high I/O (advanced), which supports ultra-high performance (advanced) DB instances.
  Changing this parameter will create a new resource.

The `backup_strategy` block supports:

* `keep_days` -
  (Optional)
  Specifies the retention days for specific backup files. The value
  range is from 0 to 732. If this parameter is not specified or set
  to 0, the automated backup policy is disabled. NOTICE:
  Primary/standby DB instances of Microsoft SQL Server do not
  support disabling the automated backup policy.

* `start_time` -
  (Required)
  Specifies the backup time window. Automated backups will be
  triggered during the backup time window. It must be a valid value in the &quot;hh:mm-HH:MM&quot;
  format. The current time is in the UTC format. The HH value must
  be 1 greater than the hh value. The values of mm and MM must be
  the same and must be set to any of the following: 00, 15, 30, or 45.
  Example value: 08:15-09:15 23:00-00:00.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `status` - Indicates the DB instance status.

* `created` - Indicates the creation time.

* `nodes` - Indicates the instance nodes information. Structure is documented below.

* `private_ips` - Indicates the private IP address list.
  It is a blank string until an ECS is created.

* `public_ips` - Indicates the public IP address list.

The `nodes` block contains:

* `availability_zone` - Indicates the AZ.

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `role` - Indicates the node type. The value can be master or slave,
  indicating the primary node or standby node respectively.

* `status` - Indicates the node status.

## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 30 minute.
- `update` - Default is 30 minute.

## Import

RDS instance can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_rds_instance.instance_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```

But due to some attrubutes missing from the API response, it's required to ignore changes as below.

```
resource "huaweicloud_rds_instance" "instance_1" {
  ...

  lifecycle {
    ignore_changes = [
      "db",
    ]
  }
}
```
