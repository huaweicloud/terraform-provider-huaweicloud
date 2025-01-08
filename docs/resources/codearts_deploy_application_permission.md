---
subcategory: "CodeArts Deploy"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_deploy_application_permission"
description: |-
  Manages a CodeArts deploy application permission resource within HuaweiCloud.
---

# huaweicloud_codearts_deploy_application_permission

Manages a CodeArts deploy application permission resource within HuaweiCloud.

-> Only when the applications using instance level permission, this resource is available.

## Example Usage

```hcl
variable "project_id" {}
variable "application_ids" {}
variable "role_id" {}

resource "huaweicloud_codearts_deploy_group_permission" "test" {
  project_id      = var.project_id
  application_ids = var.application_ids

  roles {
    role_id        = var.role_id
    can_modify     = true
    can_disable    = true
    can_delete     = true
    can_view       = true
    can_execute    = true
    can_copy       = true
    can_manage     = true
    can_create_env = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `project_id` - (Required, String, ForceNew) Specifies the project ID for CodeArts service.
  Changing this creates a new resource.

* `application_ids` - (Required, List) Specifies the application IDs.

* `roles` - (Required, List) Specifies the role permissions list.
  The [roles](#block--roles) structure is documented below.

<a name="block--roles"></a>
The `roles` block supports:

* `role_id` - (Required, String) Specifies the role ID.

* `can_copy` - (Required, Bool) Specifies whether the role has the copy permission.

* `can_create_env` - (Required, Bool) Specifies whether the role has the permission to create an environment.

* `can_delete` - (Required, Bool) Specifies whether the role has the deletion permission.

* `can_disable` - (Required, Bool) Specifies whether the role has the permission to disable application.

* `can_execute` - (Required, Bool) Specifies whether the role has the deployment permission.

* `can_manage` - (Required, Bool) Specifies whether the role has the management permission, including adding, deleting,
  modifying, querying deployment and permission modification.

* `can_modify` - (Required, Bool) Specifies whether the role has the editing permission.

* `can_view` - (Required, Bool) Specifies whether the role has the view permission.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
