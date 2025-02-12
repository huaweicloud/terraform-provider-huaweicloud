---
subcategory: "CodeArts Deploy"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_deploy_group_permission"
description: |-
  Manages a CodeArts deploy group permission resource within HuaweiCloud.
---

# huaweicloud_codearts_deploy_group_permission

Manages a CodeArts deploy group permission resource within HuaweiCloud.

## Example Usage

```hcl
variable "project_id" {}
variable "group_id" {}
variable "role_id" {}

resource "huaweicloud_codearts_deploy_group_permission" "test" {
  project_id       = var.project_id
  group_id         = var.group_id
  role_id          = var.role_id
  permission_name  = "can_add_host"
  permission_value = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `project_id` - (Required, String, ForceNew) Specifies the project ID.
  Changing this creates a new resource.

* `group_id` - (Required, String, ForceNew) Specifies the group ID.
  Changing this creates a new resource.

* `role_id` - (Required, String, ForceNew) Specifies the role ID.
  Changing this creates a new resource.

* `permission_name` - (Required, String, ForceNew) Specifies the permission name.
  Valid values are **can_view**, **can_edit**, **can_delete**, **can_add_host**, **can_manage**, and **can_copy**.

  Changing this creates a new resource.

* `permission_value` - (Optional, Bool) Specifies whether to enable the permission.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.

## Import

The CodeArts deploy group permission resource can be imported using the `project_id`, `group_id`, `role_id` and
`permission_name`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_codearts_deploy_group_permission.test <project_id>/<group_id>/<role_id>/<permission_name>
```
