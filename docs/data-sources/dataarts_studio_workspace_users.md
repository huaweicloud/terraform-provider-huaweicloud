---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_studio_workspace_users"
description: |-
  Use this data source to get the list of workspace users for DataArts Studio Management Center within HuaweiCloud.
---

# huaweicloud_dataarts_studio_workspace_users

Use this data source to get the list of workspace users for DataArts Studio Management Center within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_studio_workspace_users" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the workspace users are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID to which the users belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `users` - The list of workspace users.  
  The [users](#dataarts_studio_workspace_users_attr) structure is documented below.

<a name="dataarts_studio_workspace_users_attr"></a>
The `users` block supports:

* `id` - The ID of the IAM user to which the workspace user correspond.

* `name` - The name of the IAM user to which the workspace user correspond.

* `roles` - The role list of the workspace user.  
  The [roles](#dataarts_studio_workspace_users_roles_attr) structure is documented below.

* `created_at` - The creation time of the workspace user, in RFC3339 format.

* `updated_at` - The latest update time of the workspace user, in RFC3339 format.

<a name="dataarts_studio_workspace_users_roles_attr"></a>
The `roles` block supports:

* `id` - The role ID.

* `code` - The role code.

* `name` - The role name.

* `description` - The role description.
