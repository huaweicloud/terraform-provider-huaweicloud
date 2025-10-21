---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpcs_app"
description: |-
  Manages a CPCS application resource within HuaweiCloud.
---

# huaweicloud_cpcs_app

Manages a CPCS application resource within HuaweiCloud.

-> Currently, this resource is valid only in cn-north-9 region.

## Example Usage

```hcl
variable "app_name" {}
variable "vpc_id" {}
variable "vpc_name" {}
variable "subnet_id" {}
variable "subnet_name" {}

resource "huaweicloud_cpcs_app" "test" {
  app_name    = var.app_name
  vpc_id      = var.vpc_id
  vpc_name    = var.vpc_name
  subnet_id   = var.subnet_id
  subnet_name = var.subnet_name
  description = "test application"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `app_name` - (Required, String, NonUpdatable) Specifies the application name. The name must be unique.

* `vpc_id` - (Required, String, NonUpdatable) Specifies the ID of the VPC to which the application belongs.

* `vpc_name` - (Required, String, NonUpdatable) Specifies the name of the VPC to which the application belongs.

* `subnet_id` - (Required, String, NonUpdatable) Specifies the ID of the subnet to which the application belongs.

* `subnet_name` - (Required, String, NonUpdatable) Specifies the name of the subnet to which the application belongs.

* `description` - (Optional, String, NonUpdatable) Specifies the application description.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (same as app_id).

* `create_time` - The creation time of the application, UNIX timestamp in milliseconds.

## Import

The CPCS application resource can be imported using the `app_name`, e.g.

```bash
$ terraform import huaweicloud_cpcs_app.test <app_name>
```
