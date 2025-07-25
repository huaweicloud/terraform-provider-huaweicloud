---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_user_permissions"
description: |-
  Use this data source to get a list of CodeArts pipeline user permissions.
---

# huaweicloud_codearts_pipeline_user_permissions

Use this data source to get a list of CodeArts pipeline user permissions.

## Example Usage

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}

data "huaweicloud_codearts_pipeline_user_permissions" "test" {
  project_id  = var.codearts_project_id
  pipeline_id = var.pipeline_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `project_id` - (Required, String) Specifies the CodeArts project ID.

* `pipeline_id` - (Required, String) Specifies the pipeline ID.

* `user_name` - (Optional, String) Specifies the user name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `users` - Indicates the template list.
  The [users](#attrblock--users) structure is documented below.

<a name="attrblock--users"></a>
The `users` block supports:

* `user_id` - Indicates the user ID.

* `user_name` - Indicates the user name.

* `operation_authorize` - Indicates whether the user has the permission to authorize.

* `operation_delete` - Indicates whether the user has the permission to delete.

* `operation_execute` - Indicates whether the user has the permission to execute.

* `operation_query` - Indicates whether the user has the permission to query.

* `operation_update` - Indicates whether the user has the permission to update.

* `role_id` - Indicates the role ID.

* `role_name` - Indicates the role name.
