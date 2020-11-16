---
subcategory: "Relational Database Service (RDS)"
---

# huaweicloud\_rds\_read\_replica\_instance

Manage RDS Read Replica Instance resource.

## Example Usage

### Create a Rds read replica instance
```hcl
resource "huaweicloud_networking_secgroup" "secgroup" {
  name          = "terraform_test_sg_for_rds"
  description   = "security group for rds read replica instance test"
}

resource "huaweicloud_rds_instance" "instance" {
  name                  = "terraform_test_rds_instance"
  flavor                = "rds.pg.c2.medium"
  availability_zone     = ["{{ availability_zone }}"]
  security_group_id     = huaweicloud_networking_secgroup.secgroup.id
  vpc_id                = "{{ vpc_id }}"
  subnet_id             = "{{ subnet_id }}"
  enterprise_project_id = "{{ enterprise_project_id }}"

  db {
    password    = "Huangwei!120521"
    type        = "PostgreSQL"
    version     = "10"
    port        = "8635"
  }
  volume {
    type = "ULTRAHIGH"
    size = 50
  }
  backup_strategy {
    start_time  = "08:00-09:00"
    keep_days   = 1
  }
}

resource "huaweicloud_rds_read_replica_instance" "replica_instance" {
  name                  = "terraform_test_rds_read_replica_instance"
  flavor                = "rds.pg.c2.medium.rr"
  primary_instance_id   = huaweicloud_rds_instance.instance.id
  availability_zone     = "{{ availability_zone }}"
  enterprise_project_id = "{{ enterprise_project_id }}"
  volume {
    type = "ULTRAHIGH"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, ForceNew) The region in which to create the rds read replica instance resource. If omitted, the provider-level region will be used.

  Currently, read replicas can be created only in the same region as that of the primary DB instance.

* `name` - (Required, ForceNew) Specifies the DB instance name. The DB instance name of the same type must be unique for the same tenant. 
  
  The value must be 4 to 64 characters in length and start with a letter. It is case-sensitive and can contain only letters, digits, hyphens (-), and underscores (_).

* `flavor` - (Required) Specifies the specification code.

* `primary_instance_id` - (Required, ForceNew) Specifies the DB instance ID, which is used to create a read replica.

* `volume` - (Required, Type:List, ForceNew) Specifies the volume information. Structure is documented below.

* `availability_zone` - (Required, ForceNew) Specifies the AZ name.

* `enterprise_project_id` - (Optional, ForceNew) The enterprise project id of the read replica instance.

* `tags` - (Optional, Type:Map) A mapping of tags to assign to the RDS read replica instance. Each tag is represented by one key-value pair.

The `volume` block supports:

* `type` - (Required, ForceNew) Specifies the volume type. Its value can be any of the following and is case-sensitive: 
    - ULTRAHIGH: indicates the SSD type.
    - ULTRAHIGHPRO: indicates the ultra-high I/O.

* `disk_encryption_id` -  (Optional, ForceNew) Specifies the key ID for disk encryption.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - Indicates the instance ID.

* `status` - Indicates the instance status.

* `db` - Indicates the database information. Structure is documented below. 

* `private_ips` - Indicates the private IP address list. 

* `public_ips` - Indicates the public IP address list.

* `security_group_id` - Indicates the security group which the RDS DB instance belongs to.

* `subnet_id` - Indicates the subnet id. 

* `vpc_id` - Indicates the VPC ID. 

The `db` block supports:

* `port` - Indicates the database port information.

* `type` - Indicates the DB engine. Value: MySQL, PostgreSQL, SQLServer.

* `user_name` - Indicates the default user name of database.

* `version` - Indicates the database version.

## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 30 minute.
- `delete` - Default is 30 minute.

## Import

RDS read replica instance can be imported by `id`, e.g.

```shell
$ terraform import huaweicloud_rds_read_replica_instance.replica_instance 92302c133d13424cbe357506ce057ea5in03
```