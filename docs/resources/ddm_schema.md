---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_schema"
description: ""
---

# huaweicloud_ddm_schema

Manages a DDM schema resource within HuaweiCloud.

## Example Usage

```hcl
variable "rds_password" {}

resource "huaweicloud_vpc" "test" {
  name = "test_vpc"
  cidr = "192.168.0.0/24"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "test_subnet"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "test_secgroup"
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_ddm_engines" test {
  version = "3.0.8.5"
}

data "huaweicloud_ddm_flavors" test {
  engine_id = data.huaweicloud_ddm_engines.test.engines[0].id
  cpu_arch  = "X86"
}

resource "huaweicloud_ddm_instance" "test" {
  name              = "ddm-test"
  flavor_id         = data.huaweicloud_ddm_flavors.test.flavors[0].id
  node_num          = 2
  engine_id         = data.huaweicloud_ddm_engines.test.engines[0].id
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
}

resource "huaweicloud_rds_instance" "test" {
  name              = "rds_test"
  flavor            = "rds.mysql.n1.large.4"
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  fixed_ip          = "192.168.0.58"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  db {
    password = var.rds_password
    type     = "MySQL"
    version  = "5.7"
    port     = 3306
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

resource "huaweicloud_ddm_schema" "test"{
  instance_id  = huaweicloud_ddm_instance.test.id
  name         = "test_schema"
  shard_mode   = "single"
  shard_number = 1

  data_nodes {
    id             = huaweicloud_rds_instance.test.id
    admin_user     = "root"
    admin_password = var.rds_password
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of a DDM instance.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the DDM schema.
  An instance name starts with a letter, consists of 2 to 48 characters, and can contain only lowercase letters,
  digits, and underscores (_). Cannot contain keywords information_schema, mysql, performance_schema, or sys.

  Changing this parameter will create a new resource.

* `shard_mode` - (Required, String, ForceNew) Specifies the sharding mode of the schema. Values option: **cluster**, **single**.
  + **cluster**: indicates that the schema is in sharded mode.
  + **single**: indicates that the schema is in non-sharded mode.

  Changing this parameter will create a new resource.

* `shard_number` - (Required, Int, ForceNew) Specifies the number of shards in the same working mode.
  The value must be greater than or equal to the number of associated RDS instances and less than or equal
  to the number of associated instances multiplied by 64.

  Changing this parameter will create a new resource.

* `data_nodes` - (Required, List, ForceNew) Specifies the RDS instances associated with the schema.

  Changing this parameter will create a new resource.
  The [DataNode](#DdmSchema_DataNode) structure is documented below.

* `delete_rds_data` - (Optional, String) Specifies whether data stored on the associated DB instances is deleted.

<a name="DdmSchema_DataNode"></a>
The `DataNode` block supports:

* `id` - (Required, String) Specifies the ID of the RDS instance associated with the schema.

* `admin_user` - (Required, String) Specifies the username for logging in to the associated RDS instance.

* `admin_password` - (Required, String) Specifies the password for logging in to the associated RDS instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the schema status.

* `shards` - Indicates the sharding information of the schema.
  The [Shard](#DdmSchema_Shard) structure is documented below.

* `data_nodes` - Indicates the RDS instances associated with the schema.
  The [DataNode](#DdmSchema_DataNode) structure is documented below.

* `data_vips` - Indicates the IP address and port number for connecting to the schema.

<a name="DdmSchema_Shard"></a>
The `Shard` block supports:

* `db_slot` - Indicates the number of shards.

* `name` - Indicates the shard name.

* `status` - Indicates the shard status.

* `id` - Indicates the ID of the RDS instance where the shard is located.

<a name="DdmSchema_DataNode"></a>
The `DataNode` block supports:

* `id` - Indicates the ID of the RDS instance associated with the schema.

* `name` - Indicates the name of the associated RDS instance.

* `status` - Indicates the status of the associated RDS instance.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 30 minutes.

## Import

The DDM schema can be imported using the `<instance_id>/<schema_name>`, e.g.

```
$ terraform import huaweicloud_ddm_schema.test 80e373f9-872e-4046-aae9-ccd9ddc55511/schema_name
```
