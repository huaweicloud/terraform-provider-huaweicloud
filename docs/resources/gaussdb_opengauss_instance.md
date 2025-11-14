---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_instance"
description: |-
  GaussDB OpenGauss instance management within HuaweiCould.
---

# huaweicloud_gaussdb_opengauss_instance

GaussDB OpenGauss instance management within HuaweiCould.

## Example Usage

### Create an instance for distributed HA mode

```hcl
variable "vpc_id" {}
variable "subnet_network_id" {}
variable "security_group_id" {}
variable "instance_name" {}
variable "instance_password" {}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_network_id
  security_group_id = var.security_group_id

  flavor            = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  name              = var.instance_name
  password          = var.instance_password
  sharding_num      = 1
  coordinator_num   = 2
  availability_zone = join(",", slice(data.huaweicloud_availability_zones.test.names, 0, 3))

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "strong"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
```

### Create an instance for centralized HA mode

```hcl
variable "instance_name" {}
variable "instance_password" {}
variable "vpc_id" {}
variable "subnet_network_id" {}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_gaussdb_opengauss_instance" "instance_acc" {
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_network_id
  security_group_id = var.security_group_id
  name              = var.instance_name
  password          = var.instance_password
  flavor            = "gaussdb.opengauss.ee.m6.2xlarge.x868.ha"
  availability_zone = join(",", slice(data.huaweicloud_availability_zones.myaz.names, 0, 3))

  replica_num = 3

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "strong"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the instance.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the instance name, which can be the same as an existing instance name.
  The value must be `4` to `64` characters in length and start with a letter. It is case-sensitive and can contain only
  letters, digits, hyphens (-), and underscores (_).

* `flavor` - (Required, String) Specifies the instance specifications.

* `password` - (Required, String) Specifies the database password. The value must be `8` to `32` characters in length,
  including uppercase and lowercase letters, digits, and special characters, such as **~!@#%^*-_=+?**. You are advised
  to enter a strong password to improve security, preventing security risks such as brute force cracking.

* `availability_zone` - (Required, String, ForceNew) Specifies the availability zone information, can be three same or
  different az like **cn-north-4a,cn-north-4a,cn-north-4a**. Changing this parameter will create a new resource.

* `ha` - (Required, List, ForceNew) Specifies the HA information.
  The [object](#opengauss_ha) structure is documented below.
  Changing this parameter will create a new resource.

* `volume` - (Required, List) Specifies the volume storage information.
  The [object](#opengauss_volume) structure is documented below.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID to which the subnet belongs.
  Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the network ID of VPC subnet to which the instance belongs.
  Changing this parameter will create a new resource.

* `security_group_id` - (Optional, String, ForceNew) Specifies the security group ID to which the instance belongs.
  If the `port` parameter is specified, please ensure that the TCP ports in the inbound rule of security group
  includes the `100` ports starting with the database port.
  (For example, if the database port is `8,000`, the TCP port must include the range from `8,000` to `8,100`.)

  Changing this parameter will create a new resource.

* `port` - (Optional, String, ForceNew) Specifies the port information. Defaults to `8,000`.
  The valid values are as follows:
  + `2,378` to `2,380`
  + `4999` to `5,000`
  + `5,999` to `6,001`
  + `8,097` to `8,098`
  + `12,016` to `12,017`
  + `20,049` to `20,050`
  + `21,731` to `21,732`
  + `32,122` to `32,124`

  Changing this parameter will create a new resource.

* `configuration_id` - (Optional, String) Specifies the parameter template ID.
  Changing this parameter will create a new resource.

* `sharding_num` - (Optional, Int) Specifies the sharding number.  
  The valid value is range form `1` to `9`.

* `coordinator_num` - (Optional, Int) Specifies the coordinator number.  
  The valid value is range form `1` to `9`.
  The value must not be greater than twice value of `sharding_num`.

* `replica_num` - (Optional, Int, ForceNew) The replica number. The valid values are `2` and `3`.
  Double replicas are only available for specific users and supports only instance versions are v1.3.0 or later.
  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `time_zone` - (Optional, String, ForceNew) Specifies the time zone.
  Changing this parameter will create a new resource.

* `disk_encryption_id` - (Optional, String) Specifies the key ID for disk encryption.

* `enable_force_switch` - (Optional, Bool) Specifies whether to forcibly promote a standby node to primary.
  Defaults to **false**.

* `enable_single_float_ip` - (Optional, Bool) Specifies whether to enable single floating IP address policy, which is only
  suitable for primary/standby instances. Value options:
  + **true**: This function is enabled. Only one floating IP address is bound to the primary node of a DB instance. If a
    primary/standby fail over occurs, the floating IP address does not change.
  + **false (default value)**: The function is disabled. Each node is bound to a floating IP address. If a primary/standby
    fail over occurs, the floating IP addresses change.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the GaussDB OpenGauss instance.

* `force_import` - (Optional, Bool) Specifies whether to import the instance with the given configuration instead of
  creation. If specified, try to import the instance instead of creation if the instance already existed.

* `datastore` - (Optional, List, ForceNew) Specifies the datastore information.
  The [datastore](#opengauss_datastore) structure is documented below.
  Changing this parameter will create a new resource.

* `backup_strategy` - (Optional, List) Specifies the advanced backup policy.
  The [backup_strategy](#opengauss_backup_strategy) structure is documented below.

* `parameters` - (Optional, List) Specifies an array of one or more parameters to be set to the instance after launched.
  The [parameters](#parameters_struct) structure is documented below.

* `mysql_compatibility_port` - (Optional, String) Specifies the port for MySQL compatibility. Value range: **0** or
  **1024** to **39989**.
  + The following ports are used by the system and cannot be used: **2378**, **2379**, **2380**, **2400**, **4999**,
    **5000**, **5001**, **5100**, **5500**, **5999**, **6000**, **6001**, **6009**, **6010**, **6500**, **8015**, **8097**,
    **8098**, **8181**, **9090**, **9100**, **9180**, **9187**, **9200**, **12016**, **12017**, **20049**, **20050**,
    **21731**, **21732**, **32122**, **32123**, **32124**, **32125**, **32126**, **39001**,
    **[Database port, Database port + 10]**.
  + If the value is **0**, the MySQL compatibility port is disabled.

* `advance_features` - (Optional, List) Specifies the advanced features.
  The [advance_features](#advance_features_struct) structure is documented below.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of opengauss instance.
  The valid values are as follows:
  + **prePaid**: the yearly/monthly billing mode.
  + **postPaid**: the pay-per-use billing mode.

  Defaults to **postPaid**. Changing this parameter will create a new resource.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of opengauss instance.
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.
  Changing this parameter will create a new resource.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of opengauss instance.
  If `period_unit` is set to **month**, the value ranges from 1 to 9.
  If `period_unit` is set to **year**, the value ranges from 1 to 5.
  This parameter is mandatory if `charging_mode` is set to **prePaid**.
  Changing this parameter will create a new resource.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled.
  Valid values are **true** and **false**. Defaults to **false**.

<a name="opengauss_ha"></a>
The `ha` block supports:

* `mode` - (Required, String, ForceNew) Specifies the deployment model.
  The valid values are **enterprise** and **centralization_standard**.
  Changing this parameter will create a new resource.

* `replication_mode` - (Required, String, ForceNew) Specifies the database replication mode.
  Only **sync** is supported now. Changing this parameter will create a new resource.

* `consistency` - (Optional, String, ForceNew) Specifies the database consistency mode.
  The valid values are **strong** and **eventual**, not case-sensitive.
  Changing this parameter will create a new resource.

* `instance_mode` - (Optional, String, ForceNew) Specifies the product type of the instance. Value options:
  + **enterprise**: The instance of the enterprise edition will be created.
  + **basic**: The instance of the basic edition will be created.
  + **ecology**: The instance of the ecosystem edition will be created.

  Changing this parameter will create a new resource.

<a name="opengauss_volume"></a>
The `volume` block supports:

* `type` - (Required, String, ForceNew) Specifies the volume type. Only **ULTRAHIGH** is supported now.
  Changing this parameter will create a new resource.

* `size` - (Required, Int) Specifies the volume size (in gigabytes). The valid value is range form `40` to `4,000`.

<a name="opengauss_datastore"></a>
The `datastore` block supports:

* `engine` - (Required, String, ForceNew) Specifies the database engine. Only **GaussDB(for openGauss)** is supported
  now. Changing this parameter will create a new resource.

* `version` - (Optional, String, ForceNew) Specifies the database version. Defaults to the latest version. Please
  reference to the API docs for valid options. Changing this parameter will create a new resource.

<a name="opengauss_backup_strategy"></a>
The `backup_strategy` block supports:

* `start_time` - (Required, String) Specifies the backup time window. Automated backups will be triggered during the
  backup time window. It must be a valid value in the **hh:mm-HH:MM** format. The current time is in the UTC format. The
  **HH** value must be `1` greater than the **hh** value. The values of mm and MM must be the same and must be set to
  **00**. Example value: **08:00-09:00**, **23:00-00:00**.

* `keep_days` - (Optional, Int) Specifies the number of days to retain the generated backup files. The value ranges from
  `0` to `732`. If this parameter is set to `0`, the automated backup policy is not set.
  If this parameter is not transferred, the automated backup policy is enabled by default.

<a name="parameters_struct"></a>
The `parameters` block supports:

* `name` - (Required, String) Specifies the name of the parameter.

* `value` - (Required, String) Specifies the value of the parameter.

<a name="advance_features_struct"></a>
The `advance_features` block supports:

* `name` - (Required, String) Specifies the name of the advance feature.

* `value` - (Required, String) Specifies the value of the advance feature.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the DB instance ID.

* `status` - Indicates the DB instance status.

* `type` - Indicates the database type.

* `private_ips` - Indicates the private IP address of the DB instance.

* `public_ips` - Indicates the public IP address of the DB instance.

* `endpoints` - Indicates the connection endpoints list of the DB instance. Example: [127.0.0.1:8000].

* `db_user_name` - Indicates the default username.

* `switch_strategy` - Indicates the switch strategy.

* `balance_status` - Indicates whether the host load is balanced due to a primary/standby switchover.

* `error_log_switch_status` - Indicates whether error log collection is enabled. The value can be:
  + **ON**: enabled
  + **OFF**: disabled

* `maintenance_window` - Indicates the maintenance window.

* `nodes` - Indicates the instance nodes information. Structure is documented below.

The `nodes` block contains:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `role` - Indicates the node role.
  + **master**.
  + **slave**.

* `status` - Indicates the node status.

* `availability_zone` - Indicates the availability zone of the node.

* `private_ip` - Indicates the private IP address of the node.

* `public_ip` - Indicates the EIP that has been bound.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 120 minutes.
* `update` - Default is 150 minutes.
* `delete` - Default is 45 minutes.

## Import

OpenGaussDB instance can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_gaussdb_opengauss_instance.test <id>
```

Note that the imported state may not be identical to your resource definition, due to the attribute missing from the
API response. The missing attributes include: `password`, `ha.0.mode`, `ha.0.instance_mode`, `configuration_id`,
`disk_encryption_id`, `enable_force_switch`, `enable_single_float_ip`, `parameters`, `period_unit`, `period` and
`auto_renew`. It is generally recommended running `terraform plan` after importing a GaussDB OpenGauss instance. You can
then decide if changes should be applied to the GaussDB OpenGauss instance, or the resource definition should be updated
to align with the GaussDB OpenGauss instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  ...

  lifecycle {
    ignore_changes = [
      password, configuration_id, disk_encryption_id, enable_force_switch, enable_single_float_ip, parameters, period_unit,
      period, auto_renew,
    ]
  }
}
```
