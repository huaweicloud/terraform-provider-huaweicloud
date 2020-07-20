---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_instance"
sidebar_current: "docs-huaweicloud-resource-gaussdb-opengauss-instance"
description: |-
  GaussDB OpenGauss instance management
---

# huaweicloud\_gaussdb\_opengauss\_instance

GaussDB OpenGauss instance management within HuaweiCoud.

## Example Usage

### create a basic instance

```hcl
resource "huaweicloud_gaussdb_opengauss_instance" "instance_acc" {
  name        = "opengaussdb_instance_1"
  password    = "Test@123"
  flavor      = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  vpc_id      = var.vpc_id
  subnet_id   = var.subnet_id
  availability_zone = "cn-north-4a,cn-north-4a,cn-north-4a"

  ha {
    mode = "enterprise"
    replication_mode = "sync"
    consistency = "strong"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
  security_group_id = var.secgroup.id
  sharding_num = 1
  coordinator_num = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the instance name, which can be the same
  as an existing instance name. The value must be 4 to 64 characters in
  length and start with a letter. It is case-sensitive and can contain
  only letters, digits, hyphens (-), and underscores (_).
  Changing this parameter will create a new resource.

* `flavor` - (Required) Specifies the instance specifications. Please reference
  the API docs for valid options. Changing this parameter will create a new resource.

* `password` - (Required) Specifies the database password. The value must be 8 to 32 characters
  in length, including uppercase and lowercase letters, digits, and special characters,
  such as ~!@#%^*-_=+?

  You are advised to enter a strong password to improve security, preventing security risks
  such as brute force cracking.
  Changing this parameter will create a new resource.

* `availability_zone` -  (Required) Specifies the Availability Zone information, can be three same or
  different az like "cn-north-4a,cn-north-4a,cn-north-4a".
  Changing this parameter will create a new resource.

* `vpc_id` -  (Required) Specifies the VPC ID.
  Changing this parameter will create a new resource.

* `subnet_id` - (Required) Specifies the network ID of a subnet.
  Changing this parameter will create a new resource.

* `security_group_id` - (Required) Specifies the security group ID.
  Changing this parameter will create a new resource.

* `volume` - (Required) Specifies the volume storage information. Structure is documented below.

* `port` - (Optional) Specifies the port information. Defaults to "8000".
  Changing this parameter will create a new resource.

* `configuration_id` - (Optional) The parameter template id.
  Changing this parameter will create a new resource.

* `dsspool_id` - (Optional) The Dec dedicated storage pool id.
  Changing this parameter will create a new resource.

* `sharding_num` - (Optional) The Sharding num. Values: 1~32.

* `coordinator_num` - (Optional) The Coordinator num. Values: 1~32.

* `disk_encryption_id` - (Optional) Specifies the volume encryption KMS id.
  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional) The enterprise project id.
  Changing this parameter will create a new resource.

* `time_zone` - (Optional) Specifies the time zone. Defaults to "UTC+08:00".
  Changing this parameter will create a new resource.

* `datastore` - (Optional) Specifies the datastore information. Structure is documented below.
  Changing this parameter will create a new resource.

* `ha` - (Optional) Specifies the HA information. Structure is documented below.
  Changing this parameter will create a new resource.

The `datastore` block supports:

* `engine` - (Required) Specifies the database engine. Only "GaussDB(openGauss)" is supported now.

* `version` - (Required) Specifies the database version. Only "1.0" is supported now.


The `volume` block supports:

* `type` - (Required) Specifies the volume type. Only "ULTRAHIGH" is supported now.

* `size` - (Required) Specifies the volume size (in gigabytes) for a single DN. The value should between 40G ~ 5TB.


The `ha` block supports:

* `mode` - (Required) Specifies the database mode. Only "enterprise" is supported now.

* `replication_mode` - (Required) Specifies the database replication mode. Only "sync" is supported now.

* `consistency` - (Optional) Specifies the database consistency mode. Valid options are "strong" and "eventual".

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - Indicates the DB instance ID.
* `status` - Indicates the DB instance status.
* `region` - Indicates the region where the DB instance is deployed.
* `type` - Indicates the database type.
* `private_ips` - Indicates the private IP address of the DB instance.
* `public_ips` - Indicates the public IP address of the DB instance.
* `db_user_name` - Indicates the default username.
* `switch_strategy` - Indicates the switch strategy.
* `maintenance_window` - Indicates the maintenance window.
* `nodes` - Indicates the instance nodes information. Structure is documented below.

The `nodes` block contains:

- `id` - Indicates the node ID.
- `name` - Indicates the node name.
- `role` - Indicates the node role: master or slave.
- `status` - Indicates the node status.
- `availability_zone` - Indicates the availability zone of the node.

## Import

OpenGaussDB instance can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_gaussdb_opengauss_instance.instance_1 ee678f40-ce8e-4d0c-8221-38dead426f06
```
