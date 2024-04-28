---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_user_role_assignment"
description: ""
---

# huaweicloud_identity_user_role_assignment

Manages an IAM user role assignment within HuaweiCloud IAM.

-> **NOTE:** 1. You *must* have admin privileges to use this resource.
  <br/>2. When the resource is created, the permissions will take effect after 15 to 30 minutes.

## Example Usage

```hcl
variable "enterprise_project_id" {}
variable "user_1_password" {}

data "huaweicloud_identity_role" "test" {
  display_name = "ECS FullAccess"
}

resource "huaweicloud_identity_user" "test" {
  name        = "user_1"
  description = "A user"
  password    = var.user_1_password
}

resource "huaweicloud_identity_user_role_assignment" "test" {
  user_id               = huaweicloud_identity_user.test.id
  role_id               = data.huaweicloud_identity_role.test.id
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required, String, ForceNew) Specifies the the ID of user to assign the role to.
  Changing this parameter will create a new resource.

* `role_id` - (Required, String, ForceNew) Specifies the role to assign.
  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Required, String, ForceNew) Specifies the ID of the enterprise project
  to assign the role in. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The format is `<user_id>/<role_id>/<enterprise_project_id>`.

## Import

The role assignments can be imported using the `user_id`, `role_id` and  `enterprise_project_id`, e.g.

```bash
$ terraform import huaweicloud_identity_user_role_assignment.test <user_id>/<role_id>/<enterprise_project_id>
```
