---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_adg_instance"
description: |-
  Manages a DSC ADG (API Data Security Gateway) instance resource within HuaweiCloud.
---

# huaweicloud_dsc_adg_instance

Manages a DSC ADG (API Data Security Gateway) instance resource within HuaweiCloud.

## Example Usage

```hcl
variable "availability_zone" {}
variable "specification" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}
variable "password" {}
variable "public_ip_id" {}


resource "huaweicloud_dsc_adg_instance" "test" {
  name              = "adg_test"
  availability_zone = var.availability_zone
  type              = "ADG"
  specification     = var.specification
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.security_group_id
  admin_name        = "sysadmin"
  password          = var.password
  deploy_mode       = "CLOUD"
  mode              = "ha"
  publicip_id       = var.public_ip_id

  charge_mode   = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the ADG instance.

* `type` - (Required, String, NonUpdatable) Specifies the type of the ADG instance.
  The valid values are as follows:
  + **ADG**: API security gateway.
  + **DEG**: database security gateway.
  + **DOM**: database O&M.
  + **ENGINE**: engine.

* `specification` - (Required, String, NonUpdatable) Specifies the specification of the ADG instance.
  The value must exist in the product metadata.

* `vpc_id` - (Required, String, NonUpdatable) Specifies the VPC ID to which the ADG instance belongs.

* `subnet_id` - (Required, String, NonUpdatable) Specifies the subnet ID to which the ADG instance belongs.

* `security_group_id` - (Required, String, NonUpdatable) Specifies the security group ID of the ADG instance.

* `availability_zone` - (Required, String, NonUpdatable) Specifies the availability zone of the ADG instance.

* `deploy_mode` - (Optional, String, NonUpdatable) Specifies the deploy mode of the ADG instance.
  The valid values are as follows:
  + **CLOUD**: cloud-deployed.
  + **OUTSIDE**: off-cloud.

* `mode` - (Optional, String, NonUpdatable) Specifies the mode of the ADG instance.
  The valid values are as follows:
  + **ha**: high availability.
  + **single**: standalone.

* `password` - (Optional, String) Specifies the password of the ADG instance.

* `admin_name` - (Optional, String) Specifies the administrator name.
  The valid values are **sysadmin**, **secadmin**, and **audadmin**.

* `publicip_id` - (Optional, String) Specifies the public IP ID to bind to the ADG instance.

* `outside_ins_config` - (Optional, List, NonUpdatable) Specifies the cloud outside instance configuration.
  This parameter is mandatory for off-cloud-deployed instances.
  The [outside_ins_config](#outside_ins_config_struct) structure is documented below.

* `charge_mode` - (Required, String, NonUpdatable) Specifies the charging mode of the ADG instance.
  The valid values are as follows:
  + **prePaid**: yearly/monthly.
  + **postPaid**: pay-per-use.

  > **Note:** Currently, only **prePaid** is supported.

* `period_unit` - (Required, String, NonUpdatable) Specifies the charging period unit of the ADG instance.
  The valid values are as follows:
  + **month**
  + **year**

* `period` - (Required, Int, NonUpdatable) Specifies the charging period of the ADG instance.
  + If `period_unit` is set to **month**, the value ranges from `1` to `9`.
  + If `period_unit` is set to **year**, the value ranges from `1` to `3`.

* `auto_renew` - (Optional, String) Specifies whether auto-renew is enabled. Defaults to **false**.
  The valid values are **true** and **false**.

<a name="outside_ins_config_struct"></a>
The `outside_ins_config` block supports:

* `master_node_ip` - (Optional, String, NonUpdatable) Specifies the master node IP address.

* `slave_node_ip` - (Optional, String, NonUpdatable) Specifies the slave node IP address.

* `virtual_ip` - (Optional, String, NonUpdatable) Specifies the virtual IP address.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also the ADG instance ID.

* `status` - The status of the ADG instance.

* `public_ip` - The public IP address of the ADG instance.

* `virtual_ip` - The virtual IP address of the ADG instance.

* `version` - The version of the ADG instance.

* `create_time` - The creation time of the ADG instance.

* `started_time` - The start time of the ADG instance.

* `fail_reason` - The failure reason of the ADG instance.

* `nodes` - The node information of the ADG instance.
  The [nodes](#nodes_struct) structure is documented below.

<a name="nodes_struct"></a>
The `nodes` block supports:

* `id` - The node ID.

* `name` - The node name.

* `availability_zone` - The availability zone of the node.

* `private_ip` - The private IP address of the node.

* `role` - The role of the node.

* `status` - The status of the node.

* `vm_id` - The VM ID of the node.

* `error_reason` - The error reason of the node.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The ADG instance resource can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_dsc_adg_instance.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `password`, `admin_name`, `outside_ins_config`, `period_unit`, `period`, `auto_renew`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dsc_adg_instance" "test" {
  ...

  lifecycle {
    ignore_changes = [
      password, admin_name, outside_ins_config, period_unit, period, auto_renew, admin_password,
    ]
  }
}
```
