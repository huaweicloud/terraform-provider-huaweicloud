---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_image_permissions"
description: ""
---

# huaweicloud_swr_image_permissions

Manages a SWR image permissions within HuaweiCloud.

## Example Usage

```hcl
variable "organization_name" {}
variable "repository_name" {}

resource "huaweicloud_swr_image_permissions" "test"{
  organization = var.organization_name
  repository   = var.repository_name

  users {
    user_name  = "test_user1"
    user_id    = "5fc95b3e8fac4bce97e0e0cc8d4a3324"
    permission = "Manage"
  }
  users {
    user_id    = "8854a3426de744d7a5bcb27a171ebfb6"
    permission = "Read"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `organization` - (Required, String, ForceNew) Specifies the name of the organization.

  Changing this parameter will create a new resource.

* `repository` - (Required, String, ForceNew) Specifies the name of the repository.

  Changing this parameter will create a new resource.

* `users` - (Required, List) Specifies the users to access to the image (repository).
The [User](#SwrImagePermissions_User) structure is documented below.

<a name="SwrImagePermissions_User"></a>
The `User` block supports:

* `user_id` - (Required, String) Specifies the ID of the existing HuaweiCloud user.

* `user_name` - (Optional, String) Specifies the name of the existing HuaweiCloud user.

* `permission` - (Required, String) Specifies the user permission of the existing HuaweiCloud user.
  The values can be **Manage**, **Write** and **Read**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `self_permission` - Indicates the permission information of current user.
  The [SelfPermission](#SwrImagePermissions_SelfPermission) structure is documented below.

<a name="SwrImagePermissions_SelfPermission"></a>
The `SelfPermission` block supports:

* `user_id` - Indicates the ID of current user.

* `user_name` - Indicates the name of current user.

* `permission` - Indicates the permission of current user.

## Import

The SWR image permissions can be imported using the organization name and repository name separated by a slash, e.g.:

Only when repository name is with no slashes, can use a slash to separate.

```bash
$ terraform import huaweicloud_swr_image_permissions.test <organization_name>/<repository_name>
```

Using comma to separate is available for repository name with slashes or not.

```bash
$ terraform import huaweicloud_swr_image_permissions.test <organization_name>,<repository_name>
```
