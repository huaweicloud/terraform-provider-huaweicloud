---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_studio_workspace_user"
description: |-
  Manages a workspace user for DataArts Studio Management Center within HuaweiCloud.
---

# huaweicloud_dataarts_studio_workspace_user

Manages a workspace user for DataArts Studio Management Center within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "user_id" {}

resource "huaweicloud_dataarts_studio_workspace_user" "test" {
  workspace_id = var.workspace_id
  user_id      = var.user_id

  roles {
    id = "r00003"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the workspace user is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the user belongs.

* `user_id` - (Required, String, NonUpdatable) Specifies the ID of the IAM user to which the workspace user correspond.

* `roles` - (Required, List) Specifies the role list of the workspace user.  
  The [roles](#dataarts_studio_workspace_user_roles) structure is documented below.

<a name="dataarts_studio_workspace_user_roles"></a>
The `roles` block supports:

* `id` - (Required, String) Specifies the role ID.  
  The valid values are as follows:
  + **r00001**: administrator
  + **r00002**: developer
  + **r00003**: operator
  + **r00004**: viewer

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also same with the parameter `user_id`.

* `roles` - The role list of the workspace user.  
  The [roles](#dataarts_studio_workspace_user_roles_attr) structure is documented below.

* `created_at` - The creation time of the workspace user, in RFC3339 format.

* `updated_at` - The latest update time of the workspace user, in RFC3339 format.

<a name="dataarts_studio_workspace_user_roles_attr"></a>
The `roles` block supports:

* `code` - The role code.

* `name` - The role name.

* `description` - The role description.

## Import

The resource can be imported using `workspace_id` and `user_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dataarts_studio_workspace_user.test <workspace_id>/<user_id>
```
