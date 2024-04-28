---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_account_assignment"
description: ""
---

# huaweicloud_identitycenter_account_assignment

Manages an Identity Center account assignment resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "permission_set_id" {}
variable "principal_id" {}
variable "target_id" {}

resource "huaweicloud_identitycenter_account_assignment" "test"{
  instance_id       = var.instance_id
  permission_set_id = var.permission_set_id
  principal_id      = var.principal_id
  principal_type    = "USER"
  target_id         = var.target_id
  target_type       = "ACCOUNT"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the Identity Center instance.

  Changing this parameter will create a new resource.

* `permission_set_id` - (Required, String, ForceNew) Specifies the ID of the permission set.

  Changing this parameter will create a new resource.

* `principal_id` - (Required, String, ForceNew) Specifies the ID of the user or user group that belongs to IAM
  Identity Center.

  Changing this parameter will create a new resource.

* `principal_type` - (Required, String, ForceNew) Specifies the type of the user or user group.
  Value options: **USER**, **GROUP**.

  Changing this parameter will create a new resource.

* `target_id` - (Required, String, ForceNew) Specifies the ID of the target to be bound.

  Changing this parameter will create a new resource.

* `target_type` - (Required, String, ForceNew) Specifies the type of the target to be bound. Value options: **ACCOUNT**.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

The Identity Center account permission can be imported using the `instance_id`, `permission_set_id`,`target_id`
and `principal_id` separated by slashes, e.g.

```bash
$ terraform import huaweicloud_identitycenter_account_assignment.test <instance_id>/<permission_set_id>/<target_id>/<principal_id>
```
