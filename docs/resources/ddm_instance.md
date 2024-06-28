---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_instance"
description: ""
---

# huaweicloud_ddm_instance

Manages DDM instance resource within HuaweiCloud.

## Example Usage

```hcl
variable "flavor_id" {}
variable "engine_id" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}
variable "availability_zone" {}

resource "huaweicloud_ddm_instance" "test" {
  name              = "ddm_test"
  flavor_id         = var.flavor_id
  node_num          = 2
  engine_id         = var.engine_id
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.security_group_id
  
  availability_zones = [var.availability_zone]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the DDM instance.
  An instance name starts with a letter, consists of 4 to 64 characters, and can contain only letters,
  digits, and hyphens (-).

* `flavor_id` - (Required, String) Specifies the ID of a product.

* `node_num` - (Required, Int) Specifies the number of nodes.

* `engine_id` - (Required, String, ForceNew) Specifies the ID of an Engine.

  Changing this parameter will create a new resource.

* `availability_zones` - (Required, List, ForceNew) Specifies the list of availability zones.

  Changing this parameter will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of a VPC.

  Changing this parameter will create a new resource.

* `security_group_id` - (Required, String) Specifies the ID of a security group.

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of a subnet.

  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id.
  Value 0 indicates the default enterprise project.

* `param_group_id` - (Optional, String, ForceNew) Specifies the ID of parameter group.

  Changing this parameter will create a new resource.

* `time_zone` - (Optional, String, ForceNew) Specifies the time zone.

  Changing this parameter will create a new resource.

* `admin_user` - (Optional, String, ForceNew) Specifies the username of the administrator.
  The username starts with a letter, consists of 1 to 32 characters, and can contain only letters,
  digits, and underscores (_).

  Changing this parameter will create a new resource.

* `admin_password` - (Optional, String) Specifies the password of the administrator.
  The password consists of 8 to 32 characters, and must be a combination of uppercase letters,
  lowercase letters, digits, and the following special characters: ~!@#%^*-_=+?.

* `parameters` - (Optional, List) Specify an array of one or more parameters to be set to the instance after launched.
  The [parameters](#parameters_struct) structure is documented below.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the DDM instance.
  Valid values are **prePaid** and **postPaid**, defaults to **postPaid**.

  Changing this parameter will create a new resource.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit.
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.

  Changing this parameter will create a new resource.

* `period` - (Optional, Int, ForceNew) Specifies the charging period.
  If `period_unit` is set to **month**, the value ranges from 1 to 9.
  If `period_unit` is set to **year**, the value ranges from 1 to 3.
  This parameter is mandatory if `charging_mode` is set to **prePaid**.

  Changing this parameter will create a new resource.

* `auto_renew` - (Optional, String) Specifies whether auto-renew is enabled.
  Valid values are **true** and **false**. Defaults to **false**.

* `delete_rds_data` - (Optional, String) Specifies whether data stored on the associated DB instances is deleted.

<a name="parameters_struct"></a>
The `parameters` block supports:

* `name` - (Required, String) Specifies the parameter name. Some of them needs the instance to be restarted to take effect.

* `value` - (Required, String) Specifies the parameter value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the status of the DDM instance.

* `access_ip` - Indicates the address for accessing the DDM instance.

* `access_port` - Indicates the port for accessing the DDM instance.

* `engine_version` - Indicates the engine version.

* `nodes` - Indicates the node information.
  The [NodeInfoRef](#DdmInstance_NodeInfoRef) structure is documented below.

<a name="DdmInstance_NodeInfoRef"></a>
The `NodeInfoRef` block supports:

* `status` - Indicates the status of the DDM instance node.

* `port` - Indicates the port of the DDM instance node.

* `ip` - Indicates the IP address of the DDM instance node.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 15 minutes.
* `update` - Default is 60 minutes.
* `delete` - Default is 10 minutes.

## Import

The DDM instance can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ddm_instance.test <id>
```
