---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_permission"
description: |-
  Manages a CodeArts pipeline permission resource within HuaweiCloud.
---

# huaweicloud_codearts_pipeline_permission

Manages a CodeArts pipeline permission resource within HuaweiCloud.

## Example Usage

### Modify user permission

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}
variable "user_id" {}

resource "huaweicloud_codearts_pipeline_permission" "user" {
  project_id        = var.codearts_project_id
  pipeline_id       = var.pipeline_id
  user_id           = var.user_id
  operation_delete  = true
  operation_execute = true
  operation_query   = true
  operation_update  = true
}
```

### Modify role permission

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}
variable "role_id" {}

resource "huaweicloud_codearts_pipeline_permission" "role" {
  project_id        = var.codearts_project_id
  pipeline_id       = var.pipeline_id
  role_id           = var.role_id
  operation_delete  = true
  operation_execute = true
  operation_query   = true
  operation_update  = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `project_id` - (Required, String, NonUpdatable) Specifies the CodeArts project ID.

* `pipeline_id` - (Required, String, NonUpdatable) Specifies the pipeline ID.

* `role_id` - (Optional, Int, NonUpdatable) Specifies the role ID.

* `user_id` - (Optional, String, NonUpdatable) Specifies the user ID.

-> Only one of `role_id` and `user_id` can be specified.

* `operation_authorize` - (Optional, Bool) Specifies whether the role has the permission to authorize.
  Default to **fasle**.

* `operation_delete` - (Optional, Bool) Specifies whether the role has the permission to delete. Default to **fasle**.

* `operation_execute` - (Optional, Bool) Specifies whether the role has the permission to execute. Default to **fasle**.

* `operation_query` - (Optional, Bool) Specifies whether the role has the permission to query. Default to **fasle**.

* `operation_update` - (Optional, Bool) Specifies whether the role has the permission to update. Default to **fasle**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `role_name` - Indicates the role name.

* `user_name` - Indicates the user name.

## Import

The pipeline permission can be imported using `project_id`, `pipeline_id`, `role_id` and `user_id`, e.g.

### Import role permission

```bash
$ terraform import huaweicloud_codearts_pipeline_permission.test <project_id>/<pipeline_id>/role/<role_id>
```

### Import user permission

```bash
$ terraform import huaweicloud_codearts_pipeline_permission.test <project_id>/<pipeline_id>/user/<user_id>
```
