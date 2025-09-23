---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_studio_instance"
description: ""
---

# huaweicloud_dataarts_studio_instance

Manages DataArts Studio instance resource within HuaweiCloud.

-> Only **prePaid** charging mode is supported.

## Example Usage

```hcl
variable "availability_zone" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "secgroup_id" {}

resource "huaweicloud_dataarts_studio_instance" "my_demo" {
  name                  = "DataArts-demo"
  version               = "dayu.starter"
  vpc_id                = var.vpc_id
  subnet_id             = var.subnet_id
  security_group_id     = var.secgroup_id
  availability_zone     = var.availability_zone
  period_unit           = "month"
  period                = 1
  enterprise_project_id = "0"

  tags = {
    key = "value"
  }
}
```

<!--markdownlint-disable MD033-->

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to manage the DataArts Studio instance.
  Changing this creates a new instance.

* `name` - (Required, String, ForceNew) Specifies the DataArts Studio instance name. Changing this creates a new instance.

* `version` - (Required, String, ForceNew) Specifies the DataArts Studio version version.
  The valid values are **dayu.starter**, **dayu.nb.professional** and **dayu.nb.enterprise**.
  Changing this creates a new instance.

* `availability_zone` - (Required, String, ForceNew) Specifies the AZ name. Changing this creates a new instance.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID. Changing this creates a new instance.

* `subnet_id` - (Required, String, ForceNew) Specifies the VPC subnet ID. Changing this creates a new instance.

* `security_group_id` - (Required, String, ForceNew) Specifies the security group ID. Changing this creates a new instance.

* `period_unit` - (Required, String, ForceNew) Specifies the charging period unit of the instance.
  Valid values are **month** and **year**.
  Changing this creates a new instance.

* `period` - (Required, Int, ForceNew) Specifies the charging period of the DataArts Studio instance.
  If `period_unit` is set to **month**, the value ranges from `1` to `9`.
  If `period_unit` is set to **year**, the value ranges from `1` to `3`.
  Changing this creates a new instance.

* `auto_renew` - (Optional, String, ForceNew) Specifies whether auto renew is enabled.
  Valid values are `true` and `false`, defaults to `false`. Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id of the instance.

  -> 1. Only **one** DataArts Studio instance can be purchased in an enterprise project.
  <br/> 2. If DataArts Studio needs to communicate with other cloud services, ensure that the enterprise project
    of DataArts Studio is the same as that of other cloud services.

* `tags` - (Optional, Map, ForceNew) The key/value pairs to associate with the DataArts Studio instance.
  Changing this creates a new instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
* `charging_mode` - The charging mode. The value is `prePaid` indicates the yearly/monthly billing mode.
* `order_id` - The order ID of this DataArts Studio instance.
* `expire_days` - The expire days to renew.
* `status` - The status of this DataArts Studio instance.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `delete` - Default is 30 minutes.

## Import

DataArts Studio instances can be imported using their `id`, e.g.

```sh
terraform import huaweicloud_dataarts_studio_instance.instance e60361de2cfd42d7a6b673f0ae58db82
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `tags`, `period_unit`, `period`, `auto_renew`.
It is generally recommended running `terraform plan` after importing an instance.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dataarts_studio_instance" "instance" {
    ...

  lifecycle {
    ignore_changes = [
      tags, period_unit, period, auto_renew,
    ]
  }
}
```
