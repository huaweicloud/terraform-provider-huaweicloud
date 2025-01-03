---
subcategory: "CodeArts Deploy"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_deploy_environment_permission"
description: |-
  Manages a CodeArts deploy environment permission resource within HuaweiCloud.
---

# huaweicloud_codearts_deploy_environment_permission

Manages a CodeArts deploy environment permission resource within HuaweiCloud.

## Example Usage

```hcl
variable "application_id" {}
variable "environment_id" {}
variable "role_id" {}

resource "huaweicloud_codearts_deploy_environment_permission" "test" {
  application_id   = var.application_id
  environment_id   = var.environment_id
  role_id          = var.role_id
  permission_name  = "can_delete"
  permission_value = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `application_id` - (Required, String, ForceNew) Specifies the application ID.
  Changing this creates a new resource.

* `environment_id` - (Required, String, ForceNew) Specifies the environment ID.
  Changing this creates a new resource.

* `role_id` - (Required, String, ForceNew) Specifies the role ID.
  Changing this creates a new resource.

* `permission_name` - (Required, String, ForceNew) Specifies the permission name.
  Valid values are **can_view**, **can_edit**, **can_delete**, **can_deploy** and **can_manage**.

  Changing this creates a new resource.

* `permission_value` - (Optional, Bool) Specifies whether to enable the permission.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The CodeArts deploy environment permission resource can be imported using the `application_id`, `environment_id`,
`role_id` and `permission_name`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_codearts_deploy_environment_permission.test <app_id>/<env_id>/<role_id>/<permission_name>
```
