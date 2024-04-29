---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_instance_attach"
description: ""
---

# huaweicloud_as_instance_attach

Manages an AS instance attachment resource within HuaweiCloud.

## Example Usage

### Add an instance

```hcl
variable "scaling_group_id" {}
variable "ecs_id" {}

resource "huaweicloud_as_instance_attach" "test" {
  scaling_group_id = var.scaling_group_id
  instance_id      = var.ecs_id
}
```

### Add an instance with protection

```hcl
variable "scaling_group_id" {}
variable "ecs_id" {}

resource "huaweicloud_as_instance_attach" "test" {
  scaling_group_id = var.scaling_group_id
  instance_id      = var.ecs_id
  protected        = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `scaling_group_id` - (Required, String, ForceNew) Specifies the AS group ID.
  Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ECS instance ID.
  Changing this creates a new resource.

  <!--markdownlint-disable MD033-->
  -> Before adding instances to an AS group, ensure that the following conditions are met:
  <br/>1. The instance is not in other AS groups.
  <br/>2. The instance is in the same VPC as the AS group.
  <br/>3. The instance is in the AZs used by the AS group.
  <br/>4. After the instance is added, the total number of instances is less than or equal to the maximum number of
  instances allowed.

* `protected` - (Optional, Bool) Specifies whether the instance can be removed **automatically** from the AS group.
  Once configured, when AS automatically scales in the AS group, the instance that is protected will not be removed.

* `standby` - (Optional, Bool) Specifies whether to stop distributing traffic to the instance but do not want to remove
  it from the AS group. You can stop or restart the instance without worrying about it will be removed from the AS group.

* `append_instance` - (Optional, Bool) Specifies whether to add a new instance when the instance enter standby mode.
  This parameter takes effect only when `standby` is set to true.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of `scaling_group_id`/`instance_id`.
* `instance_name` - The ECS instance name.
* `health_status` - The instance health status. The value can be `INITIALIZING`, `NORMAL` or `ERROR`.
* `status` - The instance lifecycle status in the AS group. The value can be `INSERVICE`, `STANDBY`, `PENDING` or `REMOVING`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The AS instances can be imported by the `scaling_group_id` and `instance_id`, separated by a slash, e.g.

```shell
$ terraform import huaweicloud_as_instance_attach.test <scaling_group_id>/<instance_id>
```

Note that the imported state may not be identical to your resource definition, due to `append_instance` is missing from
the API response.

It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_as_instance_attach" "test" {
  ...

  lifecycle {
    ignore_changes = [append_instance]
  }
}
```
